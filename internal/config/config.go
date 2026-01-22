package config

import "globe-and-citizen/layer8/auth-server/pkg/utils"

type AppConfig struct {
	AppEnv string `env:"APP_ENV" required:"true"`
	Host   string `env:"HOST"`
	Port   int    `env:"PORT"`
	SPAConfig

	utils.PostgresConfig
	ScramConfig
	ZkConfig

	UserConfig
	ClientConfig
	OAuthConfig
}

type SPAConfig struct {
	StaticAssetsPath string `env:"STATIC_ASSETS_PATH" env-default:"./web/dist/assets"`
	SPAIndexPath     string `env:"SPA_INDEX_PATH" required:"true" env-default:"./web/dist/index.html"`
}

type ZkConfig struct {
	GenerateNewZkSnarksKeys bool `env:"GENERATE_NEW_ZK_SNARKS_KEYS" env-default:"true"`
}

type ScramConfig struct {
	ScramIterationCount int `env:"SCRAM_ITERATION_COUNT" env-default:"4096"`
}
