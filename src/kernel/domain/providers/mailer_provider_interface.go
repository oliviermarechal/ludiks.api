package providers

type MailerProvider interface {
	SendTemplateMail(templateId string, to string, subject string, data map[string]interface{}) error
	SendTemplateMailWithContent(to string, subject string, htmlContent string) error
}
