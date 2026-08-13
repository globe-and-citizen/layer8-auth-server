package userUC

import (
	"context"
	"errors"
	"fmt"
	"globe-and-citizen/layer8/auth-server/internal/consts"
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/internal/dto/responsedto"
	"globe-and-citizen/layer8/auth-server/internal/models/gormModels"
	"globe-and-citizen/layer8/auth-server/internal/usecases/ucerror"
	"globe-and-citizen/layer8/auth-server/pkg/scram"
	"globe-and-citizen/layer8/auth-server/pkg/utils"

	"gorm.io/gorm"
)

func (uc *UserUsecase) PrecheckRegister(
	ctx context.Context, req requestdto.UserRegisterPrecheck, iterCount int,
) (*responsedto.UserRegisterPrecheck, *ucerror.UCError) {
	registerMsg := scram.CreateServerRegisterFirstMessage(iterCount)

	user := gormModels.User{
		Username:            req.Username,
		ScramSalt:           registerMsg.Salt,
		ScramIterationCount: iterCount,
		PublicKey:           []byte{},
	}

	err := uc.postgres.CreateUser(ctx, user)
	if err != nil {
		return nil, ucerror.New(utils.StackError(fmt.Errorf("failed to create user: %w", err)), consts.ErrInternalServer)
	}

	return &responsedto.UserRegisterPrecheck{
		ServerRegisterFirstMessage: registerMsg,
	}, nil
}

func (uc *UserUsecase) Register(ctx context.Context, req requestdto.UserRegister) *ucerror.UCError {
	newUser := gormModels.User{
		Username:       req.Username,
		PublicKey:      req.PublicKey,
		ScramStoredKey: req.StoredKey,
		ScramServerKey: req.ServerKey,
	}

	err := uc.postgres.UpdateUser(ctx, newUser)
	if err != nil {
		delErr := uc.postgres.DeleteUserByUsername(ctx, req.Username)
		if delErr != nil {
			uc.logger.Error("Failed to delete user after failed update", delErr)
		}

		if errors.Is(err, context.Canceled) {
			return ucerror.New(utils.StackError(fmt.Errorf("request was canceled: %w", err)), consts.ErrRequestCanceled)
		}

		if errors.Is(err, context.DeadlineExceeded) {
			return ucerror.New(utils.StackError(fmt.Errorf("request timeout: %w", err)), consts.ErrRequestTimeout)
		}

		return ucerror.New(utils.StackError(fmt.Errorf("failed to update user: %w", err)), consts.ErrInternalServer)
	}

	return nil
}

func (uc *UserUsecase) PrecheckLogin(
	ctx context.Context, req requestdto.UserLoginPrecheck,
) (*responsedto.UserLoginPrecheck, *ucerror.UCError) {
	user, err := uc.postgres.GetUserByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ucerror.New(utils.StackError(fmt.Errorf("failed to get user: %w", err)), consts.ErrNotFound)
		}
		return nil, ucerror.New(utils.StackError(fmt.Errorf("failed to get user: %w", err)), consts.ErrInternalServer)
	}

	loginPrecheckResp := responsedto.UserLoginPrecheck{
		ServerLoginFirstMessage: scram.CreateServerLoginFirstMessage(user.ScramSalt, user.ScramIterationCount, req.ClientLoginFirstMessage),
	}

	return &loginPrecheckResp, nil
}

func (uc *UserUsecase) Login(ctx context.Context, req requestdto.UserLogin) (*responsedto.UserLogin, *ucerror.UCError) {
	user, err := uc.postgres.GetUserByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ucerror.New(utils.StackError(fmt.Errorf("failed to get user: %w", err)), consts.ErrNotFound)
		}
		return nil, ucerror.New(utils.StackError(fmt.Errorf("failed to get user: %w", err)), consts.ErrInternalServer)
	}

	tokenString, err := uc.token.GenerateUserJWTToken(*user, consts.UserLoginTokenExpiry)
	if err != nil {
		return nil, ucerror.New(utils.StackError(fmt.Errorf("error generating token: %v", err)), consts.ErrInternalServer)
	}

	serverFinalMsg, err := scram.CreateServerLoginFinalMessage(req.ClientLoginFinalMessage, req.CNonce, user.ScramSalt,
		user.ScramIterationCount, user.ScramStoredKey, user.ScramServerKey)
	if err != nil {
		return nil, ucerror.New(utils.StackError(fmt.Errorf("error creating server final message: %v", err)), consts.ErrInternalServer)
	}

	return &responsedto.UserLogin{
		ServerLoginFinalMessage: serverFinalMsg,
		Token:                   tokenString,
	}, nil
}

func (uc *UserUsecase) GetProfile(ctx context.Context, userID uint) (*responsedto.UserProfile, *ucerror.UCError) {
	user, metadata, err := uc.postgres.GetUserProfile(ctx, userID)
	if err != nil {
		return nil, ucerror.New(utils.StackError(fmt.Errorf("failed to get user profile: %w", err)), consts.ErrInternalServer)
	}

	return &responsedto.UserProfile{
		Username:              user.Username,
		DisplayName:           metadata.DisplayName,
		DisplayNameUpdatedAt:  metadata.DisplayNameUpdatedAt,
		Bio:                   metadata.Bio,
		BioUpdatedAt:          metadata.BioUpdatedAt,
		Color:                 metadata.Color,
		ColorUpdatedAt:        metadata.ColorUpdatedAt,
		EmailVerified:         metadata.IsEmailVerified,
		EmailVerifiedAt:       metadata.EmailVerifiedAt,
		PhoneNumberVerified:   metadata.IsPhoneNumberVerified,
		PhoneNumberVerifiedAt: metadata.PhoneVerifiedAt,
		PhoneLocation:         metadata.Location,
	}, nil
}

func (uc *UserUsecase) PrecheckResetPassword(
	ctx context.Context, req requestdto.UserResetPasswordPrecheck,
) (*responsedto.UserResetPasswordPrecheck, *ucerror.UCError) {
	user, err := uc.postgres.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, ucerror.New(utils.StackError(fmt.Errorf("failed to get user: %w", err)), consts.ErrNotFound)
	}

	return &responsedto.UserResetPasswordPrecheck{
		ServerRegisterFirstMessage: scram.ServerRegisterFirstMessage{
			Salt:           user.ScramSalt,
			IterationCount: user.ScramIterationCount,
		}}, nil
}

func (uc *UserUsecase) ResetPassword(ctx context.Context, req requestdto.UserResetPassword) *ucerror.UCError {
	user, err := uc.postgres.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return ucerror.New(utils.StackError(fmt.Errorf("failed to get user: %w", err)), consts.ErrNotFound)
	}

	err = utils.ValidateSignature("Sign-in with Layer8", req.Signature, user.PublicKey)
	if err != nil {
		return ucerror.New(utils.StackError(fmt.Errorf("invalid signature: %w", err)), fmt.Errorf("%w: signature is invalid", consts.ErrBadRequest))
	}

	err = uc.postgres.UpdateUserPassword(ctx, user.Username, req.StoredKey, req.ServerKey)
	if err != nil {
		return ucerror.New(utils.StackError(fmt.Errorf("failed to update user: %w", err)), consts.ErrInternalServer)
	}

	return nil
}
