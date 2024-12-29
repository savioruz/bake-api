package model

import "github.com/savioruz/bake/internal/domain/entity"

type CreateOrderRequest struct {
	UserID    string `json:"user_id" validate:"required,uuid"`
	ProductID string `json:"product_id" validate:"required,uuid"`
	Quantity  int    `json:"quantity" validate:"required,min=1"`
}

type OrderResponse struct {
	ID         string         `json:"id"`
	UserID     string         `json:"user_id"`
	ProductID  string         `json:"product_id"`
	AddressID  string         `json:"address_id"`
	Quantity   int            `json:"quantity"`
	TotalPrice float64        `json:"total_price"`
	Status     string         `json:"status"`
	CreatedAt  string         `json:"created_at"`
	UpdatedAt  string         `json:"updated_at"`
	Product    entity.Product `json:"product"`
	Address    entity.Address `json:"address"`
}

type OrderPagination struct {
	Page  int    `query:"page,omitempty" validate:"omitempty,min=1"`
	Limit int    `query:"limit,omitempty" validate:"omitempty,min=1,max=100"`
	Sort  string `query:"sort,omitempty" validate:"omitempty,oneof=id user_id product_id address_id quantity total_price status created_at updated_at ID USER_ID PRODUCT_ID ADDRESS_ID QUANTITY TOTAL_PRICE STATUS CREATED_AT UPDATED_AT"`
	Order string `query:"order,omitempty" validate:"omitempty,oneof=ASC DESC asc desc"`
}

type GetOrderRequest struct {
	ID string `param:"id" validate:"required,uuid"`
}
