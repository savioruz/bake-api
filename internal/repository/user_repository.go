package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/savioruz/bake/internal/domain/entity"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(db *sqlx.Tx, entity *entity.User) error {
	query := `INSERT INTO users (id, email, password, name, phone, role, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := db.Exec(
		query,
		entity.ID,
		entity.Email,
		entity.Password,
		entity.Name,
		entity.Phone,
		entity.Role,
		entity.CreatedAt,
		entity.UpdatedAt,
	)
	return err
}

func (r *UserRepository) GetByEmail(db *sqlx.Tx, email string) (*entity.User, error) {
	query := `SELECT * FROM users WHERE email = ?`

	var user entity.User
	err := db.Get(&user, query, email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByID(db *sqlx.Tx, id string) (*entity.User, error) {
	query := `SELECT * FROM users WHERE id = ?`

	var user entity.User
	err := db.Get(&user, query, id)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetFirst(db *sqlx.Tx, user *entity.User) error {
	query := `SELECT * FROM users LIMIT 1`

	err := db.Get(user, query)
	if err != nil {
		return err
	}

	return nil
}
