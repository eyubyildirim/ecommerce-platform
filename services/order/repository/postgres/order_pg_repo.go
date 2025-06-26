package postgres

import (
	"context"
	"database/sql"
	"ecommerce-platform/services/order/model"
	"errors"
	"log/slog"

	"github.com/go-chi/chi/middleware"
	"github.com/google/uuid"
)

type OrderPgRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

func (or *OrderPgRepository) Create(ctx context.Context, order *model.Order) error {
	repoLogger := or.logger.With("request_id", middleware.GetReqID(ctx))

	repoLogger.Info("Create started", "order", order)

	exec := `INSERT INTO orders (id, user_id, items, total_price, status) VALUES ($1, $2, $3, $4, $5) RETURNING created_at, updated_at`

	order.ID = uuid.NewString()
	repoLogger.Info("UUID generated for order", "order", order)

	createdOrder := or.db.QueryRowContext(ctx, exec, order.ID, order.UserID, order.Items, order.TotalPrice, order.Status)
	err := createdOrder.Scan(&order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		repoLogger.Error("Could not create record in database", "order", order, "error", err)
		return err
	}

	repoLogger.Info("Create successful", "order", order)

	return nil
}

func (or *OrderPgRepository) FindByID(ctx context.Context, id string) (*model.Order, error) {
	repoLogger := or.logger.With("request_id", middleware.GetReqID(ctx))

	repoLogger.Info("FindByID started", "order_id", id)

	query := `SELECT id, user_id, items, total_price, status, created_at, updated_at FROM order WHERE id = $1`

	row := or.db.QueryRowContext(ctx, query, id)

	var order model.Order
	err := row.Scan(&order.ID, &order.UserID, &order.Items, &order.TotalPrice, &order.Status, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			repoLogger.Error("No order with given id", "order_id", id, "error", err)
		}

		repoLogger.Error("Error reading database", "error", err)
		return nil, err
	}

	repoLogger.Info("FindByID successful", "order", order)

	return &order, nil
}

func (or *OrderPgRepository) UpdateStatus(ctx context.Context, id string, newStatus string) error {
	repoLogger := or.logger.With("request_id", middleware.GetReqID(ctx))

	repoLogger.Info("UpdateStatus started", "new_status", newStatus)

	query := `UPDATE orders SET status = $1 WHERE id = $2`

	res, err := or.db.Exec(query, newStatus, id)

	if err != nil {
		repoLogger.Error("Could not update database", "error", err)
		return err
	}

	rowsAff, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAff == 0 {
		repoLogger.Error("No order with given id found", "order_id", id, "error", err)
		return sql.ErrNoRows
	}

	repoLogger.Info("UpdateStatus successful", "order_id", id, "new_status", newStatus)

	return nil
}

func NewOrderPgRepository(db *sql.DB, logger *slog.Logger) (*OrderPgRepository, error) {
	if err := db.Ping(); err != nil {
		return nil, errors.New("failed to connect to the database: " + err.Error())
	}

	return &OrderPgRepository{
		db:     db,
		logger: logger.With("file", "order_pg_repo.go"),
	}, nil
}
