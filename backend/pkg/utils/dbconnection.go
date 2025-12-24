package utils

import (
	"fmt"
	"globe-and-citizen/layer8/auth-server/backend/config"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDB(postgresConfig config.PostgresConfig) *gorm.DB {
	// PostgreSQL DSN format:
	// postgres://user:password@host:port/dbname?sslmode=disable
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		postgresConfig.User,
		postgresConfig.Password,
		postgresConfig.Host,
		postgresConfig.Port,
		postgresConfig.DBName,
	)

	gormConfig := gorm.Config{
		DisableNestedTransaction: true,
		Logger:                   logger.Default.LogMode(logger.Silent),
	}

	db, err := gorm.Open(postgres.Open(dsn), &gormConfig)
	if err != nil {
		log.Fatalf("Connect to PostgreSQL server failed: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Cannot config PostgreSQL connection: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	return db
}
