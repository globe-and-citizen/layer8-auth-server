package userUC

import (
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/internal/models/gormModels"
	"time"
)

func (uc *UserUsecase) VerifyEmail(userID uint, userEmail string) error {
	user, err := uc.postgres.GetUserByID(userID)
	if err != nil {
		return err
	}

	verificationCode, err := uc.code.GenerateEmailVerificationCode(user.ScramSalt, userEmail)
	if err != nil {
		return err
	}

	err = uc.email.SendVerificationEmail(&user, userEmail, verificationCode)
	if err != nil {
		return err
	}

	err = uc.postgres.SaveEmailVerificationData(
		gormModels.EmailVerificationData{
			UserId:           user.ID,
			VerificationCode: verificationCode,
			ExpiresAt:        time.Now().Add(uc.email.GetVerificationCodeExpiry()).UTC(),
		},
	)

	return err
}

func (uc *UserUsecase) CheckEmailVerificationCode(userId uint, code string) error {
	verificationData, e := uc.postgres.GetEmailVerificationData(userId)
	if e != nil {
		return e
	}

	e = uc.email.VerifyCode(&verificationData, code)

	return e
}

func (uc *UserUsecase) SaveProofOfEmailVerification(userID uint, req requestdto.UserCheckEmailVerificationCode) (string, error) {
	user, err := uc.postgres.GetUserByID(userID)
	if err != nil {
		return "Failed to get user", err
	}

	zkProof, zkKeyPairId, err := uc.zk.GenerateProof(user.ScramSalt, req.Email, req.Code)
	if err != nil {
		return "Failed to generate zk proof of email verification", err
	}

	err = uc.postgres.SaveProofOfEmailVerification(userID, req.Code, zkProof, zkKeyPairId)
	if err != nil {
		return "Failed to save proof of the email verification procedure", err
	}

	return "Your email was successfully verified!", nil
}
