package storage

import (
	"github.com/rustoma/octo-pulse/internal/models"
)

type Store struct {
	User     UserStore
	Role     RoleStore
	Domain   DomainStore
	Category CategoryStore
	Author   AuthorStore
	Article  ArticleStore
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
}

type CategoryStore interface {
	InsertCategory(category *models.Category) (int, error)
}

type AuthorStore interface {
	InsertAuthor(author *models.Author) (int, error)
}

type ArticleStore interface {
	InsertArticle(article *models.Article) (int, error)
}
