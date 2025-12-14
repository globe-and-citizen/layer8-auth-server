package emailRepo

import (
	"globe-and-citizen/layer8/auth-server/config"
	"globe-and-citizen/layer8/auth-server/internal/models"
	"time"
)

type IEmailRepository interface {
	Send(email *models.Email) error
	SendVerificationEmail(user *models.User, userEmail string, verificationCode string) error
	VerifyCode(verificationData *models.EmailVerificationData, code string) error
	VerificationCodeValidityDuration() time.Duration
}

type EmailRepository struct {
	sender   *EmailSender
	verifier *EmailVerifier
}

func NewEmailRepository(config config.EmailConfig) *EmailRepository {

	sender := NewEmailSender(config.ApiKey, config.TemplateId)
	verifier := NewEmailVerifier(config)

	return &EmailRepository{
		sender:   sender,
		verifier: verifier,
	}
}

func (r *EmailRepository) Send(email *models.Email) error {
	return r.sender.Send(email)
}

func (r *EmailRepository) SendVerificationEmail(user *models.User, userEmail string, verificationCode string) error {
	return r.sender.Send(
		&models.Email{
			From:    r.verifier.adminEmailAddress,
			To:      userEmail,
			Subject: "Verify your email at the Layer8 service",
			Content: models.VerificationEmailContent{
				Username: user.Username,
				Code:     verificationCode,
			},
		},
	)
}

func (r *EmailRepository) VerifyCode(verificationData *models.EmailVerificationData, code string) error {
	return r.verifier.VerifyCode(verificationData, code)
}

func (r *EmailRepository) VerificationCodeValidityDuration() time.Duration {
	return r.verifier.VerificationCodeValidityDuration
}
