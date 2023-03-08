package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Xrefullx/YandexDiplom3/internal/api/consta"
	"github.com/Xrefullx/YandexDiplom3/internal/api/container"
	"github.com/Xrefullx/YandexDiplom3/internal/models"
	"github.com/Xrefullx/YandexDiplom3/internal/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
	"time"
)

func AddWithdraw(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), consta.TimeOutRequest)
	defer cancel()
	if !utils.ValidContent(c, "application/json") {
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
	var withdraw models.Withdraw
	err = json.Unmarshal(body, &withdraw)
	if err != nil {
		log.Error(consta.ErrorBody, zap.Error(err))
		c.String(http.StatusInternalServerError, consta.ErrorBody)
		return
	}
	withdraw.ProcessedAT, withdraw.UserLogin = time.Now(), user
	log.Debug("поступил запрос на списание средств",
		zap.Any("withdraw", withdraw),
		zap.String("loginUser", user))
	numberOrder, err := strconv.Atoi(withdraw.NumberOrder)
	if err != nil {
		log.Debug("ошибка преобразования номера заказа", zap.Any("withdraw", withdraw))
		c.String(http.StatusInternalServerError, "ошибка преобразования номера заказа")
		return
	}
	if !utils.LuhValid(numberOrder) {
		log.Debug(consta.ErrorNumberValidLuhn, zap.Error(err), zap.Int("numberOrder", numberOrder))
		c.String(http.StatusUnprocessableEntity, consta.ErrorNumberValidLuhn)
		return
	}
	err = storage.AddWithdraw(ctx, withdraw)
	if err != nil {
		if errors.Is(err, consta.ErrorStatusShortfallAccount) {
			c.String(http.StatusPaymentRequired, consta.ErrorStatusShortfallAccount.Error())
			return
		}
		log.Error(consta.ErrorDataBase, zap.Error(err), zap.String("func", "AddWithdraw"))
		c.String(http.StatusInternalServerError, consta.ErrorDataBase)
		return
	}
	log.Debug("списание совершено", zap.Any("withdraw", withdraw))
	c.String(http.StatusOK, "списание совершено")
}
