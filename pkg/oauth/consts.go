package oauth

type GrantType string

const (
	GrantTypeAuthorizationCode GrantType = "authorization_code"
	GrantTypePassword          GrantType = "password"
	GrantTypeClientCredentials GrantType = "client_credentials"
)

const ResponseTypeCode = "code"

const (
	ParamResponseType = "response_type"
	ParamClientID     = "client_id"
	ParamRedirectURI  = "redirect_uri"
	ParamScope        = "scope"
	ParamState        = "state"
	ParamNonce        = "nonce"
)
