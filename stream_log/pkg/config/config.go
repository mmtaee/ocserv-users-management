package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	Debug      bool
	Host       string
	Port       string
	JWTSecret  string
	WebhookApi string
}

var cfg *Config

func NewConfig(debug bool, host string, port int) *Config {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable not set")
	}
	webhookApi := os.Getenv("WEBHOOK_API")
	if webhookApi == "" {
		webhookApi = "http://ocserv:8080"
	}

	cfg = &Config{
		Debug:      debug,
		Host:       host,
		Port:       strconv.Itoa(port),
		JWTSecret:  jwtSecret,
		WebhookApi: webhookApi,
	}
	return cfg
}

func Get() *Config {
	return cfg
}
