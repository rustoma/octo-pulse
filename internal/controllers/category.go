package controllers

import (
	"fmt"
	"github.com/rustoma/octo-pulse/internal/models"
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

func (c *CategoryController) HandleGetDomainCategories(w http.ResponseWriter, r *http.Request) error {
	domainIdParam := chi.URLParam(r, "id")
	domainId, err := strconv.Atoi(domainIdParam)

	if err != nil {
		return api.Error{Err: "bad request", Status: http.StatusBadRequest}
	}

	domainCategories, err := c.categoryService.GetDomainCategories(domainId)

	if err != nil {
		return api.Error{Err: "cannot get domain categories", Status: api.HandleErrorStatus(err)}
	}

	return api.WriteJSON(w, http.StatusOK, domainCategories)
}

func (c *CategoryController) HandleCreateCategory(w http.ResponseWriter, r *http.Request) error {
	var request *models.Category

	err := api.ReadJSON(w, r, &request)
	if err != nil {
		logger.Err(err).Msg("Bad request")
		return api.Error{Err: "bad request", Status: http.StatusBadRequest}
	}

	categoryId, err := c.categoryService.CreateCategory(request)
	if err != nil {
		logger.Err(err).Send()
		return api.Error{Err: "cannot create category", Status: api.HandleErrorStatus(err)}
	}

	return api.WriteJSON(w, http.StatusOK, fmt.Sprintf("Category with ID %d was created successfully", categoryId))
}
