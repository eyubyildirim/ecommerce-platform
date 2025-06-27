package service

import (
	"context"
	"ecommerce-platform/services/order/model"
	"ecommerce-platform/services/order/repository"
	"encoding/json"
	"errors"
	"log/slog"

	pb "ecommerce-platform/pkg/grpc/inventory"

	"github.com/go-chi/chi/middleware"
)

var (
	ErrNoPrice = errors.New("Price not found for one of the items")
)

type OrderService interface {
	CreateOrder(ctx context.Context, userID string, items []model.OrderItem) (*model.Order, error)
	GetOrderByID(ctx context.Context, id string) (*model.Order, error)
	HandlePaymentSucceeded(ctx context.Context, orderID string)
	HandlePaymentFailed(ctx context.Context, orderID, reason string)
	HandleStockReserved(ctx context.Context, orderID string)
	HandleStockUnavailable(ctx context.Context, orderID string, unavailableItems *[]model.OrderItem)
}

type orderServiceImpl struct {
	orderRepo       repository.OrderRepository
	logger          *slog.Logger
	inventoryClient pb.InventoryServiceClient
}

func (or *orderServiceImpl) CreateOrder(ctx context.Context, userID string, items []model.OrderItem) (*model.Order, error) {
	serviceLogger := or.logger.With("request_id", middleware.GetReqID(ctx), "user_id", userID, "items", items)

	serviceLogger.Info("CreateOrder started")

	var productIDs []string
	for _, item := range items {
		productIDs = append(productIDs, item.ProductID)
	}

	products, err := or.inventoryClient.GetProductInfo(ctx, &pb.GetProductInfoRequest{ProductIds: productIDs})
	if err != nil {
		serviceLogger.Error("Failed to fetch product info from grpc server", "error", err)
		return nil, err
	}

	priceMap := make(map[string]float64)
	for _, prod := range products.Products {
		priceMap[prod.Id] = prod.Price
	}

	var totalPrice float64
	for i, item := range items {
		price, ok := priceMap[item.ProductID]

		if !ok {
			serviceLogger.Error("Could not fetch the price for product", "product_id", item.ProductID)
			return nil, ErrNoPrice
		}

		items[i].Price = &price
		totalPrice += price * float64(item.Quantity)
	}

	var order model.Order
	order.UserID = userID
	itemsBytes, _ := json.Marshal(items)
	order.Items = itemsBytes

	serviceLogger.Info("Set items", "items", order.Items)

	order.TotalPrice = totalPrice
	order.Status = "PENDING"

	serviceLogger.Info("Set total price and status", "total_price", order.TotalPrice, "status", order.Status)

	err = or.orderRepo.Create(ctx, &order)
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

func NewOrderService(orderRepo repository.OrderRepository, logger *slog.Logger, inventoryClient pb.InventoryServiceClient) *orderServiceImpl {
	return &orderServiceImpl{
		orderRepo:       orderRepo,
		logger:          logger.With("file", "order_service.go"),
		inventoryClient: inventoryClient,
	}
}
