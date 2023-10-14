package controllers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rustoma/octo-pulse/internal/api"
	"github.com/rustoma/octo-pulse/internal/dto"
	"github.com/rustoma/octo-pulse/internal/models"
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

func (c *ArticleController) HandleGenerateArticles(w http.ResponseWriter, r *http.Request) error {
	var request *dto.GenerateArticlesRequest
	err := api.ReadJSON(w, r, &request)
	if err != nil {
		return api.Error{Err: "bad request", Status: http.StatusBadRequest}
	}

	err = c.articleTasks.NewGenerateArticlesTask(request.DomainId, request.NumberOfArtilces, request.QuestionCategoryId)

	if err != nil {
		return api.Error{Err: err.Error(), Status: api.HandleErrorStatus(err)}
	}

	return api.WriteJSON(w, http.StatusOK, "Generate articles tasks created successfully")
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

func (c *ArticleController) HandleGetArticle(w http.ResponseWriter, r *http.Request) error {

	articleIdParam := chi.URLParam(r, "id")
	articleId, err := strconv.Atoi(articleIdParam)

	if err != nil {
		return api.Error{Err: "bad request", Status: http.StatusBadRequest}
	}

	article, err := c.articleService.GetArticle(articleId)

	if err != nil {
		return api.Error{Err: "Cannot get article", Status: api.HandleErrorStatus(err)}
	}

	return api.WriteJSON(w, http.StatusOK, article)
}

func (c *ArticleController) HandleUpdateArticle(w http.ResponseWriter, r *http.Request) error {
	var article *models.Article

	articleIdParam := chi.URLParam(r, "id")
	articleId, err := strconv.Atoi(articleIdParam)

	if err != nil {
		return api.Error{Err: "bad request", Status: http.StatusBadRequest}
	}

	err = api.ReadJSON(w, r, &article)

	if err != nil {
		return api.Error{Err: err.Error(), Status: http.StatusBadRequest}
	}

	updatedArticle, err := c.articleService.UpdateArticle(articleId, article)

	if err != nil {
		return api.Error{Err: err.Error(), Status: api.HandleErrorStatus(err)}
	}

	return api.WriteJSON(w, http.StatusOK, updatedArticle)
}
