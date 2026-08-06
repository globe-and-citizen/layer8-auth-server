package oauthUC

import (
	"fmt"
	"strings"
)

type Scope string

const MaxScopesSize = 50
const AuthorizationCodeSize = 32

const (
	// API scopes
	ScopeReadUser                Scope = "read:user"
	ScopeReadUserDisplayName     Scope = "read:user:display_name"
	ScopeReadUserColor           Scope = "read:user:color"
	ScopeReadUserBio             Scope = "read:user:bio"
	ScopeReadUserIsEmailVerified Scope = "read:user:is_email_verified"
	ScopeReadUserLocation        Scope = "read:user:location"

	// OIDC scopes are standardized, we can define custom ones for our implementation but don't redefine standard scopes
	OIDCScopeOpenID  Scope = "openid"
	OIDCScopeProfile Scope = "profile"
	OIDCScopeEmail   Scope = "email"
)

const DefaultScopes = ScopeReadUser

func isOIDCScope(s Scope) bool {
	switch s {
	case OIDCScopeOpenID, OIDCScopeProfile, OIDCScopeEmail:
		return true
	default:
		return false
	}
}

func isAPIScope(s Scope) bool {
	switch s {
	case ScopeReadUser,
		ScopeReadUserDisplayName,
		ScopeReadUserColor,
		ScopeReadUserBio,
		ScopeReadUserIsEmailVerified,
		ScopeReadUserLocation:
		return true
	default:
		return false
	}
}

func isValidScope(s Scope) bool {
	return isOIDCScope(s) || isAPIScope(s)
}

type Scopes []Scope

func (s Scopes) IsOIDC() bool {
	for _, scope := range s {
		if scope == OIDCScopeOpenID {
			return true
		}
	}
	return false
}

func (s Scopes) Strings() []string {
	strs := make([]string, len(s))
	for i, scope := range s {
		strs[i] = string(scope)
	}
	return strs
}

func (s Scopes) Contains(scope Scope) bool {
	for _, s := range s {
		if s == scope {
			return true
		}
	}
	return false
}

func (s Scopes) APIScopes() Scopes {
	apiScopes := make(Scopes, 0)
	for _, scope := range s {
		if isAPIScope(scope) {
			apiScopes = append(apiScopes, scope)
		}
	}
	return apiScopes
}

func (s Scopes) OIDCScopes() Scopes {
	oidcScopes := make(Scopes, 0)
	for _, scope := range s {
		if isOIDCScope(scope) {
			oidcScopes = append(oidcScopes, scope)
		}
	}
	return oidcScopes
}

func (s Scopes) String() string {
	return strings.Join(s.Strings(), " ")
}

// ValidateScopeStr validates the scope string, remove duplicates, and checks for OIDC compliance if needed.
// returns the parsed scopes, whether it's an OIDC request, and any error encountered.
func ValidateScopeStr(scopeStr string) (Scopes, bool, error) {
	if len(scopeStr) > MaxScopesSize {
		return nil, false, fmt.Errorf("scope too long")
	}

	if strings.TrimSpace(scopeStr) == "" {
		scopeStr = string(DefaultScopes)
	}

	// 1️⃣ Parse (space separated per spec)
	fields := strings.Fields(scopeStr)

	// 2️⃣ Deduplicate
	unique := make(map[Scope]struct{})
	scopes := make(Scopes, 0, len(fields))

	for _, f := range fields {
		s := Scope(f)

		// Basic syntax check (optional but recommended)
		if strings.Contains(f, ",") {
			return nil, false, fmt.Errorf("invalid scope delimiter")
		}

		if _, exists := unique[s]; !exists {
			unique[s] = struct{}{}
			scopes = append(scopes, s)
		}
	}

	// 3️⃣ Detect OIDC
	isOIDC := scopes.IsOIDC()

	// 4️⃣ Validate each scope
	for _, s := range scopes {
		if !isValidScope(s) {
			return nil, false, fmt.Errorf("invalid scope: %s", s)
		}

		if !isOIDC && isOIDCScope(s) {
			return nil, false, fmt.Errorf("oidc scope not allowed without openid: %s", s)
		}
	}

	return scopes, isOIDC, nil
}

// parseScopes is an unsafe version of ValidateScopeStr that assumes the input is already validated and just parses it into a Scopes slice.
// It should only be used internally after validation to avoid redundant checks.
func parseScopes(scopesStr string) Scopes {
	fields := strings.Fields(scopesStr)
	scopes := make(Scopes, 0, len(fields))

	for _, f := range fields {
		scopes = append(scopes, Scope(f))
	}

	return scopes
}

// ScopeDescriptions Scope descriptions
var ScopeDescriptions = map[Scope]string{
	ScopeReadUser:   "Read anonymized information about your account",
	OIDCScopeOpenID: "Do you agree to share your identity with the app? This is required for OIDC authentication.",
}
