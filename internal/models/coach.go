package models

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// Coach 教练模型
type Coach struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	Nickname  string     `json:"nickname" gorm:"not null"`
	AvatarURL string     `json:"avatar_url"`
	Config    string     `json:"config" gorm:"default:'{}'"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func (Coach) TableName() string {
	return "COACH"
}

// CoachAccount 教练账号模型
type CoachAccount struct {
	ProviderType string    `json:"provider_type" gorm:"primaryKey"`
	ProviderID   string    `json:"provider_id" gorm:"primaryKey"`
	ProviderArg1 int       `json:"provider_arg1"`
	ProviderArg2 string    `json:"provider_arg2"`
	ProviderArg3 string    `json:"provider_arg3"`
	CoachID      uint      `json:"coach_id" gorm:"not null"`
	CreatedAt    time.Time `json:"created_at"`
}

func (CoachAccount) TableName() string {
	return "COACH_ACCOUNT"
}

// AuthResponse represents the response for authentication operations
type AuthResponse struct {
	Token  string `json:"token"`
	Coach  Coach  `json:"coach"`
	Status string `json:"status"`
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	SecretKey     string
	TokenDuration time.Duration
}

// Claims represents the JWT claims
type Claims struct {
	CoachId string `json:"coach_id"`
	jwt.RegisteredClaims
}

// Default JWT configuration
var DefaultJWTConfig = JWTConfig{
	SecretKey:     "your-secret-key-change-in-production", // Should be loaded from environment variables
	TokenDuration: 24 * time.Hour,                         // 24 hours
}

// Helper function to generate JWT token
func GenerateJWT(coachId string) (string, error) {
	expirationTime := time.Now().Add(DefaultJWTConfig.TokenDuration)

	claims := &Claims{
		CoachId: coachId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "myapi",
			Subject:   coachId,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(DefaultJWTConfig.SecretKey))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseJWT parses and validates a JWT token
func ParseJWT(tokenString string) (*Claims, error) {
	// Remove "Bearer " prefix if present
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(DefaultJWTConfig.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// AuthMiddleware is a middleware to authenticate requests
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusOK, gin.H{"code": 401, "msg": "Authorization header is required", "data": nil})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := ParseJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 401, "msg": "Invalid or expired token", "data": nil})
			c.Abort()
			return
		}

		// Set coach ID in context
		c.Set("coachId", claims.CoachId)
		c.Next()
	}
}
