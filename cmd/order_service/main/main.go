package main

import (
	"ecommerce-platform/internal/database"
	"ecommerce-platform/services/order/api/handler"
	"ecommerce-platform/services/order/repository/postgres"
	"ecommerce-platform/services/order/service"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With("service", "order-service")

	db := database.InitDb("postgres_db", "5432", "user", "password", "ecommerce_db")
	defer db.Close()

	orderRepo, err := postgres.NewOrderPgRepository(db, logger)
	if err != nil {
		panic(err)
	}
	orderService := service.NewOrderService(orderRepo, logger)
	orderHandler := handler.NewOrderHandler(orderService, logger)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/orders", func(r chi.Router) {
		r.Post("/", orderHandler.CreateOrder)
	})

	http.ListenAndServe(":8081", r)
}
