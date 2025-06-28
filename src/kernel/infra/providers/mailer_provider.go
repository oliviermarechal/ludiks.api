package providers

import (
	"context"
	"fmt"
	"time"

	"github.com/mailersend/mailersend-go"
)

func NewMailerProvider(mailerApiKey string) *MailerProvider {
	return &MailerProvider{
		ApiKey: mailerApiKey,
	}
}

type MailerProvider struct {
	ApiKey string
}

func (m *MailerProvider) SendTemplateMail(templateId string, to string, subject string, data map[string]interface{}) error {
	ms := mailersend.NewMailersend(m.ApiKey)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	from := mailersend.From{
		Name:  "Ludiks",
		Email: "info@ludiks.io",
	}

	recipients := []mailersend.Recipient{
		{
			Email: to,
		},
	}

	personalization := []mailersend.Personalization{
		{
			Email: to,
			Data:  data,
		},
	}

	message := ms.Email.NewMessage()

	message.SetFrom(from)
	message.SetRecipients(recipients)
	message.SetSubject(subject)
	message.SetTemplateID(templateId)
	message.SetPersonalization(personalization)

	// Send the email
	_, err := ms.Email.Send(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (m *MailerProvider) SendTemplateMailWithContent(to string, subject string, htmlContent string) error {
	ms := mailersend.NewMailersend(m.ApiKey)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	from := mailersend.From{
		Name:  "Ludiks",
		Email: "info@test-nrw7gymroyng2k8e.mlsender.net",
		// Email: "info@ludiks.io",
	}

	recipients := []mailersend.Recipient{
		{
			Email: to,
		},
	}

	message := ms.Email.NewMessage()

	message.SetFrom(from)
	message.SetRecipients(recipients)
	message.SetSubject(subject)
	message.SetHTML(htmlContent)

	_, err := ms.Email.Send(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
