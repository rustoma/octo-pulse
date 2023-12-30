package controllers

import (
	"github.com/rustoma/octo-pulse/internal/api"
	"github.com/rustoma/octo-pulse/internal/services"
	"net/http"
)

type AuthorController struct {
	authorService services.AuthorService
}

func NewAuthorController(authorService services.AuthorService) *AuthorController {
	return &AuthorController{
		authorService,
	}
}

func (c *AuthorController) HandleGetAuthors(w http.ResponseWriter, r *http.Request) error {
	categories, err := c.authorService.GetAuthors()
	if err != nil {
		return api.Error{Err: "cannot get authors", Status: api.HandleErrorStatus(err)}
	}

	return api.WriteJSON(w, http.StatusOK, categories)
}
