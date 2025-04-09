package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	ServerError = iota
	AuthError
	TokenExpired
	LoginError
	RefreshTokenError
)

var HttpCode = map[uint]int{
	ServerError:       502,
	AuthError:         401,
	TokenExpired:      401,
	LoginError:        403,
	RefreshTokenError: 403,
}

var Message = map[uint]string{
	ServerError:       "服务器错误",
	AuthError:         "无权访问",
	TokenExpired:      "Token过期",
	LoginError:        "登录失败",
	RefreshTokenError: "刷新Token失败",
}

func SuccessResponse(c *gin.Context, data any) {
	c.JSON(200, gin.H{
		"msg":  "success",
		"code": 200,
		"data": data,
	})
}

func ErrorResponse(c *gin.Context, code uint, data ...any) {
	httpStatus, ok := HttpCode[code]
	if !ok {
		httpStatus = 403
	}
	msg, ok := Message[code]
	if !ok {
		msg = "未知错误"
	}
	zap.L().Error("error response", zap.Uint("code", code), zap.String("openid", c.GetString("openid")), zap.Any(msg, data))

	c.JSON(httpStatus, gin.H{
		"code": code,
		"msg":  msg,
	})
}
