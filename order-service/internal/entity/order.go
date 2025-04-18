package entity

import (
	"time"

	"github.com/google/uuid"
)

type OrderItem struct {
	ID        uuid.UUID `db:"id"`
	OrderID   uuid.UUID `db:"order_id"`
	ProductID uuid.UUID `db:"product_id"`
	Quantity  int       `db:"quantity"`
	Price     float64   `db:"price"`
}

type Order struct {
	ID          uuid.UUID    `db:"id"`
	UserID      string       `db:"user_id"`
	Status      string       `db:"status"`
	TotalAmount float64      `db:"total_amount"`
	CreatedAt   time.Time    `db:"created_at"`
	Items       []*OrderItem `db:"-"`
}