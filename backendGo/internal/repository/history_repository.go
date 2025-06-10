package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/Parron01/GerenciadorEstoque/backendGo/internal/models"
	"github.com/google/uuid"
)

// HistoryRepository defines the interface for history data operations
type HistoryRepository interface {
	Create(history *models.History) error
	CreateBatch(entries []models.History) error
	GetByBatchID(batchID string, userID int) ([]models.History, error)
	GetHistory(limit, offset int, userID int) ([]models.History, error)
	GetHistoryByEntity(entityType, entityID string, userID int) ([]models.History, error)
	GetGroupedHistoryBatches(page, pageSize int, userID int) (*models.PaginatedHistoryBatchGroups, error)
}

type historyRepository struct {
	db *sql.DB
}

// NewHistoryRepository creates a new HistoryRepository
func NewHistoryRepository(db *sql.DB) HistoryRepository {
	return &historyRepository{db: db}
}

// Create adds a new history entry to the database
func (r *historyRepository) Create(history *models.History) error {
	if history.ID == "" {
		history.ID = uuid.NewString()
	}
	if history.Date == "" {
		history.Date = time.Now().Format(time.RFC3339)
	}
	// If BatchID is not set, default it to the entry's own ID
	if history.BatchID == "" {
		history.BatchID = history.ID
	}

	// Ensure changes is valid JSON
	var js json.RawMessage = history.Changes
	if !json.Valid(js) {
		// Attempt to marshal if it's not already raw JSON string
		b, err := json.Marshal(history.Changes)
		if err != nil {
			return fmt.Errorf("history changes is not valid JSON and failed to marshal: %w", err)
		}
		js = json.RawMessage(b)
	}

	query := `INSERT INTO history (id, date, entity_type, entity_id, user_id, changes, batch_id) 
              VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.Exec(query, history.ID, history.Date, history.EntityType, history.EntityID, history.UserID, js, history.BatchID)
	if err != nil {
		return fmt.Errorf("failed to create history entry: %w", err)
	}
	return nil
}

// CreateBatch inserts multiple history entries into the database, typically within a transaction.
func (r *historyRepository) CreateBatch(entries []models.History) error {
	if len(entries) == 0 {
		return nil // No entries to insert
	}

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			tx.Rollback() // err is non-nil; don't change it
		} else {
			err = tx.Commit() // err is nil; if Commit returns error update err
		}
	}()

	stmt, err := tx.Prepare(`INSERT INTO history (id, date, entity_type, entity_id, user_id, changes, batch_id)
                             VALUES ($1, $2, $3, $4, $5, $6, $7)`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement for batch insert: %w", err)
	}
	defer stmt.Close()

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

		if _, err = stmt.Exec(entry.ID, entry.Date, entry.EntityType, entry.EntityID, entry.UserID, js, entry.BatchID); err != nil {
			return fmt.Errorf("failed to execute statement for entry %s in batch insert: %w", entry.ID, err)
		}
	}

	return nil
}

// GetByBatchID retrieves all history entries for a specific batch ID, ordered by date.
func (r *historyRepository) GetByBatchID(batchID string, userID int) ([]models.History, error) {
	var entries []models.History
	query := `SELECT id, date, entity_type, entity_id, changes, batch_id
              FROM history
              WHERE batch_id = $1 AND user_id = $2
              ORDER BY date ASC` // Order by date to maintain sequence within a batch
	rows, err := r.db.Query(query, batchID, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return []models.History{}, nil // Return empty slice if no rows found
		}
		return nil, fmt.Errorf("failed to query history entries by batch ID %s: %w", batchID, err)
	}
	defer rows.Close()

	for rows.Next() {
		var entry models.History
		if err := rows.Scan(&entry.ID, &entry.Date, &entry.EntityType, &entry.EntityID, &entry.Changes, &entry.BatchID); err != nil {
			return nil, fmt.Errorf("failed to scan history entry: %w", err)
		}
		entries = append(entries, entry)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration for batch ID %s: %w", batchID, err)
	}

	return entries, nil
}

// GetHistory retrieves a paginated list of all history entries, ordered by date descending.
func (r *historyRepository) GetHistory(limit, offset int, userID int) ([]models.History, error) {
	var entries []models.History
	query := `SELECT id, date, entity_type, entity_id, changes, batch_id
              FROM history
              WHERE user_id = $3
              ORDER BY date DESC
              LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(query, limit, offset, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query history entries: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var entry models.History
		if err := rows.Scan(&entry.ID, &entry.Date, &entry.EntityType, &entry.EntityID, &entry.Changes, &entry.BatchID); err != nil {
			return nil, fmt.Errorf("failed to scan history entry: %w", err)
		}
		entries = append(entries, entry)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration for GetHistory: %w", err)
	}
	return entries, nil
}

// GetHistoryByEntity retrieves all history entries for a specific entity, ordered by date descending.
func (r *historyRepository) GetHistoryByEntity(entityType, entityID string, userID int) ([]models.History, error) {
	var entries []models.History
	query := `SELECT id, date, entity_type, entity_id, changes, batch_id
              FROM history
              WHERE entity_type = $1 AND entity_id = $2 AND user_id = $3
              ORDER BY date DESC`
	rows, err := r.db.Query(query, entityType, entityID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query history for entity %s/%s: %w", entityType, entityID, err)
	}
	defer rows.Close()

	for rows.Next() {
		var entry models.History
		if err := rows.Scan(&entry.ID, &entry.Date, &entry.EntityType, &entry.EntityID, &entry.Changes, &entry.BatchID); err != nil {
			return nil, fmt.Errorf("failed to scan history entry for entity %s/%s: %w", entityType, entityID, err)
		}
		entries = append(entries, entry)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration for GetHistoryByEntity %s/%s: %w", entityType, entityID, err)
	}
	return entries, nil
}

// GetGroupedHistoryBatches retrieves history entries grouped by batch ID, with pagination for batches.
func (r *historyRepository) GetGroupedHistoryBatches(page, pageSize int, userID int) (*models.PaginatedHistoryBatchGroups, error) {
	var totalBatches int
	err := r.db.QueryRow("SELECT COUNT(DISTINCT batch_id) FROM history WHERE user_id = $1", userID).Scan(&totalBatches)
	if err != nil {
		if err == sql.ErrNoRows { // If no history entries at all
			totalBatches = 0
		} else {
			return nil, fmt.Errorf("failed to count total distinct batches: %w", err)
		}
	}

	if totalBatches == 0 {
		return &models.PaginatedHistoryBatchGroups{
			Groups:       []models.HistoryBatchGroup{},
			TotalBatches: 0,
			Page:         page,
			PageSize:     pageSize,
			TotalPages:   0,
		}, nil
	}

	offset := (page - 1) * pageSize

	// Step 1: Get paginated batch_ids and their first entry timestamp
	// We order by the earliest date within each batch to ensure consistent batch ordering.
	// The batch_id itself is used as a tie-breaker if timestamps are identical (though unlikely with RFC3339).
	batchQuery := `
        SELECT batch_id, MIN(date) as first_entry_date
        FROM history
        WHERE user_id = $3
        GROUP BY batch_id
        ORDER BY first_entry_date DESC, batch_id DESC
        LIMIT $1 OFFSET $2
    `
	type BatchInfo struct {
		BatchID        string `db:"batch_id"`
		FirstEntryDate string `db:"first_entry_date"`
	}
	var batchInfos []BatchInfo

	// Changed from r.db.Select to manual iteration
	rowsBatchInfo, err := r.db.Query(batchQuery, pageSize, offset, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query paginated batch IDs: %w", err)
	}
	defer rowsBatchInfo.Close()

	for rowsBatchInfo.Next() {
		var bi BatchInfo
		if err := rowsBatchInfo.Scan(&bi.BatchID, &bi.FirstEntryDate); err != nil {
			return nil, fmt.Errorf("failed to scan batch info: %w", err)
		}
		batchInfos = append(batchInfos, bi)
	}
	if err = rowsBatchInfo.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration for batch infos: %w", err)
	}

	if len(batchInfos) == 0 {
		return &models.PaginatedHistoryBatchGroups{
			Groups:       []models.HistoryBatchGroup{},
			TotalBatches: totalBatches,
			Page:         page,
			PageSize:     pageSize,
			TotalPages:   int(math.Ceil(float64(totalBatches) / float64(pageSize))),
		}, nil
	}

	groups := make([]models.HistoryBatchGroup, 0, len(batchInfos))

	// Step 2: For each batch_id, get all its records
	recordsQuery := `
        SELECT id, date, entity_type, entity_id, changes, batch_id
        FROM history
        WHERE batch_id = $1 AND user_id = $2
        ORDER BY date ASC
    ` // Order records within a batch by their date

	for _, batchInfo := range batchInfos {
		var records []models.History
		// Changed from r.db.Select to manual iteration
		rowsRecords, errRecords := r.db.Query(recordsQuery, batchInfo.BatchID, userID)
		if errRecords != nil {
			log.Printf("Error querying records for batch_id %s: %v. Skipping this batch.", batchInfo.BatchID, errRecords)
			continue
		}

		for rowsRecords.Next() {
			var record models.History
			if errScan := rowsRecords.Scan(&record.ID, &record.Date, &record.EntityType, &record.EntityID, &record.Changes, &record.BatchID); errScan != nil {
				log.Printf("Error scanning record for batch_id %s: %v. Skipping this record.", batchInfo.BatchID, errScan)
				// Decide if a single scan error should fail the whole batch or just skip the record
				continue
			}
			records = append(records, record)
		}
		if errRows := rowsRecords.Err(); errRows != nil {
			log.Printf("Error during rows iteration for records of batch_id %s: %v. Batch might be incomplete.", batchInfo.BatchID, errRows)
		}
		rowsRecords.Close() // Close rows for each batch iteration

		if len(records) > 0 {
			groups = append(groups, models.HistoryBatchGroup{
				BatchID:     batchInfo.BatchID,
				CreatedAt:   batchInfo.FirstEntryDate,
				Records:     records,
				RecordCount: len(records),
			})
		}
	}

	totalPages := int(math.Ceil(float64(totalBatches) / float64(pageSize)))
	if totalPages == 0 && totalBatches > 0 {
		totalPages = 1
	}

	return &models.PaginatedHistoryBatchGroups{
		Groups:       groups,
		TotalBatches: totalBatches,
		Page:         page,
		PageSize:     pageSize,
		TotalPages:   totalPages,
	}, nil
}
