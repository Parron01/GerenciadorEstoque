package controllers

import (
	"net/http"

	"github.com/Parron01/GerenciadorEstoque/backendGo/internal/database"
	"github.com/Parron01/GerenciadorEstoque/backendGo/internal/models"
	"github.com/gin-gonic/gin"
)

// HistoryController handles history-related requests
type HistoryController struct{}

// NewHistoryController creates a new history controller
func NewHistoryController() *HistoryController {
    return &HistoryController{}
}

// GetAll gets all history records
func (hc *HistoryController) GetAll(c *gin.Context) {
    rows, err := database.DB.Query("SELECT id, date, changes FROM history ORDER BY date DESC")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch history"})
        return
    }
    defer rows.Close()

    var history []models.History
    for rows.Next() {
        var record models.History
        if err := rows.Scan(&record.ID, &record.Date, &record.Changes); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan history"})
            return
        }
        history = append(history, record)
    }

    c.JSON(http.StatusOK, history)
}

// Create adds a new history entry
func (hc *HistoryController) Create(c *gin.Context) {
    var historyEntry models.History
    if err := c.ShouldBindJSON(&historyEntry); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid history data"})
        return
    }

    // Convert changes to JSON string for storage
    stmt, err := database.DB.Prepare("INSERT INTO history (id, date, changes) VALUES ($1, $2, $3)")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare statement"})
        return
    }
    defer stmt.Close()

    _, err = stmt.Exec(historyEntry.ID, historyEntry.Date, historyEntry.Changes)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create history entry"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "message": "Histórico adicionado com sucesso",
        "id":      historyEntry.ID,
    })
}