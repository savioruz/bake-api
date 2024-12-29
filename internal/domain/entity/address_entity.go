package entity

import "time"

type Address struct {
	ID          string    `db:"id" json:"id"`
	UserID      string    `db:"user_id" json:"user_id"`
	AddressLine string    `db:"address_line" json:"address_line"`
	City        string    `db:"city" json:"city"`
	State       string    `db:"state" json:"state"`
	PostalCode  string    `db:"postal_code" json:"postal_code"`
	Country     string    `db:"country" json:"country"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
	User        *User     `db:"-" json:"user,omitempty"`
}

func (Address) TableName() string {
	return "addresses"
}
