package middlewares

import (
	"fmt"
	"main/libs"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
		fmt.Println(roleValue,"gagal")
		if !exists {
			ctx.JSON(401, gin.H{
				"success": false,
				"message": "Unauthorized: missing role",
			})
			ctx.Abort()
			return
		}

		role, ok := roleValue.(string)
		fmt.Println(role,"gagal")
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

func AuthMiddleware(role string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.JSON(401, gin.H{"success": false, "message": "Missing or invalid Authorization header"})
			ctx.Abort()
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		secret := os.Getenv("JWT_SECRET")
		claims := &libs.UserPayload{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			ctx.JSON(401, gin.H{"success": false, "message": "Invalid token"})
			ctx.Abort()
			return
		}

		if claims.Role != role{
			ctx.JSON(401, gin.H{"success": false, "message": "Not permission"})
			ctx.Abort()
			return 
		}
		ctx.Next()
	}
} 