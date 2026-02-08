package clientUC

import (
	"context"
	"fmt"
)

func (uc *ClientUsecase) VerifyClientJWTToken(ctx context.Context, tokenString string) (string, string, error) {
	claims, err := uc.token.VerifyClientJWTToken(tokenString)
	if err != nil {
		return "", "", err
	}

	// verify client by clientID
	client, err := uc.postgres.GetClientByID(ctx, claims.ClientID)
	if err != nil {
		return "", "", fmt.Errorf("client not found: %e", err)
	}

	// verify the rest claims
	if client.Username != claims.Username {
		return "", "", fmt.Errorf("invalid client")
	}

	return claims.ClientID, claims.Username, nil
}
