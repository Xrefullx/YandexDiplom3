package container

import (
	"github.com/Xrefullx/YandexDiplom3/internal/models"
	"github.com/Xrefullx/YandexDiplom3/internal/storage"
	"go.uber.org/zap"
)

func GetLog() *zap.Logger {
	return DiContainer.Get("zap-logger").(*zap.Logger)
}

func GetStorage() storage.LoyalityStorage {
	return DiContainer.Get("storage").(storage.LoyalityStorage)
}

func GetConfig() models.Config {
	return DiContainer.Get("server-config").(models.Config)
}
