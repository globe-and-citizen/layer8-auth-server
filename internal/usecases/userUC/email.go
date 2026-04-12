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

func (uc *UserUsecase) VerifyEmail(ctx context.Context, userID uint, userEmail string) *ucerror.UCError {
	user, err := uc.postgres.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ucerror.New(fmt.Errorf("verify email user not found: %w", err), consts.ErrNotFound)
		}
		return ucerror.New(fmt.Errorf("failed to get user to verify email: %w", err), consts.ErrInternalServer)
	}

	salt := utils.GenerateRandomSalt(consts.SaltSize)

	verificationCode, err := uc.code.GenerateEmailVerificationCode(salt, userEmail)
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
			Salt:             salt,
			Email:            userEmail,
			VerificationCode: verificationCode,
			ExpiresAt:        time.Now().Add(uc.email.GetVerificationCodeExpiry()).UTC(),
		},
	)
	if err != nil {
		return ucerror.New(fmt.Errorf("failed to save verification email: %w", err), consts.ErrInternalServer)
	}

	return nil
}

func (uc *UserUsecase) CheckEmailVerificationCode(ctx context.Context, userId uint, req requestdto.UserCheckEmailVerificationCode) *ucerror.UCError {
	verificationData, err := uc.postgres.GetEmailVerificationData(ctx, userId)
	if err != nil {
		return ucerror.New(fmt.Errorf("failed to get verification data: %w", err), consts.ErrInternalServer)
	}

	err = uc.email.VerifyCode(verificationData, req.Code)
	if err != nil {
		return ucerror.New(fmt.Errorf("verification code is invalid: %w", err), consts.ErrInvalidField)
	}

	return uc.saveProofOfEmailVerification(ctx, userId, verificationData.Salt, verificationData.Email, req.Code)
}

func (uc *UserUsecase) saveProofOfEmailVerification(
	ctx context.Context,
	userID uint,
	salt string,
	email string,
	verificationCode string,
) *ucerror.UCError {
	zkProof, zkKeyPairId, err := uc.zk.GenerateProof(salt, email, verificationCode)
	if err != nil {
		return ucerror.New(fmt.Errorf("failed to generate zkproof of email verification: %w", err), consts.ErrInternalServer)
	}

	err = uc.postgres.SaveProofOfEmailVerification(ctx, userID, salt, verificationCode, zkProof, zkKeyPairId)
	if err != nil {
		return ucerror.New(fmt.Errorf("failed to save verification zkproof: %w", err), consts.ErrInternalServer)
	}

	return nil
}
