package model

import "time"

type ProductResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	Image       string  `json:"image"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type ProductQuery struct {
	ID          *string    `query:"id,omitempty" validate:"omitempty,uuid"`
	Name        *string    `query:"name,omitempty" validate:"omitempty,min=3,max=255"`
	Description *string    `query:"description,omitempty" validate:"omitempty,min=3,max=255"`
	Price       *float64   `query:"price,omitempty" validate:"omitempty,min=0"`
	Stock       *int       `query:"stock,omitempty" validate:"omitempty,min=0"`
	Image       *string    `query:"image,omitempty" validate:"omitempty,"`
	CreatedAt   *time.Time `query:"created_at,omitempty" validate:"omitempty,datetime=2006-01-02 15:04:05"`
	UpdatedAt   *time.Time `query:"updated_at,omitempty" validate:"omitempty,datetime=2006-01-02 15:04:05"`
}

type ProductPagination struct {
	Page  int    `query:"page" default:"1" validate:"numeric,omitempty,min=1"`
	Limit int    `query:"limit" default:"10" validate:"numeric,omitempty,min=1,max=100"`
	Sort  string `query:"sort" default:"created_at" validate:"omitempty,oneof=id name description price stock image created_at updated_at"`
	Order string `query:"order" default:"desc" validate:"omitempty,oneof=asc desc"`
}

type GetProductRequest struct {
	ID string `param:"id" validate:"required,uuid"`
}

type CreateProductRequest struct {
	Name        string  `json:"name" validate:"required,min=3,max=255"`
	Description string  `json:"description" validate:"required,min=3,max=255"`
	Price       float64 `json:"price" validate:"required,min=0"`
	Stock       int     `json:"stock" validate:"required,min=0"`
	Image       string  `json:"image" validate:"required,text"`
}

type UpdateProductRequest struct {
	Name        string  `json:"name" validate:"omitempty,min=3,max=255"`
	Description string  `json:"description" validate:"omitempty,min=3,max=255"`
	Price       float64 `json:"price" validate:"omitempty,min=0"`
	Stock       int     `json:"stock" validate:"omitempty,min=0"`
	Image       string  `json:"image" validate:"omitempty,text"`
}

type DeleteProductRequest struct {
	ID string `param:"id" validate:"required,uuid"`
}
