package handler

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/apple"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/google"
	"github.com/mavrk-mose/pay/config"
	"github.com/mavrk-mose/pay/internal/user/repository"
	"github.com/mavrk-mose/pay/pkg/middleware"
)

type UserHandler struct {
	repo repository.UserRepository
}

func NewUserHandler(repo repository.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

var (
	jwtSecret      []byte
	expirationTime time.Time
)

func InitAuth(cfg *config.Config) {
	var providers []goth.Provider

	// Initialize session store
	middleware.InitSessionStore(cfg)
	gothic.Store = middleware.GetSessionStore()

	jwtSecret = []byte(cfg.Server.JwtSecretKey)
	expirationTime = time.Now().Add(24 * time.Hour)

	// Register Google provider if enabled
	if cfg.OAuth.Google.Enabled {
		providers = append(providers,
			google.New(
				cfg.OAuth.Google.ClientID,
				cfg.OAuth.Google.ClientSecret,
				cfg.OAuth.Google.RedirectURL,
				"email",
				"profile",
			),
		)
	}

	if cfg.OAuth.Facebook.Enabled {
		providers = append(providers,
			facebook.New(
				cfg.OAuth.Facebook.ClientID,
				cfg.OAuth.Facebook.ClientSecret,
				cfg.OAuth.Facebook.RedirectURL,
				"email",
				"public_profile",
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

	// Set the OAuth state string from config or generate a secure random one
	stateString := cfg.OAuth.StateString
	if stateString == "" {
		stateString = middleware.GenerateSecureStateString()
	}
	gothic.SetState(stateString)
}

func BeginAuthHandler(c *gin.Context) {
	// The provider is passed as a URL parameter
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

// AuthCallbackHandler handles the callback from any provider
func (h *UserHandler) AuthCallbackHandler(c *gin.Context) {
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete user auth: " + err.Error()})
		return
	}

	user, err := h.repo.CreateOrUpdateUser(c, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create/update user: " + err.Error()})
		return
	}

	token, err := GenerateJWT(dbUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Login successful",
		"token":     token,
		"expiresIn": expirationTime,
	})
}

func GenerateJWT(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     expirationTime.Unix(),
	})

	return token.SignedString(jwtSecret)
}
