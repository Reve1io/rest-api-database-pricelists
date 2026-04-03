package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type ProductRow struct {
	Code                 string
	Name                 string
	Producer             string
	Supplier             string
	SupplierCurrency     string
	SupplierDeliveryTime string
	Quant                int
	Price                float64
	Currency             string
	QtyAvailable         int
	Moq                  int
}

type ProductRepository struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
}

func NewProductRepository(pool *pgxpool.Pool, logger *zap.Logger) *ProductRepository {
	return &ProductRepository{
		pool:   pool,
		logger: logger,
	}
}

func (r *ProductRepository) SearchByName(ctx context.Context, name string) ([]ProductRow, error) {

	query := `
	SELECT 
		p.code,
		p.name,
		pr.name,
		s.name,
		s.currency,
		s.delivery_time,
		cp.quant,
		cp.price,
		cp.currency,
		cp.qty_available,
		cp.moq
	FROM products p
	JOIN producers pr ON pr.id = p.producer_id
	JOIN current_prices cp ON cp.product_id = p.id
	JOIN suppliers s ON s.id = cp.supplier_id
	WHERE p.search_vector @@ plainto_tsquery('russian', $1)
	AND p.is_active = true
	ORDER BY cp.quant ASC;
	`

	start := time.Now()

	r.logger.Info("executing sql query",
		zap.String("search_term", name),
	)

	rows, err := r.pool.Query(ctx, query, name)
	if err != nil {
		r.logger.Error("sql query failed", zap.Error(err))
		return nil, err
	}

	r.logger.Info("sql query executed",
		zap.Duration("duration", time.Since(start)),
	)
	defer rows.Close()

	var result []ProductRow

	for rows.Next() {
		var row ProductRow
		err := rows.Scan(
			&row.Code,
			&row.Name,
			&row.Producer,
			&row.Supplier,
			&row.SupplierCurrency,
			&row.SupplierDeliveryTime,
			&row.Quant,
			&row.Price,
			&row.Currency,
			&row.QtyAvailable,
			&row.Moq,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, row)
	}

	r.logger.Info("db rows fetched",
		zap.Int("rows_count", len(result)),
	)

	return result, nil
}
