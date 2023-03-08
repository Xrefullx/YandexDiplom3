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

type Response struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

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
		c.JSON(http.StatusInternalServerError, Response{Message: consta.ErrorReadBody, Status: http.StatusInternalServerError})
		return
	}
	var withdraw models.Withdraw
	err = json.Unmarshal(body, &withdraw)
	if err != nil {
		log.Error(consta.ErrorBody, zap.Error(err))
		c.JSON(http.StatusInternalServerError, Response{Message: consta.ErrorBody, Status: http.StatusInternalServerError})
		return
	}
	withdraw.ProcessedAT, withdraw.UserLogin = time.Now(), user
	log.Debug("a request has been received for debiting funds",
		zap.Any("withdraw", withdraw),
		zap.String("loginUser", user))
	numberOrder, err := strconv.Atoi(withdraw.NumberOrder)
	if err != nil {
		log.Debug("order number conversion error", zap.Any("withdraw", withdraw))
		c.String(http.StatusInternalServerError, "order number conversion error")
		return
	}
	if !utils.LuhValid(numberOrder) {
		log.Debug(consta.ErrorNumberValidLuhn, zap.Error(err), zap.Int("numberOrder", numberOrder))
		c.JSON(http.StatusUnprocessableEntity, Response{Message: consta.ErrorNumberValidLuhn, Status: http.StatusUnprocessableEntity})
		return
	}
	err = storage.AddWithdraw(ctx, withdraw)
	if err != nil {
		if errors.Is(err, consta.ErrorStatusShortfallAccount) {
			c.JSON(http.StatusPaymentRequired, Response{Message: consta.ErrorStatusShortfallAccount.Error(), Status: http.StatusPaymentRequired})
			return
		}
		log.Error(consta.ErrorDataBase, zap.Error(err), zap.String("func", "AddWithdraw"))
		c.JSON(http.StatusInternalServerError, Response{Message: consta.ErrorDataBase, Status: http.StatusInternalServerError})
		return
	}
	log.Debug("the write-off has been completed", zap.Any("withdraw", withdraw))
	c.JSON(http.StatusOK, Response{Message: "the write-off has been completed", Status: http.StatusOK})
}
