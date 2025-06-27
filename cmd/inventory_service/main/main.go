package main

import (
	"ecommerce-platform/internal/database"
	"ecommerce-platform/services/inventory/api/handler"
	"ecommerce-platform/services/inventory/repository/postgres"
	"ecommerce-platform/services/inventory/service"
	"log/slog"
	"net"
	"net/http"
	"os"

	pb "ecommerce-platform/pkg/grpc/inventory"
	inventory_grpc "ecommerce-platform/services/inventory/grpc"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With("service", "product-service")

	db := database.InitDb("postgres_db", "5432", "user", "password", "ecommerce_db")
	defer db.Close()

	inventoryRepo, err := postgres.NewInventoryPgRepository(db, logger)
	if err != nil {
		panic(err)
	}
	inventoryService := service.NewInventoryService(inventoryRepo, logger)
	inventoryHandler := handler.NewInventoryHandler(inventoryService, logger)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	r.Route("/products", func(r chi.Router) {
		r.Post("/", inventoryHandler.AddProduct)
		r.Get("/{id}", inventoryHandler.GetPrice)
	})

	go func() {
		http.ListenAndServe(":8082", r)
	}()

	grpcPort := os.Getenv("GRPC_PORT")

	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		logger.Error("Failed to listen for gRPC", "error", err)
		os.Exit(1)
	}

	grpcServer := grpc.NewServer()
	inventoryServer := inventory_grpc.NewInventoryGRPCServer(inventoryService)
	pb.RegisterInventoryServiceServer(grpcServer, inventoryServer)

	logger.Info("gRPC inventory server is starting", "port", grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		logger.Error("gRPC server failed to start", "error", err)
		os.Exit(1)
	}
}
