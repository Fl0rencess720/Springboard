package middleware

import (
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
)

type ContextKey string

var (
	OpenidKey = ContextKey("openid")
)

type AuthClaims struct {
	Openid string `json:"openid"`
	jwt.RegisteredClaims
}

func GenAccessToken(openid string) (string, error) {
	ac := AuthClaims{
		Openid: openid,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        time.Now().String(),
			Issuer:    "Springboard",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
		},
	}
	accessSecret := viper.GetString("ACCESS_SECRET")
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, ac).SignedString([]byte(accessSecret))
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

func GenRefreshToken() (string, error) {
	rc := jwt.RegisteredClaims{
		ID:        time.Now().String(),
		Issuer:    "Springboard",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(14 * 24 * time.Hour)),
	}
	refreshSecret := viper.GetString("REFRESH_SECRET")
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, rc).SignedString([]byte(refreshSecret))
	if err != nil {
		return "", err
	}
	return refreshToken, nil
}

func GenToken(phone string) (string, string, error) {
	accessToken, err := GenAccessToken(phone)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := GenRefreshToken()
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func ParseToken(aToken string) (*AuthClaims, bool, error) {
	accessSecret := viper.GetString("ACCESS_SECRET")
	accessToken, err := jwt.ParseWithClaims(aToken, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(accessSecret), nil
	})
	if err != nil {
		return nil, false, err
	}
	if claims, ok := accessToken.Claims.(*AuthClaims); ok && accessToken.Valid {
		return claims, false, nil
	}
	return nil, true, errors.New("invalid token")
}

func RefreshToken(aToken, rToken string) (string, error) {
	accessSecret := viper.GetString("ACCESS_SECRET")
	refreshSecret := viper.GetString("REFRESH_SECRET")
	rToken = strings.TrimPrefix(rToken, "Bearer ")
	_, err := jwt.Parse(rToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(refreshSecret), nil
	})
	if err != nil {
		return "", err
	}
	var claims AuthClaims
	_, err = jwt.ParseWithClaims(aToken, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(accessSecret), nil
	})
	v, _ := err.(*jwt.ValidationError)
	if v == nil || v.Errors == jwt.ValidationErrorExpired {
		return GenAccessToken(claims.Openid)
	}
	return "", err
}

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(401, gin.H{
				"code":    401,
				"message": "miss token string",
			})
			return
		}
		parts := strings.Split(tokenString, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(401, gin.H{
				"code":    401,
				"message": "wrong token format",
			})
			return
		}
		parsedToken, isExpire, err := ParseToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{
				"code":    401,
				"message": "invalid token",
			})
			return
		}
		if isExpire {
			c.AbortWithStatusJSON(401, gin.H{
				"code":    401,
				"message": "token expired",
			})
			return
		}
		c.Set(string(OpenidKey), parsedToken.Openid)
		c.Next()
	}
}
