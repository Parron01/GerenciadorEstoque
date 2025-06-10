package repository

import (
	"database/sql"
	"fmt"

	"github.com/Parron01/GerenciadorEstoque/backendGo/internal/models"
)

type ProductRepository interface {
	GetAll(userID int) ([]models.Product, error)
	GetByID(id string, userID int) (*models.Product, error)
	Create(product *models.Product) error
	Update(product *models.Product) error
	Delete(id string, userID int) error
}

type productRepository struct {
	db             *sql.DB
	loteRepository LoteRepository // To fetch associated lotes
}

func NewProductRepository(db *sql.DB, loteRepo LoteRepository) ProductRepository {
	return &productRepository{db: db, loteRepository: loteRepo}
}

func (r *productRepository) GetAll(userID int) ([]models.Product, error) {
	rows, err := r.db.Query("SELECT id, name, unit, quantity FROM products WHERE user_id = $1", userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch products: %w", err)
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Unit, &product.Quantity); err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		// Fetch associated lotes
		lotes, err := r.loteRepository.GetByProductID(product.ID, userID)
		if err != nil {
			// Log error but continue, or return error based on policy
			fmt.Printf("Warning: failed to fetch lotes for product %s: %v\n", product.ID, err)
		}
		product.Lotes = lotes
		products = append(products, product)
	}
	return products, nil
}

func (r *productRepository) GetByID(id string, userID int) (*models.Product, error) {
	var product models.Product
	err := r.db.QueryRow("SELECT id, name, unit, quantity FROM products WHERE id = $1 AND user_id = $2", id, userID).
		Scan(&product.ID, &product.Name, &product.Unit, &product.Quantity)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Or a specific "not found" error
		}
		return nil, fmt.Errorf("product not found: %w", err)
	}

	// Fetch associated lotes
	lotes, err := r.loteRepository.GetByProductID(product.ID, userID)
	if err != nil {
		// Log error but continue, or return error based on policy
		fmt.Printf("Warning: failed to fetch lotes for product %s: %v\n", product.ID, err)
	}
	product.Lotes = lotes
	return &product, nil
}

func (r *productRepository) Create(product *models.Product) error {
	// Note: Product.Quantity will be updated by trigger if lotes are managed.
	// If creating a product without lotes, this quantity is the initial one.
	stmt, err := r.db.Prepare("INSERT INTO products (id, name, unit, quantity, user_id) VALUES ($1, $2, $3, $4, $5)")
	if err != nil {
		return fmt.Errorf("failed to prepare statement for create product: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(product.ID, product.Name, product.Unit, product.Quantity, product.UserID)
	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}
	return nil
}

func (r *productRepository) Update(product *models.Product) error {
	// Note: Product.Quantity will be updated by trigger if lotes are managed.
	// Updating product details other than quantity directly.
	stmt, err := r.db.Prepare("UPDATE products SET name = $1, unit = $2 WHERE id = $3 AND user_id = $4")
	// If quantity needs to be updatable here AND lots exist, logic is more complex.
	// For now, assuming trigger handles quantity based on lots.
	// If no lots, direct quantity update: "UPDATE products SET name = $1, unit = $2, quantity = $3 WHERE id = $4"
	if err != nil {
		return fmt.Errorf("failed to prepare statement for update product: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(product.Name, product.Unit, product.ID, product.UserID)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("product not found or no change made")
	}
	return nil
}

func (r *productRepository) Delete(id string, userID int) error {
	// Deleting a product will also delete its lotes due to ON DELETE CASCADE
	stmt, err := r.db.Prepare("DELETE FROM products WHERE id = $1 AND user_id = $2")
	if err != nil {
		return fmt.Errorf("failed to prepare statement for delete product: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("product not found")
	}
	return nil
}
