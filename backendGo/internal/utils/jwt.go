package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// JWTClaims defines the claims for our JWT tokens
type JWTClaims struct {
    ID        int    `json:"id"`
    Username  string `json:"username"`
    jwt.RegisteredClaims
}

// GenerateJWT creates a new JWT token for authentication
func GenerateJWT(userID int, username string) (string, error) {
    jwtSecret := getJWTSecret()
    expiration := getJWTExpiration()
    
    claims := JWTClaims{
        ID:       userID,
        Username: username,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(jwtSecret))
}

// VerifyJWT validates a JWT token and returns its claims
func VerifyJWT(tokenString string) (*JWTClaims, error) {
    jwtSecret := getJWTSecret()
    
    token, err := jwt.ParseWithClaims(
        tokenString, 
        &JWTClaims{}, 
        func(token *jwt.Token) (interface{}, error) {
            return []byte(jwtSecret), nil
        },
    )
    
    if err != nil {
        return nil, err
    }
    
    if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
        return claims, nil
    }
    
    return nil, jwt.ErrSignatureInvalid
}

// Helper functions to get JWT configuration
func getJWTSecret() string {
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        return "fallback_secret" // Fallback secret, same as in the Node backend
    }
    return secret
}

// getJWTExpiration returns the JWT expiration duration (default: 7 days)
func getJWTExpiration() time.Duration {
    expiration := os.Getenv("JWT_EXPIRATION")
    if expiration == "" {
        return 7 * 24 * time.Hour // Default: 7 days
    }
    
    // Parse the duration string (e.g., "7d", "24h")
    // For simplicity, we'll just use the default value here
    return 7 * 24 * time.Hour
}
