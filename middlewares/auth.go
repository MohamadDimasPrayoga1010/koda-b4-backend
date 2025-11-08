package middlewares

import (

	"github.com/gin-gonic/gin"
)
func AuthRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, exists := ctx.Get("userID")
		if !exists {
			ctx.JSON(401, gin.H{
				"success": false,
				"message": "Unauthorized: missing user info",
			})
			ctx.Abort()
			return
		}
		ctx.Set("userID", userID)
		ctx.Next()
	}
}

func AdminOnly() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		roleValue, exists := ctx.Get("userRole")
		if !exists {
			ctx.JSON(401, gin.H{
				"success": false,
				"message": "Unauthorized: missing role",
			})
			ctx.Abort()
			return
		}

		role, ok := roleValue.(string)
		if !ok || role != "admin" {
			ctx.JSON(403, gin.H{
				"success": false,
				"message": "Forbidden: admin access only",
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
