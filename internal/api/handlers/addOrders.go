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

func AddUserOrders(c *gin.Context) {
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
		log.Error(consta.ErrorUnmarshalBody, zap.Error(err))
		c.String(http.StatusInternalServerError, consta.ErrorUnmarshalBody)
		return
	}
	log.Debug("поступил номер заказа",
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
				log.Error(consta.ErrorWorkDataBase, zap.Error(errGet), zap.String("func", "GetOrder"))
				c.String(http.StatusInternalServerError, consta.ErrorWorkDataBase)
				return
			}
			if order.UserLogin == user {
				log.Debug("номер заказа уже был загружен этим пользователем", zap.Any("order", order))
				c.String(http.StatusOK, "номер заказа уже был загружен этим пользователем")
				return
			}
			log.Debug("номер заказа уже был загружен другим пользователем", zap.Any("order", order))
			c.String(http.StatusConflict, "номер заказа уже был загружен другим пользователем")
			return
		}
		log.Error(consta.ErrorWorkDataBase, zap.Error(err), zap.String("func", "AddOrder"))
		c.String(http.StatusInternalServerError, consta.ErrorWorkDataBase)
		return
	}
	log.Debug("новый номер заказа принят в обработку", zap.Any("number_order", numberOrder))
	c.String(http.StatusAccepted, "новый номер заказа принят в обработку")
}
