package model

type AddressRequest struct {
	AddressLine string `json:"address_line" validate:"required,min=5,max=255"`
	City        string `json:"city" validate:"required,min=2,max=50"`
	State       string `json:"state" validate:"required,min=2,max=50"`
	PostalCode  string `json:"postal_code" validate:"required,min=5,max=20"`
	Country     string `json:"country" validate:"required,min=2,max=50"`
}

type AddressResponse struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	AddressLine string `json:"address_line"`
	City        string `json:"city"`
	State       string `json:"state"`
	PostalCode  string `json:"postal_code"`
	Country     string `json:"country"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
