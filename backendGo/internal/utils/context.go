package utils

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// GetUserIDFromContext extracts the user ID from the Gin context
// This assumes the AuthMiddleware has set the userID in the context
func GetUserIDFromContext(c *gin.Context) (int, error) {
	userID, exists := c.Get("userID")
	if !exists {
		return 0, fmt.Errorf("user ID not found in context")
	}

	id, ok := userID.(int)
	if !ok {
		return 0, fmt.Errorf("invalid user ID format in context")
	}

	return id, nil
}

// GetUsernameFromContext extracts the username from the Gin context
func GetUsernameFromContext(c *gin.Context) (string, error) {
	username, exists := c.Get("username")
	if !exists {
		return "", fmt.Errorf("username not found in context")
	}

	name, ok := username.(string)
	if !ok {
		return "", fmt.Errorf("invalid username format in context")
	}

	return name, nil
}
