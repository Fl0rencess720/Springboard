package api

import (
	"time"

	"github.com/Fl0rencess720/Springbroad/api/portfolio"
	"github.com/Fl0rencess720/Springbroad/api/user"
	"github.com/Fl0rencess720/Springbroad/internal/controller"
	"github.com/Fl0rencess720/Springbroad/internal/middleware"
	ginZap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Init(au *controller.AuthUsecase, pu *controller.PortfolioUsecase) *gin.Engine {
	e := gin.New()
	e.Use(gin.Logger(), gin.Recovery(), ginZap.Ginzap(zap.L(), time.RFC3339, false), ginZap.RecoveryWithZap(zap.L(), false))
	auth := e.Group("/api")
	{
		auth.POST("/login", au.Login)
		auth.GET("/refresh", au.RefreshAccessToken)
	}

	app := e.Group("/api", middleware.Auth())
	{
		user.InitAPI(app.Group("/user"))
		portfolio.InitAPI(app.Group("/portfolio"), pu)
	}
	return e
}
