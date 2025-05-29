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
	RecordChange(entityType, entityID string, changeDetails interface{}, batchID ...string) error
	GetHistory(limit, offset int) ([]models.History, error)
	GetHistoryForEntity(entityType, entityID string) ([]models.History, error)
	CreateRawHistoryEntry(entry models.History) error
	CreateBatch(entries []models.History) (string, error) // New: Create batch of history entries
	GetByBatchID(batchID string) ([]models.History, error) // New: Get entries by batch ID
	GetGroupedHistory(page, pageSize int) (*models.PaginatedHistoryBatchGroups, error)
}

type historyService struct {
	repo repository.HistoryRepository
}

func NewHistoryService(repo repository.HistoryRepository) HistoryService {
	return &historyService{repo: repo}
}

// RecordChange now accepts an optional batchID
func (s *historyService) RecordChange(entityType, entityID string, changeDetails interface{}, batchID ...string) error {
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

	// Use provided batchID if available.
	// If not, the repository layer will default batch_id to entry.id if it's empty.
	if len(batchID) > 0 && batchID[0] != "" {
		entry.BatchID = batchID[0]
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
    return s.repo.Create(&entry)
}

// New method to create a batch of history entries with a common batchID
func (s *historyService) CreateBatch(entries []models.History) (string, error) {
    if len(entries) == 0 {
        return "", fmt.Errorf("no history entries provided for batch creation")
    }

    // Generate a common batch ID for all entries
    batchID := uuid.NewString()
    now := time.Now().Format(time.RFC3339)

    // Set consistent values for all entries
    for i := range entries {
        if entries[i].ID == "" {
            entries[i].ID = uuid.NewString()
        }
        if entries[i].Date == "" {
            entries[i].Date = now
        }
        entries[i].BatchID = batchID // Use the common batchID
    }

    if err := s.repo.CreateBatch(entries); err != nil {
        return "", fmt.Errorf("failed to create history batch: %w", err)
    }

    return batchID, nil
}

// New method to retrieve history entries by batchID
func (s *historyService) GetByBatchID(batchID string) ([]models.History, error) {
    if batchID == "" {
        return nil, fmt.Errorf("batch ID is required")
    }
    return s.repo.GetByBatchID(batchID)
}

func (s *historyService) GetGroupedHistory(page, pageSize int) (*models.PaginatedHistoryBatchGroups, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10 // Default page size for batches
	}
	offset := (page - 1) * pageSize

	batchIDs, firstEntryDates, totalBatches, err := s.repo.GetDistinctBatchIDs(pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get distinct batch IDs: %w", err)
	}

	groups := make([]models.HistoryBatchGroup, 0, len(batchIDs))
	for i, batchID := range batchIDs {
		records, err := s.repo.GetByBatchID(batchID)
		if err != nil {
			// Log or handle error for individual batch, maybe skip this batch
			fmt.Printf("Warning: failed to get records for batch ID %s: %v\n", batchID, err)
			continue
		}
		if len(records) > 0 {
			groups = append(groups, models.HistoryBatchGroup{
				BatchID:     batchID,
				CreatedAt:   firstEntryDates[i], // Use the fetched first entry date
				Records:     records,
				RecordCount: len(records),
			})
		}
	}

	totalPages := 0
	if totalBatches > 0 {
		totalPages = (totalBatches + pageSize - 1) / pageSize
	}

	return &models.PaginatedHistoryBatchGroups{
		Groups:       groups,
		TotalBatches: totalBatches,
		Page:         page,
		PageSize:     pageSize,
		TotalPages:   totalPages,
	}, nil
}
