package clientUC

import (
	"context"
	"fmt"
	"globe-and-citizen/layer8/auth-server/internal/consts"
	"globe-and-citizen/layer8/auth-server/internal/usecases/ucerror"
)

func (uc *ClientUsecase) MdwVerifyClientJWTToken(ctx context.Context, tokenString string) (string, string, *ucerror.UCError) {
	claims, err := uc.token.VerifyClientJWTToken(tokenString)
	if err != nil {
		return "", "", ucerror.New(fmt.Errorf("failed to verify client token: %w", err), consts.ErrUnauthorized)
	}

	// verify client by clientID
	client, err := uc.postgres.GetClientByID(ctx, claims.ClientID)
	if err != nil {
		return "", "", ucerror.New(fmt.Errorf("failed to get client: %w", err), consts.ErrInternalServer) // todo
	}

	// verify the rest claims
	if client.Username != claims.Username {
		return "", "", ucerror.New(fmt.Errorf("invalid client: username mismatch"), consts.ErrBadRequest)
	}

	return claims.ClientID, claims.Username, nil
}
