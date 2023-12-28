package controllers

import (
	"fmt"
	"github.com/rustoma/octo-pulse/internal/storage"
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

	err = c.articleTasks.NewGenerateArticlesTask(request.DomainId, request.NumberOfArtilces, request.QuestionCategoryId, request.ImagesCategory)

	if err != nil {
		return api.Error{Err: err.Error(), Status: api.HandleErrorStatus(err)}
	}

	return api.WriteJSON(w, http.StatusOK, "Generate articles tasks created successfully")
}

func (c *ArticleController) HandleGenerateDescritption(w http.ResponseWriter, r *http.Request) error {
	var request *dto.GenerateDescriptionRequest
	err := api.ReadJSON(w, r, &request)
	if err != nil {
		return api.Error{Err: "bad request", Status: http.StatusBadRequest}
	}

	pageIdParam := chi.URLParam(r, "id")
	pageId, err := strconv.Atoi(pageIdParam)

	if err != nil {
		return api.Error{Err: "bad request", Status: http.StatusBadRequest}
	}

	err = c.articleTasks.NewGenerateDescriptionTask(pageId, request.QuestionId)

	if err != nil {
		return api.Error{Err: err.Error(), Status: api.HandleErrorStatus(err)}
	}

	return api.WriteJSON(w, http.StatusOK, "")
}

func (c *ArticleController) HandleGetArticles(w http.ResponseWriter, r *http.Request) error {
	domainIdParam := r.URL.Query().Get("domainId")
	categoryIdParam := r.URL.Query().Get("categoryId")
	limitParam := r.URL.Query().Get("limit")
	offsetParam := r.URL.Query().Get("offset")
	featuredParam := r.URL.Query().Get("featured")
	slug := r.URL.Query().Get("slug")
	excludeBodyParam := r.URL.Query().Get("excludeBody")

	var filters storage.GetArticlesFilters

	if domainIdParam != "" {
		domainId, err := strconv.Atoi(domainIdParam)
		if err != nil {
			return api.Error{Err: "bad request - domainId wrong format", Status: http.StatusBadRequest}
		}

		filters.DomainId = domainId
	}

	if categoryIdParam != "" {
		categoryId, err := strconv.Atoi(categoryIdParam)
		if err != nil {
			return api.Error{Err: "bad request - categoryId wrong format", Status: http.StatusBadRequest}
		}

		filters.CategoryId = categoryId
	}

	if limitParam != "" {
		limit, err := strconv.Atoi(limitParam)
		if err != nil {
			return api.Error{Err: "bad request - limit wrong format", Status: http.StatusBadRequest}
		}

		filters.Limit = limit
	}

	if offsetParam != "" {
		offset, err := strconv.Atoi(offsetParam)
		if err != nil {
			return api.Error{Err: "bad request - offset wrong format", Status: http.StatusBadRequest}
		}

		filters.Offset = offset
	}

	if featuredParam == "true" || featuredParam == "false" {
		filters.Featured = featuredParam
	}

	if excludeBodyParam == "true" {
		filters.ExcludeBody = excludeBodyParam
	}

	if slug != "" {
		filters.Slug = slug
	}

	articles, err := c.articleService.GetArticles(&filters)

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

func (c *ArticleController) HandleDeleteArticle(w http.ResponseWriter, r *http.Request) error {
	articleIdParam := chi.URLParam(r, "id")
	articleId, err := strconv.Atoi(articleIdParam)

	if err != nil {
		return api.Error{Err: "bad request", Status: http.StatusBadRequest}
	}

	article, err := c.articleService.DeleteArticle(articleId)

	if err != nil {
		return api.Error{Err: fmt.Sprintf("cannot delete article with ID: %d , err: %+v\n", articleId, err), Status: api.HandleErrorStatus(err)}
	}

	return api.WriteJSON(w, http.StatusOK, article)
}

func (c *ArticleController) HandleRemoveDuplicatesFromArticle(w http.ResponseWriter, r *http.Request) error {
	articleIdParam := chi.URLParam(r, "id")
	articleId, err := strconv.Atoi(articleIdParam)
	if err != nil {
		return api.Error{Err: "bad request", Status: http.StatusBadRequest}
	}

	err = c.articleService.RemoveDuplicateHeadingsFromArticle(articleId)
	if err != nil {
		return api.Error{Err: fmt.Sprintf("cannot remove duplicate from article with ID: %d , err: %+v\n", articleId, err), Status: api.HandleErrorStatus(err)}
	}

	return api.WriteJSON(w, http.StatusOK, "Duplicates removed successfully!")
}
