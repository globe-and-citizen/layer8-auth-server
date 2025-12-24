package responsedto

type OauthToken struct {
	AccessToken      string `json:"access_token"`
	TokenType        string `json:"token_type"`
	ExpiresInMinutes int    `json:"expires_in_minutes"`
}
