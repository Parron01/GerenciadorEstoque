package controllers

import (
	"net/http"
	"strconv"

	"github.com/Parron01/GerenciadorEstoque/backendGo/internal/models"
	"github.com/Parron01/GerenciadorEstoque/backendGo/internal/service"
	"github.com/gin-gonic/gin"
)

// HistoryController handles history-related requests
type HistoryController struct {
	service service.HistoryService
}

// NewHistoryController creates a new history controller
func NewHistoryController(service service.HistoryService) *HistoryController {
	return &HistoryController{service: service}
}

// GetAll gets all history records with pagination
func (hc *HistoryController) GetAll(c *gin.Context) {
	limitQuery := c.DefaultQuery("limit", "20")
	offsetQuery := c.DefaultQuery("offset", "0")
	batchID := c.Query("batch_id") // Added batch_id query parameter

	limit, err := strconv.Atoi(limitQuery)
	if err != nil || limit <= 0 {
		limit = 20
	}
	offset, err := strconv.Atoi(offsetQuery)
	if err != nil || offset < 0 {
		offset = 0
	}

	var historyEntries []models.History
	var fetchErr error

	// If batch_id is provided, get by batch_id instead
	if batchID != "" {
		historyEntries, fetchErr = hc.service.GetByBatchID(batchID)
	} else {
		historyEntries, fetchErr = hc.service.GetHistory(limit, offset)
	}

	if fetchErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch history: " + fetchErr.Error()})
		return
	}
	if historyEntries == nil {
		historyEntries = []models.History{}
	}

	c.JSON(http.StatusOK, historyEntries)
}

// Create adds a new history entry (potentially from client, if design allows)
func (hc *HistoryController) Create(c *gin.Context) {
	var historyEntry models.History
	if err := c.ShouldBindJSON(&historyEntry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid history data: " + err.Error()})
		return
	}

	err := hc.service.CreateRawHistoryEntry(historyEntry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create history entry: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Histórico adicionado com sucesso",
		"id":      historyEntry.ID,
		"batch_id": historyEntry.BatchID,
	})
}

// GetHistoryForEntity retrieves history for a specific entity (product or lote)
// @Summary Get history for a specific entity
// @Description Retrieves all history records for a given entity type and ID.
// @Tags history
// @Produce json
// @Param entity_type path string true "Entity Type (e.g., product, lote)"
// @Param entity_id path string true "Entity ID"
// @Success 200 {array} models.History
// @Failure 400 {object} gin.H{"error": "message"} "Invalid entity type"
// @Failure 500 {object} gin.H{"error": "message"}
// @Router /api/history/{entity_type}/{entity_id} [get]
// @Security BearerAuth
func (hc *HistoryController) GetHistoryForEntity(c *gin.Context) {
	entityType := c.Param("entity_type")
	entityID := c.Param("entity_id")

	if entityType != service.EntityTypeProduct && entityType != service.EntityTypeLote {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity_type. Must be 'product' or 'lote'."})
		return
	}

	historyEntries, err := hc.service.GetHistoryForEntity(entityType, entityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve history for entity: " + err.Error()})
		return
	}
	if historyEntries == nil {
		historyEntries = []models.History{}
	}
	c.JSON(http.StatusOK, historyEntries)
}

// New batch operations

// CreateBatch handles the creation of multiple history entries in a single batch
// @Summary Create multiple history entries in one batch
// @Description Creates multiple history entries with a shared batch ID
// @Tags history
// @Accept json
// @Produce json
// @Param entries body []models.History true "Array of history entries (BatchID will be set automatically)"
// @Success 201 {object} gin.H{"message": "string", "batch_id": "string", "count": int}
// @Failure 400 {object} gin.H{"error": "string"}
// @Failure 500 {object} gin.H{"error": "string"}
// @Router /api/history/batch [post]
// @Security BearerAuth
func (hc *HistoryController) CreateBatch(c *gin.Context) {
	var historyEntries []models.History
	if err := c.ShouldBindJSON(&historyEntries); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid batch history data: " + err.Error()})
		return
	}

	if len(historyEntries) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Empty batch - no history entries provided"})
		return
	}

	batchID, err := hc.service.CreateBatch(historyEntries)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create history batch: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "History batch created successfully",
		"batch_id": batchID,
		"count": len(historyEntries),
	})
}

// GetByBatch retrieves all history entries for a specific batch ID
// @Summary Get history entries by batch ID
// @Description Retrieves all history entries belonging to a specific batch
// @Tags history
// @Produce json
// @Param batch_id path string true "Batch ID"
// @Success 200 {array} models.History
// @Failure 400 {object} gin.H{"error": "string"}
// @Failure 500 {object} gin.H{"error": "string"}
// @Router /api/history/batch/{batch_id} [get]
// @Security BearerAuth
func (hc *HistoryController) GetByBatch(c *gin.Context) {
	batchID := c.Param("batch_id")
	if batchID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Batch ID is required"})
		return
	}

	historyEntries, err := hc.service.GetByBatchID(batchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch history by batch ID: " + err.Error()})
		return
	}
	if historyEntries == nil {
		historyEntries = []models.History{}
	}

	c.JSON(http.StatusOK, historyEntries)
}

// GetGrouped retrieves history entries grouped by batch ID, with pagination for batches.
// @Summary Get history grouped by batch ID
// @Description Retrieves history entries, grouped by their batch ID, supporting pagination for the batches.
// @Tags history
// @Produce json
// @Param page query int false "Page number for batch pagination" default(1)
// @Param pageSize query int false "Number of batches per page" default(10)
// @Success 200 {object} models.PaginatedHistoryBatchGroups
// @Failure 400 {object} gin.H{"error": "string"} "Invalid query parameters"
// @Failure 500 {object} gin.H{"error": "string"}
// @Router /api/history/grouped [get]
// @Security BearerAuth
func (hc *HistoryController) GetGrouped(c *gin.Context) {
	pageQuery := c.DefaultQuery("page", "1")
	pageSizeQuery := c.DefaultQuery("pageSize", "10")

	page, err := strconv.Atoi(pageQuery)
	if err != nil || page <= 0 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeQuery)
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}

	paginatedGroups, err := hc.service.GetGroupedHistory(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch grouped history: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, paginatedGroups)
}

// CreateProductBatchContext godoc
// @Summary Create a product batch context history entry
// @Description Records a snapshot of a product's state (name, quantity before/after) for a given batch of operations.
// @Tags history
// @Accept json
// @Produce json
// @Param context_payload body models.ProductBatchContextChangeDetail true "Product batch context data"
// @HeaderParam X-Operation-Batch-ID header string true "Batch ID for grouping this context with other operations"
// @Success 201 {object} gin.H{"message": "Product batch context recorded"}
// @Failure 400 {object} gin.H{"error": "message"}
// @Failure 500 {object} gin.H{"error": "message"}
// @Router /api/history/product-context [post]
// @Security BearerAuth
func (hc *HistoryController) CreateProductBatchContext(c *gin.Context) {
	var payload models.ProductBatchContextChangeDetail
	operationBatchID := c.GetHeader("X-Operation-Batch-ID")

	if operationBatchID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-Operation-Batch-ID header is required"})
		return
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload: " + err.Error()})
		return
	}

	if payload.ProductID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "productId is required in payload"})
		return
	}

	// Use a distinct EntityType for these records
	err := hc.service.RecordChange(service.EntityTypeProductBatchContext, payload.ProductID, payload, operationBatchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record product batch context: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Product batch context recorded successfully"})
}