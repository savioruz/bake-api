package handler

import (
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/savioruz/bake/docs"
	"github.com/savioruz/bake/pkg/config"
)

func Handler(w http.ResponseWriter, r *http.Request) {
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
