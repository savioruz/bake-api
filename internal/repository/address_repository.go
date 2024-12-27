package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/savioruz/bake/internal/domain/entity"
)

type AddressRepository struct {
	db *sqlx.DB
}

func NewAddressRepository(db *sqlx.DB) *AddressRepository {
	return &AddressRepository{db: db}
}

func (r *AddressRepository) Create(tx *sqlx.Tx, address *entity.Address) error {
	query := `INSERT INTO addresses (id, user_id, address_line, city, state, postal_code, country) 
			  VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err := tx.Exec(
		query,
		address.ID,
		address.UserID,
		address.AddressLine,
		address.City,
		address.State,
		address.PostalCode,
		address.Country,
	)
	return err
}

func (r *AddressRepository) GetByUserID(tx *sqlx.Tx, userID string) (*entity.Address, error) {
	query := `SELECT * FROM addresses WHERE user_id = ?`

	var address entity.Address
	err := tx.Get(&address, query, userID)
	if err != nil {
		return nil, err
	}

	return &address, nil
}
