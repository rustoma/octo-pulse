package postgresstore

import (
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	lr "github.com/rustoma/octo-pulse/internal/logger"
	"github.com/rustoma/octo-pulse/internal/storage"
)

var logger *zerolog.Logger

type PostgressStore struct {
	User              storage.UserStore
	Role              storage.RoleStore
	Domain            storage.DomainStore
	Category          storage.CategoryStore
	Author            storage.AuthorStore
	Article           storage.ArticleStore
	CategoriesDomains storage.CategoriesDomainsStore
	Image             storage.ImageStorageStore
	ImageCategory     storage.ImageCategoryStore
}

func NewPostgresStorage(DB *pgxpool.Pool) *PostgressStore {
	return &PostgressStore{
		User:              NewUserStore(DB),
		Role:              NewRoleStore(DB),
		Domain:            NewDomainStore(DB),
		Category:          NewCategoryStore(DB),
		Author:            NewAuthorStore(DB),
		Article:           NewArticleStore(DB, NewCategoryStore(DB), NewImageStorageStore(DB), NewAuthorStore(DB)),
		CategoriesDomains: NewCategoriesDomainsStore(DB),
		Image:             NewImageStorageStore(DB),
		ImageCategory:     NewImageCategoryStore(DB),
	}
}

func pgQb() squirrel.StatementBuilderType {
	return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
}

func init() {
	//Init logger
	l, logFile := lr.NewLogger()
	defer logFile.Close()
	logger = l
}
