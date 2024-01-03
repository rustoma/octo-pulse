package controllers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/rustoma/octo-pulse/internal/api"
	"github.com/rustoma/octo-pulse/internal/models"
	"github.com/rustoma/octo-pulse/internal/services"
	"net/http"
	"strconv"
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

func (c *AuthorController) HandleCreateAuthor(w http.ResponseWriter, r *http.Request) error {
	var request *models.Author

	err := api.ReadJSON(w, r, &request)
	if err != nil {
		logger.Err(err).Msg("Bad request")
		return api.Error{Err: "bad request", Status: http.StatusBadRequest}
	}

	authorId, err := c.authorService.CreateAuthor(request)
	if err != nil {
		logger.Err(err).Send()
		return api.Error{Err: "cannot create domain", Status: api.HandleErrorStatus(err)}
	}

	return api.WriteJSON(w, http.StatusOK, fmt.Sprintf("Author with ID %d was created successfully", authorId))
}

func (c *AuthorController) HandleUpdateAuthor(w http.ResponseWriter, r *http.Request) error {
	var author *models.Author

	authorIdParam := chi.URLParam(r, "id")
	authorId, err := strconv.Atoi(authorIdParam)
	if err != nil {
		return api.Error{Err: "bad request", Status: http.StatusBadRequest}
	}

	err = api.ReadJSON(w, r, &author)
	if err != nil {
		return api.Error{Err: err.Error(), Status: http.StatusBadRequest}
	}

	updatedAuthor, err := c.authorService.UpdateAuthor(authorId, author)
	if err != nil {
		return api.Error{Err: err.Error(), Status: api.HandleErrorStatus(err)}
	}

	return api.WriteJSON(w, http.StatusOK, updatedAuthor)
}
