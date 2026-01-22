package requestdto

type OAuthUserLoginPrecheck struct {
	UserLoginPrecheck `json:",inline"`
}

type OAuthUserLogin struct {
	UserLogin `json:",inline"`
}

type OAuthAuthorizeContext struct {
	ClientID    string `json:"client_id"`
	RedirectURI string `json:"redirect_uri"`
	Scopes      string `json:"scopes"`
}

type OAuthAuthorizeDecision struct {
	OAuthAuthorizeContext `json:",inline"`
	ReturnResult          bool `json:"return_result" default:"false"`
	Share                 struct {
		Bio             bool `json:"bio" default:"false"`
		Color           bool `json:"color" default:"false"`
		DisplayName     bool `json:"display_name" default:"false"`
		IsEmailVerified bool `json:"is_email_verified" default:"false"`
	}
}

type OAuthAccessToken struct {
	GrantType         string `json:"grant_type" default:"authorization_code"`
	ClientID          string `json:"client_id" validate:"required"`
	ClientSecret      string `json:"client_secret" validate:"required"`
	AuthorizationCode string `json:"code" validate:"required"`
	RedirectURI       string `json:"redirect_uri"`
}

type OAuthZkMetadata struct {
	UserID uint
	Scopes string
}
