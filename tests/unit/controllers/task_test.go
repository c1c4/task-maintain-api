package controllers

import (
	"api/app/config"
	"api/app/controllers"
	"api/app/models"
	"api/app/repositories"
	"api/app/utils/error_utils"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	technician_token         = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjk5OTk5NjU1ODI4LCJwZXJtaXNzaW9ucyI6WyJjcmVhdGUiLCJ1cGRhdGUiLCJnZXRfb25lIiwibGlzdF9vd25fdGFza3MiXSwidXNlcl9pZCI6MX0.aFj60IBurKQoSPQ25jHLu2yzBrsRLQwtlvncgL2_7IY"
	manager_token            = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjk5OTQzNjU1NzM1LCJwZXJtaXNzaW9ucyI6WyJsaXN0IiwiZGVsZXRlIiwibm90aWZpZWQiXSwidXNlcl9pZCI6Mn0.x1WUovq-J4R0RhaSoXbQTvJ9pkZvO-q_oJT9HY-Pkk8"
	getTaskByIdRepository    func(id uint64) (*models.Task, error_utils.MessageErr)
	createTaskRepository     func(task *models.Task) (*models.Task, error_utils.MessageErr)
	updateTaskReposiroty     func(task *models.Task) (*models.Task, error_utils.MessageErr)
	getTasksByUserRepository func(userID uint64) []models.Task
	getTasksRepository       func() []models.Task
	deleteTasksRepository    func(id uint64) error_utils.MessageErr
	handlerCreateTask        = controllers.CreateTask
	handlerUpdateTask        = controllers.UpdateTask
	handlerGetTask           = controllers.GetTask
	handlerGetTasksByUser    = controllers.GetTasksByUser
	handlerGetAllTasks       = controllers.GetAllTasks
	handlerDeleteTasks       = controllers.DeleteTasks
)

type taskRepoMock struct{}

func (taskRepo *taskRepoMock) Create(task *models.Task) (*models.Task, error_utils.MessageErr) {
	return createTaskRepository(task)
}

func (taskRepo *taskRepoMock) Get(taskId uint64) (*models.Task, error_utils.MessageErr) {
	return getTaskByIdRepository(taskId)
}

func (taskRepo *taskRepoMock) Update(task *models.Task) (*models.Task, error_utils.MessageErr) {
	return updateTaskReposiroty(task)
}

func (taskRepo *taskRepoMock) GetAll() []models.Task {
	return getTasksRepository()
}

func (taskRepo *taskRepoMock) GetAllByUserID(userID uint64) []models.Task {
	return getTasksByUserRepository(userID)
}

func (taskRepo *taskRepoMock) Delete(taskId uint64) error_utils.MessageErr {
	return deleteTasksRepository(taskId)
}

func (taskRepo *taskRepoMock) Init() {}

func TestCreateTask_Success(t *testing.T) {
	config.SECRETKEY = "mySecretK3y"
	repositories.TaskRepo = &taskRepoMock{}

	createTaskRepository = func(task *models.Task) (*models.Task, error_utils.MessageErr) {
		return &models.Task{
			ID:        1,
			Summary:   "This is a summary test",
			UserID:    1,
			CreatedAt: tm,
			UpdatedAt: tm,
		}, nil
	}

	jsonBody := `{"summary": "This is a summary test"}`
	r := gin.Default()
	req, errRequest := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(jsonBody))
	req.Header = map[string][]string{
		"content-type":  {"application/json"},
		"Authorization": {technician_token},
	}

	if errRequest != nil {
		t.Errorf("this is the error: %v\n", errRequest)
	}

	rr := httptest.NewRecorder()
	r.POST("/tasks", handlerCreateTask)
	r.ServeHTTP(rr, req)

	var task models.Task
	err := json.Unmarshal(rr.Body.Bytes(), &task)
	assert.Nil(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Equal(t, uint64(1), task.ID)
	assert.Equal(t, "This is a summary test", task.Summary)
	assert.NotNil(t, task.CreatedAt)
	assert.NotNil(t, task.UpdatedAt)
}

func TestCreateTask_WrongJSONFormat(t *testing.T) {
	config.SECRETKEY = "mySecretK3y"
	repositories.TaskRepo = &taskRepoMock{}

	jsonBody := `{"summary": }`
	r := gin.Default()
	req, errRequest := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(jsonBody))
	req.Header = map[string][]string{
		"content-type":  {"application/json"},
		"Authorization": {technician_token},
	}

	if errRequest != nil {
		t.Errorf("this is the error: %v\n", errRequest)
	}

	rr := httptest.NewRecorder()
	r.POST("/tasks", handlerCreateTask)
	r.ServeHTTP(rr, req)

	apiErr, err := error_utils.NewApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.Equal(t, http.StatusUnprocessableEntity, apiErr.Status())
	assert.Equal(t, "it's not possible to convert the JSON into an object", apiErr.Message())
	assert.Equal(t, "invalid_request", apiErr.Error())
}

func TestCreateTask_WithoutToken(t *testing.T) {
	config.SECRETKEY = "mySecretK3y"
	repositories.TaskRepo = &taskRepoMock{}

	jsonBody := `{"summary": ""}`
	r := gin.Default()
	req, errRequest := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(jsonBody))
	req.Header = map[string][]string{
		"content-type":  {"application/json"},
		"Authorization": {""},
	}

	if errRequest != nil {
		t.Errorf("this is the error: %v\n", errRequest)
	}

	rr := httptest.NewRecorder()
	r.POST("/tasks", handlerCreateTask)
	r.ServeHTTP(rr, req)

	apiErr, err := error_utils.NewApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.Equal(t, http.StatusUnauthorized, apiErr.Status())
	assert.Equal(t, "token contains an invalid number of segments", apiErr.Message())
	assert.Equal(t, "unauthorized", apiErr.Error())
}

func TestCreateTask_WithoutSummary(t *testing.T) {
	config.SECRETKEY = "mySecretK3y"
	repositories.TaskRepo = &taskRepoMock{}

	jsonBody := `{"summary": ""}`
	r := gin.Default()
	req, errRequest := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(jsonBody))
	req.Header = map[string][]string{
		"content-type":  {"application/json"},
		"Authorization": {technician_token},
	}

	if errRequest != nil {
		t.Errorf("this is the error: %v\n", errRequest)
	}

	rr := httptest.NewRecorder()
	r.POST("/tasks", handlerCreateTask)
	r.ServeHTTP(rr, req)

	apiErr, err := error_utils.NewApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.Equal(t, http.StatusBadRequest, apiErr.Status())
	assert.Equal(t, "the field summary is required can't be empty", apiErr.Message())
	assert.Equal(t, "bad_request", apiErr.Error())
}

func TestCreateTask_Over2500Characters(t *testing.T) {
	config.SECRETKEY = "mySecretK3y"
	repositories.TaskRepo = &taskRepoMock{}

	longString := strings.Repeat("#", 2501)

	jsonBody := fmt.Sprintf(`{"summary": "%s"}`, longString)
	r := gin.Default()
	req, errRequest := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(jsonBody))
	req.Header = map[string][]string{
		"content-type":  {"application/json"},
		"Authorization": {technician_token},
	}

	if errRequest != nil {
		t.Errorf("this is the error: %v\n", errRequest)
	}

	rr := httptest.NewRecorder()
	r.POST("/tasks", handlerCreateTask)
	r.ServeHTTP(rr, req)

	apiErr, err := error_utils.NewApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.Equal(t, http.StatusBadRequest, apiErr.Status())
	assert.Equal(t, "the summary is too long need to be less or equal to 2500 characters", apiErr.Message())
	assert.Equal(t, "bad_request", apiErr.Error())
}

func TestCreateTask_WrongPermission(t *testing.T) {
	config.SECRETKEY = "mySecretK3y"
	repositories.TaskRepo = &taskRepoMock{}

	jsonBody := `{"summary": "This is a summary test"}`
	r := gin.Default()
	req, errRequest := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString(jsonBody))
	req.Header = map[string][]string{
		"content-type":  {"application/json"},
		"Authorization": {manager_token},
	}

	if errRequest != nil {
		t.Errorf("this is the error: %v\n", errRequest)
	}

	rr := httptest.NewRecorder()
	r.POST("/tasks", handlerCreateTask)
	r.ServeHTTP(rr, req)

	apiErr, err := error_utils.NewApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.Equal(t, http.StatusForbidden, apiErr.Status())
	assert.Equal(t, "The user doesn't have the right permission to create a task", apiErr.Message())
	assert.Equal(t, "forbidden", apiErr.Error())
}

func TestUpdateTask_Success(t *testing.T) {
	config.SECRETKEY = "mySecretK3y"
	repositories.TaskRepo = &taskRepoMock{}

	getTaskByIdRepository = func(id uint64) (*models.Task, error_utils.MessageErr) {
		return &models.Task{
			ID:        1,
			Summary:   "This is a summary test",
			UserID:    1,
			CreatedAt: tm,
			UpdatedAt: tm,
		}, nil
	}

	updateTaskReposiroty = func(task *models.Task) (*models.Task, error_utils.MessageErr) {
		return &models.Task{
			ID:        1,
			Summary:   "This is a summary test updated",
			UserID:    1,
			CreatedAt: tm,
			UpdatedAt: tm,
		}, nil
	}

	jsonBody := `{"id": 2, "summary": "This is a summary test", "userId": 1}`
	r := gin.Default()
	req, errRequest := http.NewRequest(http.MethodPut, "/tasks/1", bytes.NewBufferString(jsonBody))
	req.Header = map[string][]string{
		"content-type":  {"application/json"},
		"Authorization": {technician_token},
	}

	if errRequest != nil {
		t.Errorf("this is the error: %v\n", errRequest)
	}

	rr := httptest.NewRecorder()
	r.PUT("/tasks/:id", handlerUpdateTask)
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestUpdateTask_WrongJSONFormat(t *testing.T) {
	config.SECRETKEY = "mySecretK3y"
	repositories.TaskRepo = &taskRepoMock{}

	getTaskByIdRepository = func(id uint64) (*models.Task, error_utils.MessageErr) {
		return &models.Task{
			ID:        1,
			Summary:   "This is a summary test",
			UserID:    1,
			CreatedAt: tm,
			UpdatedAt: tm,
		}, nil
	}

	jsonBody := `{"id": "2", "summary": "This is a summary test", "userId": 1}`
	r := gin.Default()
	req, errRequest := http.NewRequest(http.MethodPut, "/tasks/1", bytes.NewBufferString(jsonBody))
	req.Header = map[string][]string{
		"content-type":  {"application/json"},
		"Authorization": {technician_token},
	}

	if errRequest != nil {
		t.Errorf("this is the error: %v\n", errRequest)
	}

	rr := httptest.NewRecorder()
	r.PUT("/tasks/:id", handlerUpdateTask)
	r.ServeHTTP(rr, req)

	apiErr, err := error_utils.NewApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.Equal(t, http.StatusUnprocessableEntity, apiErr.Status())
	assert.Equal(t, "it's not possible to convert the JSON into an object", apiErr.Message())
	assert.Equal(t, "invalid_request", apiErr.Error())
}

func TestUpdateTask_WithoutToken(t *testing.T) {
	config.SECRETKEY = "mySecretK3y"
	repositories.TaskRepo = &taskRepoMock{}

	jsonBody := `{"id": 2, "summary": "This is a summary test", "userId": 1}`
	r := gin.Default()
	req, errRequest := http.NewRequest(http.MethodPut, "/tasks/1", bytes.NewBufferString(jsonBody))
	req.Header = map[string][]string{
		"content-type":  {"application/json"},
		"Authorization": {""},
	}

	if errRequest != nil {
		t.Errorf("this is the error: %v\n", errRequest)
	}

	rr := httptest.NewRecorder()
	r.PUT("/tasks/:id", handlerUpdateTask)
	r.ServeHTTP(rr, req)

	apiErr, err := error_utils.NewApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.Equal(t, http.StatusUnauthorized, apiErr.Status())
	assert.Equal(t, "token contains an invalid number of segments", apiErr.Message())
	assert.Equal(t, "unauthorized", apiErr.Error())
}

func TestUpdateTask_WithoutSummary(t *testing.T) {
	config.SECRETKEY = "mySecretK3y"
	repositories.TaskRepo = &taskRepoMock{}

	getTaskByIdRepository = func(id uint64) (*models.Task, error_utils.MessageErr) {
		return &models.Task{
			ID:        1,
			Summary:   "This is a summary test",
			UserID:    1,
			CreatedAt: tm,
			UpdatedAt: tm,
		}, nil
	}

	jsonBody := `{"id": 2, "summary": "", "userId": 1}`
	r := gin.Default()
	req, errRequest := http.NewRequest(http.MethodPut, "/tasks/1", bytes.NewBufferString(jsonBody))
	req.Header = map[string][]string{
		"content-type":  {"application/json"},
		"Authorization": {technician_token},
	}

	if errRequest != nil {
		t.Errorf("this is the error: %v\n", errRequest)
	}

	rr := httptest.NewRecorder()
	r.PUT("/tasks/:id", handlerUpdateTask)
	r.ServeHTTP(rr, req)

	apiErr, err := error_utils.NewApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.Equal(t, http.StatusBadRequest, apiErr.Status())
	assert.Equal(t, "the field summary is required can't be empty", apiErr.Message())
	assert.Equal(t, "bad_request", apiErr.Error())
}

func TestUpdateTask_Over2500Characters(t *testing.T) {
	config.SECRETKEY = "mySecretK3y"
	repositories.TaskRepo = &taskRepoMock{}

	getTaskByIdRepository = func(id uint64) (*models.Task, error_utils.MessageErr) {
		return &models.Task{
			ID:        1,
			Summary:   "This is a summary test",
			UserID:    1,
			CreatedAt: tm,
			UpdatedAt: tm,
		}, nil
	}

	longString := strings.Repeat("#", 2501)

	jsonBody := fmt.Sprintf(`{"id": 2, "summary": "%s", "userId": 1}`, longString)
	r := gin.Default()
	req, errRequest := http.NewRequest(http.MethodPut, "/tasks/1", bytes.NewBufferString(jsonBody))
	req.Header = map[string][]string{
		"content-type":  {"application/json"},
		"Authorization": {technician_token},
	}

	if errRequest != nil {
		t.Errorf("this is the error: %v\n", errRequest)
	}

	rr := httptest.NewRecorder()
	r.PUT("/tasks/:id", handlerUpdateTask)
	r.ServeHTTP(rr, req)

	apiErr, err := error_utils.NewApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.Equal(t, http.StatusBadRequest, apiErr.Status())
	assert.Equal(t, "the summary is too long need to be less or equal to 2500 characters", apiErr.Message())
	assert.Equal(t, "bad_request", apiErr.Error())
}

func TestUpdateTask_WrongPermission(t *testing.T) {
	config.SECRETKEY = "mySecretK3y"
	repositories.TaskRepo = &taskRepoMock{}

	longString := strings.Repeat("#", 2501)

	jsonBody := fmt.Sprintf(`{"summary": "%s"}`, longString)
	r := gin.Default()
	req, errRequest := http.NewRequest(http.MethodPut, "/tasks/1", bytes.NewBufferString(jsonBody))
	req.Header = map[string][]string{
		"content-type":  {"application/json"},
		"Authorization": {manager_token},
	}

	if errRequest != nil {
		t.Errorf("this is the error: %v\n", errRequest)
	}

	rr := httptest.NewRecorder()
	r.PUT("/tasks/:id", handlerUpdateTask)
	r.ServeHTTP(rr, req)

	apiErr, err := error_utils.NewApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.Equal(t, http.StatusForbidden, apiErr.Status())
	assert.Equal(t, "The user doesn't have the right permission to update a task", apiErr.Message())
	assert.Equal(t, "forbidden", apiErr.Error())
}

func TestUpdateTask_InvalidID(t *testing.T) {
	config.SECRETKEY = "mySecretK3y"
	repositories.TaskRepo = &taskRepoMock{}

	jsonBody := `{"id": 1, "summary": "This is a summary test", "userId: 1}`
	r := gin.Default()
	req, errRequest := http.NewRequest(http.MethodPut, "/tasks/abs", bytes.NewBufferString(jsonBody))
	req.Header = map[string][]string{
		"content-type":  {"application/json"},
		"Authorization": {technician_token},
	}

	if errRequest != nil {
		t.Errorf("this is the error: %v\n", errRequest)
	}

	rr := httptest.NewRecorder()
	r.PUT("/tasks/:id", handlerUpdateTask)
	r.ServeHTTP(rr, req)

	apiErr, err := error_utils.NewApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.Equal(t, http.StatusBadRequest, apiErr.Status())
	assert.Equal(t, "not possible to convert abs into a number", apiErr.Message())
	assert.Equal(t, "bad_request", apiErr.Error())
}

func TestUpdateTask_DifferentUser(t *testing.T) {
	config.SECRETKEY = "mySecretK3y"
	repositories.TaskRepo = &taskRepoMock{}

	getTaskByIdRepository = func(id uint64) (*models.Task, error_utils.MessageErr) {
		return &models.Task{
			ID:        1,
			Summary:   "This is a summary test",
			UserID:    2,
			CreatedAt: tm,
			UpdatedAt: tm,
		}, nil
	}

	jsonBody := `{"id": 2, "summary": "This is a summary test", "userId": 2}`
	r := gin.Default()
	req, errRequest := http.NewRequest(http.MethodPut, "/tasks/1", bytes.NewBufferString(jsonBody))
	req.Header = map[string][]string{
		"content-type":  {"application/json"},
		"Authorization": {technician_token},
	}

	if errRequest != nil {
		t.Errorf("this is the error: %v\n", errRequest)
	}

	rr := httptest.NewRecorder()
	r.PUT("/tasks/:id", handlerUpdateTask)
	r.ServeHTTP(rr, req)

	apiErr, err := error_utils.NewApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.Equal(t, http.StatusForbidden, apiErr.Status())
	assert.Equal(t, "Not possible to update a task that does not belong to you", apiErr.Message())
	assert.Equal(t, "forbidden", apiErr.Error())
}

func TestUpdateTask_NotFoundTask(t *testing.T) {
	config.SECRETKEY = "mySecretK3y"
	repositories.TaskRepo = &taskRepoMock{}

	getTaskByIdRepository = func(id uint64) (*models.Task, error_utils.MessageErr) {
		return nil, error_utils.NewNotFoundError("no record matching given the identification")
	}

	jsonBody := `{"id": 1, "summary": "This is a summary test", "userId": 2}`
	r := gin.Default()
	req, errRequest := http.NewRequest(http.MethodPut, "/tasks/1", bytes.NewBufferString(jsonBody))
	req.Header = map[string][]string{
		"content-type":  {"application/json"},
		"Authorization": {technician_token},
	}

	if errRequest != nil {
		t.Errorf("this is the error: %v\n", errRequest)
	}

	rr := httptest.NewRecorder()
	r.PUT("/tasks/:id", handlerUpdateTask)
	r.ServeHTTP(rr, req)

	apiErr, err := error_utils.NewApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.Equal(t, http.StatusNotFound, apiErr.Status())
	assert.Equal(t, "no record matching given the identification", apiErr.Message())
	assert.Equal(t, "not_found", apiErr.Error())
}

func TestGetTask_Success(t *testing.T) {
	config.SECRETKEY = "mySecretK3y"
	repositories.TaskRepo = &taskRepoMock{}

	getTaskByIdRepository = func(id uint64) (*models.Task, error_utils.MessageErr) {
		return &models.Task{
			ID:        1,
			Summary:   "This is a summary test",
			UserID:    1,
			CreatedAt: tm,
			UpdatedAt: tm,
		}, nil
	}

	r := gin.Default()
	req, errRequest := http.NewRequest(http.MethodGet, "/tasks/1", nil)
	req.Header = map[string][]string{
		"content-type":  {"application/json"},
		"Authorization": {technician_token},
	}

	if errRequest != nil {
		t.Errorf("this is the error: %v\n", errRequest)
	}

	rr := httptest.NewRecorder()
	r.GET("/tasks/:id", handlerGetTask)
	r.ServeHTTP(rr, req)

	var task models.Task
	err := json.Unmarshal(rr.Body.Bytes(), &task)
	assert.Nil(t, err)
	assert.NotNil(t, task)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, uint64(1), task.ID)
	assert.Equal(t, "This is a summary test", task.Summary)
	assert.NotNil(t, task.CreatedAt)
	assert.NotNil(t, task.UpdatedAt)
}

func TestGetTask_WithoutToken(t *testing.T) {
	config.SECRETKEY = "mySecretK3y"
	repositories.TaskRepo = &taskRepoMock{}

	r := gin.Default()
	req, errRequest := http.NewRequest(http.MethodGet, "/tasks/1", nil)
	req.Header = map[string][]string{
		"content-type":  {"application/json"},
		"Authorization": {""},
	}

	if errRequest != nil {
		t.Errorf("this is the error: %v\n", errRequest)
	}

	rr := httptest.NewRecorder()
	r.GET("/tasks/:id", handlerGetTask)
	r.ServeHTTP(rr, req)

	apiErr, err := error_utils.NewApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.Equal(t, http.StatusUnauthorized, apiErr.Status())
	assert.Equal(t, "token contains an invalid number of segments", apiErr.Message())
	assert.Equal(t, "unauthorized", apiErr.Error())
}

func TestGetTask_WrongPermission(t *testing.T) {
	config.SECRETKEY = "mySecretK3y"
	repositories.TaskRepo = &taskRepoMock{}

	r := gin.Default()
	req, errRequest := http.NewRequest(http.MethodGet, "/tasks/1", nil)
	req.Header = map[string][]string{
		"content-type":  {"application/json"},
		"Authorization": {manager_token},
	}

	if errRequest != nil {
		t.Errorf("this is the error: %v\n", errRequest)
	}

	rr := httptest.NewRecorder()
	r.GET("/tasks/:id", handlerGetTask)
	r.ServeHTTP(rr, req)

	apiErr, err := error_utils.NewApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.Equal(t, http.StatusForbidden, apiErr.Status())
	assert.Equal(t, "The user doesn't have the right permission to get a specific task", apiErr.Message())
	assert.Equal(t, "forbidden", apiErr.Error())
}

func TestGetTask_InvalidID(t *testing.T) {
	config.SECRETKEY = "mySecretK3y"
	repositories.TaskRepo = &taskRepoMock{}

	r := gin.Default()
	req, errRequest := http.NewRequest(http.MethodGet, "/tasks/abc", nil)
	req.Header = map[string][]string{
		"content-type":  {"application/json"},
		"Authorization": {technician_token},
	}

	if errRequest != nil {
		t.Errorf("this is the error: %v\n", errRequest)
	}

	rr := httptest.NewRecorder()
	r.GET("/tasks/:id", handlerGetTask)
	r.ServeHTTP(rr, req)

	apiErr, err := error_utils.NewApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.Equal(t, http.StatusBadRequest, apiErr.Status())
	assert.Equal(t, "not possible to convert abc into a number", apiErr.Message())
	assert.Equal(t, "bad_request", apiErr.Error())
}

func TestGetTask_DifferentUser(t *testing.T) {
	config.SECRETKEY = "mySecretK3y"
	repositories.TaskRepo = &taskRepoMock{}

	getTaskByIdRepository = func(id uint64) (*models.Task, error_utils.MessageErr) {
		return &models.Task{
			ID:        1,
			Summary:   "This is a summary test",
			UserID:    2,
			CreatedAt: tm,
			UpdatedAt: tm,
		}, nil
	}

	r := gin.Default()
	req, errRequest := http.NewRequest(http.MethodGet, "/tasks/1", nil)
	req.Header = map[string][]string{
		"content-type":  {"application/json"},
		"Authorization": {technician_token},
	}

	if errRequest != nil {
		t.Errorf("this is the error: %v\n", errRequest)
	}

	rr := httptest.NewRecorder()
	r.GET("/tasks/:id", handlerGetTask)
	r.ServeHTTP(rr, req)

	apiErr, err := error_utils.NewApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.Equal(t, http.StatusForbidden, apiErr.Status())
	assert.Equal(t, "Not possible to see a task that does not belong to you", apiErr.Message())
	assert.Equal(t, "forbidden", apiErr.Error())
}

func TestGetTask_NotFoundTask(t *testing.T) {
	config.SECRETKEY = "mySecretK3y"
	repositories.TaskRepo = &taskRepoMock{}

	getTaskByIdRepository = func(id uint64) (*models.Task, error_utils.MessageErr) {
		return nil, error_utils.NewNotFoundError("no record matching given the identification")
	}

	r := gin.Default()
	req, errRequest := http.NewRequest(http.MethodGet, "/tasks/1", nil)
	req.Header = map[string][]string{
		"content-type":  {"application/json"},
		"Authorization": {technician_token},
	}

	if errRequest != nil {
		t.Errorf("this is the error: %v\n", errRequest)
	}

	rr := httptest.NewRecorder()
	r.GET("/tasks/:id", handlerGetTask)
	r.ServeHTTP(rr, req)

	apiErr, err := error_utils.NewApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.Equal(t, http.StatusNotFound, apiErr.Status())
	assert.Equal(t, "no record matching given the identification", apiErr.Message())
	assert.Equal(t, "not_found", apiErr.Error())
}

func TestGetTasksByUser_Success(t *testing.T) {
	config.SECRETKEY = "mySecretK3y"
	repositories.TaskRepo = &taskRepoMock{}

	getTasksByUserRepository = func(userId uint64) []models.Task {
		return []models.Task{
			{
				ID:        1,
				Summary:   "This is a summary test",
				UserID:    1,
				CreatedAt: tm,
				UpdatedAt: tm,
			},
			{
				ID:        2,
				Summary:   "This is a summary test 2",
				UserID:    1,
				CreatedAt: tm,
				UpdatedAt: tm,
			},
		}
	}

	r := gin.Default()
	req, errRequest := http.NewRequest(http.MethodGet, "/user_tasks", nil)
	req.Header = map[string][]string{
		"content-type":  {"application/json"},
		"Authorization": {technician_token},
	}

	if errRequest != nil {
		t.Errorf("this is the error: %v\n", errRequest)
	}

	rr := httptest.NewRecorder()
	r.GET("/user_tasks", handlerGetTasksByUser)
	r.ServeHTTP(rr, req)

	var tasks []models.Task
	err := json.Unmarshal(rr.Body.Bytes(), &tasks)
	assert.Nil(t, err)
	assert.NotNil(t, tasks)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, 2, len(tasks))
}

func TestGetTasksByUser_WithoutToken(t *testing.T) {
	config.SECRETKEY = "mySecretK3y"
	repositories.TaskRepo = &taskRepoMock{}

	r := gin.Default()
	req, errRequest := http.NewRequest(http.MethodGet, "/user_tasks", nil)
	req.Header = map[string][]string{
		"content-type":  {"application/json"},
		"Authorization": {""},
	}

	if errRequest != nil {
		t.Errorf("this is the error: %v\n", errRequest)
	}

	rr := httptest.NewRecorder()
	r.GET("/user_tasks", handlerGetTasksByUser)
	r.ServeHTTP(rr, req)

	apiErr, err := error_utils.NewApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.Equal(t, http.StatusUnauthorized, apiErr.Status())
	assert.Equal(t, "token contains an invalid number of segments", apiErr.Message())
	assert.Equal(t, "unauthorized", apiErr.Error())
}

func TestGetTasksByUser_WrongPermission(t *testing.T) {
	config.SECRETKEY = "mySecretK3y"
	repositories.TaskRepo = &taskRepoMock{}

	r := gin.Default()
	req, errRequest := http.NewRequest(http.MethodGet, "/user_tasks", nil)
	req.Header = map[string][]string{
		"content-type":  {"application/json"},
		"Authorization": {manager_token},
	}

	if errRequest != nil {
		t.Errorf("this is the error: %v\n", errRequest)
	}

	rr := httptest.NewRecorder()
	r.GET("/user_tasks", handlerGetTasksByUser)
	r.ServeHTTP(rr, req)

	apiErr, err := error_utils.NewApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.Equal(t, http.StatusForbidden, apiErr.Status())
	assert.Equal(t, "The user doesn't have the right permission to list his tasks", apiErr.Message())
	assert.Equal(t, "forbidden", apiErr.Error())
}

func TestGetTasksByUser_EmptyList(t *testing.T) {
	config.SECRETKEY = "mySecretK3y"
	repositories.TaskRepo = &taskRepoMock{}

	getTasksByUserRepository = func(userId uint64) []models.Task {
		return []models.Task{}
	}

	r := gin.Default()
	req, errRequest := http.NewRequest(http.MethodGet, "/user_tasks", nil)
	req.Header = map[string][]string{
		"content-type":  {"application/json"},
		"Authorization": {technician_token},
	}

	if errRequest != nil {
		t.Errorf("this is the error: %v\n", errRequest)
	}

	rr := httptest.NewRecorder()
	r.GET("/user_tasks", handlerGetTasksByUser)
	r.ServeHTTP(rr, req)

	var tasks []models.Task
	err := json.Unmarshal(rr.Body.Bytes(), &tasks)
	assert.Nil(t, err)
	assert.NotNil(t, tasks)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, 0, len(tasks))
}

func TestGetTasks_Success(t *testing.T) {
	config.SECRETKEY = "mySecretK3y"
	repositories.TaskRepo = &taskRepoMock{}

	getTasksRepository = func() []models.Task {
		return []models.Task{
			{
				ID:        1,
				Summary:   "This is a summary test",
				UserID:    1,
				CreatedAt: tm,
				UpdatedAt: tm,
			},
			{
				ID:        2,
				Summary:   "This is a summary test 2",
				UserID:    1,
				CreatedAt: tm,
				UpdatedAt: tm,
			},
			{
				ID:        3,
				Summary:   "This is a summary test 3",
				UserID:    2,
				CreatedAt: tm,
				UpdatedAt: tm,
			},
		}
	}

	r := gin.Default()
	req, errRequest := http.NewRequest(http.MethodGet, "/tasks", nil)
	req.Header = map[string][]string{
		"content-type":  {"application/json"},
		"Authorization": {manager_token},
	}

	if errRequest != nil {
		t.Errorf("this is the error: %v\n", errRequest)
	}

	rr := httptest.NewRecorder()
	r.GET("/tasks", handlerGetAllTasks)
	r.ServeHTTP(rr, req)

	var tasks []models.Task
	err := json.Unmarshal(rr.Body.Bytes(), &tasks)
	assert.Nil(t, err)
	assert.NotNil(t, tasks)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, 3, len(tasks))
}

func TestGetTasks_WrongPermission(t *testing.T) {
	config.SECRETKEY = "mySecretK3y"
	repositories.TaskRepo = &taskRepoMock{}

	r := gin.Default()
	req, errRequest := http.NewRequest(http.MethodGet, "/tasks", nil)
	req.Header = map[string][]string{
		"content-type":  {"application/json"},
		"Authorization": {technician_token},
	}

	if errRequest != nil {
		t.Errorf("this is the error: %v\n", errRequest)
	}

	rr := httptest.NewRecorder()
	r.GET("/tasks", handlerGetAllTasks)
	r.ServeHTTP(rr, req)

	apiErr, err := error_utils.NewApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.Equal(t, http.StatusForbidden, apiErr.Status())
	assert.Equal(t, "The user doesn't have the right permission to get all tasks", apiErr.Message())
	assert.Equal(t, "forbidden", apiErr.Error())
}

func TestGetTasks_EmptyList(t *testing.T) {
	config.SECRETKEY = "mySecretK3y"
	repositories.TaskRepo = &taskRepoMock{}

	getTasksRepository = func() []models.Task {
		return []models.Task{}
	}

	r := gin.Default()
	req, errRequest := http.NewRequest(http.MethodGet, "/tasks", nil)
	req.Header = map[string][]string{
		"content-type":  {"application/json"},
		"Authorization": {manager_token},
	}

	if errRequest != nil {
		t.Errorf("this is the error: %v\n", errRequest)
	}

	rr := httptest.NewRecorder()
	r.GET("/tasks", handlerGetAllTasks)
	r.ServeHTTP(rr, req)

	var tasks []models.Task
	err := json.Unmarshal(rr.Body.Bytes(), &tasks)
	assert.Nil(t, err)
	assert.NotNil(t, tasks)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, 0, len(tasks))
}

func TestDeleteTask_Success(t *testing.T) {
	config.SECRETKEY = "mySecretK3y"
	repositories.TaskRepo = &taskRepoMock{}

	deleteTasksRepository = func(id uint64) error_utils.MessageErr {
		return nil
	}

	r := gin.Default()
	req, errRequest := http.NewRequest(http.MethodDelete, "/tasks/1", nil)
	req.Header = map[string][]string{
		"content-type":  {"application/json"},
		"Authorization": {manager_token},
	}

	if errRequest != nil {
		t.Errorf("this is the error: %v\n", errRequest)
	}

	rr := httptest.NewRecorder()
	r.DELETE("/tasks/:id", handlerDeleteTasks)
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestDeleteTask_WrongPermission(t *testing.T) {
	config.SECRETKEY = "mySecretK3y"
	repositories.TaskRepo = &taskRepoMock{}

	r := gin.Default()
	req, errRequest := http.NewRequest(http.MethodDelete, "/tasks/1", nil)
	req.Header = map[string][]string{
		"content-type":  {"application/json"},
		"Authorization": {technician_token},
	}

	if errRequest != nil {
		t.Errorf("this is the error: %v\n", errRequest)
	}

	rr := httptest.NewRecorder()
	r.DELETE("/tasks/:id", handlerDeleteTasks)
	r.ServeHTTP(rr, req)

	apiErr, err := error_utils.NewApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.Equal(t, http.StatusForbidden, apiErr.Status())
	assert.Equal(t, "The user doesn't have the right permission to delete a specific task", apiErr.Message())
	assert.Equal(t, "forbidden", apiErr.Error())
}

func TestDeleteTask_InvalidID(t *testing.T) {
	config.SECRETKEY = "mySecretK3y"
	repositories.TaskRepo = &taskRepoMock{}

	r := gin.Default()
	req, errRequest := http.NewRequest(http.MethodDelete, "/tasks/abc", nil)
	req.Header = map[string][]string{
		"content-type":  {"application/json"},
		"Authorization": {manager_token},
	}

	if errRequest != nil {
		t.Errorf("this is the error: %v\n", errRequest)
	}

	rr := httptest.NewRecorder()
	r.DELETE("/tasks/:id", handlerDeleteTasks)
	r.ServeHTTP(rr, req)

	apiErr, err := error_utils.NewApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.Equal(t, http.StatusBadRequest, apiErr.Status())
	assert.Equal(t, "not possible to convert abc into a number", apiErr.Message())
	assert.Equal(t, "bad_request", apiErr.Error())
}

func TestDeleteTask_NotFoundTask(t *testing.T) {
	config.SECRETKEY = "mySecretK3y"
	repositories.TaskRepo = &taskRepoMock{}

	deleteTasksRepository = func(id uint64) error_utils.MessageErr {
		return error_utils.NewNotFoundError("no record matching given the identification")
	}

	r := gin.Default()
	req, errRequest := http.NewRequest(http.MethodDelete, "/tasks/1", nil)
	req.Header = map[string][]string{
		"content-type":  {"application/json"},
		"Authorization": {manager_token},
	}

	if errRequest != nil {
		t.Errorf("this is the error: %v\n", errRequest)
	}

	rr := httptest.NewRecorder()
	r.DELETE("/tasks/:id", handlerDeleteTasks)
	r.ServeHTTP(rr, req)

	apiErr, err := error_utils.NewApiErrFromBytes(rr.Body.Bytes())
	assert.Nil(t, err)
	assert.NotNil(t, apiErr)
	assert.Equal(t, http.StatusNotFound, apiErr.Status())
	assert.Equal(t, "no record matching given the identification", apiErr.Message())
	assert.Equal(t, "not_found", apiErr.Error())
}
