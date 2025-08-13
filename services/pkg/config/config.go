package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
)

type Config struct {
	Debug         bool
	Host          string
	Port          string
	JWTSecret     string
	AllowOrigins  []string
	Database      string
	Dockerized    bool
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

	host := os.Getenv("HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	allowOrigins := os.Getenv("ALLOW_ORIGINS")

	dockerizedEnv := os.Getenv("DOCKERIZED")

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable not set")
	}

	database := os.Getenv("DATABASE")
	if database == "" {
		log.Fatal("DATABASE environment variable not set")
	}

	apiURLService := os.Getenv("API_URL_SERVICE")
	if apiURLService == "" {
		apiURLService = "http://ocserv:8080"
	}

	cfg = &Config{
		Debug:         debug,
		Host:          host,
		Port:          port,
		JWTSecret:     jwtSecret,
		AllowOrigins:  strings.Split(allowOrigins, ","),
		Dockerized:    dockerizedEnv == "true",
		Database:      database,
		APIURLService: apiURLService,
	}

	log.Println("config initialized")
}

func Get() *Config {
	return cfg
}
