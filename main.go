package main

import (
	_ "encoding/hex"
	"fmt"
	"globe-and-citizen/layer8/auth-server/config"
	uH "globe-and-citizen/layer8/auth-server/internal/handlers/user"
	"globe-and-citizen/layer8/auth-server/internal/repositories/postgresRepo"
	"globe-and-citizen/layer8/auth-server/internal/repositories/tokenRepo"
	"globe-and-citizen/layer8/auth-server/internal/usecases"
	apiLog "globe-and-citizen/layer8/auth-server/utils/log"

	"github.com/gin-gonic/gin"
)

var postgresConfig config.PostgresConfig
var serverConfig config.ServerConfig

func readConfig() {

	postgresConfig = config.PostgresConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "layer8",
		Password: "layer81234",
		DBName:   "layer8db",
	}

	serverConfig = config.ServerConfig{
		Host:      "localhost",
		Port:      5002,
		JWTSecret: "5b0b18dc37004b97946367ca5d82673918a6c6e7a817bf84236abe1c0907b9bf",
	}
}

func main() {

	readConfig()

	app := gin.New()
	app.Use(apiLog.AccessLog)

	postgresdb := postgresRepo.NewPostgresRepository(postgresConfig)
	postgresdb.Migrate()

	token := tokenRepo.NewTokenRepository([]byte(serverConfig.JWTSecret))

	userUseCase := usecases.NewUserUseCase(postgresdb, token)
	userHandler := uH.NewUserHandler(app, userUseCase, config.UserConfig{})
	userHandler.RegisterHandler()

	gin.SetMode(gin.ReleaseMode)
	addr := fmt.Sprintf("%s:%d", serverConfig.Host, serverConfig.Port)
	log := apiLog.Get()
	log.Info().Msg("Server start at: http://" + addr)
	err := app.Run(addr)

	if err != nil {
		panic(err)
	}

}
