package controller

import (
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
	credentials, err := oss.GenerateAssumeRoleCredential()
	if err != nil {
		ErrorResponse(c, ServerError, err)
		return
	}
	policy, signature, err := oss.GeneratePolicyAndSignature(credentials.AccessKeyId, credentials.AccessKeySecret, credentials.SecurityToken)
	if err != nil {
		ErrorResponse(c, ServerError, err)
		return
	}
	SuccessResponse(c, gin.H{
		"accessKeyId":   credentials.AccessKeyId,
		"policy":        policy,
		"signature":     signature,
		"securityToken": credentials.SecurityToken,
		"expiration":    credentials.Expiration,
		"region":        "oss-cn-shenzhen",
		"bucket":        "springboard",
	})
}
