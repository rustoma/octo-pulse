package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gosimple/slug"
	"github.com/rustoma/octo-pulse/internal/utils"
	"math/rand"
	"os"
	"time"

	"github.com/hibiken/asynq"
	"github.com/rustoma/octo-pulse/internal/ai"
	"github.com/rustoma/octo-pulse/internal/models"
	"github.com/rustoma/octo-pulse/internal/services"
	"github.com/rustoma/octo-pulse/internal/storage"
)

const (
	TypeArticleGenerateDescription = "article:generateDescription"
	TypeArticleGenerateArticles    = "article:generateArticles"
)

type articleTasks struct {
	articleService  services.ArticleService
	domainService   services.DomainService
	scrapperService services.ScrapperService
	categoryService services.CategoryService
	imageService    services.ImageService
	ai              *ai.AI
	inspector       *asynq.Inspector
	scrapperTasks   scrapperTasks
}

func NewArticleTasks(
	articleService services.ArticleService,
	domainService services.DomainService,
	scrapperService services.ScrapperService,
	categoryService services.CategoryService,
	imageService services.ImageService,
	ai *ai.AI,
	scrapperTasks scrapperTasks,
) articleTasks {
	return articleTasks{
		articleService:  articleService,
		domainService:   domainService,
		scrapperService: scrapperService,
		categoryService: categoryService,
		imageService:    imageService,
		ai:              ai,
		scrapperTasks:   scrapperTasks,
	}
}

type DescriptionTaskPayload struct {
	ArticleId  int
	QuestionId int
}

type GenerateArticlesTaskPayload struct {
	DomainId                 int
	NumberOfArticlesToCreate int
	QuestionCategoryId       int
	ImagesCategory           int
}

func (t articleTasks) NewGenerateArticlesTask(domainId int, numberOfArticlesToCreate int, questionCategoryId int, imagesCategory int) error {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: os.Getenv("REDIS_ADDR"), Password: os.Getenv("REDIS_PASSWORD")})
	defer client.Close()

	payload, err := json.Marshal(GenerateArticlesTaskPayload{
		DomainId:                 domainId,
		NumberOfArticlesToCreate: numberOfArticlesToCreate,
		QuestionCategoryId:       questionCategoryId,
		ImagesCategory:           imagesCategory,
	})

	if err != nil {
		return err
	}

	task := asynq.NewTask(TypeArticleGenerateArticles, payload)
	info, err := client.Enqueue(task, asynq.MaxRetry(2), asynq.Timeout(2*time.Hour))

	logger.Info().Msgf("enqueued task: id=%s queue=%s", info.ID, info.Queue)

	if err != nil {
		return err
	}

	return nil
}

func (t articleTasks) NewGenerateDescriptionTask(articleId int, questionId int) error {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: os.Getenv("REDIS_ADDR"), Password: os.Getenv("REDIS_PASSWORD")})
	defer client.Close()

	payload, err := json.Marshal(DescriptionTaskPayload{ArticleId: articleId, QuestionId: questionId})
	if err != nil {
		return err
	}

	task := asynq.NewTask(TypeArticleGenerateDescription, payload)
	info, err := client.Enqueue(task, asynq.MaxRetry(2), asynq.Timeout(2*time.Hour))

	logger.Info().Msgf("enqueued task: id=%s queue=%s", info.ID, info.Queue)

	if err != nil {
		return err
	}

	return nil
}

func (t articleTasks) HandleGenerateDescription(ctx context.Context, task *asynq.Task) error {
	var payload DescriptionTaskPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	article, err := t.articleService.GetArticle(payload.ArticleId)
	question, err := t.scrapperService.GetQuestion(payload.QuestionId)

	if err != nil {
		return err
	}

	if article == nil {
		return fmt.Errorf("article with %d not found", payload.ArticleId)
	}

	if question == nil {
		return fmt.Errorf("question with %d not found", payload.QuestionId)
	}

	description, err := t.articleService.GenerateDescription(question)

	if err != nil {
		return err
	}

	article.Body = description

	readingTime := utils.CalculateReadTime(description)
	article.ReadingTime = &readingTime
	article.IsPublished = true
	article.PublicationDate = time.Now().UTC()

	_, err = t.articleService.UpdateArticle(payload.ArticleId, article)
	if err != nil {
		return err
	}

	err = t.articleService.RemoveDuplicateHeadingsFromArticle(payload.ArticleId)
	if err != nil {
		logger.Err(err).Msgf("Cannot remove duplicates for article id: %s", payload.ArticleId)
	}

	return nil
}

func (t articleTasks) HandleGenerateArticles(ctx context.Context, task *asynq.Task) error {
	var payload GenerateArticlesTaskPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	logger.Info().Interface("payload", payload).Send()

	questions, err := t.scrapperService.GetQuestions(&storage.GetQuestionsFilters{CategoryId: payload.QuestionCategoryId})

	if err != nil {
		return err
	}

	domainCategories, err := t.categoryService.GetDomainCategories(payload.DomainId)
	if err != nil {
		return err
	}

	logger.Info().Interface("Domain categories", domainCategories).Send()

	//---------------------

	createdArticles := 0
	for _, question := range questions {

		if createdArticles == payload.NumberOfArticlesToCreate {
			break
		}

		//ensures equal distribution of articles for categories
		categoriesMap := make(map[string]int, len(domainCategories))

		for _, category := range domainCategories {
			articlesFromCategory, err := t.articleService.GetArticles(&storage.GetArticlesFilters{CategoryId: category.ID, DomainId: payload.DomainId})
			logger.Info().Interface("category: ", category.ID).Send()
			if err != nil {
				logger.Err(err).Send()
			}
			logger.Info().Interface("category articles: ", len(articlesFromCategory)).Send()
			categoriesMap[category.Slug] = len(articlesFromCategory)
		}

		filteredCategories, err := filterCategoriesByEqualDistribution(domainCategories, categoriesMap)
		if err != nil {
			logger.Err(err).Send()
		}

		logger.Info().Interface("Filtered categories to which an article can be assigned: ", filteredCategories).Send()

		catgoryId, err := t.ai.ChatGPT.AssignToCategory(filteredCategories, question)
		if err != nil {
			return err
		}

		if catgoryId == 0 {
			logger.Info().Msg("There is no category that fits")
			continue
		}

		logger.Info().Interface("Assigned to category", catgoryId).Send()

		//Get random thumbnail
		var thumbnailId *int
		if payload.ImagesCategory != 0 {
			imagesFilter := &storage.GetImagesFilters{
				CategoryId: payload.ImagesCategory,
			}
			thumbnails, err := t.imageService.GetImages(imagesFilter)
			if err != nil {
				logger.Err(err).Send()
			}

			if len(thumbnails) > 0 {
				source := rand.NewSource(time.Now().UnixNano())
				random := rand.New(source)
				thumbnail := thumbnails[random.Intn(len(thumbnails))]
				thumbnailId = &thumbnail.ID
			}
		}

		article := &models.Article{
			Title:       question.Question,
			Slug:        slug.Make(question.Question),
			Body:        "Treść w przygotowaniu",
			Thumbnail:   thumbnailId,
			CategoryId:  catgoryId,
			AuthorId:    1,
			DomainId:    payload.DomainId,
			Featured:    false,
			IsSponsored: false,
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
		}

		articleId, err := t.articleService.CreateArticle(article)
		if err != nil {
			return err
		}

		logger.Info().Interface("CreatedArticle ID", articleId).Send()

		//Increase number of created articles
		createdArticles++

		err = t.scrapperTasks.NewUpdateQuestionTask(question.Id, question)
		if err != nil {
			_, _ = t.articleService.DeleteArticle(articleId)
			return err
		}

		//Generate Description For article
		_ = t.NewGenerateDescriptionTask(articleId, question.Id)
	}

	return nil
}

func filterCategoriesByEqualDistribution(categories []*models.Category, categoriesMap map[string]int) ([]*models.Category, error) {

	filteredCategoriesMap := findMaxMin(categoriesMap)

	filteredCategories := filterCategoriesByMap(categories, filteredCategoriesMap)

	return filteredCategories, nil
}

func filterCategoriesByMap(categories []*models.Category, categoriesMap map[string]int) []*models.Category {
	var filtered []*models.Category

	for _, cat := range categories {
		if _, ok := categoriesMap[cat.Slug]; ok {
			filtered = append(filtered, cat)
		}
	}

	return filtered
}
func findMaxMin(categoriesMap map[string]int) map[string]int {
	for {
		// Find categories with min and max number of articles
		maxCat := ""
		minArticles := int(^uint(0) >> 1) // Set to max possible int value initially
		maxArticles := 0

		for cat, numArticles := range categoriesMap {
			if numArticles < minArticles {
				minArticles = numArticles
			}
			if numArticles > maxArticles {
				maxArticles = numArticles
				maxCat = cat
			}
		}

		// Check the condition and remove if the difference is more than 2
		if maxArticles-minArticles > 2 {
			delete(categoriesMap, maxCat)
			findMaxMin(categoriesMap)
		} else {
			return categoriesMap
		}
	}
}
