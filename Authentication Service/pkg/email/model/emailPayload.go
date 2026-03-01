package model

type EmailPayload struct {
	Sender      Sender `json:"sender"`
	To          []To   `json:"to"`
	Subject     string `json:"subject"`
	HtmlContent string `json:"htmlContent"`
}
