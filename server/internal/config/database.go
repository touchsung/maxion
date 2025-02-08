package config

import (
	"fmt"
	"os"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func GetDatabaseConnection() (*gorm.DB, error) {
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:1433?database=%s&connection+timeout=30&encrypt=DISABLE",
		dbUser, dbPassword, dbHost, dbName)
	
	return gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
}

