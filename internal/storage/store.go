package storage

import (
	"github.com/rustoma/octo-pulse/internal/dto"
	"github.com/rustoma/octo-pulse/internal/models"
)

type Store struct {
	User              UserStore
	Role              RoleStore
	Domain            DomainStore
	Category          CategoryStore
	Author            AuthorStore
	Article           ArticleStore
	CategoriesDomains CategoriesDomainsStore
	Scrapper          ScrapperStore
	Image             ImageStorageStore
	ImageCategory     ImageCategoryStore
	BasicPage         BasicPageStore
}

type UserStore interface {
	GetUserByEmail(email string) (*models.User, error)
	UpdateRefreshToken(userId int, refreshToken string) (int, error)
	SelectUserByRefreshToken(refreshToken string) (*models.User, error)
	InsertUser(user *models.User) (int, error)
}

type RoleStore interface {
	InsertRole(role *models.Role) (int, error)
}

type DomainStore interface {
	InsertDomain(domain *models.Domain) (int, error)
	GetDomains() ([]*models.Domain, error)
	GetDomain(id int) (*models.Domain, error)
	GetDomainPublicData(id int) (*dto.DomainPublicData, error)
}

type GetCategoriesFilters struct {
	Slug string
}

type CategoryStore interface {
	InsertCategory(category *models.Category) (int, error)
	GetCategories(filters ...*GetCategoriesFilters) ([]*models.Category, error)
	GetCategory(id int) (*models.Category, error)
}

type CategoriesDomainsStore interface {
	AsignCategoryToDomain(categoryId int, domainId int) error
	GetDomainCategories(domainId int) ([]int, error)
}

type AuthorStore interface {
	InsertAuthor(author *models.Author) (int, error)
	GetAuthor(id int) (*models.Author, error)
}

type GetArticlesFilters struct {
	CategoryId int
	DomainId   int
	Limit      int
	Offset     int
	Featured   string
	Slug       string
}

type ArticleStore interface {
	InsertArticle(article *models.Article) (int, error)
	GetArticle(id int) (*models.Article, error)
	GetArticles(filters ...*GetArticlesFilters) ([]*dto.Article, error)
	UpdateArticle(id int, article *models.Article) (int, error)
	DeleteArticle(id int) (int, error)
}

type GetQuestionsFilters struct {
	CategoryId int
}

type ScrapperStore interface {
	GetQuestion(id int) (*models.Question, error)
	GetQuestions(filters ...*GetQuestionsFilters) ([]*models.Question, error)
	UpdateQuestion(id int, question *models.Question) error
}

type GetImagesFilters struct {
	CategoryId int
	Path       string
}

type ImageStorageStore interface {
	InsertImage(image *models.Image) (int, error)
	GetImage(id int) (*models.Image, error)
	GetImages(filters ...*GetImagesFilters) ([]*models.Image, error)
}

type ImageCategoryStore interface {
	InsertCategory(category *models.ImageCategory) (int, error)
}

type GetBasicPagesFilters struct {
	DomainId int
}

type GetBasicPageBySlugFilters struct {
	DomainId int
}

type BasicPageStore interface {
	InsertBasicPage(page *models.BasicPage) (int, error)
	GetBasicPages(filters ...*GetBasicPagesFilters) ([]*models.BasicPage, error)
	GetBasicPage(id int) (*models.BasicPage, error)
	GetBasicPageBySlug(slug string, filters ...*GetBasicPageBySlugFilters) (*models.BasicPage, error)
}
