package bootstrap

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/olekukonko/tablewriter"
	"ocserv-bakend/internal/models"
	"ocserv-bakend/pkg/config"
	"ocserv-bakend/pkg/database"
)

func AdminUsers() {
	config.Init(false)
	database.Connect(false)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := database.Get()

	var users []models.User
	if err := db.WithContext(ctx).
		Where("is_admin = ?", true).
		Select("uid", "username", "updated_at").
		Find(&users).Error; err != nil {
		log.Fatal(err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"UID", "Username", "Updated At"})

	for _, user := range users {
		err := table.Append([]string{
			user.UID,
			user.Username,
			user.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
		if err != nil {
			return
		}
	}
	err := table.Render()
	if err != nil {
		return
	}
}
