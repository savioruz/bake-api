package config

import (
	"github.com/savioruz/bake/pkg/jwt"
	"github.com/spf13/viper"
)

func NewJWT(viper *viper.Viper) *jwt.JWTConfig {
	return &jwt.JWTConfig{
		Secret:        viper.GetString("JWT_SECRET"),
		AccessExpiry:  viper.GetDuration("JWT_ACCESS_EXPIRY"),
		RefreshExpiry: viper.GetDuration("JWT_REFRESH_EXPIRY"),
	}
}
