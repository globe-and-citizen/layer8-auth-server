package userUC

import (
	"context"
	"errors"
	"fmt"
	"globe-and-citizen/layer8/auth-server/internal/consts"
	"globe-and-citizen/layer8/auth-server/internal/usecases/ucerror"

	"gorm.io/gorm"
)

func (uc *UserUsecase) MdwVerifyUserJWTToken(ctx context.Context, tokenString string) (uint, string, *ucerror.UCError) {
	claims, err := uc.token.VerifyUserJWTToken(tokenString)
	if err != nil {
		return 0, "", ucerror.New(fmt.Errorf("failed to verify jwt token: %w", err), consts.ErrUnauthorized)
	}

	// verify user by userID
	user, err := uc.postgres.GetUserByUsername(ctx, claims.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, "", ucerror.New(fmt.Errorf("user not found: %w", err), consts.ErrUnauthorized)
		}
		return 0, "", ucerror.New(fmt.Errorf("failed to get user: %w", err), consts.ErrInternalServer)
	}

	// verify the rest claims
	if user.Username != claims.Username {
		return 0, "", ucerror.New(fmt.Errorf("username mismatch"), consts.ErrUnauthorized)
	}

	return claims.UserID, claims.Username, nil
}
