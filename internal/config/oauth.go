package config

import "time"

type OAuthConfig struct {
	JWTSecret         string        `env:"OAUTH_JWT_SECRET"`
	CookieMaxAge      int           `env:"OAUTH_COOKIE_MAX_AGE" default:"3600"` // in seconds
	AuthzCodeSecret   string        `env:"OAUTH_AUTHZ_CODE_SECRET"`
	AuthzCodeExpiry   time.Duration `env:"OAUTH_AUTHZ_CODE_EXPIRY" default:"90s"`
	AccessTokenSecret string        `env:"OAUTH_ACCESS_TOKEN_SECRET"`
	AccessTokenExpiry time.Duration `env:"OAUTH_ACCESS_TOKEN_EXPIRY"`
}
