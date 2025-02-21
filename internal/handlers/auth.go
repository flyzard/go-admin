package handlers

import (
	"net/http"
	"strings"

	"belcamp/internal/service"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// AuthHandler is a handler for authentication
type AuthHandler struct {
	BaseHandler
	userService service.UserService
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(userService service.UserService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
	}
}

// ShowLogin renders the login page
func (h *AuthHandler) ShowLogin(c *gin.Context) {
	h.Render(c, "pages/login.html", gin.H{
		"title": "Login",
	})
}

// Login logs in a user
func (h *AuthHandler) Login(c *gin.Context) {
	var form struct {
		Email    string `form:"email" binding:"required,email"`
		Password string `form:"password" binding:"required"`
	}

	if err := c.ShouldBind(&form); err != nil {
		if c.GetHeader("HX-Request") == "true" {
			c.HTML(http.StatusBadRequest, "partials/auth/login-form.html", gin.H{
				"error": "Please fill in all fields correctly",
			})
			return
		}
		h.Render(c, "auth/login.html", gin.H{
			"error": "Please fill in all fields correctly",
		})
		return
	}

	// Clean input
	email := strings.TrimSpace(strings.ToLower(form.Email))

	// Get user by email
	user, err := h.userService.GetUserByEmail(email)
	if err != nil {
		if c.GetHeader("HX-Request") == "true" {
			c.HTML(http.StatusUnauthorized, "partials/auth/login-form.html", gin.H{
				"error": "Invalid credentials",
			})
			return
		}
		h.Render(c, "auth/login.html", gin.H{
			"error": "Invalid credentials",
		})
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password)); err != nil {
		if c.GetHeader("HX-Request") == "true" {
			c.HTML(http.StatusUnauthorized, "partials/auth/login-form.html", gin.H{
				"error": "Invalid credentials",
			})
			return
		}
		h.Render(c, "auth/login.html", gin.H{
			"error": "Invalid credentials",
		})
		return
	}

	// Set session
	session := sessions.Default(c)
	session.Set("userID", user.ID)
	session.Save()

	// Redirect based on request type
	if c.GetHeader("HX-Request") == "true" {
		c.Header("HX-Redirect", "/dashboard")
		c.Status(http.StatusOK)
		return
	}

	c.Redirect(http.StatusFound, "/dashboard")
}

// Logout logs out a user
func (h *AuthHandler) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()

	if c.GetHeader("HX-Request") == "true" {
		c.Header("HX-Redirect", "/login")
		c.Status(http.StatusOK)
		return
	}

	c.Redirect(http.StatusFound, "/login")
}
