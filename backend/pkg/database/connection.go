package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"ocserv-bakend/pkg/config"
)

var DB *gorm.DB

func Connect(debug bool) {
	cfg := config.Get()

	log.Printf("Connecting to database %s ...", cfg.Databases)
	db, err := gorm.Open(sqlite.Open(cfg.Databases), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}
	if debug {
		db = db.Debug()
	}
	DB = db
	log.Println("Connected to database")
}

func Get() *gorm.DB {
	return DB
}

func Close() {
	if DB != nil {
		sqlDB, _ := DB.DB()
		err := sqlDB.Close()
		if err != nil {
			log.Fatal("failed to close database")
		}
		log.Println("database closed successfully")
	}
}
