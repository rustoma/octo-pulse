package controllers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/rustoma/octo-pulse/internal/api"
	"net/http"
	"os"
)

type ImageController struct{}

func NewImageController() *ImageController {
	return &ImageController{}
}

func (c *ImageController) HandleGetImage(w http.ResponseWriter, r *http.Request) error {
	imagePath := chi.URLParam(r, "*")
	buf, err := os.ReadFile(os.Getenv("PATH_TO_ASSETS") + "/images/" + imagePath)

	if err != nil {
		errMessage := fmt.Sprintf("Error when serving image: %+v\n", err)
		return api.Error{Err: errMessage, Status: api.HandleErrorStatus(err)}
	}
	w.Header().Set("Content-Type", "image/jpg")
	_, err = w.Write(buf)

	if err != nil {
		errMessage := fmt.Sprintf("Error when serving image: %+v\n", err)
		return api.Error{Err: errMessage, Status: api.HandleErrorStatus(err)}
	}

	return nil

}
