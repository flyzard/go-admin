// internal/middleware/csrf.go
package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
)

func CSRF() gin.HandlerFunc {
	csrfMiddleware := csrf.Protect(
		[]byte(os.Getenv("APP_KEY")),
		csrf.Secure(os.Getenv("GIN_MODE") == "release"),
		csrf.HttpOnly(true),
	)

	return func(c *gin.Context) {
		// Wrap the ResponseWriter
		csrfMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c.Request = r
			c.Next()
		})).ServeHTTP(c.Writer, c.Request)
	}
}
