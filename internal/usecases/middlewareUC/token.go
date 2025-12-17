package middlewareUC

import (
	"globe-and-citizen/layer8/auth-server/internal/repositories/tokenRepo"
)

type IMiddlewareUsecase interface {
	VerifyUserJWTToken(tokenString string) (userID uint, err error)
	VerifyClientJWTToken(tokenString string) (clientID string, clientUsername string, err error)
}

type MiddlewareUsecase struct {
	r tokenRepo.ITokenRepository
}

func NewMiddlewareUsecase(r tokenRepo.ITokenRepository) IMiddlewareUsecase {
	return &MiddlewareUsecase{r: r}
}

func (uc *MiddlewareUsecase) VerifyUserJWTToken(tokenString string) (uint, error) {
	claims, err := uc.r.VerifyUserJWTToken(tokenString)
	if err != nil {
		return 0, err
	}

	return claims.UserID, nil
}

func (uc *MiddlewareUsecase) VerifyClientJWTToken(tokenString string) (string, string, error) {
	claims, err := uc.r.VerifyClientJWTToken(tokenString)
	if err != nil {
		return "", "", err
	}

	return claims.ClientID, claims.Username, nil
}
