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
	EntityTypeProduct = "product"
	EntityTypeLote    = "lote"
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
		processedRecords := make([]models.History, len(rawGroup.Records))
		// Temporary maps to store product data fetched once per product per batch
		productInfoCache := make(map[string]*models.Product) // Cache for product details

		// Calculate net quantity changes for each product within this batch
		productNetQuantityChangesInBatch := make(map[string]float64) // Key: ProductID

		for _, record := range rawGroup.Records {
			var productID string
			var errParse error

			if record.EntityType == EntityTypeLote {
				var loteDetail models.LoteChangeDetail
				errParse = json.Unmarshal(record.Changes, &loteDetail)
				if errParse == nil {
					productID = loteDetail.ProductID
					var netLoteChange float64 = 0
					if loteDetail.Action == "created" && loteDetail.QuantityAfter != nil {
						netLoteChange = *loteDetail.QuantityAfter
					} else if loteDetail.Action == "deleted" && loteDetail.QuantityBefore != nil {
						netLoteChange = -(*loteDetail.QuantityBefore)
					} else if loteDetail.Action == "updated" && loteDetail.QuantityAfter != nil && loteDetail.QuantityBefore != nil {
						netLoteChange = *loteDetail.QuantityAfter - *loteDetail.QuantityBefore
					}
					if productID != "" {
						productNetQuantityChangesInBatch[productID] += netLoteChange
					}
				}
			} else if record.EntityType == EntityTypeProduct {
				productID = record.EntityID
				// For direct product changes (e.g., non-lote product quantity update),
				// this might also contribute to net quantity change.
				// However, if product quantity is strictly derived from lotes via DB trigger,
				// direct quantity changes in ProductChange might be informational or for non-lote items.
				// For now, focusing productNetQuantityChangesInBatch on lote-driven changes for simplicity.
			}
		}

		// Populate processed records with context and build product summaries
		productSummaries := make(map[string]models.ProductBatchSummary)

		for j, record := range rawGroup.Records {
			processedRecord := record // Copy
			var productIDForContext string
			var errParseDetail error

			if record.EntityType == EntityTypeLote {
				var loteDetail models.LoteChangeDetail
				errParseDetail = json.Unmarshal(record.Changes, &loteDetail)
				if errParseDetail == nil {
					productIDForContext = loteDetail.ProductID
				}
			} else if record.EntityType == EntityTypeProduct {
				productIDForContext = record.EntityID
			}

			if productIDForContext != "" {
				// Fetch product from cache or DB
				product, exists := productInfoCache[productIDForContext]
				if !exists {
					fetchedProduct, errDb := s.productRepo.GetByID(productIDForContext)
					if errDb == nil && fetchedProduct != nil {
						productInfoCache[productIDForContext] = fetchedProduct
						product = fetchedProduct
					} else {
						log.Printf("WARN: GetGroupedHistory - Could not fetch product %s for context: %v", productIDForContext, errDb)
						// Add a placeholder to cache to avoid refetching on error for this batch
						productInfoCache[productIDForContext] = &models.Product{ID: productIDForContext, Name: "Unknown Product"}
						product = productInfoCache[productIDForContext]
					}
				}

				if product != nil {
					processedRecord.ProductNameContext = product.Name
					currentQty := product.Quantity // This is the product's quantity *after* all DB operations
					processedRecord.ProductCurrentTotalQuantity = &currentQty

					// Build or update product summary for this productID
					if _, summaryExists := productSummaries[productIDForContext]; !summaryExists {
						netChangeForThisProduct := productNetQuantityChangesInBatch[productIDForContext]
						productSummaries[productIDForContext] = models.ProductBatchSummary{
							ProductID:                productIDForContext,
							ProductName:              product.Name,
							TotalQuantityAfterBatch:  currentQty,
							NetQuantityChangeInBatch: netChangeForThisProduct,
							TotalQuantityBeforeBatch: currentQty - netChangeForThisProduct,
						}
					}
				}
			}
			processedRecords[j] = processedRecord
		}

		processedGroups[i] = models.HistoryBatchGroup{
			BatchID:          rawGroup.BatchID,
			CreatedAt:        rawGroup.CreatedAt,
			Records:          processedRecords,
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
