package repository

import (
	"context"
	"ecommerce-platform/services/order/model"
)

type OrderRepository interface {
	Create(ctx context.Context, order *model.Order) error
	FindByID(ctx context.Context, id string) (*model.Order, error)
	UpdateStatus(ctx context.Context, id, newStatus string) error
}
