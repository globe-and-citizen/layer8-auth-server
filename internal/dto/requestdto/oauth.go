package requestdto

type OAuthUserLoginPrecheck struct {
	UserLoginPrecheck `json:",inline"`
}

type OAuthUserLogin struct {
	UserLogin `json:",inline"`
}

// OAuthAuthorizeQueries fields are from query params
type OAuthAuthorizeQueries struct {
	ResponseType string `json:"response_type" validate:"required"`
	ClientID     string `json:"client_id" validate:"required"`
	RedirectURI  string `json:"redirect_uri" validate:"required"`
	Scopes       string `json:"scopes" validate:"required"`
	State        string `json:"state"`
	Nonce        string `json:"nonce,omitempty"`
}

type OAuthAuthorizeConsent struct {
	OAuthAuthorizeQueries `json:",inline"`
	OIDCAgreedToShare     bool `json:"oidc_agreed" default:"false"`
	Share                 struct {
		Bio             bool `json:"bio" default:"false"`
		Color           bool `json:"color" default:"false"`
		DisplayName     bool `json:"display_name" default:"false"`
		IsEmailVerified bool `json:"is_email_verified" default:"false"`
	}
}

type OAuthTokenRequest struct {
	GrantType         string `json:"grant_type" validate:"required"`
	ClientID          string `json:"client_id" validate:"required"`
	ClientSecret      string `json:"client_secret" validate:"required"`
	AuthorizationCode string `json:"code" validate:"required"`
	RedirectURI       string `json:"redirect_uri"`
}

type OAuthZkMetadata struct {
	UserID uint
	Scopes string
}
