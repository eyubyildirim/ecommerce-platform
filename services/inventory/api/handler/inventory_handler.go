package handler

import (
	"database/sql"
	"ecommerce-platform/services/inventory/service"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type AddProductRequest struct {
	Name          string  `json:"name"`
	Price         float64 `json:"price"`
	StockQuantity int     `json:"stockQuantity"`
}

type InventoryHandler struct {
	inventoryService service.InventoryService
	logger           *slog.Logger
}

func NewInventoryHandler(inventoryService service.InventoryService, logger *slog.Logger) *InventoryHandler {
	return &InventoryHandler{
		inventoryService: inventoryService,
		logger:           logger.With("file", "inventory_handler.go"),
	}
}

func (ih *InventoryHandler) AddProduct(w http.ResponseWriter, r *http.Request) {
	reqId := middleware.GetReqID(r.Context())

	reqLogger := ih.logger.With("request_id", reqId)

	reqLogger.Info("Processing new add product request")

	var req AddProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		reqLogger.Error("Invalid request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Price == 0.0 || req.StockQuantity == 0 {
		reqLogger.Error("Invalid request body", "req", req)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdProduct, err := ih.inventoryService.AddProduct(r.Context(), req.Name, req.Price, req.StockQuantity)
	if err != nil {
		reqLogger.Error("Error creating product", "error", err)
		http.Error(w, "Error creating product", http.StatusInternalServerError)
		return
	}

	reqLogger.Info("Product created successfully", "product", createdProduct)

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(createdProduct)
}

func (ih *InventoryHandler) GetPrice(w http.ResponseWriter, r *http.Request) {
	reqId := middleware.GetReqID(r.Context())
	reqLogger := ih.logger.With("request_id", reqId)

	productId := chi.URLParam(r, "id")

	reqLogger.Info("Getting price for given product id", "product_id", productId)

	price, err := ih.inventoryService.GetPrice(r.Context(), productId)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			reqLogger.Error("Product with given id not found", "error", err)
			http.Error(w, "Error creating product", http.StatusNotFound)
			return
		}
		reqLogger.Error("Error fetching price", "error", err)
		http.Error(w, "Error fetching price", http.StatusInternalServerError)
		return
	}

	reqLogger.Info("Got price for given product id", "price", price)

	priceResponse := struct {
		ProductID string  `json:"productId"`
		Price     float64 `json:"price"`
	}{
		ProductID: productId,
		Price:     price,
	}

	reqLogger.Info("Price response created", "price_response", priceResponse)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(priceResponse)
}
