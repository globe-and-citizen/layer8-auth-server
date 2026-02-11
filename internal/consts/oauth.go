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
