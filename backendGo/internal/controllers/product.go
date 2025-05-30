package controllers

import (
	"log"
	"net/http"

	"github.com/Parron01/GerenciadorEstoque/backendGo/internal/models"
	"github.com/Parron01/GerenciadorEstoque/backendGo/internal/repository"
	"github.com/Parron01/GerenciadorEstoque/backendGo/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ProductController handles product-related requests
type ProductController struct {
	repo       repository.ProductRepository
	historySvc service.HistoryService
}

// NewProductController creates a new product controller
func NewProductController(repo repository.ProductRepository, historySvc service.HistoryService) *ProductController {
	return &ProductController{repo: repo, historySvc: historySvc}
}

// GetAll returns all products
// @Summary Get all products
// @Description Retrieves a list of all products in the system.
// @Tags products
// @Produce json
// @Success 200 {array} models.Product
// @Failure 500 {object} gin.H{"error": "message"}
// @Router /api/products [get]
// @Security BearerAuth
func (pc *ProductController) GetAll(c *gin.Context) {
	products, err := pc.repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products: " + err.Error()})
		return
	}
	if products == nil {
		products = []models.Product{}
	}
	c.JSON(http.StatusOK, products)
}

// GetByID returns a specific product by ID
// @Summary Get a product by ID
// @Description Retrieves a product by its ID, including its lotes.
// @Tags products
// @Produce json
// @Param product_id path string true "Product ID"
// @Success 200 {object} models.Product
// @Failure 404 {object} gin.H{"error": "message"} "Product not found"
// @Failure 500 {object} gin.H{"error": "message"}
// @Router /api/products/{product_id} [get]
// @Security BearerAuth
func (pc *ProductController) GetByID(c *gin.Context) {
	productID := c.Param("product_id") // Changed from "id"

	product, err := pc.repo.GetByID(productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product by ID: " + err.Error()})
		return
	}
	if product == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// Create adds a new product
// @Summary Create a new product
// @Description Adds a new product to the system.
// @Tags products
// @Accept json
// @Produce json
// @Param product body models.Product true "Product data (ID is optional, will be generated if empty; Quantity is initial, will be managed by lotes if lotes are added)"
// @HeaderParam X-Operation-Batch-ID header string false "Optional Batch ID for grouping operations"
// @Success 201 {object} models.Product
// @Failure 400 {object} gin.H{"error": "message"}
// @Failure 500 {object} gin.H{"error": "message"}
// @Router /api/products [post]
// @Security BearerAuth
func (pc *ProductController) Create(c *gin.Context) {
	var product models.Product
	operationBatchID := c.GetHeader("X-Operation-Batch-ID")

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product data: " + err.Error()})
		return
	}

	if product.ID == "" {
		product.ID = uuid.NewString()
	}

	err := pc.repo.Create(&product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product: " + err.Error()})
		return
	}

	// Correctly assign pointer for QuantityAfter
	qtyAfter := product.Quantity 
	changeDetail := models.ProductChange{
		ProductID:      product.ID,
		ProductName:    product.Name,
		Action:         "created",
		QuantityAfter:  &qtyAfter, // Assign address of qtyAfter
		IsNewProduct:   true,
	}
	if err := pc.historySvc.RecordChange(service.EntityTypeProduct, product.ID, changeDetail, operationBatchID); err != nil {
		// Log or handle history recording error, but don't fail the main operation
	}

	createdProduct, fetchErr := pc.repo.GetByID(product.ID)
	if fetchErr != nil {
		c.JSON(http.StatusCreated, product)
		return
	}
	c.JSON(http.StatusCreated, createdProduct)
}

// Update modifies an existing product
// @Summary Update an existing product
// @Description Updates the details of an existing product (name, unit). Quantity is managed by lotes.
// @Tags products
// @Accept json
// @Produce json
// @Param product_id path string true "Product ID"
// @Param product body models.Product true "Product data to update (name, unit)"
// @HeaderParam X-Operation-Batch-ID header string false "Optional Batch ID for grouping operations"
// @Success 200 {object} models.Product
// @Failure 400 {object} gin.H{"error": "message"}
// @Failure 404 {object} gin.H{"error": "message"} "Product not found"
// @Failure 500 {object} gin.H{"error": "message"}
// @Router /api/products/{product_id} [put]
// @Security BearerAuth
func (pc *ProductController) Update(c *gin.Context) {
	productID := c.Param("product_id") // Changed from "id"
	// Use a struct with pointers to distinguish between omitted fields and empty strings
	var payload struct {
		Name *string `json:"name"`
		Unit *string `json:"unit"`
		// Quantity is not updated here, it's managed by lotes or a separate mechanism
	}
	operationBatchID := c.GetHeader("X-Operation-Batch-ID")

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product data: " + err.Error()})
		return
	}

	existingProduct, err := pc.repo.GetByID(productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch existing product: " + err.Error()})
		return
	}
	if existingProduct == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	productToUpdate := *existingProduct // Start with existing data
	madeChanges := false

	var changedFields []models.ChangedField

	if payload.Name != nil {
		if *payload.Name == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Product name cannot be empty"})
            return
        }
		if productToUpdate.Name != *payload.Name {
			changedFields = append(changedFields, models.ChangedField{Field: "name", OldValue: productToUpdate.Name, NewValue: *payload.Name})
			productToUpdate.Name = *payload.Name
			madeChanges = true
		}
	}

	if payload.Unit != nil {
		if *payload.Unit != "L" && *payload.Unit != "kg" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid unit value. Must be 'L' or 'kg'."})
			return
		}
		if productToUpdate.Unit != *payload.Unit {
			changedFields = append(changedFields, models.ChangedField{Field: "unit", OldValue: productToUpdate.Unit, NewValue: *payload.Unit})
			productToUpdate.Unit = *payload.Unit
			madeChanges = true
		}
	}

	if !madeChanges {
		c.JSON(http.StatusOK, existingProduct) // No changes, return existing
		return
	}
	
	// Ensure productToUpdate.Quantity is not inadvertently changed by this endpoint
	// It should retain existingProduct.Quantity as this endpoint only handles name/unit.
	// The repository's Update method should be specific about which fields it updates.
	// For safety, explicitly set it if there's any doubt:
	productToUpdate.Quantity = existingProduct.Quantity


	err = pc.repo.Update(&productToUpdate) // Pass the selectively updated product
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product: " + err.Error()})
		return
	}

	if len(changedFields) > 0 {
		changeDetail := models.ProductChange{
			ProductID:     productID,
			ProductName:   productToUpdate.Name, // Use the final name for context
			Action:        "product_details_updated",
			ChangedFields: changedFields,
		}
		if histErr := pc.historySvc.RecordChange(service.EntityTypeProduct, productID, changeDetail, operationBatchID); histErr != nil {
			log.Printf("WARN: Failed to record history for product update %s: %v", productID, histErr)
		}
	}

	finalUpdatedProduct, fetchErr := pc.repo.GetByID(productID)
	if fetchErr != nil {
		log.Printf("WARN: Failed to fetch product %s after update, returning potentially stale data: %v", productID, fetchErr)
		c.JSON(http.StatusOK, productToUpdate) // Fallback
		return
	}
	c.JSON(http.StatusOK, finalUpdatedProduct)
}

// Delete removes a product
// @Summary Delete a product
// @Description Removes a product and its associated lotes from the system.
// @Tags products
// @Produce json
// @Param product_id path string true "Product ID"
// @HeaderParam X-Operation-Batch-ID header string false "Optional Batch ID for grouping operations"
// @Success 200 {object} gin.H{"message": "Product deleted successfully"}
// @Failure 404 {object} gin.H{"error": "message"} "Product not found"
// @Failure 500 {object} gin.H{"error": "message"}
// @Router /api/products/{product_id} [delete]
// @Security BearerAuth
func (pc *ProductController) Delete(c *gin.Context) {
	productID := c.Param("product_id") // Changed from "id"
	operationBatchID := c.GetHeader("X-Operation-Batch-ID")

	existingProduct, err := pc.repo.GetByID(productID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check product existence: " + err.Error()})
        return
    }
    if existingProduct == nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
        return
    }

	err = pc.repo.Delete(productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product: " + err.Error()})
		return
	}

	// Correctly assign pointer for QuantityBefore
	qtyBefore := existingProduct.Quantity
	changeDetail := models.ProductChange{
		ProductID:        productID,
		ProductName:      existingProduct.Name,
		Action:           "deleted",
		QuantityBefore:   &qtyBefore, // Assign address of qtyBefore
		IsProductRemoval: true,
	}
	if err := pc.historySvc.RecordChange(service.EntityTypeProduct, productID, changeDetail, operationBatchID); err != nil {
		// Log or handle history recording error
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}