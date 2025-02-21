// Package middleware provides the middleware for the application.
package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware checks if user is authenticated
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("userID")

		if userID == nil {
			// If it's an HTMX request, respond accordingly
			if c.GetHeader("HX-Request") == "true" {
				c.Header("HX-Redirect", "/login")
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			// Regular browser request
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		// Set user info in context
		c.Set("userID", userID)
		c.Next()
	}
}

// NoAuthMiddleware ensures user is NOT authenticated (for login page etc.)
func NoAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("userID")

		if userID != nil {
			c.Redirect(http.StatusFound, "/dashboard")
			c.Abort()
			return
		}

		c.Next()
	}
}
