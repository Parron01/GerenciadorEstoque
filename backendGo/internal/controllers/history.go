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

	limit, err := strconv.Atoi(limitQuery)
	if err != nil || limit <= 0 {
		limit = 20
	}
	offset, err := strconv.Atoi(offsetQuery)
	if err != nil || offset < 0 {
		offset = 0
	}

	historyEntries, err := hc.service.GetHistory(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch history: " + err.Error()})
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

	if historyEntry.EntityType == "" || historyEntry.EntityID == "" {
		// This check depends on how strictly you want to enforce these fields for client-submitted history.
		// For now, we allow it, aligning with the previous direct DB insert that didn't require them.
		// However, for better data quality, requiring them is advisable.
		// c.JSON(http.StatusBadRequest, gin.H{"error": "EntityType and EntityID are required for history entries"})
		// return
	}

	err := hc.service.CreateRawHistoryEntry(historyEntry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create history entry: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Histórico adicionado com sucesso",
		"id":      historyEntry.ID, // ID is now set by service/repo if not provided
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