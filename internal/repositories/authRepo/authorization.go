package authRepo

import (
	"fmt"

	layer8Utils "github.com/globe-and-citizen/layer8-utils"
)

type IAuthorizationRepository interface {
	VerifyAuthorizationCode(clientSecret string, code string) (*layer8Utils.AuthCodeClaims, error)
}

type AuthorizationRepository struct{}

func NewAuthenticationRepository() IAuthorizationRepository {
	return &AuthorizationRepository{}
}

func (r AuthorizationRepository) VerifyAuthorizationCode(clientSecret string, code string) (*layer8Utils.AuthCodeClaims, error) {
	authClaims, err := layer8Utils.DecodeAuthCode(clientSecret, code)
	if err != nil {
		return authClaims, fmt.Errorf("failed to decode auth code: %v", err)
	}

	return authClaims, nil
}
