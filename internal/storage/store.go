package storage

import (
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
}

type CategoryStore interface {
	InsertCategory(category *models.Category) (int, error)
	GetCategories() ([]*models.Category, error)
	GetCategory(id int) (*models.Category, error)
}

type CategoriesDomainsStore interface {
	AsignCategoryToDomain(categoryId int, domainId int) error
}

type AuthorStore interface {
	InsertAuthor(author *models.Author) (int, error)
}

type ArticleStore interface {
	InsertArticle(article *models.Article) (int, error)
	GetArticle(id int) (*models.Article, error)
	GetArticles() ([]*models.Article, error)
	UpdateArticle(id int, article *models.Article) (int, error)
}

type ScrapperStore interface {
	GetQuestion(id int) (*models.Question, error)
}
