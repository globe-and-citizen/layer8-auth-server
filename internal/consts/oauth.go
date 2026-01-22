package consts

type OAuthScope string

const (
	OAuthScopeReadUser                OAuthScope = "read:user"
	OAuthScopeReadUserDisplayName     OAuthScope = "read:user:display_name"
	OAuthScopeReadUserColor           OAuthScope = "read:user:color"
	OAuthScopeReadUserBio             OAuthScope = "read:user:bio"
	OAuthScopeReadUserIsEmailVerified OAuthScope = "read:user:is_email_verified"
)

func (s OAuthScope) IsValid() bool {
	switch s {
	case OAuthScopeReadUser,
		OAuthScopeReadUserDisplayName,
		OAuthScopeReadUserColor,
		OAuthScopeReadUserBio,
		OAuthScopeReadUserIsEmailVerified:
		return true
	default:
		return false
	}
}

func OAuthScopesToStringSlice(scopes []OAuthScope) []string {
	result := make([]string, len(scopes))
	for i, scope := range scopes {
		result[i] = string(scope)
	}
	return result
}

// Scope descriptions
var ScopeDescriptions = map[OAuthScope]string{
	OAuthScopeReadUser: "read anonymized information about your account",
}

const (
	OAuthCookieName = "oauth_token"
)

type OAuthErrorCode string

const (
	OAuthErrorInvalidRedirectURI      OAuthErrorCode = "redirect_uri_mismatch"
	OAuthErrorInvalidClient           OAuthErrorCode = "invalid_client"
	OAuthErrorUnauthorizedClient      OAuthErrorCode = "unauthorized_client"
	OAuthErrorAccessDenied            OAuthErrorCode = "access_denied"
	OAuthErrorUnsupportedResponseType OAuthErrorCode = "unsupported_response_type"
	OAuthErrorInvalidScope            OAuthErrorCode = "invalid_scope"
	OAuthErrorServerError             OAuthErrorCode = "server_error"
	OAuthErrorTemporarilyUnavailable  OAuthErrorCode = "temporarily_unavailable"
	OAuthErrorInvalidAuthzCode        OAuthErrorCode = "invalid_authz_code"
)
