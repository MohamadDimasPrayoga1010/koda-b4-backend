package libs

import (
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupCORS() gin.HandlerFunc {
	envOrigins := os.Getenv("ALLOW_ORIGIN")
	origins := []string{"http://localhost:5173"} 

	if envOrigins != "" {
		origins = append(origins, strings.Split(envOrigins, ",")...)
	}

	config := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           24 * time.Hour,
		AllowOriginFunc: func(origin string) bool {
			for _, o := range origins {
				if o == origin {
					return true
				}
			}
			return false
		},
	}

	return cors.New(config)
}
