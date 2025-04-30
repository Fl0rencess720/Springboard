package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Fl0rencess720/Springboard/internal/conf"
	"github.com/Fl0rencess720/Springboard/internal/controller"
	"github.com/Fl0rencess720/Springboard/internal/data"

	"github.com/Fl0rencess720/Springboard/api"
	"github.com/Fl0rencess720/Springboard/consts"
	"github.com/Fl0rencess720/Springboard/pkgs/logger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func init() {
	conf.Init()
	logger.Init(consts.DefaultLogFilePath)
	data.Init()
}

func main() {
	srv := newSrv()
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			zap.L().Error("Server ListenAndServe", zap.Error(err))
			panic(err)
		}
	}()
	closeServer(srv, context.Background())
}

func newSrv() *http.Server {
	authRepo := data.NewAuthRepo(data.GetDB())
	portfolioRepo := data.NewPortfolioRepo(data.GetDB(), data.GetRedis())
	feedbackRepo := data.NewFeedbackRepo(data.GetDB())
	authUsecase := controller.NewAuthUsecase(authRepo)
	portfolioUsecase := controller.NewPortfolioUsecase(portfolioRepo)
	feedbackUsecase := controller.NewFeedbackUseCase(feedbackRepo)
	ossUsecase := controller.NewOSSUsecase()
	return &http.Server{
		Addr:    viper.GetString("server.port"),
		Handler: api.Init(authUsecase, portfolioUsecase, feedbackUsecase, ossUsecase),
	}
}

func closeServer(srv *http.Server, ctx context.Context) {
	defer func(l *zap.Logger) {
		logger.Sync(l)
	}(zap.L())
	srv.RegisterOnShutdown(data.Close)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Error("Server Shutdown", zap.Error(err))
	}
	zap.L().Info("Server exited")

}
