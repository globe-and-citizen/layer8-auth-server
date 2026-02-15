package oauthUC

type OAuthScope string

const (
	ScopeReadUser                OAuthScope = "read:user"
	ScopeReadUserDisplayName     OAuthScope = "read:user:display_name"
	ScopeReadUserColor           OAuthScope = "read:user:color"
	ScopeReadUserBio             OAuthScope = "read:user:bio"
	ScopeReadUserIsEmailVerified OAuthScope = "read:user:is_email_verified"
)

func (s OAuthScope) IsValid() bool {
	switch s {
	case ScopeReadUser,
		ScopeReadUserDisplayName,
		ScopeReadUserColor,
		ScopeReadUserBio,
		ScopeReadUserIsEmailVerified:
		return true
	default:
		return false
	}
}

func ScopesToStringSlice(scopes []OAuthScope) []string {
	result := make([]string, len(scopes))
	for i, scope := range scopes {
		result[i] = string(scope)
	}
	return result
}

// Scope descriptions
var ScopeDescriptions = map[OAuthScope]string{
	ScopeReadUser: "read anonymized information about your account",
}
