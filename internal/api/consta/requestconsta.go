package consta

import "time"

const (
	TimeSleepTooManyRequests          = 60 * time.Second
	TimeSleepCalculationLoyaltyPoints = 1 * time.Second
	TimeOutRequest                    = 1 * time.Second
	TimeOutStorage                    = 500 * time.Millisecond
)
