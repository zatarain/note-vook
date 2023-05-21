package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/zatarain/note-vook/mocks"
	"github.com/zatarain/note-vook/models"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/slices"
	"gorm.io/gorm"
)

func TestSignup(test *testing.T) {
	assert := assert.New(test)
	gin.SetMode(gin.TestMode)

	// Teardown test suite
	defer monkey.UnpatchAll()

	test.Run("Should create a new user", func(test *testing.T) {
		// Arrange
		server := gin.New()
		database := new(mocks.MockedDataAccessInterface)
		users := &UsersController{Database: database}
		database.
			On("Create", mock.AnythingOfType("*models.User")).
			Return(&gorm.DB{Error: nil})
		server.POST("/signup", users.Signup)
		user := Credentials{
			Nickname: "dummy-user",
			Password: "top-secret",
		}
		body, _ := json.Marshal(user)
		request, _ := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(body))
		recorder := httptest.NewRecorder()

		// Act
		server.ServeHTTP(recorder, request)

		// Assert
		assert.Equal(http.StatusCreated, recorder.Code)
		database.AssertExpectations(test)
	})

	test.Run("Should NOT create a duplicated user", func(test *testing.T) {
		// Arrange
		server := gin.New()
		database := new(mocks.MockedDataAccessInterface)
		users := &UsersController{Database: database}
		database.
			On("Create", mock.AnythingOfType("*models.User")).
			Return(&gorm.DB{Error: errors.New("User already exists")})
		server.POST("/signup", users.Signup)
		user := Credentials{
			Nickname: "dummy-user",
			Password: "top-secret",
		}
		body, _ := json.Marshal(user)
		request, _ := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(body))
		recorder := httptest.NewRecorder()

		// Act
		server.ServeHTTP(recorder, request)

		// Assert
		assert.Equal(http.StatusBadRequest, recorder.Code)
		assert.Contains(recorder.Body.String(), "User already exists")
		database.AssertExpectations(test)
	})

	test.Run("Should NOT try to create a user when unable to bind JSON", func(test *testing.T) {
		// Arrange
		server := gin.New()
		database := new(mocks.MockedDataAccessInterface)
		users := &UsersController{Database: database}
		database.
			On("Create", mock.AnythingOfType("*models.User")).
			Return(&gorm.DB{Error: nil})
		server.POST("/signup", users.Signup)
		body := bytes.NewBuffer([]byte("Malformed JSON"))
		request, _ := http.NewRequest(http.MethodPost, "/signup", body)
		recorder := httptest.NewRecorder()

		// Act
		server.ServeHTTP(recorder, request)

		// Assert
		assert.Equal(http.StatusBadRequest, recorder.Code)
		assert.Contains(recorder.Body.String(), "Failed to read input")
		database.AssertNotCalled(test, "Create", mock.AnythingOfType("*models.User"))
	})

	test.Run("Should NOT try to create a user when unable hash password", func(test *testing.T) {
		// Arrange
		server := gin.New()
		database := new(mocks.MockedDataAccessInterface)
		users := &UsersController{Database: database}
		database.
			On("Create", mock.AnythingOfType("*models.User")).
			Return(&gorm.DB{Error: nil})
		server.POST("/signup", users.Signup)
		user := Credentials{
			Nickname: "dummy-user",
			Password: "top-secret",
		}
		body, _ := json.Marshal(user)
		request, _ := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(body))
		recorder := httptest.NewRecorder()
		monkey.Patch(bcrypt.GenerateFromPassword, func([]byte, int) ([]byte, error) {
			return []byte{}, errors.New("Unable to hash")
		})

		// Act
		server.ServeHTTP(recorder, request)

		// Assert
		assert.Equal(http.StatusBadRequest, recorder.Code)
		assert.Contains(recorder.Body.String(), "Failed to create the hash for password")
		database.AssertNotCalled(test, "Create", mock.AnythingOfType("*models.User"))
	})
}

func TestLogin(test *testing.T) {
	assert := assert.New(test)
	gin.SetMode(gin.TestMode)

	CompareSuccessful := func([]byte, []byte) error {
		return nil
	}

	CompareFailure := func([]byte, []byte) error {
		return errors.New("Invalid Password")
	}

	NiceFakeToken := func(*UsersController, *models.User) (string, error) {
		return "Nice Fake Token", nil
	}

	NoToken := func(*UsersController, *models.User) (string, error) {
		return "", errors.New("No Token")
	}

	CheckCookie := func(cookie *http.Cookie) bool {
		return cookie.Name == "Authorisation"
	}

	// Teardown test suite
	defer monkey.UnpatchAll()

	test.Run("Should login the user and create the token", func(test *testing.T) {
		// Arrange
		server := gin.New()
		database := new(mocks.MockedDataAccessInterface)
		users := &UsersController{Database: database}
		anyUser := mock.AnythingOfType("*models.User")
		call := database.
			On("First", anyUser, "nickname = ?", "dummy-user").
			Return(&gorm.DB{Error: nil})
		call.RunFn = func(arguments mock.Arguments) {
			user := arguments.Get(0).(*models.User)
			user.ID = 12345
			user.Nickname = "dummy-user"
			user.Password = "top-secret"
		}

		monkey.Patch(bcrypt.CompareHashAndPassword, CompareSuccessful)
		monkey.PatchInstanceMethod(reflect.TypeOf(users), "NewToken", NiceFakeToken)
		server.POST("/login", users.Login)
		user := Credentials{
			Nickname: "dummy-user",
			Password: "top-secret",
		}
		body, _ := json.Marshal(user)
		request, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
		recorder := httptest.NewRecorder()

		// Act
		server.ServeHTTP(recorder, request)
		cookies := recorder.Result().Cookies()
		index := slices.IndexFunc(cookies, CheckCookie)

		// Assert
		database.AssertExpectations(test)
		assert.Equal(http.StatusOK, recorder.Code)
		assert.Contains(recorder.Body.String(), "Yaaay! You are logged in :)")
		require.GreaterOrEqual(test, index, 0)
		assert.Equal("Nice+Fake+Token", cookies[index].Value)
		assert.Equal(7*24*60*60, cookies[index].MaxAge)
		assert.False(cookies[index].Secure)
		assert.True(cookies[index].HttpOnly)
	})

	test.Run("Should response with internal server error when unable to generate token", func(test *testing.T) {
		// Arrange
		server := gin.New()
		database := new(mocks.MockedDataAccessInterface)
		users := &UsersController{Database: database}
		anyUser := mock.AnythingOfType("*models.User")
		call := database.
			On("First", anyUser, "nickname = ?", "dummy-user").
			Return(&gorm.DB{Error: nil})
		call.RunFn = func(arguments mock.Arguments) {
			user := arguments.Get(0).(*models.User)
			user.ID = 12345
			user.Nickname = "dummy-user"
			user.Password = "top-secret"
		}

		monkey.Patch(bcrypt.CompareHashAndPassword, CompareSuccessful)
		monkey.PatchInstanceMethod(reflect.TypeOf(users), "NewToken", NoToken)
		server.POST("/login", users.Login)
		user := Credentials{
			Nickname: "dummy-user",
			Password: "top-secret",
		}
		body, _ := json.Marshal(user)
		request, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
		recorder := httptest.NewRecorder()

		// Act
		server.ServeHTTP(recorder, request)
		cookies := recorder.Result().Cookies()
		index := slices.IndexFunc(cookies, CheckCookie)

		// Assert
		database.AssertExpectations(test)
		assert.Equal(http.StatusInternalServerError, recorder.Code)
		assert.Contains(recorder.Body.String(), "Unable to generate access token")
		require.Equal(test, index, -1)
	})

	test.Run("Should NOT try to login the user when unable to bind JSON", func(test *testing.T) {
		// Arrange
		server := gin.New()
		database := new(mocks.MockedDataAccessInterface)
		users := &UsersController{Database: database}
		anyUser := mock.AnythingOfType("*models.User")
		database.On("First", anyUser).Return(&gorm.DB{Error: nil})
		server.POST("/login", users.Login)
		body := bytes.NewBuffer([]byte("Malformed JSON"))
		request, _ := http.NewRequest(http.MethodPost, "/login", body)
		recorder := httptest.NewRecorder()

		// Act
		server.ServeHTTP(recorder, request)

		// Assert
		assert.Equal(http.StatusBadRequest, recorder.Code)
		assert.Contains(recorder.Body.String(), "Failed to read input")
		database.AssertNotCalled(test, "First", mock.AnythingOfType("*models.User"))
	})

	user := Credentials{
		Nickname: "dummy-user",
		Password: "top-secret",
	}

	InvalidNicknameOrPasswordTestcases := []struct {
		description string
		user        models.User
		compare     func([]byte, []byte) error
	}{
		{
			description: "Should NOT login the user when we didn't find nickname in database",
			user: models.User{
				ID:       0,
				Nickname: "",
				Password: "",
			},
			compare: CompareSuccessful,
		},
		{
			description: "Should NOT login the user when password doesn't match with stored hash",
			user: models.User{
				ID:       12345,
				Nickname: user.Nickname,
				Password: "secret-top",
			},
			compare: CompareFailure,
		},
		{
			description: "Should NOT login the user when either we didn't find nickname in database or password doesn't match",
			user: models.User{
				ID:       0,
				Nickname: "",
				Password: "",
			},
			compare: CompareFailure,
		},
	}

	anyUser := mock.AnythingOfType("*models.User")

	for _, testcase := range InvalidNicknameOrPasswordTestcases {
		test.Run(testcase.description, func(test *testing.T) {
			// Arrange
			server := gin.New()
			database := new(mocks.MockedDataAccessInterface)
			users := &UsersController{Database: database}
			call := database.
				On("First", anyUser, "nickname = ?", user.Nickname).
				Return(&gorm.DB{Error: nil})
			call.RunFn = func(arguments mock.Arguments) {
				user := arguments.Get(0).(*models.User)
				user.ID = testcase.user.ID
				user.Nickname = testcase.user.Nickname
				user.Password = testcase.user.Password
			}

			monkey.Patch(bcrypt.CompareHashAndPassword, testcase.compare)

			server.POST("/login", users.Login)
			body, _ := json.Marshal(user)
			request, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
			recorder := httptest.NewRecorder()

			// Act
			server.ServeHTTP(recorder, request)

			// Assert
			assert.Equal(http.StatusBadRequest, recorder.Code)
			assert.Contains(recorder.Body.String(), "Invalid nickname or password")
			database.AssertExpectations(test)
		})
	}
}

func TestAuthorise(test *testing.T) {
	assert := assert.New(test)
	gin.SetMode(gin.TestMode)
	today := time.Now()
	dummy := models.User{
		ID:        12345,
		Nickname:  "dummy-user",
		Password:  "top-secret",
		CreatedAt: today,
		UpdatedAt: today,
	}

	AuthorisedEndPointHandler := func(context *gin.Context) {
		data, exists := context.Get("user")
		user := data.(*models.User)

		assert.True(exists)
		assert.Equal(&dummy, user)
	}

	UnauthorisedEndPointHandler := func(*gin.Context) {
		assert.False(true, "This should never run otherwise the test failed!")
	}

	ValidToken := func(*UsersController, *gin.Context) (*models.User, error) {
		return &dummy, nil
	}

	InvalidToken := func(*UsersController, *gin.Context) (*models.User, error) {
		return nil, errors.New("Invalid token")
	}

	// Teardown test suite
	defer monkey.UnpatchAll()

	test.Run("Should set the user within the context and continue when token is valid", func(test *testing.T) {
		// Arrange
		server := gin.New()
		database := new(mocks.MockedDataAccessInterface)
		users := &UsersController{Database: database}
		server.GET("/", users.Authorise, AuthorisedEndPointHandler)
		monkey.PatchInstanceMethod(reflect.TypeOf(users), "ValidateToken", ValidToken)
		request, _ := http.NewRequest("GET", "/", nil)
		recorder := httptest.NewRecorder()

		// Act
		server.ServeHTTP(recorder, request)

		// Assert
		assert.Equal(http.StatusOK, recorder.Code)
	})

	test.Run("Should set the user within the context and continue when token is valid", func(test *testing.T) {
		// Arrange
		server := gin.New()
		database := new(mocks.MockedDataAccessInterface)
		users := &UsersController{Database: database}
		server.GET("/", users.Authorise, UnauthorisedEndPointHandler)
		monkey.PatchInstanceMethod(reflect.TypeOf(users), "ValidateToken", InvalidToken)
		request, _ := http.NewRequest("GET", "/", nil)
		recorder := httptest.NewRecorder()

		// Act
		server.ServeHTTP(recorder, request)

		// Assert
		assert.Equal(http.StatusUnauthorized, recorder.Code)
		assert.Contains(recorder.Body.String(), "Invalid token")
	})
}

func TestNewToken(test *testing.T) {
	assert := assert.New(test)
	users := &UsersController{SecretTokenKey: "super-secret-key"}
	FakeNow := func() time.Time {
		now, _ := time.Parse(time.DateOnly, "2021-01-01")
		return now
	}

	// Teardown test suite
	defer monkey.UnpatchAll()

	test.Run("Should generate the token", func(test *testing.T) {
		// Arrange
		user := &models.User{Nickname: "dummy-user"}
		monkey.Patch(time.Now, FakeNow)
		expiration, _ := time.Parse(time.DateOnly, "2021-01-08")

		// Act
		token, exception := users.NewToken(user)

		// Assert
		parsed, _ := jwt.Parse(token, users.Decoder)
		data, _ := parsed.Claims.(jwt.MapClaims)
		assert.Equal(user.Nickname, data["identifier"].(string))
		assert.Equal(expiration.Unix(), int64(data["expiration"].(float64)))
		assert.NotEmpty(token)
		assert.Nil(exception)
	})
}

func TestValidateToken(test *testing.T) {
	assert := assert.New(test)
	server := gin.New()
	users := &UsersController{SecretTokenKey: "super-secret-key"}
	var exception error
	var userResult *models.User
	FakeEndPoint := func(context *gin.Context) {
		userResult, exception = users.ValidateToken(context)
	}
	InvalidToken := func(string, jwt.Keyfunc, ...jwt.ParserOption) (*jwt.Token, error) {
		return &jwt.Token{Valid: false}, nil
	}
	server.GET("/", FakeEndPoint)

	// Teardown test suite
	defer monkey.UnpatchAll()

	test.Run("Should return user and non-error when user exists", func(test *testing.T) {
		// Arrange
		database := new(mocks.MockedDataAccessInterface)
		users.Database = database
		token, _ := users.NewToken(&models.User{Nickname: "dummy-user"})
		anyUser := mock.AnythingOfType("*models.User")
		call := database.
			On("First", anyUser, "nickname = ?", "dummy-user").
			Return(&gorm.DB{Error: nil})
		call.RunFn = func(arguments mock.Arguments) {
			record := arguments.Get(0).(*models.User)
			record.ID = 12345
			record.Nickname = "dummy-user"
			record.Password = "top-secret"
		}
		request, _ := http.NewRequest("GET", "/", nil)
		request.AddCookie(&http.Cookie{Name: "Authorisation", Value: token})
		recorder := httptest.NewRecorder()
		exception = nil

		// Act
		server.ServeHTTP(recorder, request)

		// Assert
		assert.Nil(exception)
		assert.NotNil(userResult)
		database.AssertExpectations(test)
	})

	test.Run("Should return error when there is no cookie", func(test *testing.T) {
		// Arrange
		request, _ := http.NewRequest("GET", "/", nil)
		recorder := httptest.NewRecorder()
		exception = nil

		// Act
		server.ServeHTTP(recorder, request)

		// Assert
		assert.Nil(userResult)
		assert.NotNil(exception)
	})

	test.Run("Should return error when detects is different algorithm", func(test *testing.T) {
		// Arrange
		token := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.MapClaims{"identifier": "token", "expiration": 4})
		signed, _ := token.SignedString([]byte(users.SecretTokenKey))
		request, _ := http.NewRequest("GET", "/", nil)
		request.AddCookie(&http.Cookie{Name: "Authorisation", Value: signed})
		recorder := httptest.NewRecorder()
		exception = nil

		// Act
		server.ServeHTTP(recorder, request)

		// Assert
		assert.Nil(userResult)
		assert.NotNil(exception)
	})

	test.Run("Should return error when detects an invalid token", func(test *testing.T) {
		// Arrange
		monkey.Patch(jwt.Parse, InvalidToken)
		defer monkey.UnpatchAll()
		request, _ := http.NewRequest("GET", "/", nil)
		request.AddCookie(&http.Cookie{Name: "Authorisation", Value: "invalid"})
		recorder := httptest.NewRecorder()
		exception = nil

		// Act
		server.ServeHTTP(recorder, request)

		// Assert
		assert.Nil(userResult)
		assert.NotNil(exception)
		assert.Contains(exception.Error(), "invalid authentication token")
	})

	test.Run("Should return error when detects an expired token", func(test *testing.T) {
		// Arrange
		token := jwt.NewWithClaims(
			jwt.SigningMethodHS256,
			jwt.MapClaims{
				"identifier": "token",
				"expiration": 0,
			},
		)
		expired, _ := token.SignedString([]byte(users.SecretTokenKey))
		request, _ := http.NewRequest("GET", "/", nil)
		request.AddCookie(&http.Cookie{Name: "Authorisation", Value: expired})
		recorder := httptest.NewRecorder()
		exception = nil

		// Act
		server.ServeHTTP(recorder, request)

		// Assert
		assert.Nil(userResult)
		assert.NotNil(exception)
		assert.Contains(exception.Error(), "expired session")
	})

	test.Run("Should return error when user doesn't exist", func(test *testing.T) {
		// Arrange
		database := new(mocks.MockedDataAccessInterface)
		users.Database = database
		token, _ := users.NewToken(&models.User{Nickname: "user-dummy"})
		anyUser := mock.AnythingOfType("*models.User")
		call := database.
			On("First", anyUser, "nickname = ?", "user-dummy").
			Return(&gorm.DB{Error: nil})
		call.RunFn = func(arguments mock.Arguments) {
			record := arguments.Get(0).(*models.User)
			record.ID = 0
			record.Nickname = "found!"
			record.Password = ""
		}
		request, _ := http.NewRequest("GET", "/", nil)
		request.AddCookie(&http.Cookie{Name: "Authorisation", Value: token})
		recorder := httptest.NewRecorder()
		exception = nil

		// Act
		server.ServeHTTP(recorder, request)

		// Assert
		database.AssertExpectations(test)
		assert.NotNil(exception)
		assert.Contains(exception.Error(), "user not found")
		assert.Nil(userResult)
	})
}

func TestDecoder(test *testing.T) {
	assert := assert.New(test)
	users := &UsersController{SecretTokenKey: "super-secret-key"}
	test.Run("Should return error when it's different algorithm", func(test *testing.T) {
		// Arrange
		token := &jwt.Token{
			Method: jwt.SigningMethodHS256,
		}

		// Act
		result, exception := users.Decoder(token)

		// Assert
		assert.NotNil(result)
		assert.Contains(string(result.([]byte)), "super-secret-key")
		assert.Nil(exception)
	})

	test.Run("Should return error when it's different algorithm", func(test *testing.T) {
		// Arrange
		token := &jwt.Token{
			Method: jwt.SigningMethodES512,
		}

		// Act
		result, exception := users.Decoder(token)

		// Assert
		assert.Nil(result)
		assert.NotNil(exception)
		assert.Contains(exception.Error(), "wrong signing method")
	})
}
