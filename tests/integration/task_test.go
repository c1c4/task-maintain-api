package integration

import (
	"api/app/models"
	"api/app/utils/error_utils"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func (s *SuiteTest) TestCreateTask_Success() {
	s.seedOneUserTech()
	jsonBody := `{"summary": "This is a summary test"}`
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/tasks", baseURL), bytes.NewBufferString(jsonBody))
	s.NoError(err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", technician_token)

	client := http.Client{}
	resp, err := client.Do(req)
	s.NoError(err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	s.NoError(err)

	var task models.Task
	errUnmarshal := json.Unmarshal(byteBody, &task)
	s.Nil(errUnmarshal)
	s.NotNil(task)
	s.Equal(http.StatusCreated, resp.StatusCode)
	s.Equal(uint64(1), task.ID)
	s.Equal("This is a summary test", task.Summary)
	s.NotNil(task.CreatedAt)
	s.NotNil(task.UpdatedAt)
}

func (s *SuiteTest) TestCreateTask_WrongJSONFormat() {
	jsonBody := `{"summary": }`
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/tasks", baseURL), bytes.NewBufferString(jsonBody))
	s.NoError(err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", technician_token)

	client := http.Client{}
	resp, err := client.Do(req)
	s.NoError(err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	s.NoError(err)

	apiErr, err := error_utils.NewApiErrFromBytes(byteBody)
	s.Nil(err)
	s.NotNil(apiErr)
	s.Equal(http.StatusUnprocessableEntity, apiErr.Status())
	s.Equal("it's not possible to convert the JSON into an object", apiErr.Message())
	s.Equal("invalid_request", apiErr.Error())
}

func (s *SuiteTest) TestCreateTask_WithoutToken() {
	jsonBody := `{"summary": ""}`
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/tasks", baseURL), bytes.NewBufferString(jsonBody))
	s.NoError(err)

	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	s.NoError(err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	s.NoError(err)

	apiErr, err := error_utils.NewApiErrFromBytes(byteBody)
	s.Nil(err)
	s.NotNil(apiErr)
	s.Equal(http.StatusUnauthorized, apiErr.Status())
	s.Equal("token contains an invalid number of segments", apiErr.Message())
	s.Equal("unauthorized", apiErr.Error())
}

func (s *SuiteTest) TestCreateTask_WithoutSummary() {
	jsonBody := `{"summary": ""}`
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/tasks", baseURL), bytes.NewBufferString(jsonBody))
	s.NoError(err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", technician_token)

	client := http.Client{}
	resp, err := client.Do(req)
	s.NoError(err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	s.NoError(err)

	apiErr, err := error_utils.NewApiErrFromBytes(byteBody)
	s.Nil(err)
	s.NotNil(apiErr)
	s.Equal(http.StatusBadRequest, apiErr.Status())
	s.Equal("the field summary is required can't be empty", apiErr.Message())
	s.Equal("bad_request", apiErr.Error())
}

func (s *SuiteTest) TestCreateTask_Over2500Characters() {
	longString := strings.Repeat("#", 2501)
	jsonBody := fmt.Sprintf(`{"summary": "%s"}`, longString)
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/tasks", baseURL), bytes.NewBufferString(jsonBody))
	s.NoError(err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", technician_token)

	client := http.Client{}
	resp, err := client.Do(req)
	s.NoError(err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	s.NoError(err)

	apiErr, err := error_utils.NewApiErrFromBytes(byteBody)
	s.Nil(err)
	s.NotNil(apiErr)
	s.Equal(http.StatusBadRequest, apiErr.Status())
	s.Equal("the summary is too long need to be less or equal to 2500 characters", apiErr.Message())
	s.Equal("bad_request", apiErr.Error())
}

func (s *SuiteTest) TestCreateTask_WrongPermission() {
	jsonBody := `{"summary": "This is a summary test"}`
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/tasks", baseURL), bytes.NewBufferString(jsonBody))
	s.NoError(err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", manager_token)

	client := http.Client{}
	resp, err := client.Do(req)
	s.NoError(err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	s.NoError(err)

	apiErr, err := error_utils.NewApiErrFromBytes(byteBody)
	s.Nil(err)
	s.NotNil(apiErr)
	s.Equal(http.StatusForbidden, apiErr.Status())
	s.Equal("The user doesn't have the right permission to create a task", apiErr.Message())
	s.Equal("forbidden", apiErr.Error())
}

func (s *SuiteTest) TestUpdateTask_Success() {
	s.seedOneUserTech()
	s.seedOneTask()
	jsonBody := `{"id": 1, "summary": "This is a summary test updated", "userId": 1}`
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/tasks/1", baseURL), bytes.NewBufferString(jsonBody))
	s.NoError(err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", technician_token)

	client := http.Client{}
	resp, err := client.Do(req)
	s.NoError(err)

	s.Equal(http.StatusNoContent, resp.StatusCode)
}

func (s *SuiteTest) TestUpdateTask_WrongJSONFormat() {
	s.seedOneUserTech()
	s.seedOneTask()
	jsonBody := `{"id": "2", "summary": "This is a summary test", "userId": 1}`
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/tasks/1", baseURL), bytes.NewBufferString(jsonBody))
	s.NoError(err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", technician_token)

	client := http.Client{}
	resp, err := client.Do(req)
	s.NoError(err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	s.NoError(err)

	apiErr, err := error_utils.NewApiErrFromBytes(byteBody)
	s.Nil(err)
	s.NotNil(apiErr)
	s.Equal(http.StatusUnprocessableEntity, apiErr.Status())
	s.Equal("it's not possible to convert the JSON into an object", apiErr.Message())
	s.Equal("invalid_request", apiErr.Error())
}

func (s *SuiteTest) TestUpdateTask_WithoutToken() {
	jsonBody := `{"id": 2, "summary": "This is a summary test", "userId": 1}`
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/tasks/1", baseURL), bytes.NewBufferString(jsonBody))
	s.NoError(err)

	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	s.NoError(err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	s.NoError(err)

	apiErr, err := error_utils.NewApiErrFromBytes(byteBody)
	s.Nil(err)
	s.NotNil(apiErr)
	s.Equal(http.StatusUnauthorized, apiErr.Status())
	s.Equal("token contains an invalid number of segments", apiErr.Message())
	s.Equal("unauthorized", apiErr.Error())
}

func (s *SuiteTest) TestUpdateTask_WithoutSummary() {
	s.seedOneUserTech()
	s.seedOneTask()
	jsonBody := `{"id": 2, "summary": "", "userId": 1}`
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/tasks/1", baseURL), bytes.NewBufferString(jsonBody))
	s.NoError(err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", technician_token)

	client := http.Client{}
	resp, err := client.Do(req)
	s.NoError(err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	s.NoError(err)

	apiErr, err := error_utils.NewApiErrFromBytes(byteBody)
	s.Nil(err)
	s.NotNil(apiErr)
	s.Equal(http.StatusBadRequest, apiErr.Status())
	s.Equal("the field summary is required can't be empty", apiErr.Message())
	s.Equal("bad_request", apiErr.Error())
}

func (s *SuiteTest) TestUpdateTask_Over2500Characters() {
	s.seedOneUserTech()
	s.seedOneTask()
	longString := strings.Repeat("#", 2501)
	jsonBody := fmt.Sprintf(`{"id": 2, "summary": "%s", "userId": 1}`, longString)
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/tasks/1", baseURL), bytes.NewBufferString(jsonBody))
	s.NoError(err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", technician_token)

	client := http.Client{}
	resp, err := client.Do(req)
	s.NoError(err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	s.NoError(err)

	apiErr, err := error_utils.NewApiErrFromBytes(byteBody)
	s.Nil(err)
	s.NotNil(apiErr)
	s.Equal(http.StatusBadRequest, apiErr.Status())
	s.Equal("the summary is too long need to be less or equal to 2500 characters", apiErr.Message())
	s.Equal("bad_request", apiErr.Error())
}

func (s *SuiteTest) TestUpdateTask_WrongPermission() {
	s.seedOneUserManager()
	jsonBody := `{"id": 2, "summary": "This is a summary test", "userId": 1}`
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/tasks/1", baseURL), bytes.NewBufferString(jsonBody))
	s.NoError(err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", manager_token)

	client := http.Client{}
	resp, err := client.Do(req)
	s.NoError(err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	s.NoError(err)

	apiErr, err := error_utils.NewApiErrFromBytes(byteBody)
	s.Nil(err)
	s.NotNil(apiErr)
	s.Equal(http.StatusForbidden, apiErr.Status())
	s.Equal("The user doesn't have the right permission to update a task", apiErr.Message())
	s.Equal("forbidden", apiErr.Error())
}

func (s *SuiteTest) TestUpdateTask_InvalidID() {
	jsonBody := `{"id": 1, "summary": "This is a summary test", "userId: 1}`
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/tasks/abs", baseURL), bytes.NewBufferString(jsonBody))
	s.NoError(err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", technician_token)

	client := http.Client{}
	resp, err := client.Do(req)
	s.NoError(err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	s.NoError(err)

	apiErr, err := error_utils.NewApiErrFromBytes(byteBody)
	s.Nil(err)
	s.NotNil(apiErr)
	s.Equal(http.StatusBadRequest, apiErr.Status())
	s.Equal("not possible to convert abs into a number", apiErr.Message())
	s.Equal("bad_request", apiErr.Error())
}

func (s *SuiteTest) TestUpdateTask_DifferentUser() {
	s.seedOneUserTech()
	s.seedOneUserManager()
	s.seedMultipleTasks()
	jsonBody := `{"id": 3, "summary": "This is a summary test", "userId": 2}`
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/tasks/3", baseURL), bytes.NewBufferString(jsonBody))
	s.NoError(err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", technician_token)

	client := http.Client{}
	resp, err := client.Do(req)
	s.NoError(err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	s.NoError(err)

	apiErr, err := error_utils.NewApiErrFromBytes(byteBody)
	s.Nil(err)
	s.NotNil(apiErr)
	s.Equal(http.StatusForbidden, apiErr.Status())
	s.Equal("Not possible to update a task that does not belong to you", apiErr.Message())
	s.Equal("forbidden", apiErr.Error())
}

func (s *SuiteTest) TestUpdateTask_NotFoundTask() {
	jsonBody := `{"id": 1, "summary": "This is a summary test", "userId": 2}`
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/tasks/99", baseURL), bytes.NewBufferString(jsonBody))
	s.NoError(err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", technician_token)

	client := http.Client{}
	resp, err := client.Do(req)
	s.NoError(err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	s.NoError(err)

	apiErr, err := error_utils.NewApiErrFromBytes(byteBody)
	s.Nil(err)
	s.NotNil(apiErr)
	s.Equal(http.StatusNotFound, apiErr.Status())
	s.Equal("no record matching given the identification", apiErr.Message())
	s.Equal("not_found", apiErr.Error())
}

func (s *SuiteTest) TestGetTask_Success() {
	s.seedOneUserTech()
	s.seedOneTask()
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/tasks/1", baseURL), nil)
	s.NoError(err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", technician_token)

	client := http.Client{}
	resp, err := client.Do(req)
	s.NoError(err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	s.NoError(err)

	var task models.Task
	errUnmarshal := json.Unmarshal(byteBody, &task)
	s.Nil(errUnmarshal)
	s.NotNil(task)
	s.Equal(http.StatusOK, resp.StatusCode)
	s.Equal(uint64(1), task.ID)
	s.Equal("This is a summary test", task.Summary)
	s.NotNil(task.CreatedAt)
	s.NotNil(task.UpdatedAt)
}

func (s *SuiteTest) TestGetTask_WithoutToken() {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/tasks/1", baseURL), nil)
	s.NoError(err)

	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	s.NoError(err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	s.NoError(err)

	apiErr, err := error_utils.NewApiErrFromBytes(byteBody)
	s.Nil(err)
	s.NotNil(apiErr)
	s.Equal(http.StatusUnauthorized, apiErr.Status())
	s.Equal("token contains an invalid number of segments", apiErr.Message())
	s.Equal("unauthorized", apiErr.Error())
}

func (s *SuiteTest) TestGetTask_WrongPermission() {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/tasks/1", baseURL), nil)
	s.NoError(err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", manager_token)

	client := http.Client{}
	resp, err := client.Do(req)
	s.NoError(err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	s.NoError(err)

	apiErr, err := error_utils.NewApiErrFromBytes(byteBody)
	s.Nil(err)
	s.NotNil(apiErr)
	s.Equal(http.StatusForbidden, apiErr.Status())
	s.Equal("The user doesn't have the right permission to get a specific task", apiErr.Message())
	s.Equal("forbidden", apiErr.Error())
}

func (s *SuiteTest) TestGetTask_InvalidID() {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/tasks/abc", baseURL), nil)
	s.NoError(err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", technician_token)

	client := http.Client{}
	resp, err := client.Do(req)
	s.NoError(err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	s.NoError(err)

	apiErr, err := error_utils.NewApiErrFromBytes(byteBody)
	s.Nil(err)
	s.NotNil(apiErr)
	s.Equal(http.StatusBadRequest, apiErr.Status())
	s.Equal("not possible to convert abc into a number", apiErr.Message())
	s.Equal("bad_request", apiErr.Error())
}

func (s *SuiteTest) TestGetTask_DifferentUser() {
	s.seedOneUserTech()
	s.seedOneUserManager()
	s.seedMultipleTasks()
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/tasks/3", baseURL), nil)
	s.NoError(err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", technician_token)

	client := http.Client{}
	resp, err := client.Do(req)
	s.NoError(err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	s.NoError(err)

	apiErr, err := error_utils.NewApiErrFromBytes(byteBody)
	s.Nil(err)
	s.NotNil(apiErr)
	s.Equal(http.StatusForbidden, apiErr.Status())
	s.Equal("Not possible to see a task that does not belong to you", apiErr.Message())
	s.Equal("forbidden", apiErr.Error())
}

func (s *SuiteTest) TestGetTask_NotFoundTask() {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/tasks/99", baseURL), nil)
	s.NoError(err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", technician_token)

	client := http.Client{}
	resp, err := client.Do(req)
	s.NoError(err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	s.NoError(err)

	apiErr, err := error_utils.NewApiErrFromBytes(byteBody)
	s.Nil(err)
	s.NotNil(apiErr)
	s.Equal(http.StatusNotFound, apiErr.Status())
	s.Equal("no record matching given the identification", apiErr.Message())
	s.Equal("not_found", apiErr.Error())
}

func (s *SuiteTest) TestGetTasksByUser_Success() {
	s.seedOneUserTech()
	s.seedOneUserManager()
	s.seedMultipleTasks()
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/user_tasks", baseURL), nil)
	s.NoError(err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", technician_token)

	client := http.Client{}
	resp, err := client.Do(req)
	s.NoError(err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	s.NoError(err)

	var tasks []models.Task
	errUnmarshal := json.Unmarshal(byteBody, &tasks)
	s.Nil(errUnmarshal)
	s.NotNil(tasks)
	s.Equal(http.StatusOK, resp.StatusCode)
	s.Equal(2, len(tasks))
}

func (s *SuiteTest) TestGetTasksByUser_WithoutToken() {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/user_tasks", baseURL), nil)
	s.NoError(err)

	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	s.NoError(err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	s.NoError(err)

	apiErr, err := error_utils.NewApiErrFromBytes(byteBody)
	s.Nil(err)
	s.NotNil(apiErr)
	s.Equal(http.StatusUnauthorized, apiErr.Status())
	s.Equal("token contains an invalid number of segments", apiErr.Message())
	s.Equal("unauthorized", apiErr.Error())
}

func (s *SuiteTest) TestGetTasksByUser_WrongPermission() {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/user_tasks", baseURL), nil)
	s.NoError(err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", manager_token)

	client := http.Client{}
	resp, err := client.Do(req)
	s.NoError(err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	s.NoError(err)

	apiErr, err := error_utils.NewApiErrFromBytes(byteBody)
	s.Nil(err)
	s.NotNil(apiErr)
	s.Equal(http.StatusForbidden, apiErr.Status())
	s.Equal("The user doesn't have the right permission to list his tasks", apiErr.Message())
	s.Equal("forbidden", apiErr.Error())
}

func (s *SuiteTest) TestGetTasksByUser_EmptyList() {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/user_tasks", baseURL), nil)
	s.NoError(err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", technician_token)

	client := http.Client{}
	resp, err := client.Do(req)
	s.NoError(err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	s.NoError(err)

	var tasks []models.Task
	errUnmarshal := json.Unmarshal(byteBody, &tasks)
	s.Nil(errUnmarshal)
	s.NotNil(tasks)
	s.Equal(http.StatusOK, resp.StatusCode)
	s.Equal(0, len(tasks))
}

func (s *SuiteTest) TestGetTasks_Success() {
	s.seedOneUserTech()
	s.seedOneUserManager()
	s.seedMultipleTasks()
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/tasks", baseURL), nil)
	s.NoError(err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", manager_token)

	client := http.Client{}
	resp, err := client.Do(req)
	s.NoError(err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	s.NoError(err)

	var tasks []models.Task
	errUnmarshal := json.Unmarshal(byteBody, &tasks)
	s.Nil(errUnmarshal)
	s.NotNil(tasks)
	s.Equal(http.StatusOK, resp.StatusCode)
	s.Equal(3, len(tasks))
}

func (s *SuiteTest) TestGetTasks_WrongPermission() {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/tasks", baseURL), nil)
	s.NoError(err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", technician_token)

	client := http.Client{}
	resp, err := client.Do(req)
	s.NoError(err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	s.NoError(err)

	apiErr, err := error_utils.NewApiErrFromBytes(byteBody)
	s.Nil(err)
	s.NotNil(apiErr)
	s.Equal(http.StatusForbidden, apiErr.Status())
	s.Equal("The user doesn't have the right permission to get all tasks", apiErr.Message())
	s.Equal("forbidden", apiErr.Error())
}

func (s *SuiteTest) TestGetTasks_EmptyList() {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/tasks", baseURL), nil)
	s.NoError(err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", manager_token)

	client := http.Client{}
	resp, err := client.Do(req)
	s.NoError(err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	s.NoError(err)

	var tasks []models.Task
	errUnmarshal := json.Unmarshal(byteBody, &tasks)
	s.Nil(errUnmarshal)
	s.NotNil(tasks)
	s.Equal(http.StatusOK, resp.StatusCode)
	s.Equal(0, len(tasks))
}

func (s *SuiteTest) TestDeleteTask_Success() {
	s.seedOneUserTech()
	s.seedOneTask()
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/tasks/1", baseURL), nil)
	s.NoError(err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", manager_token)

	client := http.Client{}
	resp, err := client.Do(req)
	s.NoError(err)

	s.Equal(http.StatusNoContent, resp.StatusCode)
}

func (s *SuiteTest) TestDeleteTask_WrongPermission() {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/tasks/1", baseURL), nil)
	s.NoError(err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", technician_token)

	client := http.Client{}
	resp, err := client.Do(req)
	s.NoError(err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	s.NoError(err)

	apiErr, err := error_utils.NewApiErrFromBytes(byteBody)
	s.Nil(err)
	s.NotNil(apiErr)
	s.Equal(http.StatusForbidden, apiErr.Status())
	s.Equal("The user doesn't have the right permission to delete a specific task", apiErr.Message())
	s.Equal("forbidden", apiErr.Error())
}

func (s *SuiteTest) TestDeleteTask_InvalidID() {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/tasks/abc", baseURL), nil)
	s.NoError(err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", manager_token)

	client := http.Client{}
	resp, err := client.Do(req)
	s.NoError(err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	s.NoError(err)

	apiErr, err := error_utils.NewApiErrFromBytes(byteBody)
	s.Nil(err)
	s.NotNil(apiErr)
	s.Equal(http.StatusBadRequest, apiErr.Status())
	s.Equal("not possible to convert abc into a number", apiErr.Message())
	s.Equal("bad_request", apiErr.Error())
}

func (s *SuiteTest) TestDeleteTask_NotFoundTask() {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/tasks/1", baseURL), nil)
	s.NoError(err)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", manager_token)

	client := http.Client{}
	resp, err := client.Do(req)
	s.NoError(err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	s.NoError(err)

	apiErr, err := error_utils.NewApiErrFromBytes(byteBody)
	s.Nil(err)
	s.NotNil(apiErr)
	s.Equal(http.StatusNotFound, apiErr.Status())
	s.Equal("no record matching given the identification", apiErr.Message())
	s.Equal("not_found", apiErr.Error())
}
