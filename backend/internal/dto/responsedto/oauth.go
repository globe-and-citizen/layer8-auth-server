package responsedto

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

type OAuthAuthorizeDecision struct {
	RedirectURI string `json:"redirect_uri"`
	Code        string `json:"code"`
}

type OAuthAccessToken struct {
	AccessToken     string `json:"access_token"`
	TokenType       string `json:"token_type"`
	ExpireInMinutes int    `json:"expire_in_minutes"`
}

type OAuthZkMetadata struct {
	IsEmailVerified bool   `json:"is_email_verified"`
	DisplayName     string `json:"display_name"`
	Color           string `json:"color"`
	Bio             string `json:"bio"`
}
