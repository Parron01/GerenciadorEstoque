package service

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Parron01/GerenciadorEstoque/backendGo/internal/models"
	"github.com/Parron01/GerenciadorEstoque/backendGo/internal/repository"
	"github.com/google/uuid"
)

const (
	EntityTypeProduct           = "product"
	EntityTypeLote              = "lote"
	EntityTypeProductBatchContext = "product_batch_context" // New entity type
)

// HistoryService defines the interface for history operations
type HistoryService interface {
	RecordChange(entityType string, entityID string, changeDetail interface{}, operationBatchIDHeader ...string) error
	GetHistory(limit, offset int) ([]models.History, error)
	GetHistoryForEntity(entityType, entityID string) ([]models.History, error)
	CreateRawHistoryEntry(entry models.History) error
	CreateBatch(entries []models.History) (string, error)
	GetByBatchID(batchID string) ([]models.History, error)
	GetGroupedHistory(page, pageSize int) (*models.PaginatedHistoryBatchGroups, error)
}

type historyService struct {
	repo        repository.HistoryRepository
	productRepo repository.ProductRepository // Added product repository
}

// NewHistoryService creates a new HistoryService
func NewHistoryService(repo repository.HistoryRepository, productRepo repository.ProductRepository) HistoryService {
	return &historyService{repo: repo, productRepo: productRepo}
}

// RecordChange creates a new history entry
func (s *historyService) RecordChange(entityType string, entityID string, changeDetail interface{}, operationBatchIDHeader ...string) error {
	jsonData, err := json.Marshal(changeDetail)
	if err != nil {
		return fmt.Errorf("failed to marshal change detail: %w", err)
	}

	batchID := ""
	if len(operationBatchIDHeader) > 0 && operationBatchIDHeader[0] != "" {
		batchID = operationBatchIDHeader[0]
	} else {
		// If no batch ID from header, generate a new one for this single operation
		// or use the history entry's own ID as its batch ID.
		// For now, let's ensure every record has a BatchID, defaulting to its own ID if not part of a larger client-defined batch.
		// This will be overridden if CreateRawHistoryEntry is called with an explicit BatchID.
	}

	historyEntry := models.History{
		ID:         uuid.NewString(),
		Date:       time.Now().Format(time.RFC3339),
		EntityType: entityType,
		EntityID:   entityID,
		Changes:    jsonData,
		BatchID:    batchID, // Will be set properly by CreateRawHistoryEntry
	}

	return s.CreateRawHistoryEntry(historyEntry)
}

// GetHistory retrieves a paginated list of all history entries
func (s *historyService) GetHistory(limit, offset int) ([]models.History, error) {
	return s.repo.GetHistory(limit, offset)
}

// GetHistoryForEntity retrieves history for a specific entity
func (s *historyService) GetHistoryForEntity(entityType, entityID string) ([]models.History, error) {
	return s.repo.GetHistoryByEntity(entityType, entityID)
}

// CreateRawHistoryEntry directly creates a history entry in the database.
// If entry.BatchID is empty, it defaults to entry.ID.
func (s *historyService) CreateRawHistoryEntry(entry models.History) error {
	if entry.ID == "" {
		entry.ID = uuid.NewString()
	}
	if entry.BatchID == "" {
		entry.BatchID = entry.ID // Default BatchID to the record's own ID if not provided
	}
	if entry.Date == "" {
		entry.Date = time.Now().Format(time.RFC3339)
	}
	return s.repo.Create(&entry)
}

// CreateBatch creates multiple history entries with a shared, new batch ID.
func (s *historyService) CreateBatch(entries []models.History) (string, error) {
	if len(entries) == 0 {
		return "", fmt.Errorf("no entries provided for batch creation")
	}
	batchID := uuid.NewString()
	for i := range entries {
		entries[i].BatchID = batchID
		if entries[i].ID == "" {
			entries[i].ID = uuid.NewString()
		}
		if entries[i].Date == "" {
			entries[i].Date = time.Now().Format(time.RFC3339)
		}
	}
	return batchID, s.repo.CreateBatch(entries)
}

// GetByBatchID retrieves all history entries for a specific batch ID.
func (s *historyService) GetByBatchID(batchID string) ([]models.History, error) {
	return s.repo.GetByBatchID(batchID)
}

// GetGroupedHistory retrieves history entries grouped by batch ID, with pagination for batches.
func (s *historyService) GetGroupedHistory(page, pageSize int) (*models.PaginatedHistoryBatchGroups, error) {
	paginatedRawGroups, err := s.repo.GetGroupedHistoryBatches(page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to get grouped history batches from repo: %w", err)
	}

	processedGroups := make([]models.HistoryBatchGroup, len(paginatedRawGroups.Groups))

	for i, rawGroup := range paginatedRawGroups.Groups {
		// Map to store product batch context data fetched once per productID per batch
		productBatchContexts := make(map[string]models.ProductBatchContextChangeDetail) // Key: ProductID
		processedRecords := make([]models.History, len(rawGroup.Records))

		// First pass: Extract ProductBatchContextChangeDetail for each product in the batch
		// These are specific history records that store the snapshot.
		for _, record := range rawGroup.Records {
			if record.EntityType == EntityTypeProductBatchContext {
				var ctxDetail models.ProductBatchContextChangeDetail
				if errUnmarshal := json.Unmarshal(record.Changes, &ctxDetail); errUnmarshal == nil {
					// Store the context detail against the product ID it refers to (which is record.EntityID for these context records)
					productBatchContexts[record.EntityID] = ctxDetail
				} else {
					log.Printf("WARN: GetGroupedHistory - Failed to unmarshal ProductBatchContextChangeDetail for EntityID %s in BatchID %s: %v", record.EntityID, rawGroup.BatchID, errUnmarshal)
				}
			}
		}
		
		productSummaries := make(map[string]models.ProductBatchSummary)
		// This map can still be useful for lote-specific net changes if needed for other purposes,
		// but the summary will primarily use the ProductBatchContext.
		// productNetQuantityChangesInBatch := make(map[string]float64) 

		// Second pass: Process records and build summaries using the extracted context
		for j, record := range rawGroup.Records {
			processedRecord := record // Copy

			// Skip populating context for the context records themselves
			if record.EntityType == EntityTypeProductBatchContext {
				processedRecords[j] = processedRecord
				continue
			}

			var productIDForContext string
			if record.EntityType == EntityTypeLote {
				var loteDetail models.LoteChangeDetail
				if json.Unmarshal(record.Changes, &loteDetail) == nil {
					productIDForContext = loteDetail.ProductID
				}
			} else if record.EntityType == EntityTypeProduct {
				productIDForContext = record.EntityID
			}

			if productIDForContext != "" {
				if ctx, ok := productBatchContexts[productIDForContext]; ok {
					processedRecord.ProductNameContext = ctx.ProductNameSnapshot
					// productCurrentTotalQuantity on individual records now means quantity *after* this batch for this product
					processedRecord.ProductCurrentTotalQuantity = &ctx.QuantityAfterBatch 
				} else {
					// Fallback if context record is missing (e.g., older data before this system, or an error saving context)
					processedRecord.ProductNameContext = "Context Unavailable"
					// Do not set ProductCurrentTotalQuantity if context is missing to avoid misleading data
				}
			}
			processedRecords[j] = processedRecord
		}

		// Build ProductSummaries using the data from productBatchContexts
		// Each product mentioned in a productBatchContexts gets a summary.
		for productID, ctx := range productBatchContexts {
			productSummaries[productID] = models.ProductBatchSummary{
				ProductID:                productID, // This is ctx.ProductID, which is the key
				ProductName:              ctx.ProductNameSnapshot,
				TotalQuantityBeforeBatch: ctx.QuantityBeforeBatch,
				TotalQuantityAfterBatch:  ctx.QuantityAfterBatch,
				// Net change is derived directly from the snapshot
				NetQuantityChangeInBatch: ctx.QuantityAfterBatch - ctx.QuantityBeforeBatch, 
			}
		}


		processedGroups[i] = models.HistoryBatchGroup{
			BatchID:          rawGroup.BatchID,
			CreatedAt:        rawGroup.CreatedAt,
			Records:          processedRecords, // These are all records, including the context ones if desired for full audit
			RecordCount:      rawGroup.RecordCount,
			ProductSummaries: productSummaries,
		}
	}

	return &models.PaginatedHistoryBatchGroups{
		Groups:       processedGroups,
		TotalBatches: paginatedRawGroups.TotalBatches,
		Page:         paginatedRawGroups.Page,
		PageSize:     paginatedRawGroups.PageSize,
		TotalPages:   paginatedRawGroups.TotalPages,
	}, nil
}
