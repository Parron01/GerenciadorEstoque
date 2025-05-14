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
}

// History represents a history entry in the database
// It corresponds to the ProductHistory interface in Node.js
type History struct {
    ID      string          `json:"id"`
    Date    string          `json:"date"`
    Changes json.RawMessage `json:"changes"` // Store as JSON string in DB
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