package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/Xrefullx/YandexDiplom3/internal/api/consta"
	"github.com/Xrefullx/YandexDiplom3/internal/api/container"
	"github.com/Xrefullx/YandexDiplom3/internal/models"
	"github.com/Xrefullx/YandexDiplom3/internal/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func Register(c *gin.Context) {
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
	log.Debug("регистрация пользователя", zap.Any("user", user))
	if user.Login == "" || user.Password == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	err := storage.Adduser(ctx, user)
	if err != nil {
		if errors.Is(err, consta.ErrorNoUNIQUE) {
			log.Debug("пользователь с таким логином уже есть", zap.Any("user", user))
			c.String(http.StatusConflict, "пользователь с таким логином уже есть")
			return
		}
		log.Error(consta.ErrorWorkDataBase, zap.Error(err), zap.String("func", "AddUser"))
		c.String(http.StatusInternalServerError, consta.ErrorWorkDataBase)
		return
	}
	//<-ctx.Done()
	fmt.Println(errors.Is(ctx.Err(), context.DeadlineExceeded))
	fmt.Println(errors.Is(ctx.Err(), context.Canceled))
	log.Debug("пользователь успешно зарегистрирован", zap.Any("user", user))
	c.Redirect(http.StatusPermanentRedirect, "/api/user/login")
}
