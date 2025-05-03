package oss

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"time"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	sts20150401 "github.com/alibabacloud-go/sts-20150401/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"go.uber.org/zap"
)

type Credentials struct {
	AccessKeyId     string `json:"AccessKeyId"`
	AccessKeySecret string `json:"AccessKeySecret"`
	SecurityToken   string `json:"SecurityToken"`
	Expiration      string `json:"Expiration"`
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
		RoleSessionName: tea.String("springboard"),
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
		Expiration:      tea.StringValue(credentials.Expiration),
	}, nil
}

func GeneratePolicyAndSignature(accessKeyID, accessKeySecret, securityToken string) (string, string, error) {
	expiration := time.Now().Add(30 * time.Minute).UTC().Format("2006-01-02T15:04:05.000Z")
	policy := map[string]interface{}{
		"expiration": expiration,
		"conditions": []interface{}{
			map[string]string{"bucket": "springboard"},
			map[string]string{"x-oss-security-token": securityToken},
		},
	}

	policyJSON, err := json.Marshal(policy)
	if err != nil {
		zap.L().Error("Failed to marshal policy", zap.Error(err))
		return "", "", err
	}
	base64Policy := base64.StdEncoding.EncodeToString(policyJSON)

	mac := hmac.New(sha1.New, []byte(accessKeySecret))
	mac.Write([]byte(base64Policy))
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	return base64Policy, signature, nil
}

func GenetrateDownloadSignedURL(creds Credentials, ossKey string) (string, error) {

	ossEndpoint := "https://oss-cn-shenzhen.aliyuncs.com"

	client, err := oss.New(ossEndpoint, creds.AccessKeyId, creds.AccessKeySecret, oss.SecurityToken(creds.SecurityToken))
	if err != nil {
		return "", fmt.Errorf("创建 OSS 客户端失败: %v", err)
	}

	bucketHandle, err := client.Bucket("springboard")
	if err != nil {
		return "", fmt.Errorf("获取 OSS Bucket 失败: %v", err)
	}

	signedURL, err := bucketHandle.SignURL(ossKey, oss.HTTPGet, 600)
	if err != nil {
		return "", fmt.Errorf("生成签名 URL 失败: %v", err)
	}
	return signedURL, nil
}

func GenetratePreviewSignedURL(creds Credentials, ossKey string) (string, error) {

	ossEndpoint := "https://oss-cn-shenzhen.aliyuncs.com"

	client, err := oss.New(ossEndpoint, creds.AccessKeyId, creds.AccessKeySecret, oss.SecurityToken(creds.SecurityToken))
	if err != nil {
		return "", fmt.Errorf("创建 OSS 客户端失败: %v", err)
	}

	bucketHandle, err := client.Bucket("springboard")
	if err != nil {
		return "", fmt.Errorf("获取 OSS Bucket 失败: %v", err)
	}

	options := []oss.Option{
		oss.ResponseContentDisposition("inline"),
	}

	signedURL, err := bucketHandle.SignURL(ossKey, oss.HTTPGet, 600, options...)
	if err != nil {
		return "", fmt.Errorf("生成签名 URL 失败: %v", err)
	}
	return signedURL, nil
}
