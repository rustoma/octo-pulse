package fixtures

import (
	"github.com/gosimple/slug"
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

func (f *fixtures) CreateRole(name string) *models.Role {
	return &models.Role{
		Name:      name,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}

func (f *fixtures) CreateDomain(name string) *models.Domain {
	return &models.Domain{
		Name:      name,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}

func (f *fixtures) CreateCategory(name string) *models.Category {
	return &models.Category{
		Name:      name,
		Slug:      slug.Make(name),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
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

func (f *fixtures) CreateArticle(title, desc string, thumbnail int, isPub bool, authorId, categoryId, domainId int, featured bool) *models.Article {
	return &models.Article{
		Title:           title,
		Slug:            slug.Make(title),
		Description:     desc,
		Thumbnail:       &thumbnail,
		PublicationDate: time.Now().UTC(),
		IsPublished:     isPub,
		AuthorId:        authorId,
		CategoryId:      categoryId,
		DomainId:        domainId,
		Featured:        featured,
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
	}
}

func (f *fixtures) CreateImage(name string, path string, size int, t string, width int, height int, alt string, categoryId int) *models.Image {
	return &models.Image{
		Name:       name,
		Path:       path,
		Size:       size,
		Type:       t,
		Width:      width,
		Height:     height,
		Alt:        alt,
		CategoryId: categoryId,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}
}

func (f *fixtures) CreateImageCategory(name string) *models.ImageCategory {
	return &models.ImageCategory{
		Name:      name,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}
