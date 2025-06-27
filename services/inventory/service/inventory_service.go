package service

import (
	"context"
	"database/sql"
	"ecommerce-platform/services/inventory/model"
	"ecommerce-platform/services/inventory/repository"
	"errors"
	"log/slog"

	"github.com/go-chi/chi/middleware"
)

type InventoryService interface {
	HandleReserveStock(ctx context.Context, orderID string, items []*model.Product) error
	AddProduct(ctx context.Context, name string, price float64, quantity int) (*model.Product, error)
	GetPrice(ctx context.Context, id string) (float64, error)
}

type inventoryServiceImpl struct {
	inventoryRepo repository.InventoryRepository
	logger        *slog.Logger
}

func (in *inventoryServiceImpl) GetPrice(ctx context.Context, id string) (float64, error) {
	serviceLogger := in.logger.With("request_id", middleware.GetReqID(ctx), "product_id", id)

	serviceLogger.Info("GetPrice started")

	product, err := in.inventoryRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			serviceLogger.Error("Product with given id could not be found")
			return 0, err
		}

		serviceLogger.Error("Could not get product", "error", err)
		return 0, err
	}

	price := product.Price

	serviceLogger.Info("GetPrice successful", "price", price)

	return price, nil
}

func (in *inventoryServiceImpl) HandleReserveStock(ctx context.Context, orderID string, items []*model.Product) error {
	panic("not implemented") // TODO: Implement
}

func (in *inventoryServiceImpl) AddProduct(ctx context.Context, name string, price float64, quantity int) (*model.Product, error) {
	serviceLogger := in.logger.With("request_id", middleware.GetReqID(ctx), "name", name, "price", price, "quantity", quantity)

	serviceLogger.Info("AddProduct started")

	product := model.Product{
		Name:          name,
		Price:         price,
		StockQuantity: quantity,
	}

	err := in.inventoryRepo.Create(ctx, &product)
	if err != nil {
		return nil, err
	}

	serviceLogger.Info("AddProduct completed successfully", "final_product", product)

	return &product, nil
}

func NewInventoryService(inventoryRepo repository.InventoryRepository, logger *slog.Logger) *inventoryServiceImpl {
	return &inventoryServiceImpl{
		inventoryRepo: inventoryRepo,
		logger:        logger.With("file", "inventory_service.go"),
	}
}
