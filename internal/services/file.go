package services

import (
	"fmt"
	"github.com/gosimple/slug"
	"github.com/rustoma/octo-pulse/internal/models"
	"github.com/rustoma/octo-pulse/internal/storage"
	"image/jpeg"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type FileService interface {
	CreateArticles(ids []int) error
	InsertJPGImagesFromDir(dirPath string, imageCategoryId int) error
	RenameFilesUsingSlug(dirPath string)
}

type fileService struct {
	articleStore  storage.ArticleStore
	domainStore   storage.DomainStore
	categoryStore storage.CategoryStore
	imageStore    storage.ImageStorageStore
}

func NewFileService(articleStore storage.ArticleStore, domainStore storage.DomainStore, categoryStore storage.CategoryStore, imageStore storage.ImageStorageStore) FileService {
	return &fileService{articleStore: articleStore, domainStore: domainStore, categoryStore: categoryStore, imageStore: imageStore}
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

func (s *fileService) InsertJPGImagesFromDir(dirPath string, imageCategoryId int) error {
	files, _ := os.ReadDir(dirPath)
	for _, imgFile := range files {

		imagesWithTheSamePath, err := s.imageStore.GetImages(&storage.GetImagesFilters{Path: filepath.Join("/", dirPath, imgFile.Name())})
		if err != nil {
			logger.Err(err).Msg("File name: " + imgFile.Name())
		}

		if len(imagesWithTheSamePath) > 0 {
			logger.Info().Msg("Image already exist on path: " + filepath.Join("/", dirPath, imgFile.Name()) + " File name: " + imgFile.Name())
			continue
		}

		if reader, err := os.Open(filepath.Join(dirPath, imgFile.Name())); err == nil {
			defer reader.Close()
			im, err := jpeg.DecodeConfig(reader)

			if err != nil {
				logger.Err(err).Msg("File name: " + imgFile.Name())
				continue
			}

			fileInfo, err := os.Stat(filepath.Join(dirPath, imgFile.Name()))
			if err != nil {
				return err
			}

			img := models.Image{
				Name:       slug.Make(imgFile.Name()),
				Path:       filepath.Join("/", dirPath, imgFile.Name()),
				Size:       int(fileInfo.Size()),
				Type:       ".jpg",
				Width:      im.Width,
				Height:     im.Height,
				Alt:        slug.Make(imgFile.Name()),
				CategoryId: imageCategoryId,
				CreatedAt:  time.Now().UTC(),
				UpdatedAt:  time.Now().UTC(),
			}

			_, err = s.imageStore.InsertImage(&img)
			if err != nil {
				return err
			}

		} else {
			return err
		}
	}

	return nil
}

func (s *fileService) RenameFilesUsingSlug(dirPath string) {

	list, err := os.ReadDir(dirPath)
	if err != nil {
		panic(err)
	}
	for _, file := range list {

		if file.Name() == ".DS_Store" {
			continue
		}

		name := file.Name()

		filename := path.Base(name)
		extension := path.Ext(name)
		filenameWithoutExt := filename[:len(filename)-len(extension)]

		newName := slug.Make(filenameWithoutExt) + extension

		err := os.Rename(filepath.Join(dirPath, name), filepath.Join(dirPath, newName))

		if err != nil {
			logger.Err(err).Send()
		}
	}
}
