package handlers

import (
	"github.com/Xrefullx/YandexDiplom3/internal/api/middleware"
	"github.com/Xrefullx/YandexDiplom3/internal/models"

	"github.com/gin-gonic/gin"
)

func Router(cfg models.Config) *gin.Engine {
	if cfg.ReleaseMOD {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(middleware.JwtValid())

	gUser := r.Group("/api/user")
	{
		gUser.POST("/register", Register)
		gUser.POST("/login", Login)
		gUser.POST("/orders", AddUserOrders)
		gUser.GET("/orders", GetUserOrders)
		gUser.GET("/balance", UserBalance)
		gUser.POST("/balance/withdraw", AddWithdraw)
		gUser.GET("/withdrawals", GetUserWithdraws)
	}
	return r
}
