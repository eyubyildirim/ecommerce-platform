package postgres

import (
	"context"
	"database/sql"
	"ecommerce-platform/services/order/model"
	"errors"
	"log/slog"

	"github.com/google/uuid"
)

type OrderPgRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

func (or *OrderPgRepository) Create(ctx context.Context, order *model.Order) error {
	exec := `INSERT INTO orders (id, user_id, items, total_price, status) VALUES ($1, $2, $3, $4, $5)`

	order.ID = uuid.NewString()

	_, err := or.db.Exec(exec, order.ID, order.UserID, order.Items, order.TotalPrice, order.Status)
	if err != nil {
		return err
	}

	return nil
}

func (or *OrderPgRepository) FindByID(ctx context.Context, id string) (*model.Order, error) {
	panic("not implemented") // TODO: Implement
}

func (or *OrderPgRepository) UpdateStatus(ctx context.Context, id string, newStatus string) error {
	panic("not implemented") // TODO: Implement
}

func NewOrderPgRepository(db *sql.DB, logger *slog.Logger) (*OrderPgRepository, error) {
	if err := db.Ping(); err != nil {
		return nil, errors.New("failed to connect to the database: " + err.Error())
	}

	return &OrderPgRepository{
		db:     db,
		logger: logger,
	}, nil
}
