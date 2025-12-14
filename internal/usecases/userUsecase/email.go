package userUsecase

import (
	"globe-and-citizen/layer8/auth-server/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/internal/models"
	"time"
)

func (uc *UserUseCase) VerifyEmail(userID uint, userEmail string) error {
	user, e := uc.postgres.FindUserByID(userID)
	if e != nil {
		return e
	}

	verificationCode, err := uc.code.GenerateVerificationCode(user.Salt, userEmail)
	if err != nil {
		return err
	}

	e = uc.email.SendVerificationEmail(&user, userEmail, verificationCode)
	if e != nil {
		return e
	}

	e = uc.postgres.SaveEmailVerificationData(
		models.EmailVerificationData{
			UserId:           user.ID,
			VerificationCode: verificationCode,
			ExpiresAt:        time.Now().Add(uc.email.VerificationCodeValidityDuration()).UTC(),
		},
	)

	return e
}

func (uc *UserUseCase) CheckEmailVerificationCode(userId uint, code string) error {
	verificationData, e := uc.postgres.GetEmailVerificationData(userId)
	if e != nil {
		return e
	}

	e = uc.email.VerifyCode(&verificationData, code)

	return e
}

func (uc *UserUseCase) SaveProofOfEmailVerification(userID uint, req requestdto.CheckEmailVerificationCode) (string, error) {
	user, err := uc.postgres.FindUserByID(userID)
	if err != nil {
		return "Failed to get user", err
	}

	zkProof, zkKeyPairId, err := uc.zk.GenerateProof(user.Salt, req.Email, req.Code)
	if err != nil {
		return "Failed to generate zk proof of email verification", err
	}

	err = uc.postgres.SaveProofOfEmailVerification(userID, req.Code, zkProof, zkKeyPairId)
	if err != nil {
		return "Failed to save proof of the email verification procedure", err
	}

	return "Your email was successfully verified!", nil
}
