package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Parron01/GerenciadorEstoque/backendGo/internal/models"
	"github.com/google/uuid"
)

type HistoryRepository interface {
	Create(entry *models.History) error
	CreateBatch(entries []models.History) error // New: Create multiple entries in a transaction
	GetAll(limit, offset int) ([]models.History, error)
	GetByEntityID(entityType, entityID string) ([]models.History, error)
	GetByBatchID(batchID string) ([]models.History, error) // New: Get by batchID
}

type historyRepository struct {
	db *sql.DB
}

func NewHistoryRepository(db *sql.DB) HistoryRepository {
	return &historyRepository{db: db}
}

func (r *historyRepository) Create(entry *models.History) error {
	if entry.ID == "" {
		entry.ID = uuid.NewString()
	}
	if entry.Date == "" {
		entry.Date = time.Now().Format(time.RFC3339) // Consistent date format
	}
	if entry.BatchID == "" {
		entry.BatchID = entry.ID // Default: use ID as batchID if not provided
	}

	// Ensure changes is valid JSON
	var js json.RawMessage = entry.Changes
	if !json.Valid(js) {
		// Attempt to marshal if it's not already raw JSON string
		b, err := json.Marshal(entry.Changes)
		if err != nil {
			return fmt.Errorf("history changes is not valid JSON and failed to marshal: %w", err)
		}
		js = json.RawMessage(b)
	}

	query := `INSERT INTO history (id, date, entity_type, entity_id, changes, batch_id) 
              VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Exec(query, entry.ID, entry.Date, entry.EntityType, entry.EntityID, js, entry.BatchID)
	if err != nil {
		return fmt.Errorf("failed to create history entry: %w", err)
	}
	return nil
}

func (r *historyRepository) CreateBatch(entries []models.History) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Will be ignored if tx.Commit() is called

	for _, entry := range entries {
		if entry.ID == "" {
			entry.ID = uuid.NewString()
		}
		if entry.Date == "" {
			entry.Date = time.Now().Format(time.RFC3339)
		}

		// Ensure changes is valid JSON
		var js json.RawMessage = entry.Changes
		if !json.Valid(js) {
			b, err := json.Marshal(entry.Changes)
			if err != nil {
				return fmt.Errorf("history changes is not valid JSON and failed to marshal: %w", err)
			}
			js = json.RawMessage(b)
		}

		query := `INSERT INTO history (id, date, entity_type, entity_id, changes, batch_id) 
                  VALUES ($1, $2, $3, $4, $5, $6)`
		_, err := tx.Exec(query, entry.ID, entry.Date, entry.EntityType, entry.EntityID, js, entry.BatchID)
		if err != nil {
			return fmt.Errorf("failed to create history entry in batch: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit batch transaction: %w", err)
	}
	return nil
}

func (r *historyRepository) GetAll(limit, offset int) ([]models.History, error) {
	query := `SELECT id, date, entity_type, entity_id, changes, batch_id 
              FROM history ORDER BY date DESC LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all history: %w", err)
	}
	defer rows.Close()

	var historyEntries []models.History
	for rows.Next() {
		var entry models.History
		if err := rows.Scan(&entry.ID, &entry.Date, &entry.EntityType, &entry.EntityID, &entry.Changes, &entry.BatchID); err != nil {
			return nil, fmt.Errorf("failed to scan history entry: %w", err)
		}
		historyEntries = append(historyEntries, entry)
	}
	return historyEntries, nil
}

func (r *historyRepository) GetByEntityID(entityType, entityID string) ([]models.History, error) {
	query := `SELECT id, date, entity_type, entity_id, changes, batch_id 
              FROM history WHERE entity_type = $1 AND entity_id = $2 ORDER BY date DESC`
	rows, err := r.db.Query(query, entityType, entityID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch history by entity id: %w", err)
	}
	defer rows.Close()

	var historyEntries []models.History
	for rows.Next() {
		var entry models.History
		if err := rows.Scan(&entry.ID, &entry.Date, &entry.EntityType, &entry.EntityID, &entry.Changes, &entry.BatchID); err != nil {
			return nil, fmt.Errorf("failed to scan history entry: %w", err)
		}
		historyEntries = append(historyEntries, entry)
	}
	return historyEntries, nil
}

func (r *historyRepository) GetByBatchID(batchID string) ([]models.History, error) {
	query := `SELECT id, date, entity_type, entity_id, changes, batch_id 
              FROM history WHERE batch_id = $1 ORDER BY date DESC`
	rows, err := r.db.Query(query, batchID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch history by batch id: %w", err)
	}
	defer rows.Close()

	var historyEntries []models.History
	for rows.Next() {
		var entry models.History
		if err := rows.Scan(&entry.ID, &entry.Date, &entry.EntityType, &entry.EntityID, &entry.Changes, &entry.BatchID); err != nil {
			return nil, fmt.Errorf("failed to scan history entry: %w", err)
		}
		historyEntries = append(historyEntries, entry)
	}
	return historyEntries, nil
}
