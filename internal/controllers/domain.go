package controllers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rustoma/octo-pulse/internal/api"
	"github.com/rustoma/octo-pulse/internal/services"
)

type DomainController struct {
	domainService services.DomainService
}

func NewDomainController(domainService services.DomainService) *DomainController {
	return &DomainController{
		domainService,
	}
}

func (c *DomainController) HandleGetDomains(w http.ResponseWriter, r *http.Request) error {
	domains, err := c.domainService.GetDomains()

	if err != nil {
		return api.Error{Err: "cannot get domains", Status: api.HandleErrorStatus(err)}
	}

	return api.WriteJSON(w, http.StatusOK, domains)
}

func (c *DomainController) HandleGetDomain(w http.ResponseWriter, r *http.Request) error {
	domainIdParam := chi.URLParam(r, "id")
	domainId, err := strconv.Atoi(domainIdParam)

	if err != nil {
		return api.Error{Err: "bad request", Status: http.StatusBadRequest}
	}

	domains, err := c.domainService.GetDomain(domainId)

	if err != nil {
		return api.Error{Err: "cannot get domain", Status: api.HandleErrorStatus(err)}
	}

	return api.WriteJSON(w, http.StatusOK, domains)
}

func (c *DomainController) HandleGetDomainPublicData(w http.ResponseWriter, r *http.Request) error {
	domainIdParam := chi.URLParam(r, "id")
	domainId, err := strconv.Atoi(domainIdParam)

	if err != nil {
		return api.Error{Err: "bad request", Status: http.StatusBadRequest}
	}

	domains, err := c.domainService.GetDomainPublicData(domainId)

	if err != nil {
		return api.Error{Err: "cannot get domain data", Status: api.HandleErrorStatus(err)}
	}

	return api.WriteJSON(w, http.StatusOK, domains)
}
