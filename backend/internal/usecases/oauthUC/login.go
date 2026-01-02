package oauthUC

import (
	"fmt"
	"globe-and-citizen/layer8/auth-server/backend/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/backend/internal/dto/responsedto"
	"globe-and-citizen/layer8/auth-server/backend/pkg/scram"
)

func (uc *OAuthUsecase) PrecheckUserLogin(req requestdto.OAuthUserLoginPrecheck) (responsedto.OAuthUserLoginPrecheck, error) {
	user, err := uc.postgres.GetUserByUsername(req.Username)
	if err != nil {
		return responsedto.OAuthUserLoginPrecheck{}, err
	}

	loginPrecheckResp := responsedto.OAuthUserLoginPrecheck{
		UserLoginPrecheck: responsedto.UserLoginPrecheck{
			ServerLoginFirstMessage: scram.CreateServerLoginFirstMessage(user.ScramSalt, user.ScramIterationCount, req.ClientLoginFirstMessage),
		},
	}

	return loginPrecheckResp, nil
}

func (uc *OAuthUsecase) UserLogin(req requestdto.OAuthUserLogin) (responsedto.OAuthUserLogin, error) {
	user, err := uc.postgres.GetUserByUsername(req.Username)
	if err != nil {
		return responsedto.OAuthUserLogin{}, err
	}

	tokenString, err := uc.token.GenerateOAuthJWTToken(user)
	if err != nil {
		return responsedto.OAuthUserLogin{}, fmt.Errorf("error generating token: %v", err)
	}

	serverFinalMsg, err := scram.CreateServerLoginFinalMessage(req.ClientLoginFinalMessage, req.CNonce, user.ScramSalt,
		user.ScramIterationCount, user.ScramStoredKey, user.ScramServerKey)
	if err != nil {
		return responsedto.OAuthUserLogin{}, fmt.Errorf("error creating server final message: %v", err)
	}

	return responsedto.OAuthUserLogin{
		UserLogin: responsedto.UserLogin{
			ServerLoginFinalMessage: serverFinalMsg,
			Token:                   tokenString,
		},
	}, nil
}
