package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	cfg := cors.Config{
		// AllowOrigins: []string{}, // or []string{"https://xxx.com"}
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"HEAD",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
			"Accept",
			"Token",
		},
		AllowCredentials: true,
	}

	return cors.New(cfg)
}