package middleware

import (
	"context"
	"net/http"
	"strings"

	e "github.com/savioruz/bake/pkg/error"
	"github.com/savioruz/bake/pkg/jwt"
	"github.com/sirupsen/logrus"
)

type AuthMiddleware struct {
	jwtService jwt.JWTService
	log        *logrus.Logger
}

func NewAuthMiddleware(jwtService jwt.JWTService, log *logrus.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
		log:        log,
	}
}

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

const (
	UserIDKey contextKey = "user_id"
	EmailKey  contextKey = "email"
	RoleKey   contextKey = "role"
)

func (m *AuthMiddleware) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			m.log.Warn("No Authorization header provided")
			e.ErrorHandler(w, r, http.StatusUnauthorized, e.ErrUnauthorized)
			return
		}

		// Check Bearer scheme
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			m.log.Warn("Invalid Authorization header format")
			e.ErrorHandler(w, r, http.StatusUnauthorized, e.ErrUnauthorized)
			return
		}

		// Validate token
		token := parts[1]
		claims, err := m.jwtService.ValidateToken(token)
		if err != nil {
			m.log.WithError(err).Warn("Invalid token")
			e.ErrorHandler(w, r, http.StatusUnauthorized, e.ErrUnauthorized)
			return
		}

		// Add claims to context
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, EmailKey, claims.Email)
		ctx = context.WithValue(ctx, RoleKey, claims.Role)

		// Call next handler with updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (m *AuthMiddleware) RequireRole(roles []string, next http.HandlerFunc) http.HandlerFunc {
	return m.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		// Get role from context
		role := r.Context().Value(RoleKey).(string)

		// Check if user's role is allowed
		allowed := false
		for _, allowedRole := range roles {
			if role == allowedRole {
				allowed = true
				break
			}
		}

		if !allowed {
			m.log.WithFields(logrus.Fields{
				"required_roles": roles,
				"user_role":      role,
			}).Warn("Insufficient permissions")
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Helper functions to get values from context
func GetUserIDFromContext(ctx context.Context) string {
	if id, ok := ctx.Value(UserIDKey).(string); ok {
		return id
	}
	return ""
}

func GetEmailFromContext(ctx context.Context) string {
	if email, ok := ctx.Value(EmailKey).(string); ok {
		return email
	}
	return ""
}

func GetRoleFromContext(ctx context.Context) string {
	if role, ok := ctx.Value(RoleKey).(string); ok {
		return role
	}
	return ""
}
