package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/savioruz/bake/internal/builder"
	"github.com/savioruz/bake/internal/handler"
	"github.com/savioruz/bake/internal/repository"
	"github.com/savioruz/bake/internal/service"
	"github.com/savioruz/bake/pkg/jwt"
	"github.com/savioruz/bake/pkg/middleware"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type BootstrapConfig struct {
	DB        *sqlx.DB
	Log       *logrus.Logger
	Validator *validator.Validate
	JWT       *jwt.JWTConfig
	Viper     *viper.Viper
}

func Bootstrap(c *BootstrapConfig) error {
	// Initialize services
	jwtService := jwt.NewJWTService(c.JWT)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtService, c.Log)

	// Initialize repositories
	userRepository := repository.NewUserRepository(c.DB)
	addressRepository := repository.NewAddressRepository(c.DB)
	productRepository := repository.NewProductRepository(c.DB)
	orderRepository := repository.NewOrderRepository(c.DB)

	// Initialize services
	userService := service.NewUserService(userRepository, addressRepository, c.DB, c.Log, c.Validator, jwtService)
	productService := service.NewProductService(productRepository, c.DB, c.Log, c.Validator)
	orderService := service.NewOrderService(orderRepository, productRepository, addressRepository, c.DB, c.Log, c.Validator)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService, c.Log)
	productHandler := handler.NewProductHandler(productService, c.Log)
	orderHandler := handler.NewOrderHandler(orderService, c.Log)

	// Initialize server
	server := NewServer(c.Viper, c.Log)

	// Register routes
	routeConfig := &builder.Config{
		AuthMiddleware: authMiddleware,
		UserHandler:    userHandler,
		ProductHandler: productHandler,
		OrderHandler:   orderHandler,
	}

	publicRoutes := builder.PublicRoutes(routeConfig)
	privateRoutes := builder.PrivateRoutes(routeConfig)
	swaggerRoutes := builder.SwaggerRoutes()
	allRoutes := make([]builder.Routes, 0, len(publicRoutes)+len(privateRoutes)+len(swaggerRoutes))
	allRoutes = append(allRoutes, publicRoutes...)
	allRoutes = append(allRoutes, privateRoutes...)
	allRoutes = append(allRoutes, swaggerRoutes...)
	server.RegisterRoutes(allRoutes)

	// Start server
	c.Log.Info("App is bootstrapped successfully")
	return server.Start()
}
