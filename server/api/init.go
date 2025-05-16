package api

import (
	"time"

	"github.com/Fl0rencess720/Springboard/api/feedback"
	"github.com/Fl0rencess720/Springboard/api/oss"
	"github.com/Fl0rencess720/Springboard/api/portfolio"
	"github.com/Fl0rencess720/Springboard/internal/controller"
	"github.com/Fl0rencess720/Springboard/internal/middleware"

	ginZap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Init(au *controller.AuthUsecase, pu *controller.PortfolioUsecase, sc *controller.FeedbackUseCase, ou *controller.OSSUsecase) *gin.Engine {
	e := gin.New()
	e.Use(gin.Logger(), gin.Recovery(), ginZap.Ginzap(zap.L(), time.RFC3339, false), ginZap.RecoveryWithZap(zap.L(), false))
	auth := e.Group("/api")
	{
		auth.POST("/login", au.Login)

		auth.POST("/register/app", au.AppRegister)
		auth.POST("/login/app", au.AppLogin)

		auth.GET("/refresh", au.RefreshAccessToken)
	}

	app := e.Group("/api", middleware.Cors(), middleware.Auth())
	{
		oss.InitAPI(app.Group("/oss"), ou)
		portfolio.InitAPI(app.Group("/portfolio"), pu)
		feedback.InitAPI(app.Group("/feedback"), sc)
	}

	return e
}
