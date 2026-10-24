package config

import (
	"fmt"
	"globe-and-citizen/layer8/auth-server/pkg/log"
	"globe-and-citizen/layer8/auth-server/pkg/otel"
	"globe-and-citizen/layer8/auth-server/pkg/utils"
	"math/big"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	AppEnv       string   `env:"APP_ENV" required:"true"`
	ServerHost   string   `env:"SERVER_HOST"`
	ServerPort   int      `env:"SERVER_PORT"`
	TokensIssuer string   `env:"TOKENS_ISSUER"`
	CORSOrigins  []string `env:"CORS_ORIGINS" envSeparator:";"`

	Frontend  SPAConfig
	LogConfig log.Config
	Postgres  utils.PostgresConfig
	Scram     ScramConfig
	ZK        ZkConfig
	Payment   PaymentConfig
	OTel      otel.Config

	Debug  DebugConfig
	User   UserConfig
	Client ClientConfig
	OAuth  OAuthConfig
	Worker WorkerConfig
}

type SPAConfig struct {
	StaticAssetsPath string `env:"STATIC_ASSETS_PATH" envDefault:"./web/dist/assets"`
	SPAIndexPath     string `env:"SPA_INDEX_PATH" required:"true" envDefault:"./web/dist/index.html"`
	BaseAPIURL       string `env:"BASE_API_URL" envDefault:""`
	ContractAddress  string `env:"CONTRACT_ADDRESS" required:"true"`
	WalletProjectID  string `env:"WALLET_PROJECT_ID" required:"true"`
}

type ZkConfig struct {
	GenerateNewZkSnarksKeys bool `env:"GENERATE_NEW_ZK_SNARKS_KEYS" envDefault:"true"`
}

type ScramConfig struct {
	ScramIterationCount int `env:"SCRAM_ITERATION_COUNT" envDefault:"4096"`
}

type PaymentConfig struct {
	Currency            string  `env:"PAYMENT_CURRENCY" required:"true" envDefault:"usdc"`
	CostPerBillingUnit  big.Int `env:"COST_PER_BILLING_UNIT" envDefault:"1"`
	BytesPerBillingUnit big.Int `env:"BYTES_PER_BILLING_UNIT" envDefault:"1000"`
}

type DebugConfig struct {
	ClientID []string `env:"DEBUG_CLIENT_ID" envSeparator:";"`
}

type WorkerConfig struct {
	UpdateUsageInterval time.Duration `env:"UPDATE_USAGE_INTERVAL" envDefault:"15m"`
	Payment             PaymentConfig
	Debug               DebugConfig
}

func LoadConfig() AppConfig {
	_ = godotenv.Load()

	cfg := AppConfig{}
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}

	if cfg.Payment.BytesPerBillingUnit.Sign() <= 0 {
		panic(fmt.Errorf("invalid BYTES_PER_BILLING_UNIT=%s: must be > 0", cfg.Payment.BytesPerBillingUnit.String()))
	}

	return cfg
}
