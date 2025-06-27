package handler

import (
	"database/sql"
	"ecommerce-platform/services/order/model"
	"ecommerce-platform/services/order/service"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type CreateOrderRequest struct {
	UserID string            `json:"userId"`
	Items  []model.OrderItem `json:"items"`
}

type OrderHandler struct {
	orderService service.OrderService
	logger       *slog.Logger
}

func NewOrderHandler(orderService service.OrderService, logger *slog.Logger) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
		logger:       logger.With("file", "order_handler.go"),
	}
}

func (oh *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	reqId := middleware.GetReqID(r.Context())

	reqLogger := oh.logger.With("request_id", reqId)

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

	reqLogger.Info("Order created successfully", "order", createdOrder)

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(createdOrder)
}

func (oh *OrderHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	reqLogger := oh.logger.With("request_id", middleware.GetReqID(r.Context()))

	orderId := chi.URLParam(r, "id")

	reqLogger.Info("Retrieving order by id", "order_id", orderId, "path", r.URL.Path)

	order, err := oh.orderService.GetOrderByID(r.Context(), orderId)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			http.Error(w, "No order with given id", http.StatusNotFound)
			return
		}

		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(order)
}
