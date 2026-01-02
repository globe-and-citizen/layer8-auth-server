package userUC

import (
	"fmt"
	"globe-and-citizen/layer8/auth-server/backend/internal/dto/requestdto"
	"globe-and-citizen/layer8/auth-server/backend/internal/models/gormModels"
	"globe-and-citizen/layer8/auth-server/backend/pkg/utils"
	"net/http"
	"time"
)

func (uc *UserUsecase) VerifyPhoneNumber(userID uint) (string, error) {
	user, err := uc.postgres.GetUserByID(userID)
	if err != nil {
		return "failed to get user", err
	}

	phoneNumber, chatID, err := uc.phone.GetPhoneNumberViaTelegramBot(user.TelegramSessionIDHash)
	if err != nil {
		return "failed to get phone number via telegram bot", err
	}

	verificationCode, err := uc.code.GeneratePhoneVerificationCode(user.ScramSalt, phoneNumber)
	if err != nil {
		return "failed to generate phone number verification code", err
	}

	err = uc.phone.SendVerificationCode(chatID, verificationCode)
	if err != nil {
		return "failed to send verification code", err
	}

	zkProof, zkPairID, err := uc.zk.GenerateProof(user.ScramSalt, phoneNumber, verificationCode)
	if err != nil {
		return "failed to generate the zk proof of phone number verification", err
	}

	verificationData := gormModels.PhoneNumberVerificationData{
		UserId:           userID,
		VerificationCode: verificationCode,
		ExpiresAt:        time.Now().UTC().Add(uc.phone.GetVerificationCodeExiry()),
		ZkProof:          zkProof,
		ZkPairID:         zkPairID,
	}

	err = uc.postgres.SavePhoneNumberVerificationData(verificationData)
	if err != nil {
		return "failed to save proof of the phone number verification into the db", err
	}

	return "Phone number successfully verified", nil
}

func (uc *UserUsecase) CheckPhoneNumberVerificationCode(userID uint, req requestdto.UserCheckPhoneNumberVerificationCode) (int, string, error) {
	verificationData, err := uc.postgres.GetPhoneNumberVerificationData(userID)
	if err != nil {
		// fixme check error types and return corresponding status codes
		return http.StatusBadRequest, "user's verification data not found", err
	}

	if req.VerificationCode != verificationData.VerificationCode {
		return http.StatusBadRequest, "failed to validate the provided verification code", fmt.Errorf("invalid verification code")
	}

	if verificationData.ExpiresAt.Before(time.Now().UTC()) {
		return http.StatusBadRequest, "failed to validate the provided verification code", fmt.Errorf("verification code is expired. Try to verify your phone number again")
	}

	err = uc.postgres.SaveProofOfPhoneNumberVerification(
		verificationData.UserId,
		verificationData.VerificationCode,
		verificationData.ZkProof,
		verificationData.ZkPairID,
	)
	if err != nil {
		return http.StatusInternalServerError, "failed to update phone number verification metadata in the db", err
	}

	return http.StatusOK, "phone number verification code successfully verified", nil
}

func (uc *UserUsecase) GenerateAndSaveTelegramSessionIDHash(userID uint) ([]byte, string, error) {
	sessionID, err := utils.GenerateTelegramSessionID()
	if err != nil {
		return []byte{}, "failed to generate telegram session id", err
	}

	sessionIDHash := utils.ComputeTelegramSessionIDHash(sessionID)

	err = uc.postgres.SaveTelegramSessionIDHash(userID, sessionIDHash[:])
	if err != nil {
		return []byte{}, "failed to save telegram session id hash", err
	}

	return sessionID, "telegram session id hash successfully generated and saved", nil
}
