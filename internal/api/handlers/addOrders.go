package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Xrefullx/YandexDiplom3/internal/api/consta"
	"github.com/Xrefullx/YandexDiplom3/internal/api/container"
	"github.com/Xrefullx/YandexDiplom3/internal/models"
	"github.com/Xrefullx/YandexDiplom3/internal/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"net/http"
	"time"
)

func AddOrders(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), consta.TimeOutRequest)
	defer cancel()
	if !utils.ValidContent(c, "text/plain") {
		return
	}
	log := container.GetLog()
	storage := container.GetStorage()
	user := c.Param("loginUser")
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Error(consta.ErrorReadBody, zap.Error(err))
		c.String(http.StatusInternalServerError, consta.ErrorReadBody)
		return
	}
	var numberOrder int
	err = json.Unmarshal(body, &numberOrder)
	if err != nil {
		log.Error(consta.ErrorBody, zap.Error(err))
		c.String(http.StatusInternalServerError, consta.ErrorBody)
		return
	}
	log.Debug("the order number has been received",
		zap.Int("numberOrder", numberOrder),
		zap.String("loginUser", user))
	if !utils.LuhValid(numberOrder) {
		log.Debug(consta.ErrorNumberValidLuhn, zap.Error(err), zap.Int("numberOrder", numberOrder))
		c.String(http.StatusUnprocessableEntity, consta.ErrorNumberValidLuhn)
		return
	}
	numberOrderStr := fmt.Sprint(numberOrder)
	err = storage.AddOrder(ctx, numberOrderStr,
		models.Order{
			NumberOrder: numberOrderStr,
			UserLogin:   user,
			Status:      consta.OrderStatusNEW,
			Uploaded:    time.Now(),
		})
	if err != nil {
		if errors.Is(err, consta.ErrorNoUNIQUE) {
			order, errGet := storage.GetOrder(ctx, numberOrderStr)
			if errGet != nil {
				log.Error(consta.ErrorDataBase, zap.Error(errGet), zap.String("func", "GetOrder"))
				c.String(http.StatusInternalServerError, consta.ErrorDataBase)
				return
			}
			if order.UserLogin == user {
				log.Debug("the order number has already been uploaded by this user", zap.Any("order", order))
				c.String(http.StatusOK, "the order number has already been uploaded by this user")
				return
			}
			log.Debug("the order number has already been uploaded by another user", zap.Any("order", order))
			c.String(http.StatusConflict, "the order number has already been uploaded by another user")
			return
		}
		log.Error(consta.ErrorDataBase, zap.Error(err), zap.String("func", "AddOrder"))
		c.String(http.StatusInternalServerError, consta.ErrorDataBase)
		return
	}
	log.Debug("the new order number has been accepted for processing", zap.Any("number_order", numberOrder))
	c.String(http.StatusAccepted, "the new order number has been accepted for processing")
}
