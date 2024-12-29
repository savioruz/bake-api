package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/savioruz/bake/internal/domain/entity"
	"github.com/savioruz/bake/internal/domain/model"
)

type OrderRepository struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(tx *sqlx.Tx, order *entity.Order) error {
	query := `INSERT INTO orders (id, user_id, product_id, address_id, quantity, total_price, status, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := tx.Exec(
		query,
		order.ID,
		order.UserID,
		order.ProductID,
		order.AddressID,
		order.Quantity,
		order.TotalPrice,
		order.Status,
		order.CreatedAt,
		order.UpdatedAt,
	)
	return err
}

func (r *OrderRepository) GetByID(tx *sqlx.Tx, id string) (*entity.Order, error) {
	query := `SELECT * FROM orders WHERE id = ?`

	var order entity.Order
	err := tx.Get(&order, query, id)

	return &order, err
}

func (r *OrderRepository) GetAll(tx *sqlx.Tx, pagination *model.OrderPagination) ([]entity.Order, int, error) {
	baseQuery := `SELECT * FROM orders`
	countQuery := `SELECT COUNT(*) FROM orders`

	baseQuery += ` ORDER BY ` + pagination.Sort + ` ` + pagination.Order

	offset := (pagination.Page - 1) * pagination.Limit
	baseQuery += ` LIMIT ? OFFSET ?`

	var total int
	if err := tx.Get(&total, countQuery); err != nil {
		return nil, 0, err
	}

	var orders []entity.Order
	err := tx.Select(&orders, baseQuery, pagination.Limit, offset)

	return orders, total, err
}
