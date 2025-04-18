package repository

import (
    "github.com/google/uuid"
    "E-Commerce/order-service/internal/entity"
)

type OrderRepository interface {
    Create(order *entity.Order) error
    Get(id uuid.UUID) (*entity.Order, error)
    UpdateStatus(id uuid.UUID, status string) error
    List(userID string, page, pageSize int) ([]*entity.Order, int, error)
}