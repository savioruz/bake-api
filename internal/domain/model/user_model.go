package model

type UserRegisterRequest struct {
	Email    string          `json:"email" validate:"required,email,min=3,max=100"`
	Password string          `json:"password" validate:"required,min=8,max=255"`
	Name     string          `json:"name" validate:"required,min=5,max=100"`
	Phone    string          `json:"phone" validate:"required,min=10,max=15"`
	Address  *AddressRequest `json:"address,omitempty"`
}

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email,min=3,max=100"`
	Password string `json:"password" validate:"required,min=8,max=255"`
}

type UserResponse struct {
	ID        string           `json:"id"`
	Email     string           `json:"email"`
	Name      string           `json:"name"`
	Phone     string           `json:"phone"`
	Role      string           `json:"role"`
	CreatedAt string           `json:"created_at"`
	UpdatedAt string           `json:"updated_at"`
	Address   *AddressResponse `json:"address,omitempty"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
