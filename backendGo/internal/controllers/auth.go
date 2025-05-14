package controllers

import (
	"net/http"
	"time"

	"github.com/Parron01/GerenciadorEstoque/backendGo/internal/config"
	"github.com/Parron01/GerenciadorEstoque/backendGo/internal/database"
	"github.com/Parron01/GerenciadorEstoque/backendGo/internal/middleware"
	"github.com/Parron01/GerenciadorEstoque/backendGo/internal/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// AuthController handles authentication-related requests
type AuthController struct{
    Config *config.Config
}

// NewAuthController creates a new auth controller
func NewAuthController(cfg *config.Config) *AuthController {
    return &AuthController{
        Config: cfg,
    }
}

// Login authenticates a user and returns a JWT token
func (ac *AuthController) Login(c *gin.Context) {
    var loginRequest struct {
        Username string `json:"username" binding:"required"`
        Password string `json:"password" binding:"required"`
    }

    if err := c.ShouldBindJSON(&loginRequest); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid login request"})
        return
    }

    // Check if user exists
    var user models.User
    var hashedPassword string
    err := database.DB.QueryRow(
        "SELECT id, username, password FROM users WHERE username = $1", 
        loginRequest.Username,
    ).Scan(&user.ID, &user.Username, &hashedPassword)

    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuário não encontrado"})
        return
    }

    // Compare password
    if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(loginRequest.Password)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Senha incorreta"})
        return
    }

    // Generate JWT token
    token, err := middleware.GenerateToken(user.Username, ac.Config)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "token": token,
        "user": gin.H{
            "id":       user.ID,
            "username": user.Username,
        },
    })
}

// Verify checks if a JWT token is valid
func (ac *AuthController) Verify(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"valid": true})
}

// Health verifies the server status
func (ac *AuthController) Health(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "status":  "online",
        "message": "Server is running",
        "time":    time.Now(),
    })
}