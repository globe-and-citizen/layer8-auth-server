package emailRepo

import (
	"fmt"
	"globe-and-citizen/layer8/auth-server/config"
	"globe-and-citizen/layer8/auth-server/internal/models"
	"log"
	"time"
)

type EmailVerifier struct {
	adminEmailAddress                string
	now                              func() time.Time
	VerificationCodeValidityDuration time.Duration
}

func NewEmailVerifier(config config.EmailConfig) *EmailVerifier {
	adminEmailAddress := fmt.Sprintf("%s@%s", config.Layer8EmailUsername, config.Layer8EmailDomain)
	verificationCodeValidityDuration, e := time.ParseDuration(config.VerificationCodeValidDuration)
	if e != nil {
		log.Fatalf("error parsing verification code validity duration: %e", e)
	}

	return &EmailVerifier{
		adminEmailAddress:                adminEmailAddress,
		now:                              time.Now,
		VerificationCodeValidityDuration: verificationCodeValidityDuration,
	}
}

func (v *EmailVerifier) VerifyCode(verificationData *models.EmailVerificationData, code string) error {
	if verificationData.ExpiresAt.Before(v.now()) {
		return fmt.Errorf(
			"the verification code is expired. Please try to run the verification process again",
		)
	}

	if code != verificationData.VerificationCode {
		return fmt.Errorf(
			"invalid verification code, expected %s, got %s",
			verificationData.VerificationCode,
			code,
		)
	}

	return nil
}
