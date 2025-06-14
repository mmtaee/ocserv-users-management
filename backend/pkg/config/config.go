package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
)

type Config struct {
	Debug         bool
	Host          string
	Port          string
	SecretKey     string
	JWTSecret     string
	AllowOrigins  []string
	Databases     string
	APIURLService string
}

var cfg *Config

func Init(debug bool) {
	if debug {
		err := godotenv.Load()
		if err != nil {
			log.Printf("Error loading .env file: %v", err)
		}
	}

	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		secretKey = "SECRET_KEY122456"
	}

	host := os.Getenv("HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	allowOrigins := os.Getenv("ALLOW_ORIGINS")

	databases := os.Getenv("DATABASES")
	if databases == "" {
		databases = "ocserv.db"
	}

	if dockerized := os.Getenv("DOCKERIZED"); dockerized == "true" {
		databases = fmt.Sprintf("/db/%s", databases)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = secretKey
	}

	apiURLService := os.Getenv("API_URL_SERVICE")
	if apiURLService == "" {
		apiURLService = "http://ocserv:8080"
	}

	cfg = &Config{
		Debug:         debug,
		Host:          host,
		Port:          port,
		SecretKey:     secretKey,
		JWTSecret:     jwtSecret,
		AllowOrigins:  strings.Split(allowOrigins, ","),
		Databases:     databases,
		APIURLService: apiURLService,
	}

	log.Println("config initialized")
}

func Get() *Config {
	return cfg
}

func Set(c *Config) {
	cfg = c
}
