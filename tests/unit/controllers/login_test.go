package controllers

import (
	"api/app/controllers"
	"api/app/models"
	"api/app/repositories"
	"api/app/utils/error_utils"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	handlerLogin = controllers.Login
)

func TestLogin_Success(t *testing.T) {
	repositories.UserRepo = &userRepoMock{}

	getUserByEmailRepository = func(email string) (*models.User, error_utils.MessageErr) {
		return &models.User{
			ID:       1,
			Name:     "Test user",
			Email:    "test@test.com",
			Password: "$2a$10$cdBTAX1B2KdXSbKaBdqY7utnuWDJHuw5V46TkzgEGrAQ4E1A6c6au",
			Type:     "Test",
		}, nil
	}

	jsonBody := `{"email": "test@test.com", "password": "123"}`
	r := gin.Default()
	req, errRequest := http.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(jsonBody))

	if errRequest != nil {
		t.Errorf("this is the error: %v\n", errRequest)
	}

	rr := httptest.NewRecorder()
	r.POST("/login", handlerLogin)
	r.ServeHTTP(rr, req)

	var authentication models.AuthenticationData
	err := json.Unmarshal(rr.Body.Bytes(), &authentication)
	assert.Nil(t, err)
	assert.NotNil(t, authentication)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "1", authentication.ID)
	assert.NotNil(t, authentication.Token)
}

func TestLogin_WrongJSONFormat(t *testing.T) {
	repositories.UserRepo = &userRepoMock{}

	jsonBody := `{"email": "test@test.com" "password": ""}`
	r := gin.Default()
	req, errRequest := http.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(jsonBody))

	if errRequest != nil {
		t.Errorf("this is the error: %v\n", errRequest)
	}

	rr := httptest.NewRecorder()
	r.POST("/login", handlerLogin)
	r.ServeHTTP(rr, req)

	apiErr, err := error_utils.NewApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.Equal(t, http.StatusUnprocessableEntity, apiErr.Status())
	assert.Equal(t, "it's not possible to convert the JSON into an object", apiErr.Message())
	assert.Equal(t, "invalid_request", apiErr.Error())
}

func TestLogin_WihtoutPassword(t *testing.T) {
	repositories.UserRepo = &userRepoMock{}

	getUserByEmailRepository = func(email string) (*models.User, error_utils.MessageErr) {
		return &models.User{
			ID:       1,
			Name:     "Test user",
			Email:    "test@test.com",
			Password: "$2a$10$cdBTAX1B2KdXSbKaBdqY7utnuWDJHuw5V46TkzgEGrAQ4E1A6c6au",
			Type:     "Test",
		}, nil
	}

	jsonBody := `{"email": "test@test.com", "password": ""}`
	r := gin.Default()
	req, errRequest := http.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(jsonBody))

	if errRequest != nil {
		t.Errorf("this is the error: %v\n", errRequest)
	}

	rr := httptest.NewRecorder()
	r.POST("/login", handlerLogin)
	r.ServeHTTP(rr, req)

	apiErr, err := error_utils.NewApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.Equal(t, http.StatusUnauthorized, apiErr.Status())
	assert.Equal(t, "crypto/bcrypt: hashedPassword is not the hash of the given password", apiErr.Message())
	assert.Equal(t, "unauthorized", apiErr.Error())
}

func TestLogin_WihtoutEmail(t *testing.T) {
	repositories.UserRepo = &userRepoMock{}

	getUserByEmailRepository = func(email string) (*models.User, error_utils.MessageErr) {
		return nil, error_utils.NewNotFoundError("no record matching given the identification")
	}

	jsonBody := `{"email": "", "password": "123"}`
	r := gin.Default()
	req, errRequest := http.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(jsonBody))

	if errRequest != nil {
		t.Errorf("this is the error: %v\n", errRequest)
	}

	rr := httptest.NewRecorder()
	r.POST("/login", handlerLogin)
	r.ServeHTTP(rr, req)

	apiErr, err := error_utils.NewApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.Equal(t, http.StatusNotFound, apiErr.Status())
	assert.Equal(t, "no record matching given the identification", apiErr.Message())
	assert.Equal(t, "not_found", apiErr.Error())
}
