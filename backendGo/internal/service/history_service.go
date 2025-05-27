package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Parron01/GerenciadorEstoque/backendGo/internal/models"
	"github.com/Parron01/GerenciadorEstoque/backendGo/internal/repository"
	"github.com/google/uuid"
)

const (
	EntityTypeProduct = "product"
	EntityTypeLote    = "lote"
)

type HistoryService interface {
	RecordChange(entityType, entityID string, changeDetails interface{}) error
	GetHistory(limit, offset int) ([]models.History, error)
	GetHistoryForEntity(entityType, entityID string) ([]models.History, error)
	CreateRawHistoryEntry(entry models.History) error
}

type historyService struct {
	repo repository.HistoryRepository
}

func NewHistoryService(repo repository.HistoryRepository) HistoryService {
	return &historyService{repo: repo}
}

func (s *historyService) RecordChange(entityType, entityID string, changeDetails interface{}) error {
	changesJSON, err := json.Marshal(changeDetails)
	if err != nil {
		return fmt.Errorf("failed to marshal change details to JSON: %w", err)
	}

	entry := models.History{
		ID:         uuid.NewString(),
		Date:       time.Now().Format(time.RFC3339), // Or your preferred date format
		EntityType: entityType,
		EntityID:   entityID,
		Changes:    changesJSON,
	}

	return s.repo.Create(&entry)
}

func (s *historyService) GetHistory(limit, offset int) ([]models.History, error) {
	if limit <= 0 {
		limit = 20 // Default limit
	}
	if offset < 0 {
		offset = 0 // Default offset
	}
	return s.repo.GetAll(limit, offset)
}

func (s *historyService) GetHistoryForEntity(entityType, entityID string) ([]models.History, error) {
	return s.repo.GetByEntityID(entityType, entityID)
}

// CreateRawHistoryEntry allows creating a history entry directly, typically used by HistoryController
func (s *historyService) CreateRawHistoryEntry(entry models.History) error {
    if entry.ID == "" {
        entry.ID = uuid.NewString()
    }
    if entry.Date == "" {
        entry.Date = time.Now().Format(time.RFC3339)
    }
    // Ensure EntityType and EntityID are provided if this method is used
    if entry.EntityType == "" || entry.EntityID == "" {
        // This might be relaxed if the `changes` blob contains enough info
        // and a pre-processing step populates these fields.
        // For now, require them for raw entries.
        // However, the original controller allowed creating history without these.
        // Let's align with the original controller's flexibility for now.
    }
    return s.repo.Create(&entry)
}
