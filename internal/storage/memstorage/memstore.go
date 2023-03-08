package memstorage

import (
	"context"
	"github.com/Xrefullx/YandexDiplom3/internal/api/consta"
	"github.com/Xrefullx/YandexDiplom3/internal/models"

	"sync"

	"github.com/google/uuid"
)

// add
type MemStorage struct {
	userCash     map[uuid.UUID]models.User
	orderCash    map[string]models.Order
	withdrawCash map[uuid.UUID]models.Withdraw
	mu           *sync.RWMutex
}

func New() (*MemStorage, error) {
	return &MemStorage{
		userCash:     make(map[uuid.UUID]models.User),
		orderCash:    make(map[string]models.Order),
		withdrawCash: make(map[uuid.UUID]models.Withdraw),
		mu:           new(sync.RWMutex),
	}, nil
}

func (MS *MemStorage) Ping() error {
	return nil
}

func (MS *MemStorage) Close() error {
	return nil
}

func (MS *MemStorage) Adduser(ctx context.Context, user models.User) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	MS.mu.Lock()
	defer MS.mu.Unlock()
	for _, v := range MS.userCash {
		if v.Login == user.Login {
			return consta.ErrorNoUNIQUE
		}
	}
	MS.userCash[uuid.New()] = user
	select {
	default:
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}

func (MS *MemStorage) Authentication(_ context.Context, user models.User) (bool, error) {
	MS.mu.RLock()
	defer MS.mu.RUnlock()
	for _, v := range MS.userCash {
		if v.Login == user.Login && v.Password == user.Password {
			return true, nil
		}
	}
	return false, nil
}

func (MS *MemStorage) GetOrder(_ context.Context, numberOrder string) (models.Order, error) {
	MS.mu.RLock()
	defer MS.mu.RUnlock()

	return MS.orderCash[numberOrder], nil
}

func (MS *MemStorage) GetOrders(_ context.Context, userLogin string) ([]models.Order, error) {
	MS.mu.RLock()
	defer MS.mu.RUnlock()
	var orders []models.Order
	for _, v := range MS.orderCash {
		if v.UserLogin == userLogin {
			orders = append(orders, v)
		}
	}
	return orders, nil
}

func (MS *MemStorage) AddOrder(_ context.Context, numberOrder string, order models.Order) error {
	MS.mu.Lock()
	defer MS.mu.Unlock()
	if MS.orderCash[numberOrder].NumberOrder == numberOrder {
		return consta.ErrorNoUNIQUE
	}
	MS.orderCash[numberOrder] = order
	return nil
}

func (MS *MemStorage) GetOrdersProcess(_ context.Context) ([]models.Order, error) {
	MS.mu.RLock()
	defer MS.mu.RUnlock()
	var orders []models.Order
	for _, v := range MS.orderCash {
		if v.Status == consta.OrderStatusPROCESSING ||
			v.Status == consta.OrderStatusNEW ||
			v.Status == consta.OrderStatusREGISTERED ||
			v.Status == consta.OrderStatusInvalid {
			orders = append(orders, v)
		}
	}
	return orders, nil
}

func (MS *MemStorage) UpdateOrder(_ context.Context, Loyality models.Loyalty) error {
	MS.mu.Lock()
	defer MS.mu.Unlock()
	order := MS.orderCash[Loyality.NumberOrder]
	order.Status, order.Accrual = Loyality.Status, Loyality.Accrual
	MS.orderCash[Loyality.NumberOrder] = order
	return nil
}

func (MS *MemStorage) GetUserBalance(_ context.Context, userLogin string) (float64, float64, error) {
	MS.mu.RLock()
	defer MS.mu.RUnlock()
	pointsSUM := 0.0
	pointsSPEND := 0.0
	for _, order := range MS.orderCash {
		if order.UserLogin == userLogin {
			pointsSUM += order.Accrual
		}
	}
	for _, withdraw := range MS.withdrawCash {
		if withdraw.UserLogin == userLogin {
			pointsSPEND += withdraw.Sum
		}
	}
	return pointsSUM, pointsSPEND, nil
}

func (MS *MemStorage) AddWithdraw(ctx context.Context, withdraw models.Withdraw) error {
	orderSum, wSum, err := MS.GetUserBalance(ctx, withdraw.UserLogin)
	if err != nil {
		return err
	}
	if orderSum < wSum+withdraw.Sum {
		return consta.ErrorStatusShortfallAccount
	}
	MS.mu.Lock()
	defer MS.mu.Unlock()
	MS.withdrawCash[uuid.New()] = withdraw
	return nil
}

func (MS *MemStorage) GetWithdraws(_ context.Context, userLogin string) ([]models.Withdraw, error) {
	MS.mu.RLock()
	defer MS.mu.RUnlock()
	var withdraws []models.Withdraw
	for _, v := range MS.withdrawCash {
		if v.UserLogin == userLogin {
			withdraws = append(withdraws, v)
		}
	}
	return withdraws, nil
}
