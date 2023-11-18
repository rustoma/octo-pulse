package controllers

import (
	"github.com/rustoma/octo-pulse/internal/api"
	"github.com/rustoma/octo-pulse/internal/dto"
	"github.com/rustoma/octo-pulse/internal/services"
	"net/http"
)

type FileController struct {
	fileService services.FileService
}

func NewFileController(fileService services.FileService) *FileController {
	return &FileController{
		fileService,
	}
}

func (c *FileController) HandleCreateArticles(w http.ResponseWriter, r *http.Request) error {
	var request *dto.CreateArticlesRequest

	err := api.ReadJSON(w, r, &request)
	if err != nil {
		return api.Error{Err: "bad login request", Status: http.StatusBadRequest}
	}

	err = c.fileService.CreateArticles(request.Ids)
	if err != nil {
		return api.Error{Err: "cannot create article files", Status: api.HandleErrorStatus(err)}
	}

	return api.WriteJSON(w, http.StatusOK, "Files are created successfully")
}
