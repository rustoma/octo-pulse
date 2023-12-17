package services

import (
	"fmt"
	"github.com/gosimple/slug"
	"github.com/rustoma/octo-pulse/internal/storage"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type FileService interface {
	CreateArticles(ids []int) error
}

type fileService struct {
	articleStore  storage.ArticleStore
	domainStore   storage.DomainStore
	categoryStore storage.CategoryStore
}

func NewFileService(articleStore storage.ArticleStore, domainStore storage.DomainStore, categoryStore storage.CategoryStore) FileService {
	return &fileService{articleStore: articleStore, domainStore: domainStore, categoryStore: categoryStore}
}

func (s *fileService) ConvertHTMLToDocx(htmlPath, docxPath string) error {
	cmd := exec.Command("pandoc", "-s", htmlPath, "-o", docxPath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to execute pandoc command: %w", err)
	}

	return nil
}

func (s *fileService) CreateHtmlFile(htmlFilePath string, content string) error {
	f, err := os.Create(htmlFilePath)
	defer f.Close()
	if err != nil {
		return err
	}

	htmlEntry := "<!DOCTYPE html><html lang=\"pl\"><head><meta charset=\"UTF-8\">\n\t<meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"></head><body>\n"
	htmlEnd := "\n</body>\n\n</html>"
	htmlContent := fmt.Sprintf("%s%s%s", htmlEntry, content, htmlEnd)

	_, err = f.WriteString(strings.TrimSpace(htmlContent))
	if err != nil {
		return err
	}

	return nil
}

func (s *fileService) CreateArticles(ids []int) error {

	for _, id := range ids {
		article, err := s.articleStore.GetArticle(id)
		if err != nil {
			return err
		}

		domain, err := s.domainStore.GetDomain(article.DomainId)
		if err != nil {
			return err
		}

		category, err := s.categoryStore.GetCategory(article.CategoryId)
		if err != nil {
			return err
		}

		if err := os.MkdirAll(filepath.Join("assets", "articles", slug.Make(domain.Name), slug.Make(category.Name)), os.ModePerm); err != nil {
			logger.Err(err).Send()
			return err
		}

		htmlFilePath := "temp.html"
		docxFilePath := filepath.Join("assets", "articles", slug.Make(domain.Name), slug.Make(category.Name), fmt.Sprintf("%s.docx", slug.Make(article.Title)))

		err = s.CreateHtmlFile(htmlFilePath, article.Body)
		if err != nil {
			logger.Err(err).Send()
			return err
		}

		if err := s.ConvertHTMLToDocx(htmlFilePath, docxFilePath); err != nil {
			logger.Err(err).Send()
			return err
		}

		err = os.Remove(htmlFilePath)
		if err != nil {
			logger.Err(err).Send()
			return err
		}
	}

	return nil
}
