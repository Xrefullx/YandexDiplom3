package models

type Loyalty struct {
	Status      string  `json:"status"`
	Accrual     float64 `json:"accrual,omitempty"`
	NumberOrder string  `json:"order"`
}
