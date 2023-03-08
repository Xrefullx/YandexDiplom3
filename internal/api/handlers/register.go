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
		log.Error(consta.ErrorBody, zap.Error(err))
		c.String(http.StatusInternalServerError, consta.ErrorBody)
		return
	}
	log.Debug("register user", zap.Any("user", user))
	if user.Login == "" || user.Password == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	err := storage.Adduser(ctx, user)
	if err != nil {
		if errors.Is(err, consta.ErrorNoUNIQUE) {
			log.Debug("a user with this username already exists", zap.Any("user", user))
			c.String(http.StatusConflict, "a user with this username already exists")
			return
		}
		log.Error(consta.ErrorDataBase, zap.Error(err), zap.String("func", "Adduser"))
		c.String(http.StatusInternalServerError, consta.ErrorDataBase)
		return
	}
	//<-ctx.Done()
	fmt.Println(errors.Is(ctx.Err(), context.DeadlineExceeded))
	fmt.Println(errors.Is(ctx.Err(), context.Canceled))
	log.Debug("the user has been successfully registered", zap.Any("user", user))
	c.Redirect(http.StatusPermanentRedirect, "/api/user/login")
}
