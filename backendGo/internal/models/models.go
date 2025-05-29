package models

import (
	"encoding/json"
	"time"
)

// User represents a user in the system
type User struct {
    ID        int       `json:"id"`
    Username  string    `json:"username"`
    Password  string    `json:"-"` // Password is not serialized to JSON
    CreatedAt time.Time `json:"created_at"`
}

// Product matches the Product interface from the Node.js backend
type Product struct {
    ID       string  `json:"id"`
    Name     string  `json:"name"`
    Unit     string  `json:"unit"`
    Quantity float64 `json:"quantity"`
    Lotes    []Lote  `json:"lotes,omitempty"` // Added: Lotes associated with the product
}

// Lote represents a batch of a product
type Lote struct {
    ID            string    `json:"id"`                       // UUID
    ProductID     string    `json:"product_id"`               // FK to Product.ID
    Quantity      float64   `json:"quantity" binding:"required,gt=0"`
    DataValidade  string    `json:"data_validade" binding:"required"` // YYYY-MM-DD
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
}

// History represents a history entry in the database
// It corresponds to the ProductHistory interface in Node.js
type History struct {
    ID         string          `json:"id" db:"id"`
    Date       string          `json:"date" db:"date"` // Consider using time.Time and handling format in marshalling/unmarshalling
    EntityType string          `json:"entityType" db:"entity_type"`
    EntityID   string          `json:"entityId" db:"entity_id"`
    Changes    json.RawMessage `json:"changes" db:"changes"` // Storing as raw JSON
    BatchID    string          `json:"batchId" db:"batch_id"` // New field for grouping history entries
}

// HistoryBatchGroup represents a collection of history records for a single batch operation.
type HistoryBatchGroup struct {
	BatchID     string    `json:"batchId"`
	CreatedAt   string    `json:"createdAt"` // Timestamp of the first entry in the batch, for ordering
	Records     []History `json:"records"`
	RecordCount int       `json:"recordCount"`
}

// PaginatedHistoryBatchGroups represents a paginated response for grouped history.
type PaginatedHistoryBatchGroups struct {
	Groups       []HistoryBatchGroup `json:"groups"`
	TotalBatches int                 `json:"totalBatches"`
	Page         int                 `json:"page"`
	PageSize     int                 `json:"pageSize"`
	TotalPages   int                 `json:"totalPages"`
}

// LoginRequest represents a login request
type LoginRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

// LoginResponse represents the response after a successful login
type LoginResponse struct {
    Token string `json:"token"`
}

// ProductChange matches the ProductChange interface from the Node.js backend
type ProductChange struct {
    ProductID       string  `json:"productId"`
    ProductName     string  `json:"productName"`
    Action          string  `json:"action"`
    QuantityChanged float64 `json:"quantityChanged"`
    QuantityBefore  float64 `json:"quantityBefore"`
    QuantityAfter   float64 `json:"quantityAfter"`
    IsNewProduct    bool    `json:"isNewProduct,omitempty"`
    IsProductRemoval bool   `json:"isProductRemoval,omitempty"`
}

// LoteChangeDetail describes changes made to a Lote for history records
type LoteChangeDetail struct {
	LoteID          string    `json:"loteId"`
	ProductID       string    `json:"productId"`
	Action          string    `json:"action"` // e.g., "created", "updated", "deleted"
	QuantityChanged *float64  `json:"quantityChanged,omitempty"`
	QuantityBefore  *float64  `json:"quantityBefore,omitempty"`
	QuantityAfter   *float64  `json:"quantityAfter,omitempty"`
	DataValidade    *string   `json:"dataValidade,omitempty"`    // Current value after change
	DataValidadeOld *string   `json:"dataValidadeOld,omitempty"` // Previous value if updated
	DataValidadeNew *string   `json:"dataValidadeNew,omitempty"` // New value if updated
}