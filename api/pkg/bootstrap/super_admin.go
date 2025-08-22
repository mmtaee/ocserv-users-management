package bootstrap

import (
	"api/internal/models"
	"api/pkg/config"
	"api/pkg/crypto"
	"api/pkg/database"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func CreateSuperAdmin(username, password string, debug bool) {
	cfg := config.NewConfig(debug, "", 0)
	database.Connect(cfg)
	Migrate()

	if username == "" || password == "" {
		log.Println("You need to specify username and password")
		os.Exit(1)
	}

	db := database.Get()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var existing models.User

	if err := db.WithContext(ctx).
		Where("LOWER(username) = ?", strings.ToLower(username)).
		First(&existing).Error; err == nil {
		log.Printf("User with username '%s' already exists (ID: %d)\n", existing.Username, existing.ID)
		return
	}

	cryptoRepo := crypto.NewCustomPassword()
	passwd := cryptoRepo.CreatePassword(password)

	user := models.User{
		Username: strings.ToLower(username),
		Password: passwd.Hash,
		Salt:     passwd.Salt,
		IsAdmin:  true,
	}

	err := db.WithContext(ctx).Create(&user).Error
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("User created with ID: %d\n", user.ID)
}
