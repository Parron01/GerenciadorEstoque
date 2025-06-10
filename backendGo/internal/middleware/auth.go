package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Parron01/GerenciadorEstoque/backendGo/internal/config"
	"github.com/Parron01/GerenciadorEstoque/backendGo/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Claims represents the JWT claims structure
type Claims struct {
    UserID   int    `json:"userID"`
    Username string `json:"username"`
    jwt.RegisteredClaims
}

// GenerateToken generates a new JWT token for a user
func GenerateToken(userID int, username string, cfg *config.Config) (string, error) {
    expirationTime := time.Now().Add(cfg.JWT.Expiration)

    claims := &Claims{
        UserID:   userID,
        Username: username,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expirationTime),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString([]byte(cfg.JWT.Secret))

    return tokenString, err
}

// AuthMiddleware is a middleware function for authentication
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }

        // Extract the token from the Authorization header
        // Format: "Bearer {token}"
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
            c.Abort()
            return
        }

        tokenString := parts[1]

        // Parse and validate the token
        claims := &Claims{}
        token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }
            return []byte(cfg.JWT.Secret), nil
        })

        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        if !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        // Set username and userID in context for later use
        c.Set("userID", claims.UserID)
        c.Set("username", claims.Username)
        c.Next()
    }
}

// AuthenticateToken verifies that the request includes a valid JWT token
func AuthenticateToken(c *gin.Context) {
    // Skip authentication in development mode
    if os.Getenv("GO_ENV") == "development" {
        c.Next()
        return
    }

    authHeader := c.GetHeader("Authorization")
    if authHeader == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"message": "Token de autenticação não fornecido"})
        c.Abort()
        return
    }

    // Extract the token from the Authorization header
    // Format: "Bearer {token}"
    parts := strings.Split(authHeader, " ")
    if len(parts) != 2 || parts[0] != "Bearer" {
        c.JSON(http.StatusUnauthorized, gin.H{"message": "Formato de token inválido"})
        c.Abort()
        return
    }

    token := parts[1]

    // Verify the token
    claims, err := utils.VerifyJWT(token)
    if err != nil {
        c.JSON(http.StatusForbidden, gin.H{"message": "Token inválido ou expirado"})
        c.Abort()
        return
    }

    // Store user claims in context for use by subsequent handlers
    c.Set("userID", claims.ID)
    c.Set("username", claims.Username)

    c.Next()
}