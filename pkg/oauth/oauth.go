package oauth

import (
	"fmt"
	"globe-and-citizen/layer8/auth-server/pkg/utils"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
)

type AuthorizationCodeClaims struct {
	ClientID    string            `json:"cid"`
	UserID      uint              `json:"uid"`
	RedirectURI string            `json:"ruri"`
	Scopes      string            `json:"scp"`
	ExpiresAt   int64             `json:"exp"`
	HeaderMap   map[string]string `json:"hmap"`
	jwt.RegisteredClaims
}

type AuthURL struct {
	URL          string
	State        string
	Code         string
	CodeVerifier string
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
	clientID string,
	code string,
	clientRedirectURI string,
	scopes []string,
) (string, error) {
	state, stateErr := utils.GenerateRandomBase64String(24) // todo why 24?
	if stateErr != nil {
		return "", fmt.Errorf("could not generate random state: %v", stateErr)
	}

	config := oauth2.Config{
		ClientID:    clientID,
		RedirectURL: clientRedirectURI,
		Scopes:      scopes,
	}

	authURL := AuthURL{
		URL: config.AuthCodeURL(
			state,
			oauth2.SetAuthURLParam("code", code),
		),
		Code:  code,
		State: state,
	}

	return authURL.String(), nil
}

func GenerateAuthorizationCode(
	clientID string,
	clientSecret string,
	clientRedirectURI string,
	scopes []string,
	userID uint,
	expiry time.Duration,
) (string, error) {
	claims := AuthorizationCodeClaims{
		ClientID:    clientID,
		UserID:      userID,
		RedirectURI: clientRedirectURI,
		Scopes:      strings.Join(scopes, ","), // fixme: should the separator be a space?
		ExpiresAt:   time.Now().Add(expiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	code, err := token.SignedString([]byte(clientSecret))
	if err != nil {
		return "", fmt.Errorf("could not generate auth code: %s", err)
	}

	return code, nil
}

func VerifyAuthorizationCode(clientSecret string, code string) (*AuthorizationCodeClaims, error) {
	authClaims, err := DecodeAuthCode(clientSecret, code)
	if err != nil {
		return authClaims, fmt.Errorf("failed to decode auth code: %v", err)
	}

	return authClaims, nil
}

func DecodeAuthCode(secret, code string) (*AuthorizationCodeClaims, error) {
	token, err := jwt.ParseWithClaims(code, &AuthorizationCodeClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not decode auth code: %s", err)
	}

	claims, ok := token.Claims.(*AuthorizationCodeClaims)
	if !ok {
		return nil, fmt.Errorf("could not decode auth code: %s", err)
	}

	//if claims.ExpiresAt < time.Now().Unix() {
	//	return nil, fmt.Errorf("could not decode auth code: token expired")
	//}
	return claims, nil
}
