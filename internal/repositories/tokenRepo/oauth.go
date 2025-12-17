package tokenRepo

import (
	"globe-and-citizen/layer8/auth-server/internal/consts"
	"globe-and-citizen/layer8/auth-server/internal/models"
	"globe-and-citizen/layer8/auth-server/internal/models/gormModels"
	"time"

	layer8Utils "github.com/globe-and-citizen/layer8-utils"
	"github.com/golang-jwt/jwt/v5"
)

func (t TokenRepository) GenerateClientOauthJWTToken(client gormModels.Client, authClaims layer8Utils.AuthCodeClaims) (string, error) {
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
