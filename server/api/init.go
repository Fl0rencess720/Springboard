package api

import (
	"time"

	"github.com/Fl0rencess720/Springbroad/api/user"
	"github.com/Fl0rencess720/Springbroad/internal/controller"
	"github.com/Fl0rencess720/Springbroad/internal/middleware"
	ginZap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Init() *gin.Engine {
	e := gin.New()
	e.Use(gin.Logger(), gin.Recovery(), ginZap.Ginzap(zap.L(), time.RFC3339, false), ginZap.RecoveryWithZap(zap.L(), false))
	basic := e.Group("/api")
	{
		basic.POST("/login", controller.Login)
		basic.GET("/refresh", controller.RefreshAccessToken)
	}
	app := e.Group("/api", middleware.Auth())
	{
		user.InitAPI(app.Group("/user"))
	}
	return e
}
