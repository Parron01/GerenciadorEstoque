package service

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Parron01/GerenciadorEstoque/backendGo/internal/models"
	"github.com/Parron01/GerenciadorEstoque/backendGo/internal/repository"
)

type LoteService interface {
	CreateLote(productID string, loteReq models.Lote, userID int, operationBatchID string) (*models.Lote, error)
	GetLotesByProductID(productID string, userID int) ([]models.Lote, error)
	GetLoteByID(loteID string, userID int) (*models.Lote, error)
	UpdateLote(loteID string, loteReq models.Lote, userID int, operationBatchID string) (*models.Lote, error)
	DeleteLote(loteID string, userID int, operationBatchID string) error
}

type loteService struct {
	loteRepo    repository.LoteRepository
	productRepo repository.ProductRepository // To check if product exists
	historySvc  HistoryService
	db          *sql.DB // For transactions
}

func NewLoteService(loteRepo repository.LoteRepository, productRepo repository.ProductRepository, historySvc HistoryService, db *sql.DB) LoteService {
	return &loteService{
		loteRepo:    loteRepo,
		productRepo: productRepo,
		historySvc:  historySvc,
		db:          db,
	}
}

func (s *loteService) CreateLote(productID string, loteReq models.Lote, userID int, operationBatchID string) (*models.Lote, error) {
	// Check if product exists
	product, err := s.productRepo.GetByID(productID, userID)
	if err != nil {
		return nil, fmt.Errorf("error checking product existence: %w", err)
	}
	if product == nil {
		return nil, fmt.Errorf("product with ID %s not found", productID)
	}

	// Validate DataValidade format (YYYY-MM-DD)
	if _, err := time.Parse("2006-01-02", loteReq.DataValidade); err != nil {
		return nil, fmt.Errorf("invalid data_validade format, expected YYYY-MM-DD: %w", err)
	}

	newLote := models.Lote{
		ProductID:    productID,
		UserID:       userID,
		Quantity:     loteReq.Quantity,
		DataValidade: loteReq.DataValidade,
	}

	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback if not committed

	if err := s.loteRepo.Create(tx, &newLote); err != nil {
		return nil, fmt.Errorf("failed to create lote in repository: %w", err)
	}

	// Record history with operationBatchID if provided, otherwise a new batchID (its own ID) will be used by historySvc or repo.
	changeDetail := models.LoteChangeDetail{
		LoteID:       newLote.ID,
		ProductID:    productID,
		Action:       "created",
		QuantityAfter: &newLote.Quantity,
		DataValidade: &newLote.DataValidade,
	}
	if err := s.historySvc.RecordChange(EntityTypeLote, newLote.ID, changeDetail, userID, operationBatchID); err != nil {
		// Log error, but don't fail the primary operation for history recording failure
		fmt.Printf("Warning: failed to record history for lote creation %s: %v\n", newLote.ID, err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &newLote, nil
}

func (s *loteService) GetLotesByProductID(productID string, userID int) ([]models.Lote, error) {
	return s.loteRepo.GetByProductID(productID, userID)
}

func (s *loteService) GetLoteByID(loteID string, userID int) (*models.Lote, error) {
	lote, err := s.loteRepo.GetByID(loteID, userID)
    if err != nil {
        return nil, fmt.Errorf("failed to get lote by ID from repository: %w", err)
    }
    if lote == nil {
        return nil, nil // Not found
    }
    return lote, nil
}

func (s *loteService) UpdateLote(loteID string, loteReq models.Lote, userID int, operationBatchID string) (*models.Lote, error) {
	existingLote, err := s.loteRepo.GetByID(loteID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch existing lote: %w", err)
	}
	if existingLote == nil {
		return nil, fmt.Errorf("lote with ID %s not found", loteID)
	}

    // Validate DataValidade format (YYYY-MM-DD)
	if _, err := time.Parse("2006-01-02", loteReq.DataValidade); err != nil {
		return nil, fmt.Errorf("invalid data_validade format, expected YYYY-MM-DD: %w", err)
	}

	originalQuantity := existingLote.Quantity
	originalDataValidade := existingLote.DataValidade

	existingLote.Quantity = loteReq.Quantity
	existingLote.DataValidade = loteReq.DataValidade
	// ProductID should not change during an update of a lote

	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if err := s.loteRepo.Update(tx, existingLote); err != nil {
		return nil, fmt.Errorf("failed to update lote in repository: %w", err)
	}

	// Record history with operationBatchID
	changeDetail := models.LoteChangeDetail{
		LoteID:          loteID,
		ProductID:       existingLote.ProductID,
		Action:          "updated",
		QuantityBefore:  &originalQuantity,
		QuantityAfter:   &existingLote.Quantity,
		DataValidadeOld: &originalDataValidade,
		DataValidadeNew: &existingLote.DataValidade,
	}
	if existingLote.Quantity != originalQuantity {
		qtyChanged := existingLote.Quantity - originalQuantity
		changeDetail.QuantityChanged = &qtyChanged
	}
	if err := s.historySvc.RecordChange(EntityTypeLote, loteID, changeDetail, userID, operationBatchID); err != nil {
		fmt.Printf("Warning: failed to record history for lote update %s: %v\n", loteID, err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return existingLote, nil
}

func (s *loteService) DeleteLote(loteID string, userID int, operationBatchID string) error {
	existingLote, err := s.loteRepo.GetByID(loteID, userID)
	if err != nil {
		return fmt.Errorf("failed to fetch lote for deletion: %w", err)
	}
	if existingLote == nil {
		return fmt.Errorf("lote with ID %s not found", loteID)
	}

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if err := s.loteRepo.Delete(tx, loteID, userID); err != nil {
		return fmt.Errorf("failed to delete lote in repository: %w", err)
	}

	// Record history with operationBatchID
	changeDetail := models.LoteChangeDetail{
		LoteID:         loteID,
		ProductID:      existingLote.ProductID,
		Action:         "deleted",
		QuantityBefore: &existingLote.Quantity,
		DataValidade:   &existingLote.DataValidade,
	}
	if err := s.historySvc.RecordChange(EntityTypeLote, loteID, changeDetail, userID, operationBatchID); err != nil {
		fmt.Printf("Warning: failed to record history for lote deletion %s: %v\n", loteID, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}
