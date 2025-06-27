package main

import (
	"ecommerce-platform/internal/database"
	"ecommerce-platform/services/order/api/handler"
	"ecommerce-platform/services/order/repository/postgres"
	"ecommerce-platform/services/order/service"
	"log/slog"
	"net/http"
	"os"

	pb "ecommerce-platform/pkg/grpc/inventory"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With("service", "order-service")

	inventorySvcAddr := os.Getenv("INVENTORY_SERVICE_GRPC_ADDR") // e.g., "inventory-service:9090"

	// Create the gRPC connection to the inventory service
	conn, err := grpc.NewClient(inventorySvcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error("Failed to connect to inventory service", "error", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Create the gRPC client from the connection
	inventoryClient := pb.NewInventoryServiceClient(conn)
	db := database.InitDb("postgres_db", "5432", "user", "password", "ecommerce_db")
	defer db.Close()

	orderRepo, err := postgres.NewOrderPgRepository(db, logger)
	if err != nil {
		panic(err)
	}
	orderService := service.NewOrderService(orderRepo, logger, inventoryClient)
	orderHandler := handler.NewOrderHandler(orderService, logger)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	r.Route("/orders", func(r chi.Router) {
		r.Post("/", orderHandler.CreateOrder)
		r.Get("/{id}", orderHandler.GetOrderByID)
	})

	http.ListenAndServe(":8081", r)
}
