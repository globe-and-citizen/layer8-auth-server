package models

import (
	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	Username string `json:"username"`
	UserID   uint   `json:"user_id"`
	jwt.RegisteredClaims
}

type ClientClaims struct {
	Username string `json:"username"`
	ClientID string `json:"user_id"`
	jwt.RegisteredClaims
}

// OAuthAuthenticationClaims represents the claims of an OAuth authentication token
type OAuthAuthenticationClaims struct {
	jwt.RegisteredClaims
}

type ClientAccessTokenClaims struct {
	UserID uint   `json:"user_id"`
	Scopes string `json:"scopes"`
	jwt.RegisteredClaims
}
