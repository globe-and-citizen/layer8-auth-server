package oauthUC

import "fmt"

func (uc *OAuthUsecase) VerifyOAuthJWTToken(tokenString string) (uint, string, error) {
	claims, err := uc.token.VerifyOAuthJWTToken(tokenString)
	if err != nil {
		return 0, "", err
	}

	// verify user by username
	user, err := uc.postgres.GetUserByUsername(claims.Subject)
	if err != nil {
		return 0, "", fmt.Errorf("user not found: %e", err)
	}

	// todo verify the rest claims

	return user.ID, user.Username, nil
}

func (uc *OAuthUsecase) VerifyAccessToken(tokenString string) (uint, string, error) {
	claims, err := uc.token.ParseOAuthAccessToken(tokenString)
	if err != nil {
		return 0, "", err
	}

	// validate clientID
	client, err := uc.postgres.GetClientByID(claims.Subject)
	if err != nil {
		return 0, "", fmt.Errorf("client not found: %e", err)
	}

	err = uc.token.VerifyOAuthAccessToken(tokenString, []byte(client.Secret))
	if err != nil {
		return 0, "", fmt.Errorf("invalid token: %w", err)
	}

	return uint(claims.UserID), claims.Scopes, nil
}
