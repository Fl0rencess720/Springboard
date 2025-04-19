package oss

import (
	"github.com/Fl0rencess720/Springbroad/internal/controller"
	"github.com/gin-gonic/gin"
)

func InitAPI(group *gin.RouterGroup, ou *controller.OSSUsecase) {
	group.GET("/sts", ou.GetCredentials)
}
