package config

import "time"

type PostgresConfig struct {
	Host     string `env:"DB_HOST" env-default:"localhost"`
	Port     int    `env:"DB_PORT" env-default:"5432"`
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
	DBName   string `env:"DB_NAME"`
}

type ServerConfig struct {
	Host      string `env:"HOST"`
	Port      int    `env:"PORT"`
	JWTSecret string `env:"JWT_SECRET"`
}

type UserConfig struct {
	ScramIterationCount int `env:"SCRAM_ITERATION_COUNT" env-default:"4096"`
}

type EmailConfig struct {
	ApiKey                 string        `env:"MAILER_SEND_API_KEY"`
	TemplateId             string        `env:"MAILER_SEND_TEMPLATE_ID"`
	Layer8EmailUsername    string        `env:"LAYER8_EMAIL_USERNAME"`
	Layer8EmailDomain      string        `env:"LAYER8_EMAIL_DOMAIN"`
	VerificationCodeExpiry time.Duration `env:"EMAIL_VERIFICATION_CODE_EXPIRY" default:"15m"`
}

type ZkConfig struct {
	GenerateNewZkSnarksKeys bool `env:"GENERATE_NEW_ZK_SNARKS_KEYS" default:"false"`
}

type PhoneConfig struct {
	TelegramApiKey         string        `env:"TELEGRAM_API_KEY"`
	VerificationCodeExpiry time.Duration `env:"PHONE_VERIFICATION_CODE_EXPIRY" default:"10m"`
}
