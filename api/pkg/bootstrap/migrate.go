package bootstrap

import (
	"api/internal/models"
	"api/pkg/audit_log"
	"api/pkg/database"
	"fmt"
	"log"
)

var tables = []interface{}{
	&models.System{},
	&models.User{},
	&models.UserToken{},
	&models.OcservGroup{},
	&models.OcservUser{},
	&models.OcservUserTrafficStatistics{},
	&audit_log.AuditLog{},
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
