package emailRepo

import (
	"globe-and-citizen/layer8/auth-server/config"
	"globe-and-citizen/layer8/auth-server/internal/models"
	"globe-and-citizen/layer8/auth-server/internal/models/gormModels"
	"time"
)

type IEmailRepository interface {
	SendVerificationEmail(user *gormModels.User, userEmail string, verificationCode string) error
	VerifyCode(verificationData *gormModels.EmailVerificationData, code string) error
	GetVerificationCodeExpiry() time.Duration // todo this sounds silly?
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

func (r *EmailRepository) SendVerificationEmail(user *gormModels.User, userEmail string, verificationCode string) error {
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

func (r *EmailRepository) VerifyCode(verificationData *gormModels.EmailVerificationData, code string) error {
	return r.verifier.VerifyCode(verificationData, code)
}

func (r *EmailRepository) GetVerificationCodeExpiry() time.Duration {
	return r.verifier.VerificationCodeExpiry
}
