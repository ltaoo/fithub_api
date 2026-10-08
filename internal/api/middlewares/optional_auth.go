package middlewares

import (
	"github.com/gin-gonic/gin"
	"myapi/config"
	"myapi/internal/models"
	"strings"
	"time"
)

// OptionalAuthMiddleware keeps valid identity for catalog reads and otherwise treats the request as a guest.
func OptionalAuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		parts := strings.Fields(ctx.GetHeader("Authorization"))
		if len(parts) == 2 && parts[0] == "Bearer" {
			claims, err := models.ParseJWT(parts[1], cfg.TokenSecretKey)
			if err == nil && claims.Id != 0 && claims.ExpiresAt > float64(time.Now().Unix()) {
				ctx.Set("id", claims.Id)
			}
		}
		ctx.Next()
	}
}
