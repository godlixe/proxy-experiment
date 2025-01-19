package internal

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v TimeZone=Asia/Jakarta", dbHost, dbUser, dbPass, dbName, dbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db
}

func CloseDB(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	sqlDB.Close()
}
