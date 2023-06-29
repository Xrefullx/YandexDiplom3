package middleware

import (
	"fmt"
	"github.com/Xrefullx/YandexDiplom3/internal/api/container"
	"github.com/Xrefullx/YandexDiplom3/internal/models"
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gin-gonic/gin"
	"github.com/zhashkevych/auth/pkg/auth"
	"net/http"
	"strings"
)

func JwtValid() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/api/user/register" || c.Request.URL.Path == "/api/user/login" {
			return
		}
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		header := strings.Split(authHeader, " ")
		if len(header) != 2 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if header[0] != "Bearer" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		login, err := parseToken(header[1],
			[]byte(container.DiContainer.Get("server-config").(models.Config).SecretKey),
		)
		if err != nil {
			status := http.StatusBadRequest
			if err == auth.ErrInvalidAccessToken {
				status = http.StatusUnauthorized
			}

			c.AbortWithStatus(status)
			return
		}
		c.AddParam("loginUser", login)
	}
}

func parseToken(accessToken string, signingKey []byte) (string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return signingKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
		return claims.Login, nil
	}

	return "", auth.ErrInvalidAccessToken
}
