package dto

type SendEmailRequest struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
}
