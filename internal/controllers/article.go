package controllers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rustoma/octo-pulse/internal/api"
	"github.com/rustoma/octo-pulse/internal/services"
	"github.com/rustoma/octo-pulse/internal/tasks"
)

type ArticleController struct {
	articleService services.ArticleService
	articleTasks   tasks.ArticleTasker
}

func NewArticleController(articleService services.ArticleService, articleTasks tasks.ArticleTasker) *ArticleController {
	return &ArticleController{
		articleService: articleService,
		articleTasks:   articleTasks,
	}
}

func (c *ArticleController) HandleGenerateDescritption(w http.ResponseWriter, r *http.Request) error {
	pageIdParam := chi.URLParam(r, "id")
	pageId, err := strconv.Atoi(pageIdParam)

	if err != nil {
		return api.Error{Err: "bad request", Status: http.StatusBadRequest}
	}

	err = c.articleTasks.NewGenerateDescriptionTask(pageId)

	if err != nil {
		return api.Error{Err: err.Error(), Status: api.HandleErrorStatus(err)}
	}

	return api.WriteJSON(w, http.StatusOK, "")
}

func (c *ArticleController) HandleGetArticles(w http.ResponseWriter, r *http.Request) error {

	articles, err := c.articleService.GetArticles()

	if err != nil {
		return api.Error{Err: "Cannot get articles", Status: api.HandleErrorStatus(err)}
	}

	return api.WriteJSON(w, http.StatusOK, articles)
}
