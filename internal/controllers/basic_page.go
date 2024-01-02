package controllers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/rustoma/octo-pulse/internal/api"
	"github.com/rustoma/octo-pulse/internal/models"
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

func (c *BasicPageController) HandleGetBasicPage(w http.ResponseWriter, r *http.Request) error {
	idParam := chi.URLParam(r, "id")

	pageId, err := strconv.Atoi(idParam)
	if err != nil {
		return api.Error{Err: "bad request", Status: http.StatusBadRequest}
	}

	page, err := c.basicPageService.GetBasicPage(pageId)
	if err != nil {
		return api.Error{Err: "cannot get the page", Status: api.HandleErrorStatus(err)}
	}

	return api.WriteJSON(w, http.StatusOK, page)
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

	pages, err := c.basicPageService.GetBasicPageBySlug(slug, &filters)
	if err != nil {
		return api.Error{Err: "cannot get the page", Status: api.HandleErrorStatus(err)}
	}

	return api.WriteJSON(w, http.StatusOK, pages)
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

func (c *BasicPageController) HandleCreateBasicPage(w http.ResponseWriter, r *http.Request) error {
	var request *models.BasicPage

	err := api.ReadJSON(w, r, &request)
	if err != nil {
		logger.Err(err).Msg("Bad request")
		return api.Error{Err: "bad request", Status: http.StatusBadRequest}
	}

	pageId, err := c.basicPageService.CreateBasicPage(request)
	if err != nil {
		logger.Err(err).Send()
		return api.Error{Err: "cannot create basic page", Status: api.HandleErrorStatus(err)}
	}

	return api.WriteJSON(w, http.StatusOK, fmt.Sprintf("Basic page with ID %d was created successfully", pageId))
}
