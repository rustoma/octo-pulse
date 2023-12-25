package controllers

import (
	"github.com/rustoma/octo-pulse/internal/api"
	"github.com/rustoma/octo-pulse/internal/dto"
	"github.com/rustoma/octo-pulse/internal/services"
	"net/http"
)

type EmailController struct {
	emailService services.EmailService
}

func NewEmailController(emailService services.EmailService) *EmailController {
	return &EmailController{
		emailService,
	}
}

func (c *EmailController) HandleSendEmail(w http.ResponseWriter, r *http.Request) error {
	var request *dto.SendEmailRequest

	err := api.ReadJSON(w, r, &request)
	if err != nil {
		return api.Error{Err: "bad send email request", Status: http.StatusBadRequest}
	}

	email := services.Email{
		Subject: request.Subject,
		Body:    request.Body,
	}

	err = c.emailService.Send(email)
	if err != nil {
		return api.Error{Err: "sending email failed", Status: api.HandleErrorStatus(err)}
	}

	return api.WriteJSON(w, http.StatusOK, "Email was sent successfully.")
}
