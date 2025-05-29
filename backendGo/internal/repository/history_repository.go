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
	GetDistinctBatchIDs(limit, offset int) (batchIDs []string, firstEntryDates []string, totalBatches int, err error)
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
		entry.Date = time.Now().Format(time.RFC3339)
	}
	// If BatchID is not set, default it to the entry's own ID
	if entry.BatchID == "" {
		entry.BatchID = entry.ID
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
		// BatchID should be pre-set by the service for all entries in a batch.
		// If not, this indicates a logic error upstream or a different use case.
		if entry.BatchID == "" {
			// Fallback, though ideally the service layer ensures this.
			// For CreateBatch, a common BatchID is expected to be set by the service.
			// If individual entries don't have it, they might get their own ID as BatchID.
			// This specific CreateBatch is for when the service *has* defined a common BatchID.
			return fmt.Errorf("entry in batch is missing a BatchID (ID: %s)", entry.ID)
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
              FROM history 
              WHERE batch_id = $1 
              ORDER BY date ASC` // Order records within the batch by date
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

func (r *historyRepository) GetDistinctBatchIDs(limit, offset int) (batchIDs []string, firstEntryDates []string, totalBatches int, err error) {
	countQuery := `
        SELECT COUNT(*)
        FROM (
            SELECT 1
            FROM history
            WHERE batch_id IS NOT NULL AND batch_id <> ''
            GROUP BY batch_id
        ) AS distinct_count;
    `
	err = r.db.QueryRow(countQuery).Scan(&totalBatches)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("failed to count distinct batch_ids: %w", err)
	}

	if totalBatches == 0 {
		return []string{}, []string{}, 0, nil
	}

	batchQuery := `
        SELECT batch_id, first_entry_date
        FROM (
            SELECT batch_id, MIN(date) as first_entry_date
            FROM history
            WHERE batch_id IS NOT NULL AND batch_id <> ''
            GROUP BY batch_id
            ORDER BY first_entry_date DESC
        ) AS distinct_batches
        LIMIT $1 OFFSET $2;
    `
	rows, err := r.db.Query(batchQuery, limit, offset)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("failed to query distinct batch_ids: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var batchID, firstEntryDate string
		if err := rows.Scan(&batchID, &firstEntryDate); err != nil {
			return nil, nil, 0, fmt.Errorf("failed to scan distinct batch_id: %w", err)
		}
		batchIDs = append(batchIDs, batchID)
		firstEntryDates = append(firstEntryDates, firstEntryDate)
	}
	if err = rows.Err(); err != nil {
		return nil, nil, 0, fmt.Errorf("error iterating distinct batch_id rows: %w", err)
	}

	return batchIDs, firstEntryDates, totalBatches, nil
}
