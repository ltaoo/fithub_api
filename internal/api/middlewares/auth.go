package middlewares

import (
	"net/http"
	"strings"

	"myapi/internal/models"
	"myapi/pkg/logger"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware checks for a valid JWT token in the Authorization header
func AuthMiddleware(logger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusOK, gin.H{"code": 401, "msg": "Authorization header required", "data": nil})
			c.Abort()
			return
		}

		// Extract the token from the Authorization header
		// Format: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusOK, gin.H{"code": 401, "msg": "Invalid Authorization header format", "data": nil})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := models.ParseJWT(tokenString)
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
