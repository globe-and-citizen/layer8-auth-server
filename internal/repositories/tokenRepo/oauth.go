package tokenRepo

import (
	"fmt"
	"globe-and-citizen/layer8/auth-server/internal/consts"
	"globe-and-citizen/layer8/auth-server/internal/models"
	"globe-and-citizen/layer8/auth-server/internal/models/gormModels"
	"globe-and-citizen/layer8/auth-server/pkg/oauth"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (t TokenRepository) GenerateOAuthJWTToken(user gormModels.User) (string, error) {
	claims := &jwt.RegisteredClaims{
		Subject:   user.Username, // The value was originally user.ID; it I changed it to Username to avoid type conversion overhead
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(t.oauthJWTSecret)
	if err != nil {
		return "", fmt.Errorf("could not generate oauth token: %s", err)
	}

	return tokenString, nil
}

func (t TokenRepository) VerifyOAuthJWTToken(tokenString string) (models.OAuthAuthenticationClaims, error) {
	claims := &models.OAuthAuthenticationClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return t.oauthJWTSecret, nil
	})
	if err != nil {
		return models.OAuthAuthenticationClaims{}, err
	}

	if !token.Valid {
		return models.OAuthAuthenticationClaims{}, fmt.Errorf("invalid token")
	}

	return *claims, nil
}

func (t TokenRepository) GenerateOAuthAccessToken(client gormModels.Client, authClaims oauth.AuthorizationCodeClaims) (string, error) {
	claims := models.ClientAccessTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Globe and Citizen",
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			Subject:   client.ID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(consts.AccessTokenValidityMinutes * time.Minute).UTC()),
		},
		Scopes: authClaims.Scopes,
		UserID: authClaims.UserID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(client.Secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (t TokenRepository) ParseOAuthAccessToken(tokenString string) (*models.ClientAccessTokenClaims, error) {
	claims := &models.ClientAccessTokenClaims{}
	parser := jwt.NewParser()
	_, _, err := parser.ParseUnverified(tokenString, claims)
	if err != nil {
		return nil, err
	}

	return claims, nil
}

func (t TokenRepository) VerifyOAuthAccessToken(tokenString string, clientSecret []byte) error {
	claims := &models.ClientAccessTokenClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return clientSecret, nil
	})
	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}
