package integration

import (
	"api/app/models"
	"api/app/utils/error_utils"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (s *SuiteTest) TestCreateUser_Success() {
	jsonBody := `{"name": "Tech user", "email": "tech@tech.com", "password": "123", "type": "Technician"}`
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/users", baseURL), bytes.NewBufferString(jsonBody))
	s.NoError(err)

	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	s.NoError(err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	s.NoError(err)

	var user models.User
	errUnmarshal := json.Unmarshal(byteBody, &user)
	s.Nil(errUnmarshal)
	s.NotNil(user)
	s.Equal(http.StatusCreated, resp.StatusCode)
	s.Equal(uint64(1), user.ID)
	s.NotEqual("123", user.Password)
	resp.Body.Close()
}

func (s *SuiteTest) TestCreateUser_WrongJSONFormat() {
	jsonBody := `{"name": "Test user", "email": "test@test.com", "password": "" "type": "Manager"}`
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/users", baseURL), bytes.NewBufferString(jsonBody))
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
	s.Equal(http.StatusUnprocessableEntity, apiErr.Status())
	s.Equal("it's not possible to convert the JSON into an object", apiErr.Message())
	s.Equal("invalid_request", apiErr.Error())
}

func (s *SuiteTest) TestCreateUser_UniqueEmailError() {
	s.seedOneUserTech()
	jsonBody := `{"name": "Tech user", "email": "tech@tech.com", "password": "123", "type": "Technician"}`
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/users", baseURL), bytes.NewBufferString(jsonBody))
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
	s.Equal(http.StatusInternalServerError, apiErr.Status())
	s.Equal("error when processing request Error 1062: Duplicate entry 'tech@tech.com' for key 'users.email'", apiErr.Message())
	s.Equal("server_error", apiErr.Error())
}

func (s *SuiteTest) TestCreateUser_WihtoutPassword() {
	jsonBody := `{"name": "Test user", "email": "test@test.com", "password": "", "type": "Manager"}`
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/users", baseURL), bytes.NewBufferString(jsonBody))
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
	s.Equal(http.StatusBadRequest, apiErr.Status())
	s.Equal("the field password is required can't be empty", apiErr.Message())
	s.Equal("bad_request", apiErr.Error())
}

func (s *SuiteTest) TestCreateUser_WihtoutEmail() {
	jsonBody := `{"name": "Test user", "email": "", "password": "123", "type": "Manager"}`
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/users", baseURL), bytes.NewBufferString(jsonBody))
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
	s.Equal(http.StatusBadRequest, apiErr.Status())
	s.Equal("the field email is required can't be empty", apiErr.Message())
	s.Equal("bad_request", apiErr.Error())
}

func (s *SuiteTest) TestCreateUser_InvalidEmail() {
	jsonBody := `{"name": "Test user", "email": "email", "password": "123", "type": "Manager"}`
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/users", baseURL), bytes.NewBufferString(jsonBody))
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
	s.Equal(http.StatusBadRequest, apiErr.Status())
	s.Equal("the email is invalid", apiErr.Message())
	s.Equal("bad_request", apiErr.Error())
}

func (s *SuiteTest) TestCreateUser_InvalidUserType() {
	jsonBody := `{"name": "Test user", "email": "email@email.com", "password": "123", "type": "Test"}`
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/users", baseURL), bytes.NewBufferString(jsonBody))
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
	s.Equal(http.StatusBadRequest, apiErr.Status())
	s.Equal("the field type should be Technician or Manager", apiErr.Message())
	s.Equal("bad_request", apiErr.Error())
}

func (s *SuiteTest) TestCreateUser_MultipleErrors() {
	jsonBody := `{"name": "", "email": "email@email.com", "password": "123", "type": ""}`
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/users", baseURL), bytes.NewBufferString(jsonBody))
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
	s.Equal(http.StatusBadRequest, apiErr.Status())
	s.Equal(
		"the field name is required can't be empty\nthe field type is required can't be empty",
		apiErr.Message(),
	)
	s.Equal("bad_request", apiErr.Error())
}
