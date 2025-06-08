package gamification_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ludiks/tests/e2e/testutil"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type GenerateCircuitStepTestSuite struct {
	suite.Suite
	env       *testutil.TestEnv
	token     string
	circuitID string
}

func TestGenerateCircuitStepTestSuite(t *testing.T) {
	suite.Run(t, new(GenerateCircuitStepTestSuite))
}

func (s *GenerateCircuitStepTestSuite) SetupSuite() {
	s.env = testutil.NewTestEnv()
	token, _ := s.env.GenerateTestUserAndToken()
	s.token = token

	projectID := uuid.New().String()
	s.circuitID = uuid.New().String()
	s.env.ExecSQL(`
		INSERT INTO projects (id, name) VALUES ($1, $2);
	`, projectID, "Project 1")
	s.env.ExecSQL(`
		INSERT INTO circuits (id, name, type, project_id, active) VALUES ($1, $2, $3, $4, true);
	`, s.circuitID, "Circuit 1", "actions", projectID)
}

func (s *GenerateCircuitStepTestSuite) TearDownSuite() {
	s.env.Cleanup()
}

func (s *GenerateCircuitStepTestSuite) SetupTest() {
	s.env.ExecSQL(`DELETE FROM steps;`)
}

func (s *GenerateCircuitStepTestSuite) TestGenerateCircuitActionsLinearSteps() {
	router := s.env.NewRouter()

	requestBody := map[string]interface{}{
		"circuitType":   "actions",
		"curve":         "linear",
		"numberOfSteps": 10,
		"maxValue":      100,
		"eventName":     "test",
	}
	jsonBody, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/api/circuits/"+s.circuitID+"/generate-steps", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	s.Equal(http.StatusCreated, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	steps := response["steps"].([]interface{})

	s.Equal(10, len(steps))
	s.Equal(float64(100), steps[9].(map[string]interface{})["completionThreshold"].(float64))
	s.Equal(float64(50), steps[4].(map[string]interface{})["completionThreshold"].(float64))
	s.Equal("actions", response["type"].(string))
}

func (s *GenerateCircuitStepTestSuite) TestGenerateCircuitActionsExponentialSteps() {
	router := s.env.NewRouter()

	requestBody := map[string]interface{}{
		"circuitType":   "actions",
		"curve":         "exponential",
		"numberOfSteps": 10,
		"maxValue":      100,
		"exponent":      2.0,
		"eventName":     "test",
	}
	jsonBody, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/api/circuits/"+s.circuitID+"/generate-steps", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	s.Equal(http.StatusCreated, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	steps := response["steps"].([]interface{})

	s.Equal(10, len(steps))

	expectedValues := map[int]int{
		0: 1,   // Premier palier
		2: 9,   // 3ème palier
		4: 25,  // 5ème palier
		6: 49,  // 7ème palier
		8: 81,  // 9ème palier
		9: 100, // Dernier palier
	}

	for index, expectedValue := range expectedValues {
		actualValue := int(steps[index].(map[string]interface{})["completionThreshold"].(float64))
		s.Equal(expectedValue, actualValue, "Valeur incorrecte à l'étape %d", index)
	}

	s.Equal("actions", response["type"].(string))
}

func (s *GenerateCircuitStepTestSuite) TestGenerateCircuitActionsLogarithmicSteps() {
	router := s.env.NewRouter()

	requestBody := map[string]interface{}{
		"circuitType":   "actions",
		"curve":         "logarithmic",
		"numberOfSteps": 10,
		"maxValue":      100,
		"exponent":      2.0,
		"eventName":     "test",
	}
	jsonBody, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/api/circuits/"+s.circuitID+"/generate-steps", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	s.Equal(http.StatusCreated, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	steps := response["steps"].([]interface{})

	s.Equal(10, len(steps))

	expectedValues := map[int]int{
		0: 32,  // Premier palier
		2: 55,  // 3ème palier
		4: 71,  // 5ème palier
		6: 84,  // 7ème palier
		8: 95,  // 9ème palier
		9: 100, // Dernier palier
	}

	for index, expectedValue := range expectedValues {
		actualValue := int(steps[index].(map[string]interface{})["completionThreshold"].(float64))
		s.Equal(expectedValue, actualValue, "Valeur incorrecte à l'étape %d", index)
	}

	s.Equal("actions", response["type"].(string))
}
