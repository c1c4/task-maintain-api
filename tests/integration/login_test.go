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

func (s *SuiteTest) TestLogin_Success() {
	s.seedOneUserTech()
	jsonBody := `{"email": "tech@tech.com", "password": "123"}`
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/login", baseURL), bytes.NewBufferString(jsonBody))
	s.NoError(err)

	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	s.NoError(err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	s.NoError(err)

	var authentication models.AuthenticationData
	errUnmarshal := json.Unmarshal(byteBody, &authentication)
	s.Nil(errUnmarshal)
	s.NotNil(authentication)
	s.Equal(http.StatusOK, resp.StatusCode)
	s.Equal("1", authentication.ID)
	s.NotNil(authentication.Token)
	resp.Body.Close()
}

func (s *SuiteTest) TestLogin_WrongJSONFormat() {
	jsonBody := `{"email": "tech@tech.com" "password": ""}`
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/login", baseURL), bytes.NewBufferString(jsonBody))
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

func (s *SuiteTest) TestLogin_WihtoutPassword() {
	s.seedOneUserTech()
	jsonBody := `{"email": "tech@tech.com", "password": ""}`
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/login", baseURL), bytes.NewBufferString(jsonBody))
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
	s.Equal("crypto/bcrypt: hashedPassword is not the hash of the given password", apiErr.Message())
	s.Equal("unauthorized", apiErr.Error())
}

func (s *SuiteTest) TestLogin_WihtoutEmail() {
	jsonBody := `{"email": "", "password": "123"}`
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/login", baseURL), bytes.NewBufferString(jsonBody))
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
	s.Equal(http.StatusNotFound, apiErr.Status())
	s.Equal("no record matching given the identification", apiErr.Message())
	s.Equal("not_found", apiErr.Error())
}
