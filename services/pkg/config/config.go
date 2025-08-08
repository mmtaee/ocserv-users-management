package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
)

type Config struct {
	Debug        bool
	Host         string
	Port         string
	JWTSecret    string
	AllowOrigins []string
	Dockerized   bool
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

	cfg = &Config{
		Debug:        debug,
		Host:         host,
		Port:         port,
		JWTSecret:    jwtSecret,
		AllowOrigins: strings.Split(allowOrigins, ","),
		Dockerized:   dockerizedEnv == "true",
	}

	log.Println("config initialized")
}

func Get() *Config {
	return cfg
}

func Set(c *Config) {
	cfg = c
}
