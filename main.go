package main

import (
	_ "encoding/hex"
	"fmt"
	"globe-and-citizen/layer8/auth-server/entity"
	"log"
	"time"

	apiLog "globe-and-citizen/layer8/auth-server/utils/log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var postgresConfig entity.PostgresConfig
var serverConfig entity.ServerConfig

func readConfig() {

	postgresConfig = entity.PostgresConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "layer8",
		Password: "layer81234",
		DBName:   "layer8db",
	}

	serverConfig = entity.ServerConfig{
		Host:      "localhost",
		Port:      5002,
		JWTSecret: "5b0b18dc37004b97946367ca5d82673918a6c6e7a817bf84236abe1c0907b9bf",
	}
}

func SetupDatabase() *gorm.DB {
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

	// Auto migrate tables
	err = db.AutoMigrate( /*add tables here*/ )
	if err != nil {
		log.Fatalf("Cannot migrate tables: %v", err)
	}

	return db
}

func main() {

	readConfig()

	_ = SetupDatabase()

	app := gin.New()
	app.Use(apiLog.AccessLog)

	gin.SetMode(gin.ReleaseMode)
	addr := fmt.Sprintf("%s:%d", serverConfig.Host, serverConfig.Port)
	log := apiLog.Get()
	log.Info().Msg("Server start at: http://" + addr)
	err := app.Run(addr)

	if err != nil {
		panic(err)
	}

}
