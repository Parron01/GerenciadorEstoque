package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Parron01/GerenciadorEstoque/backendGo/internal/models"
	"github.com/google/uuid"
)

type LoteRepository interface {
	Create(tx *sql.Tx, lote *models.Lote) error
	GetByID(id string, userID int) (*models.Lote, error)
	GetByProductID(productID string, userID int) ([]models.Lote, error)
	Update(tx *sql.Tx, lote *models.Lote) error
	Delete(tx *sql.Tx, id string, userID int) error
	CountByProductID(productID string, userID int) (int, error)
}

type loteRepository struct {
	db *sql.DB
}

func NewLoteRepository(db *sql.DB) LoteRepository {
	return &loteRepository{db: db}
}

func (r *loteRepository) Create(tx *sql.Tx, lote *models.Lote) error {
	lote.ID = uuid.NewString()
	lote.CreatedAt = time.Now()
	lote.UpdatedAt = time.Now()

	query := `INSERT INTO product_lots (id, product_id, user_id, quantity, data_validade, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7)`
	
	var err error
	if tx != nil {
		_, err = tx.Exec(query, lote.ID, lote.ProductID, lote.UserID, lote.Quantity, lote.DataValidade, lote.CreatedAt, lote.UpdatedAt)
	} else {
		_, err = r.db.Exec(query, lote.ID, lote.ProductID, lote.UserID, lote.Quantity, lote.DataValidade, lote.CreatedAt, lote.UpdatedAt)
	}

	if err != nil {
		return fmt.Errorf("failed to create lote: %w", err)
	}
	return nil
}

func (r *loteRepository) GetByID(id string, userID int) (*models.Lote, error) {
	lote := &models.Lote{}
	query := `SELECT id, product_id, user_id, quantity, data_validade, created_at, updated_at 
              FROM product_lots WHERE id = $1 AND user_id = $2`
	err := r.db.QueryRow(query, id, userID).Scan(&lote.ID, &lote.ProductID, &lote.UserID, &lote.Quantity, &lote.DataValidade, &lote.CreatedAt, &lote.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Or a specific "not found" error
		}
		return nil, fmt.Errorf("failed to get lote by id: %w", err)
	}
	return lote, nil
}

func (r *loteRepository) GetByProductID(productID string, userID int) ([]models.Lote, error) {
	rows, err := r.db.Query(`SELECT id, product_id, user_id, quantity, data_validade, created_at, updated_at 
                             FROM product_lots WHERE product_id = $1 AND user_id = $2 ORDER BY data_validade ASC`, productID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get lotes by product id: %w", err)
	}
	defer rows.Close()

	var lotes []models.Lote
	for rows.Next() {
		var lote models.Lote
		if err := rows.Scan(&lote.ID, &lote.ProductID, &lote.UserID, &lote.Quantity, &lote.DataValidade, &lote.CreatedAt, &lote.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan lote: %w", err)
		}
		lotes = append(lotes, lote)
	}
	return lotes, nil
}

func (r *loteRepository) Update(tx *sql.Tx, lote *models.Lote) error {
	lote.UpdatedAt = time.Now()
	query := `UPDATE product_lots SET quantity = $1, data_validade = $2, updated_at = $3
              WHERE id = $4 AND product_id = $5`
	
	var result sql.Result
	var err error

	if tx != nil {
		result, err = tx.Exec(query, lote.Quantity, lote.DataValidade, lote.UpdatedAt, lote.ID, lote.ProductID)
	} else {
		result, err = r.db.Exec(query, lote.Quantity, lote.DataValidade, lote.UpdatedAt, lote.ID, lote.ProductID)
	}

	if err != nil {
		return fmt.Errorf("failed to update lote: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("lote not found or no change made") // Or a specific error
	}
	return nil
}

func (r *loteRepository) Delete(tx *sql.Tx, id string, userID int) error {
	query := `DELETE FROM product_lots WHERE id = $1 AND user_id = $2`
	var result sql.Result
	var err error

	if tx != nil {
		result, err = tx.Exec(query, id, userID)
	} else {
		result, err = r.db.Exec(query, id, userID)
	}

	if err != nil {
		return fmt.Errorf("failed to delete lote: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("lote not found") // Or a specific error
	}
	return nil
}

func (r *loteRepository) CountByProductID(productID string, userID int) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM product_lots WHERE product_id = $1 AND user_id = $2`
	err := r.db.QueryRow(query, productID, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count lotes by product id: %w", err)
	}
	return count, nil
}

// GetDB returns the underlying sql.DB instance for transaction management
func (r *loteRepository) GetDB() *sql.DB {
    return r.db
}

// Ensure loteRepository implements LoteRepository and has GetDB
var _ LoteRepository = &loteRepository{}

func GetLoteRepositoryDB(repo LoteRepository) *sql.DB {
	if r, ok := repo.(*loteRepository); ok {
		return r.GetDB()
	}
	return nil // Or panic, or handle error
}
