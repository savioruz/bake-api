package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/savioruz/bake/internal/domain/entity"
	"github.com/savioruz/bake/internal/domain/model"
)

type ProductRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetAll(tx *sqlx.Tx, pagination *model.ProductPagination) ([]entity.Product, int, error) {
	baseQuery := `SELECT * FROM products`
	countQuery := `SELECT COUNT(*) FROM products`

	baseQuery += ` ORDER BY ` + pagination.Sort + ` ` + pagination.Order

	offset := (pagination.Page - 1) * pagination.Limit
	baseQuery += ` LIMIT ? OFFSET ?`

	var total int
	if err := tx.Get(&total, countQuery); err != nil {
		return nil, 0, err
	}

	var products []entity.Product
	err := tx.Select(&products, baseQuery, pagination.Limit, offset)

	return products, total, err
}

func (r *ProductRepository) Search(tx *sqlx.Tx, query *model.ProductQuery, pagination *model.ProductPagination) ([]entity.Product, int, error) {
	baseQuery := `SELECT * FROM products WHERE 1=1`
	countQuery := `SELECT COUNT(*) FROM products WHERE 1=1`

	args := []interface{}{}

	if query.ID != nil {
		baseQuery += ` AND id = ?`
		countQuery += ` AND id = ?`
		args = append(args, *query.ID)
	}
	if query.Name != nil {
		baseQuery += ` AND name LIKE ?`
		countQuery += ` AND name LIKE ?`
		args = append(args, "%"+*query.Name+"%")
	}
	if query.Description != nil {
		baseQuery += ` AND description LIKE ?`
		countQuery += ` AND description LIKE ?`
		args = append(args, "%"+*query.Description+"%")
	}
	if query.Price != nil {
		baseQuery += ` AND price = ?`
		countQuery += ` AND price = ?`
		args = append(args, *query.Price)
	}
	if query.Stock != nil {
		baseQuery += ` AND stock = ?`
		countQuery += ` AND stock = ?`
		args = append(args, *query.Stock)
	}

	baseQuery += ` ORDER BY ` + pagination.Sort + ` ` + pagination.Order

	offset := (pagination.Page - 1) * pagination.Limit
	baseQuery += ` LIMIT ? OFFSET ?`

	paginationArgs := append(args, pagination.Limit, offset)

	var total int
	if err := tx.Get(&total, countQuery, args...); err != nil {
		return nil, 0, err
	}

	var products []entity.Product
	err := tx.Select(&products, baseQuery, paginationArgs...)

	return products, total, err
}

func (r *ProductRepository) GetByID(tx *sqlx.Tx, id string) (*entity.Product, error) {
	query := `SELECT * FROM products WHERE id = ?`

	var product entity.Product
	err := tx.Get(&product, query, id)

	return &product, err
}

func (r *ProductRepository) Create(tx *sqlx.Tx, product *entity.Product) error {
	query := `INSERT INTO products (id, name, description, price, stock, image, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := tx.Exec(
		query,
		product.ID,
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
		product.Image,
		product.CreatedAt,
		product.UpdatedAt,
	)
	return err
}

func (r *ProductRepository) Update(tx *sqlx.Tx, product *entity.Product) error {
	query := `UPDATE products SET name = ?, description = ?, price = ?, stock = ?, image = ?, updated_at = ? WHERE id = ?`

	_, err := tx.Exec(
		query,
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
		product.Image,
		product.UpdatedAt,
		product.ID,
	)
	return err
}

func (r *ProductRepository) Delete(tx *sqlx.Tx, id string) error {
	query := `DELETE FROM products WHERE id = ?`
	_, err := tx.Exec(query, id)
	return err
}
