package bootstrap

import (
	"fmt"
	"log"
	"ocserv-bakend/internal/models"
	"ocserv-bakend/pkg/database"
)

var tables = []interface{}{
	&models.System{},
	&models.User{},
	&models.UserToken{},
	&models.OcservGroup{},
	&models.OcservUser{},
	&models.OcservUserTrafficStatistics{},
}

func Migrate() {
	log.Println("starting migrations...")
	engine := database.Get()
	err := engine.AutoMigrate(tables...)
	if err != nil {
		log.Fatalln(fmt.Sprintf("error sync tables: %v", err))
	}
	log.Println("migrating tables successfully")
}
