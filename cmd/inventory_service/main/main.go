package main

import (
	"ecommerce-platform/internal/database"
	"ecommerce-platform/services/inventory/api/handler"
	"ecommerce-platform/services/inventory/repository/postgres"
	"ecommerce-platform/services/inventory/service"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/lib/pq"
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

	http.ListenAndServe(":8082", r)
}
