package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tclutin/shoppinglist-api/internal/domain/product"
)

type ProductRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{db: db}
}

func (p *ProductRepository) Create(ctx context.Context, product product.Product) (uint64, error) {
	panic("implement me")
}

func (p *ProductRepository) GetCategories(ctx context.Context) ([]product.Category, error) {
	sql := `SELECT * FROM public.categories`

	rows, err := p.db.Query(ctx, sql)
	if err != nil {
		return nil, err
	}

	categories, err := pgx.CollectRows(rows, pgx.RowToStructByName[product.Category])
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (p *ProductRepository) GetProductsByCategoryId(ctx context.Context, categoryID uint64) ([]product.ProductName, error) {
	sql := `SELECT * FROM public.product_names WHERE category_id = $1`

	rows, err := p.db.Query(ctx, sql, categoryID)
	if err != nil {
		return nil, err
	}

	products, err := pgx.CollectRows(rows, pgx.RowToStructByName[product.ProductName])
	if err != nil {
		return nil, err
	}

	return products, nil
}
