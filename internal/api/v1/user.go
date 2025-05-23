package handler

import (
	"github.com/mavrk-mose/pay/internal/user/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/apple"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/google"
	"github.com/mavrk-mose/pay/config"
	"github.com/mavrk-mose/pay/internal/api/middleware"
	"github.com/mavrk-mose/pay/internal/user/models"
	"github.com/mavrk-mose/pay/internal/user/repository"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(db *sqlx.DB) *UserHandler {
	return &UserHandler{service: repository.NewUserService(db)}
}

var (
	expirationTime time.Time
)

func InitAuth(cfg *config.Config) {
	var providers []goth.Provider

	middleware.InitSessionStore(cfg)
	gothic.Store = middleware.GetSessionStore()

	expirationTime = time.Now().Add(24 * time.Hour)

	if cfg.OAuth.Google.Enabled {
		providers = append(providers,
			google.New(
				cfg.OAuth.Google.ClientID,
				cfg.OAuth.Google.ClientSecret,
				cfg.OAuth.Google.RedirectURL,
			),
		)
	}

	if cfg.OAuth.Facebook.Enabled {
		providers = append(providers,
			facebook.New(
				cfg.OAuth.Facebook.ClientID,
				cfg.OAuth.Facebook.ClientSecret,
				cfg.OAuth.Facebook.RedirectURL,
			),
		)
	}

	if cfg.OAuth.Apple.Enabled {
		providers = append(providers,
			apple.New(
				cfg.OAuth.Apple.ClientID,
				cfg.OAuth.Apple.ClientSecret,
				cfg.OAuth.Apple.RedirectURL,
				nil,
				apple.ScopeName,
				apple.ScopeEmail,
			),
		)
	}

	goth.UseProviders(providers...)
}

func BeginAuthHandler(c *gin.Context) {
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func (h *UserHandler) AuthCallbackHandler(c *gin.Context) {
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete user auth: " + err.Error()})
		return
	}

	token, err := h.service.RegisterUser(c.Request.Context(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create/update user: " + err.Error()})
		return
	}

	c.Redirect(http.StatusFound, "http://localhost:8080/dashboard")

	c.JSON(http.StatusOK, gin.H{
		"message":   "Login successful",
		"token":     token,
		"expiresIn": expirationTime,
	})
}

func (h *UserHandler) LogoutHandler(c *gin.Context) {
	err := gothic.Logout(c.Writer, c.Request)
	if err != nil {
		return
	}
	c.Writer.Header().Set("Set-Cookie", "token=; Path=/; Expires=Thu, 01 Jan 1970 00:00:00 GMT; Secure")
	c.Redirect(http.StatusFound, "http://localhost:8080/auth/login")
	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successful",
	})
}

// ListUsers godoc
// @Summary      List users
// @Description  Retrieves users based on filters
// @Tags         users
// @Produce      json
// @Param        name  query  string  false  "User name"
// @Success      200  {array}  models.User
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	var filter models.UserFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid filter parameters"})
		return
	}

	users, err := h.service.ListUsers(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// AssignRole godoc
// @Summary      Assign role to user
// @Description  Assigns a role to a user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        userID  path  string  true  "User ID"
// @Param        role    body  string  true  "Role"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/users/{userID}/role [post]
func (h *UserHandler) AssignRole(c *gin.Context) {
	userID := c.Param("userID")
	var request struct {
		Role string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role is required"})
		return
	}

	if err := h.service.AssignRole(c.Request.Context(), userID, request.Role); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign role: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role assigned successfully"})
}

// RevokeRole removes a role from a user
func (h *UserHandler) RevokeRole(c *gin.Context) {
	userID := c.Param("userID")

	if err := h.service.RevokeRole(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to revoke role: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role revoked successfully"})
}

// BanUser bans a user and records the reason
func (h *UserHandler) BanUser(c *gin.Context) {
	userID := c.Param("userID")
	var request struct {
		Reason string `json:"reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ban reason is required"})
		return
	}

	if err := h.service.BanUser(c.Request.Context(), userID, request.Reason); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to ban user: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User banned successfully"})
}

// UnbanUser removes the ban from a user
func (h *UserHandler) UnbanUser(c *gin.Context) {
	userID := c.Param("userID")

	if err := h.service.UnbanUser(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unban user: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User unbanned successfully"})
}
