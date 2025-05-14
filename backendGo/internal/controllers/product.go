package controllers

import (
	"net/http"

	"github.com/Parron01/GerenciadorEstoque/backendGo/internal/database"
	"github.com/Parron01/GerenciadorEstoque/backendGo/internal/models"
	"github.com/gin-gonic/gin"
)

// ProductController handles product-related requests
type ProductController struct{}

// NewProductController creates a new product controller
func NewProductController() *ProductController {
    return &ProductController{}
}

// GetAll returns all products
func (pc *ProductController) GetAll(c *gin.Context) {
    rows, err := database.DB.Query("SELECT id, name, unit, quantity FROM products")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
        return
    }
    defer rows.Close()

    var products []models.Product
    for rows.Next() {
        var product models.Product
        if err := rows.Scan(&product.ID, &product.Name, &product.Unit, &product.Quantity); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan product"})
            return
        }
        products = append(products, product)
    }

    c.JSON(http.StatusOK, products)
}

// GetByID returns a specific product by ID
func (pc *ProductController) GetByID(c *gin.Context) {
    id := c.Param("id")

    var product models.Product
    err := database.DB.QueryRow("SELECT id, name, unit, quantity FROM products WHERE id = $1", id).
        Scan(&product.ID, &product.Name, &product.Unit, &product.Quantity)
    
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
        return
    }

    c.JSON(http.StatusOK, product)
}

// Create adds a new product
func (pc *ProductController) Create(c *gin.Context) {
    var product models.Product
    if err := c.ShouldBindJSON(&product); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product data"})
        return
    }

    stmt, err := database.DB.Prepare("INSERT INTO products (id, name, unit, quantity) VALUES ($1, $2, $3, $4)")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare statement"})
        return
    }
    defer stmt.Close()

    _, err = stmt.Exec(product.ID, product.Name, product.Unit, product.Quantity)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
        return
    }

    c.JSON(http.StatusCreated, product)
}

// Update modifies an existing product
func (pc *ProductController) Update(c *gin.Context) {
    id := c.Param("id")
    
    var product models.Product
    if err := c.ShouldBindJSON(&product); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product data"})
        return
    }

    // Ensure the ID in URL matches the product ID
    if id != product.ID {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID mismatch"})
        return
    }

    stmt, err := database.DB.Prepare("UPDATE products SET name = $1, unit = $2, quantity = $3 WHERE id = $4")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare statement"})
        return
    }
    defer stmt.Close()

    result, err := stmt.Exec(product.Name, product.Unit, product.Quantity, id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
        return
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil || rowsAffected == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
        return
    }

    c.JSON(http.StatusOK, product)
}

// Delete removes a product
func (pc *ProductController) Delete(c *gin.Context) {
    id := c.Param("id")

    stmt, err := database.DB.Prepare("DELETE FROM products WHERE id = $1")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare statement"})
        return
    }
    defer stmt.Close()

    result, err := stmt.Exec(id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
        return
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil || rowsAffected == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}