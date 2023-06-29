package storage

import (
	"context"
	"github.com/Xrefullx/YandexDiplom3/internal/models"
)

type LoyalityStorage interface {
	Ping() error
	Close() error
	Adduser(ctx context.Context, user models.User) error
	Authentication(ctx context.Context, user models.User) (bool, error)
	GetOrder(ctx context.Context, numberOrder string) (models.Order, error)
	GetOrders(ctx context.Context, userLogin string) ([]models.Order, error)
	AddOrder(ctx context.Context, numberOrder string, order models.Order) error
	UpdateOrder(ctx context.Context, loyaltyPoint models.Loyalty) error
	GetOrdersProcess(ctx context.Context) ([]models.Order, error)
	GetUserBalance(ctx context.Context, userLogin string) (float64, float64, error)
	AddWithdraw(ctx context.Context, withdraw models.Withdraw) error
	GetWithdraws(ctx context.Context, userLogin string) ([]models.Withdraw, error)
}
