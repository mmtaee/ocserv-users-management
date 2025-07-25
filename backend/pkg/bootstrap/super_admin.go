package bootstrap

import (
	"context"
	"fmt"
	"log"
	"ocserv-bakend/internal/models"
	"ocserv-bakend/pkg/config"
	"ocserv-bakend/pkg/crypto"
	"ocserv-bakend/pkg/database"
	"os"
	"strings"
	"time"
)

func CreateSuperAdmin(username, password string, debug bool) {
	config.Init(debug)
	database.Connect(debug)

	if username == "" || password == "" {
		log.Println("You need to specify username and password")
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	cryptoRepo := crypto.NewCustomPassword()
	passwd := cryptoRepo.CreatePassword(password)

	user := models.User{
		Username: strings.ToLower(username),
		Password: passwd.Hash,
		Salt:     passwd.Salt,
		IsAdmin:  true,
	}

	db := database.Get()
	err := db.WithContext(ctx).Create(&user).Error
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("User created with ID: %d\n", user.ID)
}
