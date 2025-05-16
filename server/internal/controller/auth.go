package controller

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Fl0rencess720/Springboard/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/thedevsaddam/gojsonq"
)

type AuthRepo interface {
	RegisterAppUser(username, password string) error
	VerifyLogin(username, password string) error
}

type AuthUsecase struct {
	repo AuthRepo
}

type AppRegisterLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewAuthUsecase(repo AuthRepo) *AuthUsecase {
	return &AuthUsecase{repo: repo}
}

func (s *AuthUsecase) Login(c *gin.Context) {
	code := c.Query("code")
	url := "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code "
	url = fmt.Sprintf(url, viper.GetString("APP_ID"), viper.GetString("APP_SECRET"), code)
	resp, err := http.Get(url)
	if err != nil {
		ErrorResponse(c, LoginError, err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ErrorResponse(c, LoginError, err)
		return
	}
	json := gojsonq.New().FromString(string(body)).Find("openid")
	if json == nil {
		ErrorResponse(c, LoginError, errors.New("openid not found"))
		return
	}
	openId := json.(string)
	accessToken, refreshToken, err := middleware.GenToken(openId)
	if err != nil {
		ErrorResponse(c, LoginError, err)
		return
	}
	SuccessResponse(c, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (s *AuthUsecase) AppRegister(c *gin.Context) {
	var req AppRegisterLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, ServerError, err)
		return
	}
	accessToken, refreshToken, err := middleware.GenToken(req.Username)
	if err != nil {
		ErrorResponse(c, LoginError, err)
		return
	}
	if err := s.repo.RegisterAppUser(req.Username, req.Password); err != nil {
		ErrorResponse(c, ServerError, err)
		return
	}
	SuccessResponse(c, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (s *AuthUsecase) AppLogin(c *gin.Context) {
	var req AppRegisterLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ErrorResponse(c, ServerError, err)
		return
	}
	accessToken, refreshToken, err := middleware.GenToken(req.Username)
	if err != nil {
		ErrorResponse(c, LoginError, err)
		return
	}
	if err := s.repo.VerifyLogin(req.Username, req.Password); err != nil {
		ErrorResponse(c, LoginError, err)
		return
	}
	SuccessResponse(c, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (s *AuthUsecase) RefreshAccessToken(c *gin.Context) {
	refreshToken := c.Query("refresh_token")
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		ErrorResponse(c, RefreshTokenError, errors.New("miss token string"))
		return
	}
	parts := strings.Split(tokenString, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		ErrorResponse(c, RefreshTokenError, errors.New("wrong token format"))
		return
	}
	accessToken, err := middleware.RefreshToken(parts[1], refreshToken)
	if err != nil {
		ErrorResponse(c, RefreshTokenError, err)
		return
	}
	SuccessResponse(c, gin.H{
		"access_token": accessToken,
	})
}

func MD5(input string) string {
	hash := md5.Sum([]byte(input))
	return hex.EncodeToString(hash[:])
}
