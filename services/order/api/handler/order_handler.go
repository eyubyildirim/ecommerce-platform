package handler

import (
	"ecommerce-platform/services/order/model"
	"ecommerce-platform/services/order/service"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

type CreateOrderRequest struct {
	UserID string             `json:"userId"`
	Items  *[]model.OrderItem `json:"items"`
}

type OrderHandler struct {
	orderService service.OrderService
	logger       *slog.Logger
}

func NewOrderHandler(orderService service.OrderService, logger *slog.Logger) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
		logger:       logger,
	}
}

func (oh *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {

	reqLogger := oh.logger.With("request_id", middleware.GetReqID(r.Context()))

	reqLogger.Info("Processing new order request")

	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		reqLogger.Error("Invalid request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdOrder, err := oh.orderService.CreateOrder(r.Context(), req.UserID, req.Items)
	if err != nil {
		reqLogger.Error("Error creating order")
		http.Error(w, "Error creating order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(createdOrder)
}
