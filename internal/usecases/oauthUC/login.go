package oauthUC

import (
	"fmt"
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	responsedto2 "globe-and-citizen/layer8/auth-server/internal/dto/responsedto"
	"globe-and-citizen/layer8/auth-server/pkg/scram"
)

func (uc *OAuthUsecase) PrecheckUserLogin(req requestdto.OAuthUserLoginPrecheck) (responsedto2.OAuthUserLoginPrecheck, error) {
	user, err := uc.postgres.GetUserByUsername(req.Username)
	if err != nil {
		return responsedto2.OAuthUserLoginPrecheck{}, err
	}

	loginPrecheckResp := responsedto2.OAuthUserLoginPrecheck{
		UserLoginPrecheck: responsedto2.UserLoginPrecheck{
			ServerLoginFirstMessage: scram.CreateServerLoginFirstMessage(user.ScramSalt, user.ScramIterationCount, req.ClientLoginFirstMessage),
		},
	}

	return loginPrecheckResp, nil
}

func (uc *OAuthUsecase) UserLogin(req requestdto.OAuthUserLogin) (responsedto2.OAuthUserLogin, error) {
	user, err := uc.postgres.GetUserByUsername(req.Username)
	if err != nil {
		return responsedto2.OAuthUserLogin{}, err
	}

	tokenString, err := uc.token.GenerateOAuthJWTToken(user)
	if err != nil {
		return responsedto2.OAuthUserLogin{}, fmt.Errorf("error generating token: %v", err)
	}

	serverFinalMsg, err := scram.CreateServerLoginFinalMessage(req.ClientLoginFinalMessage, req.CNonce, user.ScramSalt,
		user.ScramIterationCount, user.ScramStoredKey, user.ScramServerKey)
	if err != nil {
		return responsedto2.OAuthUserLogin{}, fmt.Errorf("error creating server final message: %v", err)
	}

	return responsedto2.OAuthUserLogin{
		UserLogin: responsedto2.UserLogin{
			ServerLoginFinalMessage: serverFinalMsg,
			Token:                   tokenString,
		},
	}, nil
}
