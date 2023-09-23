package controllers

import (
	"net/http"

	"github.com/rustoma/octo-pulse/internal/api"
	"github.com/rustoma/octo-pulse/internal/services"
)

type CategoryController struct {
	categoryService services.CategoryService
}

func NewCategoryController(categoryService services.CategoryService) *CategoryController {
	return &CategoryController{
		categoryService,
	}
}

func (c *CategoryController) HandleGetCategories(w http.ResponseWriter, r *http.Request) error {
	categories, err := c.categoryService.GetCategories()

	if err != nil {
		return api.Error{Err: "cannot get categories", Status: http.StatusInternalServerError}
	}

	return api.WriteJSON(w, http.StatusOK, categories)
}
