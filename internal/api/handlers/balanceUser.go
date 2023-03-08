package handlers

import (
	"context"
	"github.com/Xrefullx/YandexDiplom3/internal/api/consta"
	"github.com/Xrefullx/YandexDiplom3/internal/api/container"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func UserBalance(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), consta.TimeOutRequest)
	defer cancel()
	log := container.GetLog()
	storage := container.GetStorage()
	user := c.Param("loginUser")
	log.Debug("поступил запрос на проверку баланса",
		zap.String("loginUser", user))

	sum, spent, err := storage.GetUserBalance(ctx, user)
	if err != nil {
		log.Error(consta.ErrorDataBase, zap.Error(err), zap.String("func", "GetUserBalance"))
		c.String(http.StatusInternalServerError, consta.ErrorDataBase)
		return
	}
	response := map[string]float64{
		"current":   sum - spent,
		"withdrawn": spent,
	}
	log.Debug("баланс пользователя", zap.String("loginUser", user),
		zap.Float64("sum", sum),
		zap.Float64("spent", spent),
		zap.Float64("current", sum-spent),
	)
	c.JSONP(http.StatusOK, response)
}
