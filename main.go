package main

import (
	_ "encoding/hex"
	"fmt"
	"globe-and-citizen/layer8/auth-server/config"
	"globe-and-citizen/layer8/auth-server/internal/handlers/tokenHandler"
	"globe-and-citizen/layer8/auth-server/internal/handlers/userHandler"
	"globe-and-citizen/layer8/auth-server/internal/models/gormModels"
	"globe-and-citizen/layer8/auth-server/internal/repositories/codeGenRepo"
	"globe-and-citizen/layer8/auth-server/internal/repositories/emailRepo"
	"globe-and-citizen/layer8/auth-server/internal/repositories/phoneRepo"
	"globe-and-citizen/layer8/auth-server/internal/repositories/postgresRepo"
	"globe-and-citizen/layer8/auth-server/internal/repositories/tokenRepo"
	"globe-and-citizen/layer8/auth-server/internal/repositories/zkRepo"
	"globe-and-citizen/layer8/auth-server/internal/usecases/tokenUsecase"
	"globe-and-citizen/layer8/auth-server/internal/usecases/userUsecase"
	"globe-and-citizen/layer8/auth-server/pkg/code"
	apiLog "globe-and-citizen/layer8/auth-server/pkg/log"
	"globe-and-citizen/layer8/auth-server/pkg/utils"
	"globe-and-citizen/layer8/auth-server/pkg/zk"
	"log"
	"os"
	"time"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/constraint"
	"github.com/gin-gonic/gin"
)

var postgresConfig config.PostgresConfig
var serverConfig config.ServerConfig
var emailConfig config.EmailConfig
var zkConfig config.ZkConfig
var phoneConfig config.PhoneConfig

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

	emailConfig = config.EmailConfig{
		VerificationCodeExpiry: time.Minute * 10,
	}

	zkConfig = config.ZkConfig{
		GenerateNewZkSnarksKeys: true,
	}

	phoneConfig = config.PhoneConfig{
		TelegramApiKey:         os.Getenv("TELEGRAM_API_KEY"),
		VerificationCodeExpiry: time.Minute * 10,
	}

	// TODO: read from env variables or config files
}

func zkSetup(postgresRepository postgresRepo.IPostgresRepository, zkConfig config.ZkConfig) zk.IProofProcessor {
	var cs constraint.ConstraintSystem
	var zkKeyPairId uint
	var provingKey groth16.ProvingKey
	var verifyingKey groth16.VerifyingKey
	var err error

	if zkConfig.GenerateNewZkSnarksKeys {
		cs, provingKey, verifyingKey = zk.RunZkSnarksSetup()

		zkKeyPairId, err = postgresRepository.SaveZkSnarksKeyPair(
			gormModels.ZkSnarksKeyPair{
				ProvingKey:   utils.WriteBytes(provingKey),
				VerifyingKey: utils.WriteBytes(verifyingKey),
			},
		)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		zkSnarksKeyPair, err := postgresRepository.GetLatestZkSnarksKeys()
		if err != nil {
			log.Fatalf("Error while reading zk-snarks keys from the database: %e", err)
		}

		cs = zk.GenerateConstraintSystem()
		zkKeyPairId = zkSnarksKeyPair.ID

		// Empty proving key initialised with elliptic curve id
		provingKey = groth16.NewProvingKey(ecc.BN254)
		// Deserialize proving key representation bytes from db into the provingKey object
		utils.ReadBytes[groth16.ProvingKey](provingKey, zkSnarksKeyPair.ProvingKey)

		// Empty verifying key initialised with elliptic curve id
		verifyingKey = groth16.NewVerifyingKey(ecc.BN254)
		// Deserialize verifying key representation bytes from db into the verifyingKey object
		utils.ReadBytes[groth16.VerifyingKey](verifyingKey, zkSnarksKeyPair.VerifyingKey)
	}

	return zk.NewProofProcessor(cs, zkKeyPairId, provingKey, verifyingKey)
}

func main() {

	readConfig()

	app := gin.New()
	app.Use(apiLog.AccessLog)
	router := app.Group("/api/v1")

	postgresRepository := postgresRepo.NewPostgresRepository(postgresConfig)
	postgresRepository.Migrate()

	tokenRepository := tokenRepo.NewTokenRepository([]byte(serverConfig.JWTSecret))
	tokenUC := tokenUsecase.NewTokenUseCase(tokenRepository)
	tokenH := tokenHandler.NewTokenHandler(tokenUC)

	emailRepository := emailRepo.NewEmailRepository(emailConfig)
	codeGenRepository := codeGenRepo.NewCodeGenerateRepository(code.NewMIMCCodeGenerator())
	zkRepository := zkRepo.NewZkRepository(zkSetup(postgresRepository, zkConfig))
	phoneRepository := phoneRepo.NewPhoneRepository(phoneConfig)

	userUC := userUsecase.NewUserUsecase(
		postgresRepository,
		tokenRepository,
		emailRepository,
		codeGenRepository,
		zkRepository,
		phoneRepository,
	)
	userH := userHandler.NewUserHandler(router, userUC, config.UserConfig{})
	userH.RegisterHandler(tokenH.UserAuthentication)

	gin.SetMode(gin.ReleaseMode)
	addr := fmt.Sprintf("%s:%d", serverConfig.Host, serverConfig.Port)
	log := apiLog.Get()
	log.Info().Msg("Server start at: http://" + addr)
	err := app.Run(addr)

	if err != nil {
		panic(err)
	}

}
