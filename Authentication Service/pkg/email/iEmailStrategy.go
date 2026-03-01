package email

import "Authentication_Service/pkg/email/model"

type IEmailStrategy interface {
	SendEmail(to model.To, subject string, body string) error
}
