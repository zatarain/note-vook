package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/zatarain/note-vook/models"
	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Nickname string
	Password string
}

type UsersController struct {
	Database       models.DataAccessInterface
	SecretTokenKey string
}

type TokenMaker interface {
	NewToken(*gin.Context, *models.User) (string, error)
}

type Authoriser interface {
	ValidateToken(*gin.Context) (*models.User, error)
}

func (credentials *Credentials) HashPassword() error {
	hash, exception := bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost)
	credentials.Password = string(hash)
	return exception
}

func getCredentialsFromRequest(context *gin.Context) *Credentials {
	var credentials Credentials

	// Trying to bind input from JSON
	if binding := context.BindJSON(&credentials); binding != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"summary": "Failed to read input",
			"details": binding.Error(),
		})
		return nil
	}

	return &credentials
}

func (users *UsersController) Signup(context *gin.Context) {
	credentials := getCredentialsFromRequest(context)
	if credentials == nil {
		return
	}

	// Trying to crete a hash for password
	if exception := credentials.HashPassword(); exception != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"summary": "Failed to create the hash for password",
			"details": exception.Error(),
		})
		return
	}

	// Insert user into the database table users
	user := models.User{
		Nickname: credentials.Nickname,
		Password: credentials.Password,
	}
	inserting := users.Database.Create(&user).Error
	if inserting != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"summary": "Failed to insert user into table users",
			"details": inserting.Error(),
		})
		return
	}

	context.JSON(http.StatusCreated, gin.H{
		"summary": "User successfully created",
		"details": user.String(),
	})
}

func (users *UsersController) NewToken(user *models.User) (string, error) {
	// Create the token for user that last for 7 days
	data := jwt.MapClaims{
		"identifier": user.Nickname,
		"expiration": time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, data)

	// Signing the token with secret key
	return token.SignedString([]byte(users.SecretTokenKey))
}

func (users *UsersController) Login(context *gin.Context) {
	credentials := getCredentialsFromRequest(context)
	if credentials == nil {
		return
	}

	// Checking the credentials
	user := &models.User{}
	users.Database.First(user, "nickname = ?", credentials.Nickname)
	failed := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(credentials.Password),
	)
	if user.ID == 0 || failed != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"summary": "Invalid nickname or password",
		})
		return
	}

	// Generate JWT Token and send it in the Cookies
	token, exception := users.NewToken(user)
	if exception != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"summary": "Unable to generate access token",
			"details": exception.Error(),
		})
		return
	}

	// Send cookie to the client
	context.SetSameSite(http.SameSiteLaxMode)
	context.SetCookie("Authorisation", token, 7*24*60*60, "", "", false, true)
	context.JSON(http.StatusOK, gin.H{"summary": "Yaaay! You are logged in :)"})
}

func (users *UsersController) Decoder(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("wrong signing method: %v", token.Header["alg"])
	}
	return []byte(users.SecretTokenKey), nil
}

func (users *UsersController) ValidateToken(context *gin.Context) (*models.User, error) {
	// Retrieving the Authorisation cookie
	cookie, exception := context.Cookie("Authorisation")
	if exception != nil {
		return nil, exception
	}

	// Decoding the token using the secret key
	token, exception := jwt.Parse(cookie, users.Decoder)
	if exception != nil {
		return nil, exception
	}

	// Validating token consistentcy and retrieving the claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
		return nil, errors.New("invalid authentication token")
	}

	// Checking expiration date/time
	expiration := int64(claims["expiration"].(float64))
	if time.Now().Unix() > expiration {
		return nil, errors.New("expired session")
	}

	// Looking for the user nickname
	user := &models.User{}
	users.Database.First(user, "nickname = ?", claims["identifier"].(string))
	if user.ID == 0 {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (users *UsersController) Authorise(context *gin.Context) {
	user, exception := users.ValidateToken(context)
	if exception != nil {
		context.AbortWithStatusJSON(
			http.StatusUnauthorized,
			gin.H{
				"summary": "Unauthorised",
				"details": exception.Error(),
			},
		)
	}

	// Attach user to context, allow access and continue
	context.Set("user", user)
	context.Next()
}
