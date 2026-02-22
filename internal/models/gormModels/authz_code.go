package gormModels

import (
	"time"
)

type OAuthAuthorizationCode struct {
	ID          uint      `gorm:"column:id;primaryKey;autoIncrement"`
	Code        string    `gorm:"column:code;uniqueIndex;not null"`
	ClientID    string    `gorm:"column:client_id;not null"`
	UserID      uint      `gorm:"column:user_id;not null"`
	RedirectURI string    `gorm:"column:redirect_uri;not null"`
	Scopes      string    `gorm:"column:scopes;not null"`
	Nonce       string    `gorm:"column:nonce"` // if oidc
	ExpiresAt   int64     `gorm:"column:expires_at;not null"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (OAuthAuthorizationCode) TableName() string {
	return "oauth_authorization_codes"
}
