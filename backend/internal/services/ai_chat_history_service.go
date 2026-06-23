package services

import (
	"github.com/anraaa/visual-mesin/internal/models"
	"github.com/anraaa/visual-mesin/internal/repository"
)

type AiChatHistoryService struct {
	repo *repository.AiChatHistoryRepository
}

func NewAiChatHistoryService(repo *repository.AiChatHistoryRepository) *AiChatHistoryService {
	return &AiChatHistoryService{repo: repo}
}

func (s *AiChatHistoryService) GetByID(id uint) (*models.AiChatHistory, error) {
	return s.repo.FindByID(id)
}

func (s *AiChatHistoryService) Create(h *models.AiChatHistory) error {
	return s.repo.Create(h)
}

func (s *AiChatHistoryService) Update(h *models.AiChatHistory) error {
	return s.repo.Update(h)
}

func (s *AiChatHistoryService) ListBySession(sessionID string, page, limit int) ([]models.AiChatHistory, int64, error) {
	return s.repo.ListBySession(sessionID, page, limit)
}

func (s *AiChatHistoryService) ListSessions(userID uint) ([]models.AiChatHistory, error) {
	return s.repo.ListSessionsByUser(userID)
}

func (s *AiChatHistoryService) DeleteSession(sessionID string) error {
	return s.repo.DeleteSession(sessionID)
}
