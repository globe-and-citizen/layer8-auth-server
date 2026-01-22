package config

import "time"

type UserConfig struct {
	JWTSecret string `env:"USER_JWT_SECRET"`
	EmailConfig
	PhoneConfig
	ScramConfig // todo remove
}

type EmailConfig struct {
	ApiKey                 string        `env:"MAILER_SEND_API_KEY"`
	TemplateId             string        `env:"MAILER_SEND_TEMPLATE_ID"`
	Layer8EmailUsername    string        `env:"LAYER8_EMAIL_USERNAME"`
	Layer8EmailDomain      string        `env:"LAYER8_EMAIL_DOMAIN"`
	VerificationCodeExpiry time.Duration `env:"EMAIL_VERIFICATION_CODE_EXPIRY" default:"15m"`
}

type PhoneConfig struct {
	TelegramApiKey         string        `env:"TELEGRAM_API_KEY"`
	VerificationCodeExpiry time.Duration `env:"PHONE_VERIFICATION_CODE_EXPIRY" default:"10m"`
}
