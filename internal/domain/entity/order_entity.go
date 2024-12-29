package entity

import "time"

type Order struct {
	ID         string    `db:"id"`
	UserID     string    `db:"user_id"`
	ProductID  string    `db:"product_id"`
	AddressID  string    `db:"address_id"`
	Quantity   int       `db:"quantity"`
	TotalPrice float64   `db:"total_price"`
	Status     string    `db:"status"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}
