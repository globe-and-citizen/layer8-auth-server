package emailRepo

import (
	"context"
	"fmt"
	"globe-and-citizen/layer8/auth-server/backend/internal/consts"
	"globe-and-citizen/layer8/auth-server/backend/internal/models"
	"net/http"

	"github.com/mailersend/mailersend-go"
)

type EmailSender struct {
	apiKey     string
	templateId string
}

func NewEmailSender(apiKey, templateId string) *EmailSender {
	return &EmailSender{
		apiKey:     apiKey,
		templateId: templateId,
	}
}

func (s *EmailSender) Send(email *models.Email) error {
	mailerSendClient := mailersend.NewMailersend(s.apiKey)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, consts.EmailSendTimeout)
	defer cancel()

	from := mailersend.From{
		Name:  "Layer8 team",
		Email: email.From,
	}
	to := mailersend.Recipient{
		Name:  email.Content.Username,
		Email: email.To,
	}

	//personalization := []mailersend.Personalization{
	//	{
	//		Email: email.To,
	//		Data: map[string]interface{}{
	//			"code": email.Content.Code,
	//			"user": email.Content.Username,
	//		},
	//	},
	//}

	message := mailerSendClient.Email.NewMessage()
	message.SetFrom(from)
	message.SetRecipients([]mailersend.Recipient{to})
	message.SetSubject(email.Subject)
	message.SetText(
		fmt.Sprintf(
			"Hi, %s!\nYour verification code is: %s",
			email.Content.Username,
			email.Content.Code,
		),
	)
	//message.SetTemplateID(ms.templateId)
	//message.SetPersonalization(personalization)

	response, e := mailerSendClient.Email.Send(ctx, message)
	if e != nil {
		return fmt.Errorf("error while sending a verification email via MailerSend: %e", e)
	}

	if response.StatusCode != http.StatusAccepted {
		return fmt.Errorf(
			"failed to send a verification email, status code %d",
			response.StatusCode,
		)
	}

	return nil
}
