package main

import (
	_ "encoding/hex"
	"fmt"
	"globe-and-citizen/layer8/auth-server/entity"
	"globe-and-citizen/layer8/auth-server/internal/repositories/postgresRepo"
	apiLog "globe-and-citizen/layer8/auth-server/utils/log"

	"github.com/gin-gonic/gin"
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

func main() {

	readConfig()

	repo := postgresRepo.NewPostgresRepository(postgresConfig)
	repo.Migrate()

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
