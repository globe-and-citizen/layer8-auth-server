package main

import (
	"context"
	_ "encoding/hex"
	"fmt"
	"globe-and-citizen/layer8/auth-server/internal/config"
	"globe-and-citizen/layer8/auth-server/internal/handlers/clientH"
	"globe-and-citizen/layer8/auth-server/internal/handlers/oauthH"
	"globe-and-citizen/layer8/auth-server/internal/handlers/userH"
	"globe-and-citizen/layer8/auth-server/internal/models/gormModels"
	"globe-and-citizen/layer8/auth-server/internal/repositories/codeGenRepo"
	"globe-and-citizen/layer8/auth-server/internal/repositories/emailRepo"
	"globe-and-citizen/layer8/auth-server/internal/repositories/ethRepo"
	"globe-and-citizen/layer8/auth-server/internal/repositories/influxdbRepo"
	"globe-and-citizen/layer8/auth-server/internal/repositories/phoneRepo"
	"globe-and-citizen/layer8/auth-server/internal/repositories/postgresRepo"
	"globe-and-citizen/layer8/auth-server/internal/repositories/tokenRepo"
	"globe-and-citizen/layer8/auth-server/internal/repositories/zkRepo"
	"globe-and-citizen/layer8/auth-server/internal/usecases/clientUC"
	"globe-and-citizen/layer8/auth-server/internal/usecases/oauthUC"
	"globe-and-citizen/layer8/auth-server/internal/usecases/userUC"
	"globe-and-citizen/layer8/auth-server/internal/usecases/workerUC"
	"globe-and-citizen/layer8/auth-server/pkg/code"
	"globe-and-citizen/layer8/auth-server/pkg/eth"
	"globe-and-citizen/layer8/auth-server/pkg/ginUtils"
	"globe-and-citizen/layer8/auth-server/pkg/log"
	"globe-and-citizen/layer8/auth-server/pkg/utils"
	"globe-and-citizen/layer8/auth-server/pkg/zk"
	"os"
	"os/signal"
	"time"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/constraint"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	appConfig := config.LoadConfig()
	logger := log.NewLogger(appConfig.LogConfig)

	app := gin.New()
	app.Use(ginUtils.RequestID, gin.Recovery(), ginUtils.AccessLog(logger))
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*.layer8proxy.net", "localhost:*"}, // Vue dev server
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	app.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	app.GET("/config.js", serveFrontendConfig(appConfig.SPAConfig))
	// Serve static assets
	app.Static("/assets", appConfig.StaticAssetsPath)
	// Serve SPA index.html for all other routes (to support client-side routing)
	app.NoRoute(func(c *gin.Context) {
		c.File(appConfig.SPAIndexPath)
	})

	postgresRepository := postgresRepo.NewPostgresRepository(appConfig.PostgresConfig)
	//postgresRepository.Migrate()
	tokenRepository := tokenRepo.NewTokenRepository(
		appConfig.ServerDNS,
		[]byte(appConfig.UserConfig.JWTSecret),
		[]byte(appConfig.ClientConfig.JWTSecret),
		[]byte(appConfig.OAuthConfig.JWTSecret),
	)
	emailRepository := emailRepo.NewEmailRepository(appConfig.EmailConfig)
	codeGenRepository := codeGenRepo.NewCodeGenerateRepository(code.NewMIMCCodeGenerator())
	zkRepository := zkRepo.NewZkRepository(zkSetup(postgresRepository, appConfig.ZkConfig))
	phoneRepository := phoneRepo.NewPhoneRepository(appConfig.PhoneConfig)
	influxdbRepository := influxdbRepo.NewInfluxdbRepository(appConfig.InfluxDB2Config)
	err := influxdbRepository.IsConnected(&gin.Context{})
	if err != nil {
		panic(err)
	}

	client, err := eth.ConnectToEthereum(appConfig.Web3Config.WebsocketRPCURL)
	if err != nil {
		panic(fmt.Errorf("failed to connect to %s: %w", appConfig.Web3Config.WebsocketRPCURL, err))
	}
	defer eth.CloseEthereumConnection(client)
	ethRepository := ethRepo.NewEthereumRepository(logger, client, appConfig.Web3Config)

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
	oauthUsecase := oauthUC.NewOAuthUsecase(appConfig.OAuthConfig, postgresRepository, tokenRepository)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	workerUsecase := workerUC.NewWorkerUsecase(logger, ctx, postgresRepository, influxdbRepository, ethRepository)

	go func() {
		ticker := time.NewTicker(appConfig.UpdateUsageInterval)

		for currTime := range ticker.C {
			logger.Info(fmt.Sprintf("Update usage balance with interval: %s", appConfig.UpdateUsageInterval))
			err = workerUsecase.UpdateUsageBalance(appConfig.BillingRatePerByte, currTime)
			if err != nil {
				logger.Error("Error while updating usage balance", err)
			}
		}
	}()
	go workerUsecase.ListenToEthereumEvents()

	apiGroup := app.Group("/api/v1")
	userHandler := userH.NewUserHandler(logger, apiGroup, userUsecase, appConfig.UserConfig)
	userHandler.RegisterAPIs()
	clientHandler := clientH.NewClientHandler(logger, apiGroup, config.ClientConfig{}, clientUsecase)
	clientHandler.RegisterAPIs()
	oauthHandler := oauthH.NewOAuthHandler(logger, apiGroup, config.OAuthConfig{CookieMaxAge: 3600}, oauthUsecase)
	oauthHandler.RegisterAPIs()

	gin.SetMode(gin.ReleaseMode)
	addr := fmt.Sprintf("%s:%d", appConfig.ServerHost, appConfig.ServerPort)
	logger.Info("Server start at: http://" + addr)
	err = app.Run(addr)
	if err != nil {
		panic(err)
	}
}

func serveFrontendConfig(cfg config.SPAConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/javascript")
		c.String(200, `
window.__APP_CONFIG__ = {
  BASE_API_URL: %q,
  CONTRACT_ADDRESS: %q,
  WALLET_PROJECT_ID: %q
};
`, cfg.BaseAPIURL, cfg.ContractAddress, cfg.WalletProjectID)
	}
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
			panic(err)
		}
	} else {
		zkSnarksKeyPair, err := postgresRepository.GetLatestZkSnarksKeys()
		if err != nil {
			panic(fmt.Errorf("get latest zk snarks keys failed: %w", err))
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
