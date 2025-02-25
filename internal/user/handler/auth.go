package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/dgrijalva/jwt-go"
	"github.com/markbates/goth/providers/google"
	"github.com/mavrk-mose/pay/internal/user/repository"
	"net/http"
)

type UserHandler struct {
	repo repository.UserRepository
}

func NewUserHandler(repo repository.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

var (
	oauthStateString = "}4PYRBlq{~m7)@wt%7jHfjo]8QyHaL6QxkwoB" // Change this to a secure random string in production read from yaml config
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	expirationTime := time.Now().Add(24 * time.Hour)
)

func InitGoogleOAuth(clientID, clientSecret, redirectURL string) {
	goth.UseProviders(
		google.New(clientID, clientSecret, redirectURL, "email", "profile"),
	)
}

func HandleGoogleLogin(c *gin.Context) {
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func (h *UserHandler) HandleGoogleCallback(c *gin.Context, db *sqlx.DB) {
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to complete user auth"})
		return
	}

	dbUser, err := h.repo.CreateOrUpdateUser(db, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create/update user"})
		return
	}

	token, err := GenerateJWT(dbUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate JWT"})
		return
	}


	c.JSON(http.StatusOK, gin.H{
		"message": "login successful", 
		"token": token,
		"expiresIn": expirationTime
	})
}


func GenerateJWT(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     expirationTime.Unix(),
	})

	return token.SignedString(jwtSecret)
}