package controllers

import (
	"github.com/rustoma/octo-pulse/internal/api"
	"github.com/rustoma/octo-pulse/internal/services"
	"net/http"
)

type ScrapperController struct {
	scrapperService services.ScrapperService
}

func NewScrapperController(scrapperService services.ScrapperService) *ScrapperController {
	return &ScrapperController{
		scrapperService,
	}
}

func (c *ScrapperController) HandleGetQuestionCategories(w http.ResponseWriter, r *http.Request) error {
	categories, err := c.scrapperService.GetQuestionCategories()

	if err != nil {
		return api.Error{Err: "cannot get question categories", Status: api.HandleErrorStatus(err)}
	}

	return api.WriteJSON(w, http.StatusOK, categories)
}
