package config

import (
	"math/big"
	"time"
)

type ClientConfig struct {
	JWTSecret           string        `env:"CLIENT_JWT_SECRET"`
	UpdateUsageInterval time.Duration `env:"UPDATE_USAGE_INTERVAL" default:"15m"`
	BillingRatePerByte  *big.Int      `env:"BILLING_RATE_PER_BYTE" default:"1"`
	InfluxDB2Config
	Web3Config
	ScramConfig // todo remove
}

type InfluxDB2Config struct {
	Url      string `env:"INFLUXDB_URL" default:"http://localhost:8086"`
	Username string `env:"INFLUXDB_USERNAME" default:"admin"`
	Password string `env:"INFLUXDB_PASSWORD" default:"somethingthatyoudontknow"`
	Org      string `env:"INFLUXDB_ORG" default:"layer8"`
	Bucket   string `env:"INFLUXDB_BUCKET" default:"layer8"`
	Token    string `env:"INFLUXDB_TOKEN" default:"DEFAULT_TOKEN_FOR_TESTING"`
}

type ExternalConfig struct {
	// todo add authenticate forward proxy config
}

type Web3Config struct {
	WebsocketRPCURL     string `env:"WEB3_WS_RPC_URL"`
	PaymentContractAddr string `env:"WEB3_PAYMENT_CONTRACT_ADDRESS"`
	PaymentContractABI  string `env:"WEB3_PAYMENT_CONTRACT_ABI"`
	// other contracts addr and abi
}
