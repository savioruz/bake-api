package builder

import (
	"net/http"
	"strings"

	"github.com/savioruz/bake/internal/handler"
	"github.com/savioruz/bake/pkg/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

const (
	APIPrefix = "/api/v1"
)

type Routes struct {
	Method  string
	Path    string
	Handler http.HandlerFunc
	Roles   []string
}

type Config struct {
	AuthMiddleware *middleware.AuthMiddleware
	UserHandler    *handler.UserHandler
	ProductHandler *handler.ProductHandler
}

// Helper function to prefix routes with /api/v1
func prefixRoute(path string) string {
	return strings.TrimRight(APIPrefix, "/") + "/" + strings.TrimLeft(path, "/")
}

func PublicRoutes(c *Config) []Routes {
	return []Routes{
		{
			Method: http.MethodGet,
			Path:   "/",
			Handler: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("Hello World"))
			},
		},
		{
			Method:  http.MethodPost,
			Path:    prefixRoute("/users"),
			Handler: c.UserHandler.Register,
		},
		{
			Method:  http.MethodPost,
			Path:    prefixRoute("/users/login"),
			Handler: c.UserHandler.Login,
		},
		{
			Method:  http.MethodPost,
			Path:    prefixRoute("/users/refresh"),
			Handler: c.UserHandler.RefreshToken,
		},
		{
			Method: http.MethodGet,
			Path:   prefixRoute("/docs/*"),
			Handler: httpSwagger.Handler(
				httpSwagger.URL(prefixRoute("/docs/doc.json")),
			),
		},
		{
			Method:  http.MethodGet,
			Path:    prefixRoute("/products"),
			Handler: c.ProductHandler.GetAll,
		},
		{
			Method:  http.MethodGet,
			Path:    prefixRoute("/products/search"),
			Handler: c.ProductHandler.Search,
		},
		{
			Method:  http.MethodGet,
			Path:    prefixRoute("/products/{id}"),
			Handler: c.ProductHandler.GetByID,
		},
	}
}

func PrivateRoutes(c *Config) []Routes {
	return []Routes{
		{
			Method:  http.MethodGet,
			Path:    prefixRoute("/users/me"),
			Handler: c.AuthMiddleware.RequireAuth(c.UserHandler.Me),
			Roles:   []string{"user", "admin"},
		},
		{
			Method:  http.MethodPost,
			Path:    prefixRoute("/products"),
			Handler: c.AuthMiddleware.RequireRole([]string{"admin"}, c.ProductHandler.Create),
		},
		{
			Method:  http.MethodPut,
			Path:    prefixRoute("/products/{id}"),
			Handler: c.AuthMiddleware.RequireRole([]string{"admin"}, c.ProductHandler.Update),
		},
		{
			Method:  http.MethodDelete,
			Path:    prefixRoute("/products/{id}"),
			Handler: c.AuthMiddleware.RequireRole([]string{"admin"}, c.ProductHandler.Delete),
		},
	}
}

func SwaggerRoutes() []Routes {
	return []Routes{
		{
			Method: http.MethodGet,
			Path:   "/docs/",
			Handler: httpSwagger.Handler(
				httpSwagger.URL("/docs/doc.json"),
				httpSwagger.DeepLinking(true),
				httpSwagger.DocExpansion("list"),
			),
		},
	}
}
