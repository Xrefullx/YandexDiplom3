package models

import "time"

type Order struct {
	UserLogin   string
	Status      string    `json:"status"`
	Accrual     float64   `json:"accrual,omitempty"`
	Uploaded    time.Time `json:"uploaded_at"`
	NumberOrder string    `json:"number"`
}
