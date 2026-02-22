package oauth

import (
	"net/url"

	"golang.org/x/oauth2"
)

type AuthURL struct {
	URL   string
	State string
	Code  string
}

// String returns the URL to redirect the user to for authentication
func (u *AuthURL) String() string {
	urlc, _ := url.Parse(u.URL)
	q := urlc.Query()
	redirectURI := q.Get("redirect_uri")
	q.Del("client_id")
	q.Del("redirect_uri")
	urlc.RawQuery = q.Encode()
	return redirectURI + urlc.String()
}

func GenerateAuthURL(
	clientID, code, clientRedirectURI string, scopes []string, state, nonce string,
) (string, error) {
	config := oauth2.Config{
		ClientID:    clientID,
		RedirectURL: clientRedirectURI,
		Scopes:      scopes,
	}
	var opts = make([]oauth2.AuthCodeOption, 0)
	if nonce != "" {
		opts = append(opts, oauth2.SetAuthURLParam("nonce", nonce))
	}

	authURL := AuthURL{
		URL:   config.AuthCodeURL(state, opts...),
		Code:  code,
		State: state,
	}

	return authURL.String(), nil
}
