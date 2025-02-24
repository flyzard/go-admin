package handlers

import (
	"net/http"
	"strings"

	"belcamp/internal/service"
	"belcamp/internal/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
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
	h.Render(c, "auth.login", gin.H{
		"title": "Login",
	}, "")
}

// Login logs in a user
func (h *AuthHandler) Login(c *gin.Context) {
	var form struct {
		Email    string `form:"email" binding:"required,email"`
		Password string `form:"password" binding:"required"`
	}

	if err := c.ShouldBind(&form); err != nil {
		h.Render(c, "auth.login", gin.H{
			"error": "Please fill in all fields correctly",
		}, "")
		return
	}

	// Clean input
	email := strings.TrimSpace(strings.ToLower(form.Email))

	// Get user by email
	user, err := h.userService.GetUserByEmail(email)
	if err != nil {
		h.Render(c, "auth.login", gin.H{
			"error": "Invalid credentials",
		}, "")
		return
	}

	if !utils.CheckPassword(form.Password, user.Password) {
		h.Render(c, "auth.login", gin.H{
			"title": "Login",
			"error": "Invalid credentials",
		}, "")
		return
	}

	// Set session
	session := sessions.Default(c)
	session.Set("userID", user.ID)
	if err := session.Save(); err != nil {
		h.Render(c, "auth.login", gin.H{
			"error": "Failed to save session",
		}, "")
		return
	}

	c.Redirect(http.StatusFound, "/")
}

// Logout logs out a user
func (h *AuthHandler) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	if err := session.Save(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Redirect(http.StatusFound, "/login")
}
