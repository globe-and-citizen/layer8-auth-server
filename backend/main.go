package main

import (
	"context"
	_ "encoding/hex"
	"fmt"
	"globe-and-citizen/layer8/auth-server/backend/config"
	"globe-and-citizen/layer8/auth-server/backend/internal/handlers/clientH"
	"globe-and-citizen/layer8/auth-server/backend/internal/handlers/oauthH"
	"globe-and-citizen/layer8/auth-server/backend/internal/handlers/userH"
	"globe-and-citizen/layer8/auth-server/backend/internal/models/gormModels"
	"globe-and-citizen/layer8/auth-server/backend/internal/repositories/codeGenRepo"
	"globe-and-citizen/layer8/auth-server/backend/internal/repositories/emailRepo"
	"globe-and-citizen/layer8/auth-server/backend/internal/repositories/ethRepo"
	"globe-and-citizen/layer8/auth-server/backend/internal/repositories/influxdbRepo"
	"globe-and-citizen/layer8/auth-server/backend/internal/repositories/phoneRepo"
	"globe-and-citizen/layer8/auth-server/backend/internal/repositories/postgresRepo"
	"globe-and-citizen/layer8/auth-server/backend/internal/repositories/tokenRepo"
	"globe-and-citizen/layer8/auth-server/backend/internal/repositories/zkRepo"
	"globe-and-citizen/layer8/auth-server/backend/internal/usecases/clientUC"
	"globe-and-citizen/layer8/auth-server/backend/internal/usecases/oauthUC"
	"globe-and-citizen/layer8/auth-server/backend/internal/usecases/userUC"
	"globe-and-citizen/layer8/auth-server/backend/internal/usecases/workerUC"
	"globe-and-citizen/layer8/auth-server/backend/pkg/code"
	"globe-and-citizen/layer8/auth-server/backend/pkg/eth"
	apiLog "globe-and-citizen/layer8/auth-server/backend/pkg/log"
	"globe-and-citizen/layer8/auth-server/backend/pkg/utils"
	"globe-and-citizen/layer8/auth-server/backend/pkg/zk"
	"log"
	"math/big"
	"os"
	"os/signal"
	"time"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/constraint"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	zlog "github.com/rs/zerolog/log"
)

var postgresConfig config.PostgresConfig
var serverConfig config.ServerConfig
var emailConfig config.EmailConfig
var zkConfig config.ZkConfig
var phoneConfig config.PhoneConfig
var userConfig config.UserConfig
var clientConfig config.ClientConfig
var influxdbConfig config.InfluxDB2Config
var ethConfig config.EthereumConfig

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
		Port:      5001,
		JWTSecret: "5b0b18dc37004b97946367ca5d82673918a6c6e7a817bf84236abe1c0907b9bf",
	}

	influxdbConfig = config.InfluxDB2Config{
		Url:         "http://localhost:8086",
		TelegrafURL: "http://host.docker.internal:8086",
		Username:    "admin",
		Password:    "somethingthatyoudontknow",
		Org:         "layer8",
		Bucket:      "layer8",
		Token:       "DEFAULT_TOKEN_FOR_TESTING",
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

	userConfig = config.UserConfig{
		ScramIterationCount: 4096,
	}

	clientConfig = config.ClientConfig{
		ScramIterationCount: 4096,
		StatsUpdateInterval: time.Second * 30,
		BillingRatePerByte:  big.NewInt(10000000),
	}

	ethConfig = config.EthereumConfig{
		WebsocketRPCURL:     "wss://eth-sepolia.g.alchemy.com/v2/3nAO4RbKoWXgNDwEFWTH5",
		PaymentContractAddr: "0xB3c1aD831A2655D46f0f1Fc0907c543b9bc184E1",
		PaymentContractABI:  "../smart-contract/artifacts/contracts/l8TrafficPayment.sol/L8TrafficPayment.json",
	}

	// TODO: read from env variables or config files
}

func main() {
	readConfig()

	app := gin.Default()
	app.Use(apiLog.AccessLog)
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Vue dev server
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Serve static assets
	app.Static("/assets", "../frontend/dist/assets")
	// SPA fallback
	app.NoRoute(func(c *gin.Context) {
		c.File("../frontend/dist/index.html")
	})

	postgresRepository := postgresRepo.NewPostgresRepository(postgresConfig)
	postgresRepository.Migrate()
	tokenRepository := tokenRepo.NewTokenRepository([]byte(serverConfig.JWTSecret), []byte(serverConfig.JWTSecret), []byte(serverConfig.JWTSecret))
	emailRepository := emailRepo.NewEmailRepository(emailConfig)
	codeGenRepository := codeGenRepo.NewCodeGenerateRepository(code.NewMIMCCodeGenerator())
	zkRepository := zkRepo.NewZkRepository(zkSetup(postgresRepository, zkConfig))
	phoneRepository := phoneRepo.NewPhoneRepository(phoneConfig)
	influxdbRepository := influxdbRepo.NewInfluxdbRepository(influxdbConfig)
	err := influxdbRepository.IsConnected(&gin.Context{})
	if err != nil {
		panic(err)
	}

	client, err := eth.ConnectToEthereum(ethConfig.WebsocketRPCURL)
	if err != nil {
		log.Fatalf("Connect to %s failed: %w\n", ethConfig.WebsocketRPCURL, err)
	}
	defer eth.CloseEthereumConnection(client)
	ethRepository := ethRepo.NewEthereumRepository(client, ethConfig)

	userUsecase := userUC.NewUserUsecase(
		postgresRepository,
		tokenRepository,
		emailRepository,
		codeGenRepository,
		zkRepository,
		phoneRepository,
	)
	clientUsecase := clientUC.NewClientUsecase(
		postgresRepository,
		tokenRepository,
		influxdbRepository,
	)
	oauthUsecase := oauthUC.NewOAuthUsecase(postgresRepository, tokenRepository)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	workerUsecase := workerUC.NewWorkerUsecase(ctx, postgresRepository, influxdbRepository, ethRepository)

	go func() {
		ticker := time.NewTicker(clientConfig.StatsUpdateInterval)

		for currTime := range ticker.C {
			zlog.Info().Msgf("Started usage balance updater with interval: %s", clientConfig.StatsUpdateInterval)
			err = workerUsecase.UpdateUsageBalance(clientConfig.BillingRatePerByte, currTime)
			if err != nil {
				zlog.Error().Err(err).Msg("Error while updating usage balance")
			}
		}
	}()
	go workerUsecase.ListenToEthereumEvents()

	apiGroup := app.Group("/api/v1")
	userHandler := userH.NewUserHandler(apiGroup, userUsecase, userConfig)
	userHandler.RegisterAPIs()
	clientHandler := clientH.NewClientHandler(apiGroup, config.ClientConfig{}, clientUsecase)
	clientHandler.RegisterAPIs()
	oauthHandler := oauthH.NewOAuthHandler(apiGroup, config.OAuthConfig{CookieMaxAge: 3600}, oauthUsecase)
	oauthHandler.RegisterAPIs()

	gin.SetMode(gin.ReleaseMode)
	addr := fmt.Sprintf("%s:%d", serverConfig.Host, serverConfig.Port)
	log := apiLog.Get()
	log.Info().Msg("Server start at: http://" + addr)
	err = app.Run(addr)
	if err != nil {
		panic(err)
	}
}

func zkSetup(postgresRepository postgresRepo.IPostgresRepository, zkConfig config.ZkConfig) zk.IProofProcessor {
	log := apiLog.Get()
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
			log.Fatal().Err(err)
		}
	} else {
		zkSnarksKeyPair, err := postgresRepository.GetLatestZkSnarksKeys()
		if err != nil {
			log.Fatal().Errs("Error while reading zk-snarks keys from the database", []error{err})
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
