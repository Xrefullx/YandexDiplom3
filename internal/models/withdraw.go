package models

import "time"

type Withdraw struct {
	UserLogin   string
	NumberOrder string    `json:"order"`
	Sum         float64   `json:"sum"`
	ProcessedAT time.Time `json:"processed_at"`
}
