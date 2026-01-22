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
	log2 "globe-and-citizen/layer8/auth-server/pkg/log"
	"globe-and-citizen/layer8/auth-server/pkg/utils"
	zk2 "globe-and-citizen/layer8/auth-server/pkg/zk"
	"log"
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

func main() {
	appConfig := config.LoadConfig()

	app := gin.Default()
	app.Use(log2.AccessLog)
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Vue dev server
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	app.GET("/config.js", serveFrontendConfig(appConfig.SPAConfig))
	// Serve static assets
	app.Static("/assets", appConfig.StaticAssetsPath)
	// Serve SPA index.html for all other routes (to support client-side routing)
	app.NoRoute(func(c *gin.Context) {
		c.File(appConfig.SPAIndexPath)
	})

	postgresRepository := postgresRepo.NewPostgresRepository(appConfig.PostgresConfig)
	postgresRepository.Migrate()
	tokenRepository := tokenRepo.NewTokenRepository(
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
		log.Fatalf("Connect to %s failed: %w\n", appConfig.Web3Config.WebsocketRPCURL, err)
	}
	defer eth.CloseEthereumConnection(client)
	ethRepository := ethRepo.NewEthereumRepository(client, appConfig.Web3Config)

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
		ticker := time.NewTicker(appConfig.UpdateUsageInterval)

		for currTime := range ticker.C {
			zlog.Info().Msgf("Started usage balance updater with interval: %s", appConfig.UpdateUsageInterval)
			err = workerUsecase.UpdateUsageBalance(appConfig.BillingRatePerByte, currTime)
			if err != nil {
				zlog.Error().Err(err).Msg("Error while updating usage balance")
			}
		}
	}()
	go workerUsecase.ListenToEthereumEvents()

	apiGroup := app.Group("/api/v1")
	userHandler := userH.NewUserHandler(apiGroup, userUsecase, appConfig.UserConfig)
	userHandler.RegisterAPIs()
	clientHandler := clientH.NewClientHandler(apiGroup, config.ClientConfig{}, clientUsecase)
	clientHandler.RegisterAPIs()
	oauthHandler := oauthH.NewOAuthHandler(apiGroup, config.OAuthConfig{CookieMaxAge: 3600}, oauthUsecase)
	oauthHandler.RegisterAPIs()

	gin.SetMode(gin.ReleaseMode)
	addr := fmt.Sprintf("%s:%d", appConfig.Host, appConfig.Port)
	log := log2.Get()
	log.Info().Msg("Server start at: http://" + addr)
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
  CONTRACT_ADDRESS: %q,
  WALLET_PROJECT_ID: %q
};
`, cfg.ContractAddress, cfg.WalletProjectID)
	}
}

func zkSetup(postgresRepository postgresRepo.IPostgresRepository, zkConfig config.ZkConfig) zk2.IProofProcessor {
	log := log2.Get()
	var cs constraint.ConstraintSystem
	var zkKeyPairId uint
	var provingKey groth16.ProvingKey
	var verifyingKey groth16.VerifyingKey
	var err error

	if zkConfig.GenerateNewZkSnarksKeys {
		cs, provingKey, verifyingKey = zk2.RunZkSnarksSetup()

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

		cs = zk2.GenerateConstraintSystem()
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

	return zk2.NewProofProcessor(cs, zkKeyPairId, provingKey, verifyingKey)
}
