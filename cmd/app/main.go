package main

import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/savioruz/bake/docs"
	"github.com/savioruz/bake/pkg/config"
)

// @title Bake API
// @version 0.1
// @description This is an auto-generated API Docs.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email jakueenak@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	viper := config.NewViper()
	log := config.NewLogrus(viper)
	db := config.NewDB(viper, log)
	validator := config.NewValidator()
	jwt := config.NewJWT(viper)

	err := config.Bootstrap(&config.BootstrapConfig{
		Viper:     viper,
		Log:       log,
		DB:        db,
		Validator: validator,
		JWT:       jwt,
	})
	if err != nil {
		log.Fatalf("Failed to bootstrap app: %v", err)
	}
}
