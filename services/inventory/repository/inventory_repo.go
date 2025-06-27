package repository

import (
	"context"
	"ecommerce-platform/services/inventory/model"
)

type InventoryRepository interface {
	Create(ctx context.Context, product *model.Product) error
	FindManyByIDs(ctx context.Context, ids []string) ([]*model.Product, error)
	FindByID(ctx context.Context, id string) (*model.Product, error)
	UpdateStockQuantity(ctx context.Context, id string, change int) error
}
