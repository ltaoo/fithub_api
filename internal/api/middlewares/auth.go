package middlewares

import (
	"net/http"
	"strings"
	"time"

	"myapi/config"
	"myapi/internal/models"
	"myapi/pkg/logger"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware checks for a valid JWT token in the Authorization header
func AuthMiddleware(logger *logger.Logger, config *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth_header := c.GetHeader("Authorization")
		if auth_header == "" {
			c.JSON(http.StatusOK, gin.H{"code": 401, "msg": "缺少登录凭证", "data": nil})
			c.Abort()
			return
		}
		// fmt.Println("process in auth middleware")
		// fmt.Println(auth_header)
		// Extract the token from the Authorization header
		// Format: "Bearer <token>"
		parts := strings.Split(auth_header, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusOK, gin.H{"code": 401, "msg": "非法凭证请重新登录", "data": nil})
			c.Abort()
			return
		}

		token_str := strings.TrimPrefix(auth_header, "Bearer ")
		claims, err := models.ParseJWT(token_str, config.TokenSecretKey)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 401, "msg": "凭证失效请重新登录", "data": nil})
			c.Abort()
			return
		}
		// 检查过期时间
		if claims.ExpiresAt < float64(time.Now().Unix()) {
			c.JSON(http.StatusOK, gin.H{"code": 401, "msg": "凭证过期请重新登录", "data": nil})
			c.Abort()
			return
		}
		if claims.Id == 0 {
			c.JSON(http.StatusOK, gin.H{"code": 401, "msg": "请先登录", "data": nil})
			c.Abort()
			return
		}
		c.Set("id", claims.Id)
		c.Next()
	}
}
