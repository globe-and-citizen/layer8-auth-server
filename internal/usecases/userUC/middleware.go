package userUC

import (
	"context"
	"fmt"
)

func (uc *UserUsecase) VerifyUserJWTToken(ctx context.Context, tokenString string) (uint, string, error) {
	claims, err := uc.token.VerifyUserJWTToken(tokenString)
	if err != nil {
		return 0, "", err
	}

	// verify user by userID
	user, err := uc.postgres.GetUserByUsername(ctx, claims.Username)
	if err != nil {
		return 0, "", fmt.Errorf("user not found: %e", err)
	}

	// verify the rest claims
	if user.Username != claims.Username {
		return 0, "", fmt.Errorf("invalid user")
	}

	return claims.UserID, claims.Username, nil
}
