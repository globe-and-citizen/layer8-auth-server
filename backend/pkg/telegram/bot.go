package telegram

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

type TelegramBot struct {
	baseURL string
}

func NewTelegramBot(baseURL string) *TelegramBot {
	return &TelegramBot{
		baseURL: baseURL,
	}
}

func (t TelegramBot) Start(telegramSessionIDHash []byte) (telegramUserID int64, err error) {
	var offset int64 = 0

	for i := 0; i < 500; i++ {
		if i > 0 {
			time.Sleep(500 * time.Millisecond)
		}

		var updates []MessageUpdateDTO
		updates, err = t.RefreshMessages(offset)
		if err != nil {
			log.Printf("failed to refresh messages from Telegram: %v\n", err)
			time.Sleep(2 * time.Second)
			continue
		}

		foundStartMessage := false

		for _, u := range updates {
			offset = u.UpdateID + 1

			if u.Message == nil || u.Message.Chat.Type != "private" {
				continue
			}

			msg := u.Message

			if !strings.HasPrefix(msg.Text, "/start ") {
				continue
			}

			sessionID := msg.Text[7:]
			if sessionID == "" {
				continue
			}

			var sessionIDBytes []byte
			sessionIDBytes, err = base64.RawURLEncoding.DecodeString(sessionID)
			if err != nil {
				log.Printf("failed to decode session id %s, skipping\n", sessionID)
				continue
			}

			sessionIDHash := sha256.Sum256(sessionIDBytes)

			if bytes.Equal(telegramSessionIDHash, sessionIDHash[:]) {
				log.Printf("Found matching session id\n")

				// Send a reply keyboard that *requests contact*.
				kb := ReplyKeyboardMarkup{
					Keyboard: [][]KeyboardButton{
						{
							{Text: "Share my phone", RequestContact: true},
						},
					},
					ResizeKeyboard:  true,
					OneTimeKeyboard: true,
				}

				err = t.SendMessage(SendMessageRequestDTO{
					ChatID:      msg.Chat.ID,
					Text:        "Hi! In order for us to verify your phone number, please tap the button below to allow Telegram sharing your phone number with us.",
					ParseMode:   "Markdown",
					ReplyMarkup: kb,
				})
				if err != nil {
					log.Printf("sendMessage error: %v", err)
				}

				foundStartMessage = true
				telegramUserID = msg.From.ID
				break
			}
		}

		if foundStartMessage {
			break
		}
	}

	return telegramUserID, err
}

func (t TelegramBot) WaitForContactShare(telegramUserID int64) (phoneNumber string, chatID int64, err error) {
	var offset int64 = 0

	for i := 0; i < 500; i++ {
		if i > 0 {
			time.Sleep(400 * time.Millisecond)
		}

		var updates []MessageUpdateDTO
		updates, err = t.RefreshMessages(offset)
		if err != nil {
			log.Printf("failed to refresh messages from Telegram: %v\n", err)
			time.Sleep(2 * time.Second)
			continue
		}

		for _, u := range updates {
			offset = u.UpdateID + 1

			if u.Message == nil || u.Message.Chat.Type != "private" {
				continue
			}

			msg := u.Message

			if msg.Contact == nil {
				continue
			}

			// Validate the contact belongs to the sender.
			if msg.From == nil || msg.Contact.UserID != msg.From.ID || msg.From.ID != telegramUserID {
				continue
			}

			phoneNumber = strings.TrimSpace(msg.Contact.PhoneNumber)
			return phoneNumber, msg.Chat.ID, nil
		}

		err = fmt.Errorf("failed to wait for contact share from Telegram")
	}

	return "", 0, err
}

func (t TelegramBot) SendVerificationCode(chatID int64, verificationCode string) error {
	return t.SendMessage(SendMessageRequestDTO{
		ChatID: chatID,
		Text: fmt.Sprintf(
			"Thanks! Your verification code is: %s. You can go back to the Layer8 user portal now.",
			verificationCode,
		),
		ReplyMarkup: ReplyKeyboardRemove{
			RemoveKeyboard: true,
		},
	})
}

func (t TelegramBot) RefreshMessages(offset int64) ([]MessageUpdateDTO, error) {
	requestUrl := fmt.Sprintf("%s/getUpdates?timeout=2&limit=50", t.baseURL)
	if offset > 0 {
		requestUrl += "&offset=" + strconv.FormatInt(offset, 10)
	}

	req, err := http.NewRequest(http.MethodGet, requestUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create getUpdates request to the Telegram API")
	}

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response GetUpdatesResponseDTO
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode getUpdates response from the Telegram API")
	}

	if !response.Ok {
		return nil, fmt.Errorf("received not-ok response from the getUpdates endpoint of Telegram API, result: %v", response.Result)
	}

	return response.Result, nil
}

func (t TelegramBot) SendMessage(request SendMessageRequestDTO) error {
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal the send message request: %v", err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/sendMessage", t.baseURL),
		bytes.NewReader(requestBytes),
	)
	if err != nil {
		return fmt.Errorf("failed to create the sendMessage request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	var response SendMessageResponseDTO
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}
	if !response.Ok {
		return fmt.Errorf("received not-ok response from the sendMessage endpoint of Telegram API")
	}

	return nil
}
