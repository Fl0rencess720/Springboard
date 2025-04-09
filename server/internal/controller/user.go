package controller

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Fl0rencess720/Springbroad/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/thedevsaddam/gojsonq"
	"go.uber.org/zap"
)

func Login(c *gin.Context) {
	code := c.Query("code")
	url := "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code "
	url = fmt.Sprintf(url, viper.GetString("APP_ID"), viper.GetString("APP_SECRET"), code)
	resp, err := http.Get(url)
	if err != nil {
		zap.L().Error("get wechat openid error", zap.Error(err))
		c.JSON(401, gin.H{
			"code":    401,
			"message": "login failed",
		})
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		zap.L().Error("read wechat openid error", zap.Error(err))
		c.JSON(401, gin.H{
			"code":    401,
			"message": "login failed",
		})
		return
	}
	json := gojsonq.New().FromString(string(body)).Find("openid")
	openId := json.(string)
	accessToken, refreshToken, err := middleware.GenToken(openId)
	if err != nil {
		zap.L().Error("gen token error", zap.Error(err))
		c.JSON(401, gin.H{
			"code":    401,
			"message": "login failed",
		})
		return
	}
	c.JSON(200, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func RefreshAccessToken(c *gin.Context) {
	refreshToken := c.Query("refresh_token")
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(401, gin.H{
			"code":    401,
			"message": "miss token string",
		})
		return
	}
	parts := strings.Split(tokenString, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(401, gin.H{
			"code":    401,
			"message": "wrong token format",
		})
		return
	}
	accessToken, err := middleware.RefreshToken(parts[1], refreshToken)
	if err != nil {
		zap.L().Error("refresh token error", zap.Error(err))
		c.JSON(401, gin.H{
			"code":    401,
			"message": "refresh token error",
		})
		return
	}
	c.JSON(200, gin.H{
		"access_token": accessToken,
	})
}
