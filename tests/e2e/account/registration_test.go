package account_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ludiks/tests/e2e/testutil"

	"github.com/stretchr/testify/suite"
)

type RegistrationTestSuite struct {
	suite.Suite
	env *testutil.TestEnv
}

func TestAccountSuite(t *testing.T) {
	suite.Run(t, new(RegistrationTestSuite))
}

func (s *RegistrationTestSuite) SetupSuite() {
	s.env = testutil.NewTestEnv()
}

func (s *RegistrationTestSuite) TearDownSuite() {
	s.env.Cleanup()
}

func (s *RegistrationTestSuite) TestRegistration() {
	router := s.env.NewRouter()

	requestBody := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	jsonBody, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/api/accounts/registration", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	s.Equal(http.StatusCreated, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	s.Equal("test@example.com", response["user"].(map[string]interface{})["email"])
	s.NotNil(response["token"])
}
