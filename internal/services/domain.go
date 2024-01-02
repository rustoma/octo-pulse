package services

import (
	"github.com/rustoma/octo-pulse/internal/dto"
	"github.com/rustoma/octo-pulse/internal/models"
	"github.com/rustoma/octo-pulse/internal/storage"
	"github.com/rustoma/octo-pulse/internal/validator"
)

type DomainService interface {
	GetDomains() ([]*models.Domain, error)
	GetDomain(id int) (*models.Domain, error)
	GetDomainPublicData(id int) (*dto.DomainPublicData, error)
	CreateDomain(domain *models.Domain) (int, error)
}

type domainService struct {
	domainStore     storage.DomainStore
	domainValidator validator.DomainValidatorer
}

func NewDomainService(domainStore storage.DomainStore, domainValidator validator.DomainValidatorer) DomainService {
	return &domainService{domainStore: domainStore, domainValidator: domainValidator}
}

func (s *domainService) GetDomains() ([]*models.Domain, error) {
	return s.domainStore.GetDomains()
}

func (s *domainService) GetDomain(id int) (*models.Domain, error) {
	return s.domainStore.GetDomain(id)
}

func (s *domainService) CreateDomain(domain *models.Domain) (int, error) {
	err := s.domainValidator.Validate(domain)
	if err != nil {
		logger.Err(err).Send()
		return 0, err
	}

	return s.domainStore.InsertDomain(domain)
}

func (s *domainService) GetDomainPublicData(id int) (*dto.DomainPublicData, error) {
	return s.domainStore.GetDomainPublicData(id)
}
