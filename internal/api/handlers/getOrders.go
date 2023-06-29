package handlers

import (
	"context"
	"github.com/Xrefullx/YandexDiplom3/internal/api/consta"
	"github.com/Xrefullx/YandexDiplom3/internal/api/container"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func GetUserOrders(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), consta.TimeOutRequest)
	defer cancel()
	log := container.GetLog()
	storage := container.GetStorage()
	user := c.Param("loginUser")
	log.Debug("поступил запрос на показ заказов",
		zap.String("loginUser", user))

	orders, err := storage.GetOrders(ctx, user)
	if err != nil {
		log.Error(consta.ErrorWorkDataBase, zap.Error(err), zap.String("func", "GetManyOrders"))
		c.String(http.StatusInternalServerError, consta.ErrorWorkDataBase)
		return
	}
	if len(orders) == 0 {
		log.Debug("нет данных для ответа", zap.String("loginUser", user))
		c.String(http.StatusNoContent, "нет данных для ответа")
		return
	}
	c.JSON(http.StatusOK, orders)
}
