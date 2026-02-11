package oauthUC

import (
	"context"
	"errors"
	"fmt"
	"globe-and-citizen/layer8/auth-server/internal/consts"
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/internal/dto/responsedto"
	"globe-and-citizen/layer8/auth-server/internal/usecases/ucerror"
	"globe-and-citizen/layer8/auth-server/pkg/scram"

	"gorm.io/gorm"
)

func (uc *OAuthUsecase) PrecheckUserLogin(ctx context.Context, req requestdto.OAuthUserLoginPrecheck) (responsedto.OAuthUserLoginPrecheck, *ucerror.UCError) {
	user, err := uc.postgres.GetUserByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return responsedto.OAuthUserLoginPrecheck{}, ucerror.New(fmt.Errorf("user not found: %w", err), consts.ErrBadRequest)
		}
		return responsedto.OAuthUserLoginPrecheck{}, ucerror.New(fmt.Errorf("failed to get user by username: %w", err), consts.ErrInternalServer)
	}

	loginPrecheckResp := responsedto.OAuthUserLoginPrecheck{
		UserLoginPrecheck: responsedto.UserLoginPrecheck{
			ServerLoginFirstMessage: scram.CreateServerLoginFirstMessage(user.ScramSalt, user.ScramIterationCount, req.ClientLoginFirstMessage),
		},
	}

	return loginPrecheckResp, nil
}

func (uc *OAuthUsecase) UserLogin(ctx context.Context, req requestdto.OAuthUserLogin) (responsedto.OAuthUserLogin, *ucerror.UCError) {
	user, err := uc.postgres.GetUserByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return responsedto.OAuthUserLogin{}, ucerror.New(fmt.Errorf("user not found: %w", err), consts.ErrBadRequest)

		}
		return responsedto.OAuthUserLogin{}, ucerror.New(fmt.Errorf("failed to get user by username: %w", err), consts.ErrInternalServer)
	}

	tokenString, err := uc.token.GenerateOAuthJWTToken(user)
	if err != nil {
		return responsedto.OAuthUserLogin{}, ucerror.New(fmt.Errorf("error generating login token: %w", err), consts.ErrInternalServer)
	}

	serverFinalMsg, err := scram.CreateServerLoginFinalMessage(req.ClientLoginFinalMessage, req.CNonce, user.ScramSalt,
		user.ScramIterationCount, user.ScramStoredKey, user.ScramServerKey)
	if err != nil {
		return responsedto.OAuthUserLogin{}, ucerror.New(fmt.Errorf("error creating server final message: %w", err), consts.ErrInternalServer)
	}

	return responsedto.OAuthUserLogin{
		UserLogin: responsedto.UserLogin{
			ServerLoginFinalMessage: serverFinalMsg,
			Token:                   tokenString,
		},
	}, nil
}
