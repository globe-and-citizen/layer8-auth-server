package tokenUsecase

import (
	"globe-and-citizen/layer8/auth-server/internal/repositories/tokenRepo"
)

type ITokenUseCase interface {
	VerifyUserJWTToken(tokenString string) (uint, error)
}

type TokenUseCase struct {
	r tokenRepo.ITokenRepository
}

func NewTokenUseCase(r tokenRepo.ITokenRepository) ITokenUseCase {
	return &TokenUseCase{r: r}
}

func (uc *TokenUseCase) VerifyUserJWTToken(tokenString string) (uint, error) {
	claims, err := uc.r.VerifyUserJWTToken(tokenString)
	if err != nil {
		return 0, err
	}

	return claims.UserID, nil
}
