package userUC

import (
	"context"
	"errors"
	"fmt"
	"globe-and-citizen/layer8/auth-server/internal/consts"
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/internal/models/gormModels"
	"globe-and-citizen/layer8/auth-server/internal/usecases/ucerror"
	"time"

	"gorm.io/gorm"
)

func (uc *UserUsecase) VerifyEmail(ctx context.Context, userID uint, userEmail string) *ucerror.UCError {
	user, err := uc.postgres.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ucerror.New(fmt.Errorf("verify email user not found: %w", err), consts.ErrNotFound)
		}
		return ucerror.New(fmt.Errorf("failed to get user to verify email: %w", err), consts.ErrInternalServer)
	}

	verificationCode, err := uc.code.GenerateEmailVerificationCode(user.ScramSalt, userEmail)
	if err != nil {
		return ucerror.New(fmt.Errorf("failed to generate verification code: %w", err), consts.ErrInternalServer)
	}

	err = uc.email.SendVerificationEmail(user, userEmail, verificationCode)
	if err != nil {
		return ucerror.New(fmt.Errorf("failed to send verification email: %w", err), consts.ErrInternalServer)
	}

	err = uc.postgres.SaveEmailVerificationData(
		ctx,
		gormModels.EmailVerificationData{
			UserId:           user.ID,
			VerificationCode: verificationCode,
			ExpiresAt:        time.Now().Add(uc.email.GetVerificationCodeExpiry()).UTC(),
		},
	)
	return ucerror.New(fmt.Errorf("failed to save verification email: %w", err), consts.ErrInternalServer)
}

func (uc *UserUsecase) CheckEmailVerificationCode(ctx context.Context, userId uint, code string) *ucerror.UCError {
	verificationData, err := uc.postgres.GetEmailVerificationData(ctx, userId)
	if err != nil {
		return ucerror.New(fmt.Errorf("failed to get verification data: %w", err), consts.ErrInternalServer)
	}

	err = uc.email.VerifyCode(verificationData, code)
	return ucerror.New(fmt.Errorf("error verify verification code: %w", err), consts.ErrInternalServer)
}

func (uc *UserUsecase) SaveProofOfEmailVerification(
	ctx context.Context,
	userID uint,
	req requestdto.UserCheckEmailVerificationCode,
) *ucerror.UCError {
	user, err := uc.postgres.GetUserByID(ctx, userID)
	if err != nil { // userID must be valid because it should've been checked in the authentication middleware
		return ucerror.New(fmt.Errorf("failed to get user by id: %w", err), consts.ErrInternalServer)
	}

	zkProof, zkKeyPairId, err := uc.zk.GenerateProof(user.ScramSalt, req.Email, req.Code)
	if err != nil {
		return ucerror.New(fmt.Errorf("failed to generate zkproof of email verification: %w", err), consts.ErrInternalServer)
	}

	err = uc.postgres.SaveProofOfEmailVerification(ctx, userID, req.Code, zkProof, zkKeyPairId)
	if err != nil {
		return ucerror.New(fmt.Errorf("failed to save verification zkproof: %w", err), consts.ErrInternalServer)
	}

	return nil
}
