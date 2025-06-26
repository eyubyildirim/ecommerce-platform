package service

import (
	"context"
	"ecommerce-platform/services/order/model"
	"ecommerce-platform/services/order/repository"
	"encoding/json"
	"log/slog"

	"github.com/go-chi/chi/middleware"
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
	serviceLogger := or.logger.With("request_id", middleware.GetReqID(ctx), "user_id", userID, "items", items)

	serviceLogger.Info("CreateOrder started")

	var order model.Order
	order.UserID = userID
	itemsBytes, _ := json.Marshal(items)
	order.Items = itemsBytes

	order.TotalPrice = 50
	order.Status = "PENDING"

	serviceLogger.Info("Set total price and status", "total_price", order.TotalPrice, "status", order.Status)

	err := or.orderRepo.Create(ctx, &order)
	if err != nil {
		serviceLogger.Error("CreateOrder failed")
		return nil, err
	}

	serviceLogger.Info("CreateOrder completed successfully", "order", order)
	return &order, nil
}

func (or *orderServiceImpl) GetOrderByID(ctx context.Context, id string) (*model.Order, error) {
	serviceLogger := or.logger.With("request_id", middleware.GetReqID(ctx), "order_id", id)

	serviceLogger.Info("GetOrderByID started")

	order, err := or.orderRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	serviceLogger.Info("GetOrderByID completed successfully")

	return order, nil
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
		logger:    logger.With("file", "order_service.go"),
	}
}
