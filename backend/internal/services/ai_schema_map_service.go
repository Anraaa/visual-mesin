package services

import (
	"github.com/anraaa/visual-mesin/internal/models"
	"github.com/anraaa/visual-mesin/internal/repository"
)

type AiSchemaMapService struct {
	repo *repository.AiSchemaMapRepository
}

func NewAiSchemaMapService(repo *repository.AiSchemaMapRepository) *AiSchemaMapService {
	return &AiSchemaMapService{repo: repo}
}

func (s *AiSchemaMapService) List(page, limit int) ([]models.AiSchemaMap, int64, error) {
	return s.repo.List(page, limit)
}

func (s *AiSchemaMapService) GetByID(id uint) (*models.AiSchemaMap, error) {
	return s.repo.FindByID(id)
}

func (s *AiSchemaMapService) Create(req *models.AiSchemaMapRequest) (*models.AiSchemaMap, error) {
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	m := &models.AiSchemaMap{
		IntentName:      req.IntentName,
		Keywords:        req.Keywords,
		TablesInvolved:  req.TablesInvolved,
		SchemaContext:   req.SchemaContext,
		FewShotExamples: req.FewShotExamples,
		Description:     req.Description,
		IsActive:        isActive,
	}

	err := s.repo.Create(m)
	return m, err
}

func (s *AiSchemaMapService) Update(id uint, req *models.AiSchemaMapUpdateRequest) (*models.AiSchemaMap, error) {
	m, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.IntentName != nil {
		m.IntentName = *req.IntentName
	}
	if req.Keywords != nil {
		m.Keywords = *req.Keywords
	}
	if req.TablesInvolved != nil {
		m.TablesInvolved = *req.TablesInvolved
	}
	if req.SchemaContext != nil {
		m.SchemaContext = req.SchemaContext
	}
	if req.FewShotExamples != nil {
		m.FewShotExamples = req.FewShotExamples
	}
	if req.Description != nil {
		m.Description = req.Description
	}
	if req.IsActive != nil {
		m.IsActive = *req.IsActive
	}

	err = s.repo.Update(m)
	return m, err
}

func (s *AiSchemaMapService) Delete(id uint) error {
	return s.repo.Delete(id)
}
