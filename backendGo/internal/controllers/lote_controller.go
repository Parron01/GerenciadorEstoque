package controllers

import (
	"fmt"
	"net/http"

	"github.com/Parron01/GerenciadorEstoque/backendGo/internal/models"
	"github.com/Parron01/GerenciadorEstoque/backendGo/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type LoteController struct {
	service   service.LoteService
	validator *validator.Validate
}

func NewLoteController(service service.LoteService) *LoteController {
	return &LoteController{
		service:   service,
		validator: validator.New(),
	}
}

// CreateLote godoc
// @Summary Create a new lote for a product
// @Description Adds a new lote to a specified product. The sum of lote quantities will update the product's total quantity.
// @Tags lotes
// @Accept json
// @Produce json
// @Param product_id path string true "Product ID"
// @Param lote body models.Lote true "Lote data (quantity, data_validade)"
// @HeaderParam X-Operation-Batch-ID header string false "Optional Batch ID for grouping operations"
// @Success 201 {object} models.Lote
// @Failure 400 {object} gin.H{"error": "message"}
// @Failure 404 {object} gin.H{"error": "message"} "Product not found"
// @Failure 500 {object} gin.H{"error": "message"}
// @Router /api/products/{product_id}/lotes [post]
// @Security BearerAuth
func (lc *LoteController) CreateLote(c *gin.Context) {
	productID := c.Param("product_id")
	var loteReq models.Lote
	operationBatchID := c.GetHeader("X-Operation-Batch-ID")

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	if err := c.ShouldBindJSON(&loteReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}
	
	// Validate specific fields if needed using validator
    err := lc.validator.StructExcept(loteReq, "ID", "ProductID", "CreatedAt", "UpdatedAt")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed: " + err.Error()})
        return
    }

	createdLote, err := lc.service.CreateLote(productID, loteReq, userID.(int), operationBatchID)
	if err != nil {
		// Basic error type checking, can be more granular
		if err.Error() == fmt.Sprintf("product with ID %s not found", productID) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if _, ok := err.(validator.ValidationErrors); ok {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Validation error: " + err.Error()})
        } else if err.Error() == "invalid data_validade format, expected YYYY-MM-DD" {
             c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        } else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create lote: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, createdLote)
}

// GetLotesForProduct godoc
// @Summary Get all lotes for a specific product
// @Description Retrieves a list of all lotes associated with a product ID.
// @Tags lotes
// @Produce json
// @Param product_id path string true "Product ID"
// @Success 200 {array} models.Lote
// @Failure 500 {object} gin.H{"error": "message"}
// @Router /api/products/{product_id}/lotes [get]
// @Security BearerAuth
func (lc *LoteController) GetLotesForProduct(c *gin.Context) {
	productID := c.Param("product_id")
	
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	lotes, err := lc.service.GetLotesByProductID(productID, userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve lotes: " + err.Error()})
		return
	}
	if lotes == nil {
		lotes = []models.Lote{} // Return empty array instead of null
	}
	c.JSON(http.StatusOK, lotes)
}

// UpdateLote godoc
// @Summary Update an existing lote
// @Description Updates the quantity or expiration date of a specific lote.
// @Tags lotes
// @Accept json
// @Produce json
// @Param lote_id path string true "Lote ID"
// @Param lote body models.Lote true "Lote data to update (quantity, data_validade)"
// @HeaderParam X-Operation-Batch-ID header string false "Optional Batch ID for grouping operations"
// @Success 200 {object} models.Lote
// @Failure 400 {object} gin.H{"error": "message"}
// @Failure 404 {object} gin.H{"error": "message"} "Lote not found"
// @Failure 500 {object} gin.H{"error": "message"}
// @Router /api/lotes/{lote_id} [put]
// @Security BearerAuth
func (lc *LoteController) UpdateLote(c *gin.Context) {
	loteID := c.Param("lote_id")
	var loteReq models.Lote
	operationBatchID := c.GetHeader("X-Operation-Batch-ID")

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	if err := c.ShouldBindJSON(&loteReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

    err := lc.validator.StructExcept(loteReq, "ID", "ProductID", "CreatedAt", "UpdatedAt")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed: " + err.Error()})
        return
    }

	updatedLote, err := lc.service.UpdateLote(loteID, loteReq, userID.(int), operationBatchID)
	if err != nil {
		if err.Error() == fmt.Sprintf("lote with ID %s not found", loteID) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if _, ok := err.(validator.ValidationErrors); ok {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Validation error: " + err.Error()})
        } else if err.Error() == "invalid data_validade format, expected YYYY-MM-DD" {
             c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        }else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update lote: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, updatedLote)
}

// DeleteLote godoc
// @Summary Delete a lote
// @Description Removes a lote by its ID.
// @Tags lotes
// @Produce json
// @Param lote_id path string true "Lote ID"
// @HeaderParam X-Operation-Batch-ID header string false "Optional Batch ID for grouping operations"
// @Success 200 {object} gin.H{"message": "Lote deleted successfully"}
// @Failure 404 {object} gin.H{"error": "message"} "Lote not found"
// @Failure 500 {object} gin.H{"error": "message"}
// @Router /api/lotes/{lote_id} [delete]
// @Security BearerAuth
func (lc *LoteController) DeleteLote(c *gin.Context) {
	loteID := c.Param("lote_id")
	operationBatchID := c.GetHeader("X-Operation-Batch-ID")

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	err := lc.service.DeleteLote(loteID, userID.(int), operationBatchID)
	if err != nil {
		if err.Error() == fmt.Sprintf("lote with ID %s not found", loteID) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete lote: " + err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Lote deleted successfully"})
}
