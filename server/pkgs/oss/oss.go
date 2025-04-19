package oss

import (
	"os"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	sts20150401 "github.com/alibabacloud-go/sts-20150401/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"go.uber.org/zap"
)

type Credentials struct {
	AccessKeyId     string `json:"AccessKeyId"`
	AccessKeySecret string `json:"AccessKeySecret"`
	SecurityToken   string `json:"SecurityToken"`
}

func GenerateAssumeRoleCredential() (Credentials, error) {
	accessKeyId := os.Getenv("OSSAccessKeyId")
	accessKeySecret := os.Getenv("OSSAccessKeySecret")
	roleArn := os.Getenv("OSSRoleArn")
	config := &openapi.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
	}
	// Endpoint 请参考 https://api.aliyun.com/product/Sts
	config.Endpoint = tea.String("sts.cn-hangzhou.aliyuncs.com")
	client, err := sts20150401.NewClient(config)
	if err != nil {
		zap.L().Error("Failed to create STS client", zap.Error(err))
		return Credentials{}, err
	}

	request := &sts20150401.AssumeRoleRequest{
		DurationSeconds: tea.Int64(3600),
		RoleArn:         tea.String(roleArn),
		RoleSessionName: tea.String("springbroad"),
	}
	response, err := client.AssumeRoleWithOptions(request, &util.RuntimeOptions{})
	if err != nil {
		zap.L().Error("Failed to assume role", zap.Error(err))
		return Credentials{}, err
	}
	credentials := response.Body.Credentials
	return Credentials{
		AccessKeyId:     tea.StringValue(credentials.AccessKeyId),
		AccessKeySecret: tea.StringValue(credentials.AccessKeySecret),
		SecurityToken:   tea.StringValue(credentials.SecurityToken),
	}, nil
}
