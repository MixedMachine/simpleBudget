package store

import (
	"github.com/mixedmachine/simple-budget-app/models"

	log "github.com/sirupsen/logrus"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"fmt"
	"os"
)

type SQLType string

const (
	SQLITE   SQLType = "sqlite"
	POSTGRES SQLType = "postgres"
)

func InitializeSQL(sqlType SQLType) *gorm.DB {
	var err error
	var DB *gorm.DB

	switch sqlType {
	case SQLITE:
		DB, err = gorm.Open(sqlite.Open("budget.db"), &gorm.Config{})
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

	if err != nil {
		log.Fatal("failed to connect database")
	}

	DB.AutoMigrate(&models.Income{})
	DB.AutoMigrate(&models.Expense{})
	DB.AutoMigrate(&models.Allocation{})

	return DB
}
