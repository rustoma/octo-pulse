package controllers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/rustoma/octo-pulse/internal/api"
	"github.com/rustoma/octo-pulse/internal/services"
	"github.com/rustoma/octo-pulse/internal/storage"
	"net/http"
	"os"
	"strconv"
)

type ImageController struct {
	imageService services.ImageService
}

func NewImageController(imageService services.ImageService) *ImageController {
	return &ImageController{
		imageService: imageService,
	}
}

func (c *ImageController) HandleGetImageByPath(w http.ResponseWriter, r *http.Request) error {
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

func (c *ImageController) HandleGetImages(w http.ResponseWriter, r *http.Request) error {
	categoryIdParam := r.URL.Query().Get("categoryId")
	limitParam := r.URL.Query().Get("limit")
	offsetParam := r.URL.Query().Get("offset")
	pathParam := r.URL.Query().Get("path")

	var filters storage.GetImagesFilters

	if limitParam != "" {
		limit, err := strconv.Atoi(limitParam)
		if err != nil {
			return api.Error{Err: "bad request - limit wrong format", Status: http.StatusBadRequest}
		}

		filters.Limit = limit
	}

	if offsetParam != "" {
		offset, err := strconv.Atoi(offsetParam)
		if err != nil {
			return api.Error{Err: "bad request - offset wrong format", Status: http.StatusBadRequest}
		}

		filters.Offset = offset
	}

	if categoryIdParam != "" {
		categoryId, err := strconv.Atoi(categoryIdParam)
		if err != nil {
			return api.Error{Err: "bad request - categoryId wrong format", Status: http.StatusBadRequest}
		}

		filters.CategoryId = categoryId
	}

	if pathParam != "" {
		filters.Path = pathParam
	}

	articles, err := c.imageService.GetImages(&filters)
	if err != nil {
		return api.Error{Err: "Cannot get images", Status: api.HandleErrorStatus(err)}
	}

	return api.WriteJSON(w, http.StatusOK, articles)
}

func (c *ImageController) HandleGetImage(w http.ResponseWriter, r *http.Request) error {
	imageIdParam := chi.URLParam(r, "id")

	imageId, err := strconv.Atoi(imageIdParam)
	if err != nil {
		return api.Error{Err: "bad request", Status: http.StatusBadRequest}
	}

	image, err := c.imageService.GetImage(imageId)
	if err != nil {
		return api.Error{Err: "Cannot get image", Status: api.HandleErrorStatus(err)}
	}

	return api.WriteJSON(w, http.StatusOK, image)

}

func (c *ImageController) HandleGetImageCategories(w http.ResponseWriter, r *http.Request) error {
	categories, err := c.imageService.GetImageCategories()
	if err != nil {
		return api.Error{Err: "Cannot get image categories", Status: api.HandleErrorStatus(err)}
	}

	return api.WriteJSON(w, http.StatusOK, categories)
}

func (c *ImageController) HandleUploadImage(w http.ResponseWriter, r *http.Request) error {
	categoryIdParam := chi.URLParam(r, "id")
	categoryId, err := strconv.Atoi(categoryIdParam)
	if err != nil {
		return api.Error{Err: "bad request", Status: http.StatusBadRequest}
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		logger.Err(err).Send()
		return api.Error{Err: "Uploading file is too big", Status: api.HandleErrorStatus(err)}
	}

	file, handler, err := r.FormFile("image")
	if err != nil {
		logger.Err(err).Send()
		return api.Error{Err: "Error Retrieving the File", Status: api.HandleErrorStatus(err)}
	}
	defer file.Close()

	_, err = c.imageService.UploadImage(file, handler, categoryId)
	if err != nil {
		return api.Error{Err: err.Error(), Status: api.HandleErrorStatus(err)}
	}

	return api.WriteJSON(w, http.StatusOK, "Image uploaded successfully")
}
