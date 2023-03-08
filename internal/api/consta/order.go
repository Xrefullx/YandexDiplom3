package consta

const (
	OrderStatusPROCESSING = "PROCESSING" // — расчёт начисления в процессе;
	OrderStatusREGISTERED = "REGISTERED" // — заказ зарегистрирован, но начисление не рассчитано;
	OrderStatusNEW        = "NEW"        // — заказ загружен в систему, но не попал в обработку;
	OrderStatusINVALID    = "INVALID"
)
