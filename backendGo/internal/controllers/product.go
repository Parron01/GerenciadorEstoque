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
func (pc *ProductController) GetByID(c *gin.Context) {
	id := c.Param("id")

	product, err := pc.repo.GetByID(id)
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
	}

	createdProduct, fetchErr := pc.repo.GetByID(product.ID)
	if fetchErr != nil {
		c.JSON(http.StatusCreated, product)
		return
	}
	c.JSON(http.StatusCreated, createdProduct)
}

// Update modifies an existing product
func (pc *ProductController) Update(c *gin.Context) {
	id := c.Param("id")
	var productUpdates models.Product
	if err := c.ShouldBindJSON(&productUpdates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product data: " + err.Error()})
		return
	}

	productUpdates.ID = id

	existingProduct, err := pc.repo.GetByID(id)
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
			ProductID:     id,
			ProductName:   productUpdates.Name,
			Action:        "updated",
		}
		if err := pc.historySvc.RecordChange(service.EntityTypeProduct, id, changeDetail); err != nil {
		}
	}

	updatedProduct, fetchErr := pc.repo.GetByID(id)
    if fetchErr != nil {
        c.JSON(http.StatusOK, productUpdates) // Fallback
        return
    }
	c.JSON(http.StatusOK, updatedProduct)
}

// Delete removes a product
func (pc *ProductController) Delete(c *gin.Context) {
	id := c.Param("id")

	existingProduct, err := pc.repo.GetByID(id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check product existence: " + err.Error()})
        return
    }
    if existingProduct == nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
        return
    }

	err = pc.repo.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product: " + err.Error()})
		return
	}

	changeDetail := models.ProductChange{
		ProductID:        id,
		ProductName:      existingProduct.Name,
		Action:           "deleted",
		QuantityBefore:   existingProduct.Quantity,
		IsProductRemoval: true,
	}
	if err := pc.historySvc.RecordChange(service.EntityTypeProduct, id, changeDetail); err != nil {
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}