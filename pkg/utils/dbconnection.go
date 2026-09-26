package utils

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	gormLogger "gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/tracing"
)

type PostgresConfig struct {
	Host        string `env:"DB_HOST" env-default:"localhost"`
	Port        int    `env:"DB_PORT" env-default:"5432"`
	User        string `env:"DB_USER"`
	Password    string `env:"DB_PASSWORD"`
	DBName      string `env:"DB_NAME"`
	SSLMode     string `env:"DB_SSL_MODE" required:"true"`
	MigratePath string `env:"DB_MIGRATE_PATH" default:"./migrations"`
}

func ConnectDB(postgresConfig PostgresConfig) *gorm.DB {
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
		Logger:                   gormLogger.Default.LogMode(gormLogger.Silent),
	}

	db, err := gorm.Open(postgres.Open(dsn), &gormConfig)
	if err != nil {
		log.Fatalf("Connect to PostgreSQL server failed: %v", err)
	}

	if err := db.Use(tracing.NewPlugin()); err != nil {
		log.Fatalf("Enable OpenTelemetry GORM tracing failed: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Cannot configure PostgreSQL connection: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db
}
