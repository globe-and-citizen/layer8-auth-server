package config

import (
	"globe-and-citizen/layer8/auth-server/pkg/log"
	"globe-and-citizen/layer8/auth-server/pkg/otel"
	"globe-and-citizen/layer8/auth-server/pkg/utils"
	"math/big"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	AppEnv      string   `env:"APP_ENV" required:"true"`
	ServerHost  string   `env:"SERVER_HOST"`
	ServerPort  int      `env:"SERVER_PORT"`
	ServerDNS   string   `env:"SERVER_DNS"`
	CORSOrigins []string `env:"CORS_ORIGINS" envSeparator:";"`
	SPAConfig

	LogConfig log.Config

	utils.PostgresConfig
	ScramConfig
	ZkConfig
	PaymentConfig PaymentConfig

	User       UserConfig
	Client     ClientConfig
	OAuth      OAuthConfig
	Worker     WorkerConfig
	OTelConfig otel.Config
}

type SPAConfig struct {
	StaticAssetsPath string `env:"STATIC_ASSETS_PATH" env-default:"./web/dist/assets"`
	SPAIndexPath     string `env:"SPA_INDEX_PATH" required:"true" env-default:"./web/dist/index.html"`
	BaseAPIURL       string `env:"BASE_API_URL" env-default:""`
	ContractAddress  string `env:"CONTRACT_ADDRESS" required:"true"`
	WalletProjectID  string `env:"WALLET_PROJECT_ID" required:"true"`
}

type ZkConfig struct {
	GenerateNewZkSnarksKeys bool `env:"GENERATE_NEW_ZK_SNARKS_KEYS" env-default:"true"`
}

type ScramConfig struct {
	ScramIterationCount int `env:"SCRAM_ITERATION_COUNT" env-default:"4096"`
}

type PaymentConfig struct {
	Currency           string  `env:"PAYMENT_CURRENCY" required:"true" env-default:"usdc"`
	BillingRatePerByte big.Int `env:"BILLING_RATE_PER_BYTE" default:"1"`
}

type WorkerConfig struct {
	UpdateUsageInterval time.Duration `env:"UPDATE_USAGE_INTERVAL" env-default:"15m"`
	Payment             PaymentConfig
}

func LoadConfig() AppConfig {
	_ = godotenv.Load()

	cfg := AppConfig{}
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}

	return cfg
}
