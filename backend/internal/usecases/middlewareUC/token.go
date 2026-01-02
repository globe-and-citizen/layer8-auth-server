package middlewareUC

import (
	"fmt"
	"globe-and-citizen/layer8/auth-server/backend/internal/repositories/postgresRepo"
	"globe-and-citizen/layer8/auth-server/backend/internal/repositories/tokenRepo"
)

type IMiddlewareUsecase interface {
	VerifyUserJWTToken(tokenString string) (userID uint, userUsername string, err error)
	VerifyClientJWTToken(tokenString string) (clientID string, clientUsername string, err error)
	VerifyOAuthJWTToken(tokenString string) (userID uint, userUsername string, err error)
	VerifyAccessToken(tokenString string) (userID uint, scopes string, err error)
}

type MiddlewareUsecase struct {
	token   tokenRepo.ITokenRepository
	postres postgresRepo.IPostgresRepository
}

func NewMiddlewareUsecase(token tokenRepo.ITokenRepository, postgres postgresRepo.IPostgresRepository) IMiddlewareUsecase {
	return &MiddlewareUsecase{
		token:   token,
		postres: postgres,
	}
}

func (uc *MiddlewareUsecase) VerifyUserJWTToken(tokenString string) (uint, string, error) {
	claims, err := uc.token.VerifyUserJWTToken(tokenString)
	if err != nil {
		return 0, "", err
	}

	// verify user by userID
	user, err := uc.postres.GetUserByUsername(claims.Username)
	if err != nil {
		return 0, "", fmt.Errorf("user not found: %e", err)
	}

	// verify the rest claims
	if user.Username != claims.Username {
		return 0, "", fmt.Errorf("invalid user")
	}

	return claims.UserID, claims.Username, nil
}

func (uc *MiddlewareUsecase) VerifyClientJWTToken(tokenString string) (string, string, error) {
	claims, err := uc.token.VerifyClientJWTToken(tokenString)
	if err != nil {
		return "", "", err
	}

	// verify client by clientID
	client, err := uc.postres.GetClientByID(claims.ClientID)
	if err != nil {
		return "", "", fmt.Errorf("client not found: %e", err)
	}

	// verify the rest claims
	if client.Username != claims.Username {
		return "", "", fmt.Errorf("invalid client")
	}

	return claims.ClientID, claims.Username, nil
}

func (uc *MiddlewareUsecase) VerifyOAuthJWTToken(tokenString string) (uint, string, error) {
	claims, err := uc.token.VerifyOAuthJWTToken(tokenString)
	if err != nil {
		return 0, "", err
	}

	// verify user by username
	user, err := uc.postres.GetUserByUsername(claims.Subject)
	if err != nil {
		return 0, "", fmt.Errorf("user not found: %e", err)
	}

	// todo verify the rest claims

	return user.ID, user.Username, nil
}

func (uc *MiddlewareUsecase) VerifyAccessToken(tokenString string) (uint, string, error) {
	claims, err := uc.token.ParseOAuthAccessToken(tokenString)
	if err != nil {
		return 0, "", err
	}

	// validate clientID
	client, err := uc.postres.GetClientByID(claims.Subject)
	if err != nil {
		return 0, "", fmt.Errorf("client not found: %e", err)
	}

	err = uc.token.VerifyOAuthAccessToken(tokenString, []byte(client.Secret))
	if err != nil {
		return 0, "", fmt.Errorf("invalid token: %w", err)
	}

	return uint(claims.UserID), claims.Scopes, nil
}
