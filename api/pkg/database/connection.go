package database

import (
	"api/pkg/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"path/filepath"
)

var DB *gorm.DB

func Connect(cfg *config.Config) {
	dbPath := "./db"
	if cfg.Debug {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		dbPath = filepath.Join(home, "ocserv_db")
	}

	err := os.MkdirAll(dbPath, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	dbPath = filepath.Join(dbPath, "ocserv.db")
	log.Printf("Connecting to database %s ...", dbPath)
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}
	if cfg.Debug {
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
