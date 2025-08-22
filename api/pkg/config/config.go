package config

import (
	"log"
	"os"
	"strings"
)

type Config struct {
	Debug        bool
	Host         string
	Port         int
	SecretKey    string
	JWTSecret    string
	AllowOrigins []string
	WebhookApi   string
}

var cfg *Config

func NewConfig(debug bool, host string, port int) *Config {
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		secretKey = "SECRET_KEY122456"
	}

	allowOrigins := os.Getenv("ALLOW_ORIGINS")

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable not set")
	}

	webhookApi := os.Getenv("WEBHOOK_API")
	if webhookApi == "" {
		webhookApi = "http://ocserv:8080"
	}

	cfg = &Config{
		Debug:        debug,
		Host:         host,
		Port:         port,
		SecretKey:    secretKey,
		JWTSecret:    jwtSecret,
		AllowOrigins: strings.Split(allowOrigins, ","),
		WebhookApi:   webhookApi,
	}
	return cfg
}

func Get() *Config {
	return cfg
}
