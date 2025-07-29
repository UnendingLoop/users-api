package config

import (
	"log"

	"github.com/UnendingLoop/users-api/cmd/internal/model"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var models []any = []any{
	&model.User{},
	&model.Friendship{},
}

func ConnectDB(path string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		log.Fatalf("Cannot open db: %v", err)
	}

	if err := db.AutoMigrate(models...); err != nil {
		log.Fatalf("Failed to migrate: %v", err)
	}

	return db
}
func ConnectPostgres(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Cannot open db: %v", err)
	}

	if err := db.AutoMigrate(models...); err != nil {
		log.Fatalf("Failed to migrate: %v", err)
	}
	return db
}
