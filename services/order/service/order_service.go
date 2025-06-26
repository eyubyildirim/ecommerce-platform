package service

import (
	"context"
	"ecommerce-platform/services/order/model"
	"ecommerce-platform/services/order/repository"
	"log/slog"
)

type OrderService interface {
	CreateOrder(ctx context.Context, userID string, items *[]model.OrderItem) (*model.Order, error)
	GetOrderByID(ctx context.Context, id string) (*model.Order, error)
	HandlePaymentSucceeded(ctx context.Context, orderID string)
	HandlePaymentFailed(ctx context.Context, orderID, reason string)
	HandleStockReserved(ctx context.Context, orderID string)
	HandleStockUnavailable(ctx context.Context, orderID string, unavailableItems *[]model.OrderItem)
}

type orderServiceImpl struct {
	orderRepo repository.OrderRepository
	logger    *slog.Logger
}

func (or *orderServiceImpl) CreateOrder(ctx context.Context, userID string, items *[]model.OrderItem) (*model.Order, error) {
	productPrice := float64(23.99)
	return &model.Order{
		ID:     "test",
		UserID: "test-user",
		Items: &[]model.OrderItem{
			{
				ProductID: "test-product-id",
				Quantity:  2,
				Price:     &productPrice,
			},
		},
		TotalPrice: 47.98,
		Status:     "PENDING",
	}, nil
}

func (or *orderServiceImpl) GetOrderByID(ctx context.Context, id string) (*model.Order, error) {
	panic("not implemented") // TODO: Implement
}

func (or *orderServiceImpl) HandlePaymentSucceeded(ctx context.Context, orderID string) {
	panic("not implemented") // TODO: Implement
}

func (or *orderServiceImpl) HandlePaymentFailed(ctx context.Context, orderID string, reason string) {
	panic("not implemented") // TODO: Implement
}

func (or *orderServiceImpl) HandleStockReserved(ctx context.Context, orderID string) {
	panic("not implemented") // TODO: Implement
}

func (or *orderServiceImpl) HandleStockUnavailable(ctx context.Context, orderID string, unavailableItems *[]model.OrderItem) {
	panic("not implemented") // TODO: Implement
}

func NewOrderService(orderRepo repository.OrderRepository, logger *slog.Logger) *orderServiceImpl {
	return &orderServiceImpl{
		orderRepo: orderRepo,
		logger:    logger,
	}
}
