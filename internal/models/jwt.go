package models

import (
	"time"

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

type OAuthAccessTokenClaims struct {
	UserID uint   `json:"user_id"`
	Scopes string `json:"scopes"`
	jwt.RegisteredClaims
}

type OAuthIDTokenClaims struct {
	jwt.RegisteredClaims
	AuthTime time.Time `json:"auth_time,omitempty"`
	Nonce    string    `json:"nonce,omitempty"`
	OIDCUserProfile
}

type OIDCUserProfile struct {
	DisplayName string `json:"display_name"`
	Bio         string `json:"bio"`
}
