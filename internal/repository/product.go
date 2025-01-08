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
	sql := `INSERT INTO public.products (group_id, product_name_id, price, status, quantity, added_by, bought_by, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING product_id`

	row := p.db.QueryRow(ctx, sql,
		product.GroupID,
		product.ProductNameID,
		product.Price,
		product.Status,
		product.Quantity,
		product.AddedBy,
		product.BoughtBy,
		product.CreatedAt)

	var productID uint64
	if err := row.Scan(&productID); err != nil {
		return 0, err
	}

	return productID, nil
}

func (p *ProductRepository) Update(ctx context.Context, product product.Product) error {
	sql := `UPDATE public.products
			SET price = $1,
			    quantity = $2,
			    status = $3,
			    bought_by = $4
			WHERE product_id = $5`

	_, err := p.db.Exec(ctx, sql, product.Price, product.Quantity, product.Status, product.BoughtBy, product.ProductID)

	return err
}

func (p *ProductRepository) Delete(ctx context.Context, productID uint64) error {
	sql := `DELETE FROM public.products WHERE product_id = $1`

	_, err := p.db.Exec(ctx, sql, productID)

	return err
}

func (p *ProductRepository) GetById(ctx context.Context, productID uint64) (product.Product, error) {
	sql := `SELECT * FROM public.products WHERE product_id = $1`

	row := p.db.QueryRow(ctx, sql, productID)

	var product product.Product
	err := row.Scan(
		&product.ProductID,
		&product.GroupID,
		&product.ProductNameID,
		&product.Price,
		&product.Status,
		&product.Quantity,
		&product.AddedBy,
		&product.BoughtBy,
		&product.CreatedAt)

	if err != nil {
		return product, err
	}

	return product, nil
}

func (p *ProductRepository) GetByProductNameId(ctx context.Context, productNameID uint64) (product.ProductName, error) {
	sql := `SELECT * FROM public.product_names WHERE product_name_id = $1`

	row := p.db.QueryRow(ctx, sql, productNameID)

	var productName product.ProductName
	err := row.Scan(&productName.ProductNameID, &productName.CategoryID, &productName.Name)

	if err != nil {
		return product.ProductName{}, err
	}

	return productName, nil
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
