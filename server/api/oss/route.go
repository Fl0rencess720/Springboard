package oss

import (
	"github.com/Fl0rencess720/Springboard/internal/controller"
	"github.com/gin-gonic/gin"
)

func InitAPI(group *gin.RouterGroup, ou *controller.OSSUsecase) {
	// group.GET("/sts", ou.GetCredentials)
	group.GET("/sts/upload", ou.GetUploadSignedUrl)
	group.GET("/sts/preview", ou.GetPreviewSignedUrl)
}
