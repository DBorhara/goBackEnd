package model

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	OrderId     uint64         `json:"order_id"`
	CustomerID  uuid.UUI       `json:"customer_id"`
	Products    []OrderProduct `json:"products"`
	CreatedAt   *time.Time     `json:"created_at"`
	ShippedAt   *time.Time     `json:"shipped_at"`
	CompletedAt *time.Time     `json:"completed_at"`
}

type OrderProduct struct {
	ProducID uuid.UUID
	Quantuty uint64
	Price    float64
}
