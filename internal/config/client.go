package config

type ClientConfig struct {
	JWTSecret string `env:"CLIENT_JWT_SECRET"`
	Payment   PaymentConfig
	InfluxDB2Config
	Web3Config  Web3Config
	ScramConfig // todo remove
}

type InfluxDB2Config struct {
	Url    string `env:"INFLUXDB_URL" envDefault:"http://localhost:8086"`
	Org    string `env:"INFLUXDB_ORG" envDefault:"layer8org"`
	Bucket string `env:"INFLUXDB_BUCKET" envDefault:"layer8bucket"`
	Token  string `env:"INFLUXDB_TOKEN" envDefault:"default_token"`
}

type ExternalConfig struct {
	// todo add authenticate forward proxy config
}

type Web3Config struct {
	WebsocketRPCURL string `env:"WEB3_WS_RPC_URL" required:"true"`

	TrafficPaidEnabled  bool   `env:"WEB3_TRAFFIC_PAID_ENABLED" envDefault:"false"`
	PaymentContractAddr string `env:"WEB3_PAYMENT_CONTRACT_ADDRESS" required:"true"`
	PaymentContractABI  string `env:"WEB3_PAYMENT_CONTRACT_ABI" required:"true"`

	// CNC payment
	CNCPaymentEnabled bool   `env:"WEB3_CNC_PAYMENT_ENABLED" envDefault:"false"`
	CNCContractAddr   string `env:"WEB3_CNC_CONTRACT_ADDRESS" required:"true"`
	CNCEventABI       string `env:"WEB3_CNC_EVENT_ABI" required:"true"`

	// other contracts addr and abi
}
