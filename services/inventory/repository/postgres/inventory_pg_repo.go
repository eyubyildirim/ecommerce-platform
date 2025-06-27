package postgres

import (
	"context"
	"database/sql"
	"ecommerce-platform/services/inventory/model"
	"errors"
	"log/slog"

	"github.com/go-chi/chi/middleware"
	"github.com/google/uuid"
)

type InventoryPgRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

func (in *InventoryPgRepository) FindByID(ctx context.Context, id string) (*model.Product, error) {
	repoLogger := in.logger.With("request_id", middleware.GetReqID(ctx))

	repoLogger.Info("FindByID started", "product_id", id)

	query := "SELECT id, name, price, stock_quantity, created_at, updated_at FROM products WHERE id = $1"

	repoLogger.Info("Finding product", "product_id", id)

	var product model.Product
	row := in.db.QueryRowContext(ctx, query, id)

	err := row.Scan(&product.ID, &product.Name, &product.Price, &product.StockQuantity, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		repoLogger.Error("Error finding product", "error", err)
		return nil, err
	}

	repoLogger.Info("FindByID successful", "product", product)

	return &product, nil
}

func (in *InventoryPgRepository) Create(ctx context.Context, product *model.Product) error {
	repoLogger := in.logger.With("request_id", middleware.GetReqID(ctx))

	repoLogger.Info("Create started", "input_product", product)

	query := `INSERT INTO products (id, name, price, stock_quantity) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
	row := in.db.QueryRowContext(ctx, query, uuid.NewString(), product.Name, product.Price, product.StockQuantity)

	err := row.Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		repoLogger.Error("Could not create product", "error", err)
		return err
	}

	repoLogger.Info("Create successful", "output_product", product)

	return nil
}

func (in *InventoryPgRepository) FindManyByIDs(ctx context.Context, ids []string) ([]*model.Product, error) {
	repoLogger := in.logger.With("request_id", middleware.GetReqID(ctx))

	repoLogger.Info("FindManyByIDs started", "id_list", ids)

	query := "SELECT id, name, price, stock_quantity, created_at, updated_at FROM products WHERE id = $1"

	var products []*model.Product
	for _, id := range ids {
		var product model.Product
		repoLogger.Info("Finding product", "product_id", id)
		row := in.db.QueryRowContext(ctx, query, id)
		err := row.Scan(&product.ID, &product.Name, &product.Price, &product.StockQuantity, &product.CreatedAt, &product.UpdatedAt)
		if err != nil {
			repoLogger.Error("Error finding product", "error", err)
			return nil, err
		}

		products = append(products, &product)
	}

	repoLogger.Info("FindManyByIDs successful", "products", products)

	return products, nil
}

func (in *InventoryPgRepository) UpdateStockQuantity(ctx context.Context, id string, change int) error {
	repoLogger := in.logger.With("request_id", middleware.GetReqID(ctx))

	repoLogger.Info("UpdateStockQuantity started", "product_id", id, "stock_change", change)

	query := `UPDATE products SET stock = stock + $1 WHERE id = $2`

	res, err := in.db.ExecContext(ctx, query, change, id)
	if err != nil {
		repoLogger.Info("Error updating stock quantity", "error", err)
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected < 1 {
		repoLogger.Error("Could not find product with given id", "product_id", id)
		return sql.ErrNoRows
	}

	repoLogger.Info("UpdateStockQuantity successful")

	return nil
}

func NewInventoryPgRepository(db *sql.DB, logger *slog.Logger) (*InventoryPgRepository, error) {
	if err := db.Ping(); err != nil {
		return nil, errors.New("failed to connect to the database: " + err.Error())
	}

	return &InventoryPgRepository{
		db:     db,
		logger: logger.With("file", "inventory_pg_repo.go"),
	}, nil
}
