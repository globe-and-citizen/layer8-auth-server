package phoneRepo

import (
	"fmt"
	"globe-and-citizen/layer8/auth-server/internal/config"
	"globe-and-citizen/layer8/auth-server/pkg/telegram"
	"globe-and-citizen/layer8/auth-server/pkg/utils"
	"net/url"
	"time"
)

type IPhoneRepository interface {
	GetPhoneNumberViaTelegramBot(telegramSessionIDHash []byte) (phoneNumber string, chatID int64, err error)
	SendVerificationCode(chatID int64, verificationCode string) error
	GetVerificationCodeExpiry() time.Duration
}

type PhoneRepository struct {
	telegramBot            *telegram.TelegramBot
	verificationCodeExpiry time.Duration
}

func NewPhoneRepository(config config.PhoneConfig) IPhoneRepository {
	baseURL := fmt.Sprintf("https://api.telegram.org/bot%s", url.PathEscape(config.TelegramApiKey))

	return &PhoneRepository{
		telegramBot:            telegram.NewTelegramBot(baseURL),
		verificationCodeExpiry: config.VerificationCodeExpiry,
	}
}

func (r *PhoneRepository) GetPhoneNumberViaTelegramBot(telegramSessionIDHash []byte) (phoneNumber string, chatID int64, err error) {
	var telegramUserID int64
	telegramUserID, err = r.telegramBot.Start(telegramSessionIDHash)
	if err != nil {
		err = utils.StackError(err)
		return
	}

	phoneNumber, chatID, err = r.telegramBot.WaitForContactShare(telegramUserID)
	err = utils.StackError(err)
	return
}

func (r *PhoneRepository) SendVerificationCode(chatID int64, verificationCode string) error {
	return utils.StackError(r.telegramBot.SendVerificationCode(chatID, verificationCode))
}

func (r *PhoneRepository) GetVerificationCodeExpiry() time.Duration {
	return r.verificationCodeExpiry
}
