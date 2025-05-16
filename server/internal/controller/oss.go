package controller

import (
	"context"

	"github.com/Fl0rencess720/Springboard/pkgs/oss"
	"github.com/gin-gonic/gin"
)

type OSSRepo interface {
}

type OSSUsecase struct {
}

func NewOSSUsecase() *OSSUsecase {
	return &OSSUsecase{}
}

func (uc *OSSUsecase) GetCredentials(c *gin.Context) {
	credentials, err := oss.GenerateAssumeRoleCredential(context.TODO())
	if err != nil {
		ErrorResponse(c, ServerError, err)
		return
	}
	SuccessResponse(c, credentials)
}

func (uc *OSSUsecase) GetPreviewSignedUrl(c *gin.Context) {
	ossKey := c.Query("ossKey")
	previewUrl, err := oss.PresignPreviewUrl(ossKey)
	if err != nil {
		ErrorResponse(c, ServerError, err)
		return
	}
	SuccessResponse(c, gin.H{
		"previewUrl": previewUrl,
	})
}

func (uc *OSSUsecase) GetUploadSignedUrl(c *gin.Context) {
	filename := c.Query("filename")
	contentType := c.DefaultQuery("contentType", "application/octet-stream")
	objectkey := oss.GenerateUniqueKey(filename)
	uploadUrl, err := oss.PresignUploadUrl(objectkey, contentType)
	if err != nil {
		ErrorResponse(c, ServerError, err)
		return
	}
	SuccessResponse(c, gin.H{
		"uploadUrl": uploadUrl,
		"ossKey":    objectkey,
	})
}
