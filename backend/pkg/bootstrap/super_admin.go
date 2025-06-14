package bootstrap

import (
	"context"
	"fmt"
	"log"
	"ocserv-bakend/internal/models"
	"ocserv-bakend/pkg/config"
	"ocserv-bakend/pkg/database"
	"os"
	"time"
)

func CreateSuperAdmin(username, password string) {
	config.Init(false)
	database.Connect(false)

	if username == "" || password == "" {
		log.Println("You need to specify username and password")
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	db := database.Get()
	user := models.User{
		Username: username,
		Password: password,
		IsAdmin:  true,
	}
	err := db.WithContext(ctx).Create(&user).Error
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("User created with ID: %d\n", user.ID)
}
