package controllers

import (
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
// @Success 201 {object} models.Product
// @Failure 400 {object} gin.H{"error": "message"}
// @Failure 500 {object} gin.H{"error": "message"}
// @Router /api/products [post]
// @Security BearerAuth
func (pc *ProductController) Create(c *gin.Context) {
	var product models.Product
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

	changeDetail := models.ProductChange{
		ProductID:      product.ID,
		ProductName:    product.Name,
		Action:         "created",
		QuantityAfter:  product.Quantity,
		IsNewProduct:   true,
	}
	if err := pc.historySvc.RecordChange(service.EntityTypeProduct, product.ID, changeDetail); err != nil {
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
// @Success 200 {object} models.Product
// @Failure 400 {object} gin.H{"error": "message"}
// @Failure 404 {object} gin.H{"error": "message"} "Product not found"
// @Failure 500 {object} gin.H{"error": "message"}
// @Router /api/products/{product_id} [put]
// @Security BearerAuth
func (pc *ProductController) Update(c *gin.Context) {
	productID := c.Param("product_id") // Changed from "id"
	var productUpdates models.Product
	if err := c.ShouldBindJSON(&productUpdates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product data: " + err.Error()})
		return
	}

	productUpdates.ID = productID

	existingProduct, err := pc.repo.GetByID(productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch existing product: " + err.Error()})
		return
	}
	if existingProduct == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	originalName := existingProduct.Name
	originalUnit := existingProduct.Unit

	err = pc.repo.Update(&productUpdates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product: " + err.Error()})
		return
	}

	if originalName != productUpdates.Name || originalUnit != productUpdates.Unit {
		changeDetail := models.ProductChange{
			ProductID:     productID,
			ProductName:   productUpdates.Name,
			Action:        "updated",
			// Add more details if needed, e.g., old name/unit
		}
		if err := pc.historySvc.RecordChange(service.EntityTypeProduct, productID, changeDetail); err != nil {
			// Log or handle history recording error
		}
	}

	updatedProduct, fetchErr := pc.repo.GetByID(productID)
    if fetchErr != nil {
        c.JSON(http.StatusOK, productUpdates) // Fallback
        return
    }
	c.JSON(http.StatusOK, updatedProduct)
}

// Delete removes a product
// @Summary Delete a product
// @Description Removes a product and its associated lotes from the system.
// @Tags products
// @Produce json
// @Param product_id path string true "Product ID"
// @Success 200 {object} gin.H{"message": "Product deleted successfully"}
// @Failure 404 {object} gin.H{"error": "message"} "Product not found"
// @Failure 500 {object} gin.H{"error": "message"}
// @Router /api/products/{product_id} [delete]
// @Security BearerAuth
func (pc *ProductController) Delete(c *gin.Context) {
	productID := c.Param("product_id") // Changed from "id"

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

	changeDetail := models.ProductChange{
		ProductID:        productID,
		ProductName:      existingProduct.Name,
		Action:           "deleted",
		QuantityBefore:   existingProduct.Quantity, // Quantity before deletion (sum of its lotes)
		IsProductRemoval: true,
	}
	if err := pc.historySvc.RecordChange(service.EntityTypeProduct, productID, changeDetail); err != nil {
		// Log or handle history recording error
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}