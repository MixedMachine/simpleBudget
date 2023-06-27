package store

import (
	"github.com/mixedmachine/simple-budget-app/internal/models"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"fmt"
	"os"
)

type SQLType string

const (
	SQLITE      SQLType = "sqlite"
	POSTGRES    SQLType = "postgres"
	SQLITE_FILE         = "budget.db"
)

func InitializeSQL(sqlType SQLType, dbLocation string) *gorm.DB {
	var err error
	var DB *gorm.DB

	switch sqlType {
	case SQLITE:
		DB, err = gorm.Open(sqlite.Open(dbLocation), &gorm.Config{})
	case POSTGRES:
		DB, err = gorm.Open(postgres.Open(fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			"budget",
			os.Getenv("DB_PORT"),
		)), &gorm.Config{})
	}

	if err != nil || DB == nil {
		log.Fatal("failed to connect database")
	}

	// TODO: Make this more dynamic
	// IDEA: Use argument to pass in models and iterate through them
	DB.AutoMigrate(&models.Income{})
	DB.AutoMigrate(&models.Expense{})
	DB.AutoMigrate(&models.Allocation{})
	DB.AutoMigrate(&models.Notes{})

	return DB
}
