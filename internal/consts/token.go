package consts

import "time"

const (
	TokenTypeBearer        = "Bearer"
	UserLoginTokenExpiry   = 60 * time.Minute
	ClientLoginTokenExpiry = 60 * time.Minute
	OAuthLoginTokenExpiry  = 10 * time.Minute
)
