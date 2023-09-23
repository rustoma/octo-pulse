package controllers

import (
	"net/http"

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
		return api.Error{Err: "cannot get domains", Status: http.StatusInternalServerError}
	}

	return api.WriteJSON(w, http.StatusOK, domains)
}
