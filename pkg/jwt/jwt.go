package jwt

type JWTService interface {
	GenerateAccessToken(userID, email, role string) (string, error)
	GenerateRefreshToken(userID, email, role string) (string, error)
	ValidateToken(tokenString string) (*JWTClaims, error)
}
