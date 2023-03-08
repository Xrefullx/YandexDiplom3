package service

import (
	"context"
	"encoding/json"
	"github.com/Xrefullx/YandexDiplom3/internal/api/consta"
	"github.com/Xrefullx/YandexDiplom3/internal/api/container"
	"github.com/Xrefullx/YandexDiplom3/internal/models"
	"io"
	"net/http"
	"net/url"
	"time"
)

func CalculationLoyaltyPoints(ctx context.Context) error {
	storage := container.GetStorage()
	cfg := container.GetConfig()
	ctxStorage, cancelStorage := context.WithTimeout(ctx, consta.TimeOutStorage)
	orders, err := storage.GetOrdersProcess(ctxStorage)
	cancelStorage()
	if err != nil {
		return err
	}
	for _, order := range orders {
		ctx, cancel := context.WithTimeout(ctx, consta.TimeOutRequest)
		joinPath, errJoin := url.JoinPath(cfg.AccrualAddress, "/api/orders/", order.NumberOrder)
		if errJoin != nil {
			cancel()
			return errJoin
		}
		r, errGet := http.Get(joinPath)
		if errGet != nil {
			cancel()
			return errGet
		}
		switch r.StatusCode {
		case http.StatusTooManyRequests:
			time.Sleep(consta.TimeSleepTooManyRequests)
		case http.StatusInternalServerError:
			_, errRead := io.ReadAll(r.Body)
			if errRead != nil {
				cancel()
				return errRead
			}
			errClose := r.Body.Close()
			if errClose != nil {
				cancel()
				return errClose
			}
		case http.StatusOK:
			body, errRead := io.ReadAll(r.Body)
			if errRead != nil {
				cancel()
				return errRead
			}
			var loyalty models.Loyalty
			if errUnmarshal := json.Unmarshal(body, &loyalty); errUnmarshal != nil {
				cancel()
				return errUnmarshal
			}
			errClose := r.Body.Close()
			if errClose != nil {
				cancel()
				return errClose
			}
			if loyalty.Status == consta.OrderStatusREGISTERED {
				loyalty.Status = consta.OrderStatusPROCESSING
			}
			if order.Status != loyalty.Status {
				if err = storage.UpdateOrder(ctx, loyalty); err != nil {
					cancel()
					return err
				}
			}
		}
		cancel()
	}
	return nil
}
