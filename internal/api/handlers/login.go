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
		log.Error(consta.ErrorUnmarshalBody, zap.Error(err))
		c.String(http.StatusInternalServerError, consta.ErrorUnmarshalBody)
		return
	}
	log.Debug("авторизация пользователя", zap.Any("user", user))
	if user.Login == "" || user.Password == "" {
		log.Debug("не валидные логин или пароль", zap.Any("user", user))
		c.String(http.StatusBadRequest, "не валидные логин или пароль")
		return
	}
	authenticationUser, err := storage.Authentication(ctx, user)
	if err != nil {
		log.Error(consta.ErrorWorkDataBase, zap.Error(err))
		c.String(http.StatusInternalServerError, consta.ErrorWorkDataBase)
		return
	}
	if !authenticationUser {
		log.Debug("пароль или логин не верный", zap.Any("user", user))
		c.String(http.StatusUnauthorized, "пароль или логин не верный")
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &models.Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(time.Hour * 100)),
			IssuedAt:  jwt.At(time.Now())},
		Login: user.Login,
	})
	log.Debug("пользователь успешно авторизовался",
		zap.Any("user", user),
		zap.Any("token", token))
	accessToken, err := token.SignedString([]byte(container.GetConfig().SecretKey))
	if err != nil {
		log.Error("ошибка генерация токена", zap.Error(err))
		c.String(http.StatusInternalServerError, "ошибка генерация токена")
		return
	}
	c.Header("Authorization", "Bearer "+accessToken)
}
