package middlewares

import (
	"net/http"
	"strings"
	"time"

	"myapi/internal/models"
	"myapi/pkg/logger"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware checks for a valid JWT token in the Authorization header
func AuthMiddleware(logger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth_header := c.GetHeader("Authorization")
		if auth_header == "" {
			c.JSON(http.StatusOK, gin.H{"code": 401, "msg": "Authorization header required", "data": nil})
			c.Abort()
			return
		}

		// fmt.Println("process in auth middleware")
		// fmt.Println(auth_header)
		// Extract the token from the Authorization header
		// Format: "Bearer <token>"
		parts := strings.Split(auth_header, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusOK, gin.H{"code": 401, "msg": "Invalid Authorization header format", "data": nil})
			c.Abort()
			return
		}

		token_str := strings.TrimPrefix(auth_header, "Bearer ")
		claims, err := models.ParseJWT(token_str)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 401, "msg": "Invalid or expired token", "data": nil})
			c.Abort()
			return
		}
		// 检查过期时间
		if claims.ExpiresAt < float64(time.Now().Unix()) {
			c.JSON(http.StatusOK, gin.H{"code": 401, "msg": "Token has expired", "data": nil})
			c.Abort()
			return
		}
		c.Set("id", claims.Id)
		c.Next()
	}
}
