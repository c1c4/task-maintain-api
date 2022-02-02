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
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	tm                       = time.Now()
	getUserByEmailRepository func(email string) (*models.User, error_utils.MessageErr)
	getUserByIdRepository    func(id uint64) (*models.User, error_utils.MessageErr)
	createUserRepository     func(user *models.User) (*models.User, error_utils.MessageErr)
	handlerCreateUser        = controllers.CreateUser
)

type userRepoMock struct{}

func (userRepo *userRepoMock) GetByEmail(email string) (*models.User, error_utils.MessageErr) {
	return getUserByEmailRepository(email)
}

func (userRepo *userRepoMock) Create(user *models.User) (*models.User, error_utils.MessageErr) {
	return createUserRepository(user)
}

func (userRepo *userRepoMock) Get(userId uint64) (*models.User, error_utils.MessageErr) {
	return getUserByIdRepository(userId)
}

func (userRepo *userRepoMock) Init() {}

func TestCreateUser_Success(t *testing.T) {
	repositories.UserRepo = &userRepoMock{}

	createUserRepository = func(user *models.User) (*models.User, error_utils.MessageErr) {
		return &models.User{
			ID:        1,
			Name:      "Test user",
			Email:     "test@test.com",
			Password:  "$2a$10$cdBTAX1B2KdXSbKaBdqY7utnuWDJHuw5V46TkzgEGrAQ4E1A6c6au",
			Type:      "Test",
			Tasks:     []models.Task{},
			CreatedAt: tm,
			UpdatedAt: tm,
		}, nil
	}

	jsonBody := `{"name": "Test user", "email": "test@test.com", "password": "123", "type": "Manager"}`
	r := gin.Default()
	req, errRequest := http.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(jsonBody))

	if errRequest != nil {
		t.Errorf("this is the error: %v\n", errRequest)
	}

	rr := httptest.NewRecorder()
	r.POST("/users", handlerCreateUser)
	r.ServeHTTP(rr, req)

	var user models.User
	err := json.Unmarshal(rr.Body.Bytes(), &user)
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Equal(t, uint64(1), user.ID)
	assert.Equal(t, "Test user", user.Name)
	assert.Equal(t, "test@test.com", user.Email)
	assert.NotEqual(t, "123", user.Password)
	assert.NotNil(t, user.CreatedAt)
	assert.NotNil(t, user.UpdatedAt)
}

func TestCreateUser_WrongJSONFormat(t *testing.T) {
	repositories.UserRepo = &userRepoMock{}

	jsonBody := `{"name": "Test user", "email": "test@test.com", "password": "" "type": "Manager"}`
	r := gin.Default()
	req, errRequest := http.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(jsonBody))

	if errRequest != nil {
		t.Errorf("this is the error: %v\n", errRequest)
	}

	rr := httptest.NewRecorder()
	r.POST("/users", handlerCreateUser)
	r.ServeHTTP(rr, req)

	apiErr, err := error_utils.NewApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.Equal(t, http.StatusUnprocessableEntity, apiErr.Status())
	assert.Equal(t, "it's not possible to convert the JSON into an object", apiErr.Message())
	assert.Equal(t, "invalid_request", apiErr.Error())
}

func TestCreateUser_WihtoutPassword(t *testing.T) {
	repositories.UserRepo = &userRepoMock{}

	jsonBody := `{"name": "Test user", "email": "test@test.com", "password": "", "type": "Manager"}`
	r := gin.Default()
	req, errRequest := http.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(jsonBody))

	if errRequest != nil {
		t.Errorf("this is the error: %v\n", errRequest)
	}

	rr := httptest.NewRecorder()
	r.POST("/users", handlerCreateUser)
	r.ServeHTTP(rr, req)

	apiErr, err := error_utils.NewApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.Equal(t, http.StatusBadRequest, apiErr.Status())
	assert.Equal(t, "the field password is required can't be empty", apiErr.Message())
	assert.Equal(t, "bad_request", apiErr.Error())
}

func TestCreateUser_WihtoutEmail(t *testing.T) {
	repositories.UserRepo = &userRepoMock{}

	jsonBody := `{"name": "Test user", "email": "", "password": "123", "type": "Manager"}`
	r := gin.Default()
	req, errRequest := http.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(jsonBody))

	if errRequest != nil {
		t.Errorf("this is the error: %v\n", errRequest)
	}

	rr := httptest.NewRecorder()
	r.POST("/users", handlerCreateUser)
	r.ServeHTTP(rr, req)

	apiErr, err := error_utils.NewApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.Equal(t, http.StatusBadRequest, apiErr.Status())
	assert.Equal(t, "the field email is required can't be empty", apiErr.Message())
	assert.Equal(t, "bad_request", apiErr.Error())
}

func TestCreateUser_InvalidEmail(t *testing.T) {
	repositories.UserRepo = &userRepoMock{}

	jsonBody := `{"name": "Test user", "email": "email", "password": "123", "type": "Manager"}`
	r := gin.Default()
	req, errRequest := http.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(jsonBody))

	if errRequest != nil {
		t.Errorf("this is the error: %v\n", errRequest)
	}

	rr := httptest.NewRecorder()
	r.POST("/users", handlerCreateUser)
	r.ServeHTTP(rr, req)

	apiErr, err := error_utils.NewApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.Equal(t, http.StatusBadRequest, apiErr.Status())
	assert.Equal(t, "the email is invalid", apiErr.Message())
	assert.Equal(t, "bad_request", apiErr.Error())
}

func TestCreateUser_InvalidUserType(t *testing.T) {
	repositories.UserRepo = &userRepoMock{}

	jsonBody := `{"name": "Test user", "email": "email@email.com", "password": "123", "type": "Test"}`
	r := gin.Default()
	req, errRequest := http.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(jsonBody))

	if errRequest != nil {
		t.Errorf("this is the error: %v\n", errRequest)
	}

	rr := httptest.NewRecorder()
	r.POST("/users", handlerCreateUser)
	r.ServeHTTP(rr, req)

	apiErr, err := error_utils.NewApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.Equal(t, http.StatusBadRequest, apiErr.Status())
	assert.Equal(t, "the field type should be Technician or Manager", apiErr.Message())
	assert.Equal(t, "bad_request", apiErr.Error())
}

func TestCreateUser_MultipleErrors(t *testing.T) {
	repositories.UserRepo = &userRepoMock{}

	jsonBody := `{"name": "", "email": "email@email.com", "password": "123", "type": ""}`
	r := gin.Default()
	req, errRequest := http.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(jsonBody))

	if errRequest != nil {
		t.Errorf("this is the error: %v\n", errRequest)
	}

	rr := httptest.NewRecorder()
	r.POST("/users", handlerCreateUser)
	r.ServeHTTP(rr, req)

	apiErr, err := error_utils.NewApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.Equal(t, http.StatusBadRequest, apiErr.Status())
	assert.Equal(
		t,
		"the field name is required can't be empty\nthe field type is required can't be empty",
		apiErr.Message(),
	)
	assert.Equal(t, "bad_request", apiErr.Error())
}
