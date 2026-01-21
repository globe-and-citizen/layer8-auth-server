package config

type AppConfig struct {
	AppEnv string `env:"APP_ENV" required:"true"`
	Host   string `env:"HOST"`
	Port   int    `env:"PORT"`
	SPAConfig

	PostgresConfig
	ScramConfig
	ZkConfig

	UserConfig
	ClientConfig
	OAuthConfig
}

type SPAConfig struct {
	StaticAssetsPath string `env:"STATIC_ASSETS_PATH" env-default:"./frontend/dist/assets"`
	SPAIndexPath     string `env:"SPA_INDEX_PATH" required:"true" env-default:"../frontend/dist/index.html"`
}

type PostgresConfig struct {
	Host        string `env:"DB_HOST" env-default:"localhost"`
	Port        int    `env:"DB_PORT" env-default:"5432"`
	User        string `env:"DB_USER"`
	Password    string `env:"DB_PASSWORD"`
	DBName      string `env:"DB_NAME"`
	SSLMode     string `env:"DB_SSL_MODE" required:"true"`
	MigratePath string `env:"DB_MIGRATE_PATH" default:"./migrations"`
}

type ZkConfig struct {
	GenerateNewZkSnarksKeys bool `env:"GENERATE_NEW_ZK_SNARKS_KEYS" env-default:"true"`
}

type ScramConfig struct {
	ScramIterationCount int `env:"SCRAM_ITERATION_COUNT" env-default:"4096"`
}
