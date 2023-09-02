package fixtures

import (
	"time"

	"github.com/rustoma/octo-pulse/internal/models"
	"github.com/rustoma/octo-pulse/internal/services"
)

type fixtures struct {
	authService services.AuthService
}

func NewFixtures(authService services.AuthService) *fixtures {
	return &fixtures{
		authService,
	}
}

func (f *fixtures) CreateUser(email string, pass string, roleId int) *models.User {

	hashedPass, _ := f.authService.HashPassword(pass)

	user := &models.User{
		Email:        email,
		RefreshToken: "",
		PasswordHash: hashedPass,
		RoleId:       roleId,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	return user
}

func (f *fixtures) CreateRole(roleName string) *models.Role {
	return &models.Role{
		RoleName:  roleName,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}

func (f *fixtures) CreateDomain(domainName string) *models.Domain {
	return &models.Domain{
		DomainName: domainName,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}
}

func (f *fixtures) CreateCategory(categoryName string) *models.Category {
	return &models.Category{
		CategoryName: categoryName,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}
}

func (f *fixtures) CreateAuthor(fn, ln, desc, imageUrl string) *models.Author {
	return &models.Author{
		FirstName:   fn,
		LastName:    ln,
		Description: desc,
		ImageUrl:    imageUrl,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
}

func (f *fixtures) CreateArticle(title, desc, imageUrl string, isPub bool, authorId, categoryId, domainId int) *models.Article {
	return &models.Article{
		Title:           title,
		Description:     desc,
		ImageUrl:        imageUrl,
		PublicationDate: time.Now().UTC(),
		IsPublished:     isPub,
		AuthorId:        authorId,
		CategoryId:      categoryId,
		DomainId:        domainId,
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
	}
}
