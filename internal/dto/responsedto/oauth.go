package responsedto

import "time"

type OAuthUserLoginPrecheck struct {
	UserLoginPrecheck `json:",inline"`
}

type OAuthUserLogin struct {
	UserLogin `json:",inline"`
}

type OAuthAuthorizeContext struct {
	ClientName string                 `json:"client_name"`
	Scopes     []OAuthAuthorizeScopes `json:"scopes"`
}

type OAuthAuthorizeScopes struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type OAuthAuthorizeConsent struct {
	RedirectURI string `json:"redirect_uri"`
	Code        string `json:"code"`
}

type OAuthTokenRequest struct {
	AccessToken      string `json:"access_token"`
	TokenType        string `json:"token_type"`
	ExpiresInMinutes int    `json:"expires_in_minutes"`
	IDToken          string `json:"id_token,omitempty"`
	RefreshToken     string `json:"refresh_token,omitempty"`
}

type OAuthZkMetadata struct {
	IsEmailVerified        *bool      `json:"is_email_verified,omitempty"`
	EmailVerifiedAt        *time.Time `json:"email_verified_at,omitempty"`
	DisplayName            *string    `json:"display_name,omitempty"`
	DisplayNameUpdatedAt   *time.Time `json:"display_name_updated_at,omitempty"`
	FavoriteColor          *string    `json:"favorite_color,omitempty"`
	FavoriteColorUpdatedAt *time.Time `json:"favorite_color_updated_at,omitempty"`
	Bio                    *string    `json:"bio,omitempty"`
	BioUpdatedAt           *time.Time `json:"bio_updated_at,omitempty"`
	Location               *string    `json:"location,omitempty"`
	LocationUpdatedAt      *time.Time `json:"location_updated_at,omitempty"`
}
