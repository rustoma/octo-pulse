package controllers

import (
	"github.com/rustoma/octo-pulse/internal/storage"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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
	slug := r.URL.Query().Get("slug")

	var filters storage.GetCategoriesFilters

	if slug != "" {
		filters.Slug = slug
	}

	categories, err := c.categoryService.GetCategories(&filters)

	if err != nil {
		return api.Error{Err: "cannot get categories", Status: api.HandleErrorStatus(err)}
	}

	return api.WriteJSON(w, http.StatusOK, categories)
}

func (c *CategoryController) HandleGetCategory(w http.ResponseWriter, r *http.Request) error {
	categoryIdParam := chi.URLParam(r, "id")
	categoryId, err := strconv.Atoi(categoryIdParam)

	if err != nil {
		return api.Error{Err: "bad request", Status: http.StatusBadRequest}
	}

	category, err := c.categoryService.GetCategory(categoryId)

	if err != nil {
		return api.Error{Err: "cannot get category", Status: api.HandleErrorStatus(err)}
	}

	return api.WriteJSON(w, http.StatusOK, category)
}
