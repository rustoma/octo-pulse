package services

import (
	"github.com/rustoma/octo-pulse/internal/models"
	"github.com/rustoma/octo-pulse/internal/storage"
)

type DomainService interface {
	GetDomains() ([]*models.Domain, error)
}

type domainService struct {
	domainStore storage.DomainStore
}

func NewDomainService(domainStore storage.DomainStore) DomainService {
	return &domainService{domainStore: domainStore}
}

func (s *domainService) GetDomains() ([]*models.Domain, error) {
	return s.domainStore.GetDomains()
}
