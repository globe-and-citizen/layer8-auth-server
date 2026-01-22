package config

import "time"

type OAuthConfig struct {
	JWTSecret       string        `env:"OAUTH_JWT_SECRET"`
	CookieMaxAge    int           `env:"OAUTH_COOKIE_MAX_AGE" default:"3600"` // in seconds
	AuthzCodeExpiry time.Duration `env:"OAUTH_AUTHZ_CODE_EXPIRY" default:"10m"`
}
