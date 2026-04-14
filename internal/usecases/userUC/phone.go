package userUC

import (
	"context"
	"errors"
	"fmt"
	"globe-and-citizen/layer8/auth-server/internal/consts"
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/internal/models/gormModels"
	"globe-and-citizen/layer8/auth-server/internal/usecases/ucerror"
	"globe-and-citizen/layer8/auth-server/pkg/utils"
	"time"

	"gorm.io/gorm"
)

func (uc *UserUsecase) VerifyPhoneNumber(ctx context.Context, userID uint) *ucerror.UCError {
	user, err := uc.postgres.GetUserByID(ctx, userID)
	if err != nil { // todo handle ctx canceled error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ucerror.New(fmt.Errorf("user not found: %w", err), consts.ErrNotFound)
		}
		return ucerror.New(fmt.Errorf("failed to get user by ID: %w", err), consts.ErrInternalServer)
	}

	phoneNumber, chatID, err := uc.phone.GetPhoneNumberViaTelegramBot(user.TelegramSessionIDHash)
	if err != nil {
		return ucerror.New(fmt.Errorf("failed to get phone number: %w", err), consts.ErrInternalServer)
	}

	salt := utils.GenerateRandomSalt(consts.SaltSize)
	verificationCode, err := uc.code.GeneratePhoneVerificationCode(salt, phoneNumber)
	if err != nil {
		return ucerror.New(fmt.Errorf("failed to generate verification code: %w", err), consts.ErrInternalServer)
	}

	err = uc.phone.SendVerificationCode(chatID, verificationCode)
	if err != nil {
		return ucerror.New(fmt.Errorf("failed to send verification code: %w", err), consts.ErrInternalServer)
	}

	verificationData := gormModels.PhoneVerificationData{
		UserId:           userID,
		Salt:             salt,
		PhoneNumber:      phoneNumber,
		VerificationCode: verificationCode,
		ExpiresAt:        time.Now().UTC().Add(uc.phone.GetVerificationCodeExpiry()),
	}

	err = uc.postgres.SavePhoneVerificationData(ctx, verificationData)
	if err != nil {
		return ucerror.New(fmt.Errorf("failed to save zkproof of phone number verification: %w", err), consts.ErrInternalServer)
	}

	return nil
}

func (uc *UserUsecase) CheckPhoneVerificationCode(
	ctx context.Context,
	userID uint,
	req requestdto.UserCheckPhoneVerificationCode,
) *ucerror.UCError {
	verificationData, err := uc.postgres.GetPhoneVerificationData(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ucerror.New(fmt.Errorf("phone number verification data not found: %w", err), consts.ErrNotFound)
		}
		return ucerror.New(fmt.Errorf("failed to get phone number verification data: %w", err), consts.ErrInternalServer)
	}

	if req.Code != verificationData.VerificationCode {
		return ucerror.New(fmt.Errorf("phone number verification code does not match"), consts.ErrBadRequest)
	}

	if verificationData.ExpiresAt.Before(time.Now().UTC()) {
		return ucerror.New(fmt.Errorf("phone number verification expired"), fmt.Errorf("%w: verification code is expired", consts.ErrBadRequest))
	}

	zkProof, zkPairID, err := uc.zk.GenerateProof(verificationData.Salt, verificationData.PhoneNumber, req.Code)
	if err != nil {
		return ucerror.New(fmt.Errorf("failed to generate zkproof of phone number verification: %w", err), consts.ErrInternalServer)
	}

	err = uc.postgres.SaveProofOfPhoneVerification(
		ctx,
		verificationData.UserId,
		verificationData.Salt,
		verificationData.VerificationCode,
		zkProof,
		zkPairID,
	)
	if err != nil {
		return ucerror.New(fmt.Errorf("failed to save proof of phone number verification: %w", err), consts.ErrInternalServer)
	}

	return nil
}

func (uc *UserUsecase) GenerateAndSaveTelegramSessionIDHash(ctx context.Context, userID uint) ([]byte, *ucerror.UCError) {
	sessionID, err := utils.GenerateTelegramSessionID()
	if err != nil {
		return nil, ucerror.New(fmt.Errorf("failed to generate telegram session id: %w", err), consts.ErrInternalServer)
	}

	sessionIDHash := utils.ComputeTelegramSessionIDHash(sessionID)

	err = uc.postgres.SaveTelegramSessionIDHash(ctx, userID, sessionIDHash[:])
	if err != nil {
		return nil, ucerror.New(fmt.Errorf("failed to save telegram session id hash: %w", err), consts.ErrInternalServer)
	}

	return sessionID, nil
}
