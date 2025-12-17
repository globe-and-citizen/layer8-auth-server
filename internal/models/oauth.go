package models

import "github.com/golang-jwt/jwt/v5"

type ClientAccessTokenClaims struct {
	UserID int64
	Scopes string
	jwt.RegisteredClaims
}
