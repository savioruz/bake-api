package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTConfig struct {
	Secret        string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
}

type JWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type JWTServiceImpl struct {
	secretKey     []byte
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

func NewJWTService(config *JWTConfig) *JWTServiceImpl {
	return &JWTServiceImpl{
		secretKey:     []byte(config.Secret),
		accessExpiry:  config.AccessExpiry,
		refreshExpiry: config.RefreshExpiry,
	}
}

func (s *JWTServiceImpl) GenerateAccessToken(userID, email, role string) (string, error) {
	return s.generateToken(userID, email, role, s.accessExpiry)
}

func (s *JWTServiceImpl) GenerateRefreshToken(userID, email, role string) (string, error) {
	return s.generateToken(userID, email, role, s.refreshExpiry)
}

func (s *JWTServiceImpl) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (s *JWTServiceImpl) generateToken(userID, email, role string, expiry time.Duration) (string, error) {
	claims := &JWTClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}
