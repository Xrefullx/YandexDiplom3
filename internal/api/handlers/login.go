package handlers

import (
	"context"
	"github.com/Xrefullx/YandexDiplom3/internal/api/consta"
	"github.com/Xrefullx/YandexDiplom3/internal/api/container"
	"github.com/Xrefullx/YandexDiplom3/internal/models"
	"github.com/Xrefullx/YandexDiplom3/internal/utils"
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func Login(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), consta.TimeOutRequest)
	defer cancel()
	if !utils.ValidContent(c, "application/json") {
		return
	}
	log := container.GetLog()
	storage := container.GetStorage()
	var user models.User
	if err := c.Bind(&user); err != nil {
		log.Error(consta.ErrorBody, zap.Error(err))
		c.String(http.StatusInternalServerError, consta.ErrorBody)
		return
	}
	log.Debug("user authorization", zap.Any("user", user))
	if user.Login == "" || user.Password == "" {
		log.Debug("invalid username or password", zap.Any("user", user))
		c.String(http.StatusBadRequest, "invalid username or password")
		return
	}
	authenticationUser, err := storage.Authentication(ctx, user)
	if err != nil {
		log.Error(consta.ErrorDataBase, zap.Error(err))
		c.String(http.StatusInternalServerError, consta.ErrorDataBase)
		return
	}
	if !authenticationUser {
		log.Debug("the password or login is not correct", zap.Any("user", user))
		c.String(http.StatusUnauthorized, "the password or login is not correct")
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &models.Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(time.Hour * 100)),
			IssuedAt:  jwt.At(time.Now())},
		Login: user.Login,
	})
	log.Debug("the user has successfully logged in",
		zap.Any("user", user),
		zap.Any("token", token))
	accessToken, err := token.SignedString([]byte(container.GetConfig().SecretKey))
	if err != nil {
		log.Error("error token", zap.Error(err))
		c.String(http.StatusInternalServerError, "error token")
		return
	}
	c.Header("Authorization", "Bearer "+accessToken)
}
