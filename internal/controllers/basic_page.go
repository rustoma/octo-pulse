package controllers

import (
	"github.com/go-chi/chi/v5"
	"github.com/rustoma/octo-pulse/internal/api"
	"github.com/rustoma/octo-pulse/internal/services"
	"github.com/rustoma/octo-pulse/internal/storage"
	"net/http"
	"strconv"
)

type BasicPageController struct {
	basicPageService services.BasicPageService
}

func NewBasicPageController(basicPageService services.BasicPageService) *BasicPageController {
	return &BasicPageController{
		basicPageService,
	}
}

func (c *BasicPageController) HandleGetBasicPageBySlug(w http.ResponseWriter, r *http.Request) error {
	slug := chi.URLParam(r, "slug")
	domainIdParam := r.URL.Query().Get("domainId")

	var filters storage.GetBasicPageBySlugFilters

	if domainIdParam != "" {
		domainId, err := strconv.Atoi(domainIdParam)
		if err != nil {
			return api.Error{Err: "bad request - domainId wrong format", Status: http.StatusBadRequest}
		}

		filters.DomainId = domainId
	}

	domains, err := c.basicPageService.GetBasicPageBySlug(slug, &filters)

	if err != nil {
		return api.Error{Err: "cannot get the page", Status: api.HandleErrorStatus(err)}
	}

	return api.WriteJSON(w, http.StatusOK, domains)
}

func (c *BasicPageController) HandleGetBasicPages(w http.ResponseWriter, r *http.Request) error {
	domainIdParam := r.URL.Query().Get("domainId")

	var filters storage.GetBasicPagesFilters

	if domainIdParam != "" {
		domainId, err := strconv.Atoi(domainIdParam)
		if err != nil {
			return api.Error{Err: "bad request - domainId wrong format", Status: http.StatusBadRequest}
		}

		filters.DomainId = domainId
	}

	domains, err := c.basicPageService.GetBasicPages(&filters)

	if err != nil {
		return api.Error{Err: "cannot get the basic pages", Status: api.HandleErrorStatus(err)}
	}

	return api.WriteJSON(w, http.StatusOK, domains)
}
