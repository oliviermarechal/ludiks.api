package account_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ludiks/tests/e2e/testutil"

	"github.com/stretchr/testify/suite"
)

type CreateProjectTestSuite struct {
	suite.Suite
	env *testutil.TestEnv
}

func TestCreateProjectTestSuite(t *testing.T) {
	suite.Run(t, new(CreateProjectTestSuite))
}

func (s *CreateProjectTestSuite) SetupSuite() {
	s.env = testutil.NewTestEnv()
}

func (s *CreateProjectTestSuite) TearDownSuite() {
	s.env.Cleanup()
}

func (s *CreateProjectTestSuite) TestCreateProject() {
	router := s.env.NewRouter()

	requestBody := map[string]string{
		"name": "Mon Projet",
	}

	token, _ := s.env.GenerateTestUserAndToken()
	req := s.env.AuthenticatedRequest("POST", "/api/projects", requestBody, token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	s.Equal(http.StatusCreated, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	s.Equal("Mon Projet", response["name"])
}
