package oss

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	sts20150401 "github.com/alibabacloud-go/sts-20150401/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"go.uber.org/zap"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
)

// type Credentials struct {
// 	AccessKeyId     string `json:"AccessKeyId"`
// 	AccessKeySecret string `json:"AccessKeySecret"`
// 	SecurityToken   string `json:"SecurityToken"`
// }

var (
	region     string = "cn-shenzhen" // SDK 会在前面添加 "oss-" 前缀
	bucketName string = "springboard"
	cfg        *oss.Config
)

func init() {
	// SDK 会自动调用传入的函数刷新 credential
	cfg = oss.LoadDefaultConfig().
		WithCredentialsProvider(
			credentials.NewCredentialsFetcherProvider(
				credentials.CredentialsFetcherFunc(GenerateAssumeRoleCredential),
			),
		).
		WithRegion(region)
}

func GenerateAssumeRoleCredential(ctx context.Context) (credentials.Credentials, error) {
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
		return credentials.Credentials{}, fmt.Errorf("failed to create STS client: %w", err)
	}

	request := &sts20150401.AssumeRoleRequest{
		DurationSeconds: tea.Int64(3600),
		RoleArn:         tea.String(roleArn),
		RoleSessionName: tea.String("springboard"),
	}
	response, err := client.AssumeRoleWithOptions(request, &util.RuntimeOptions{})
	if err != nil {
		zap.L().Error("Failed to assume role", zap.Error(err))
		return credentials.Credentials{}, err
	}
	stsRespCredentials := response.Body.Credentials

	expirationStr := tea.StringValue(stsRespCredentials.Expiration)
	if expirationStr == "" {
		zap.L().Error("Empty Expiration string from STS", zap.Error(err))
		return credentials.Credentials{}, fmt.Errorf("STS response contained empty Expiration string")
	}

	expiresTime, parseErr := time.Parse(time.RFC3339, expirationStr)
	if parseErr != nil {
		zap.L().Error("Failed to parse expiration time from STS", zap.String("expiration", expirationStr), zap.Error(parseErr))
		return credentials.Credentials{}, fmt.Errorf("failed to parse expiration time '%s': %w", expirationStr, parseErr)
	}

	return credentials.Credentials{
		AccessKeyID:     tea.StringValue(stsRespCredentials.AccessKeyId),
		AccessKeySecret: tea.StringValue(stsRespCredentials.AccessKeySecret),
		SecurityToken:   tea.StringValue(stsRespCredentials.SecurityToken),
		Expires:         &expiresTime,
	}, nil
}

func PresignPreviewUrl(objectkey string) (string, error) {
	client := oss.NewClient(cfg)

	log.Println("[PresignPreview Url] objectkey: ", objectkey)

	result, err := client.Presign(context.TODO(), &oss.GetObjectRequest{
		Bucket: oss.Ptr(bucketName),
		Key:    oss.Ptr(objectkey),
	}, oss.PresignExpires(30*time.Minute))

	if err != nil {
		zap.L().Error("failed to get object presign: ", zap.Error(err))
		return "", fmt.Errorf("failed to get object presign: %w", err)
	}

	return result.URL, nil
}

func PresignUploadUrl(objectkey string, contentType string) (string, error) {
	client := oss.NewClient(cfg)

	result, err := client.Presign(context.TODO(), &oss.PutObjectRequest{
		Bucket:      oss.Ptr(bucketName),
		Key:         oss.Ptr(objectkey),
		ContentType: oss.Ptr(contentType),
	}, oss.PresignExpires(30*time.Minute))

	if err != nil {
		zap.L().Error("failed to put object presign: ", zap.Error(err))
		return "", fmt.Errorf("failed to put object presign: %w", err)
	}

	if len(result.SignedHeaders) > 0 {
		log.Printf("signed headers:\n")
		for k, v := range result.SignedHeaders {
			log.Printf("%v: %v\n", k, v)
		}
	}

	return result.URL, nil
}

func GenerateUniqueKey(filename string) string {
	timestamp := time.Now().Unix()
	return fmt.Sprintf("uploads/%d_%s", timestamp, filename)
}
