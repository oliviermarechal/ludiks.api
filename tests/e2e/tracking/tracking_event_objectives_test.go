package tracking_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ludiks/tests/e2e/testutil"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type TrackingEventObjectivesTestSuite struct {
	suite.Suite
	env       *testutil.TestEnv
	apiKey    string
	projectID string
	circuitID string
	endUserID string
	router    *gin.Engine
}

func TestTrackingEventObjectivesSuite(t *testing.T) {
	suite.Run(t, new(TrackingEventObjectivesTestSuite))
}

func (s *TrackingEventObjectivesTestSuite) SetupTest() {
	s.env = testutil.NewTestEnv()

	s.env.ExecSQL(`DELETE FROM user_step_progressions`)
	s.env.ExecSQL(`DELETE FROM user_circuit_progressions`)
	s.env.ExecSQL(`DELETE FROM steps`)
	s.env.ExecSQL(`DELETE FROM circuits`)
	s.env.ExecSQL(`DELETE FROM end_users`)
	s.env.ExecSQL(`DELETE FROM api_keys`)
	s.env.ExecSQL(`DELETE FROM projects`)

	s.projectID = uuid.New().String()
	s.circuitID = uuid.New().String()
	s.endUserID = uuid.New().String()
	apiKeyValue := "api_key_value"

	s.env.ExecSQL(`
		INSERT INTO projects (id, name) VALUES ($1, $2);
	`, s.projectID, "Project Test")

	s.env.ExecSQL(`
		INSERT INTO api_keys (value, name, project_id) VALUES ($1, $2, $3);
	`, apiKeyValue, "test-api-key", s.projectID)
	s.apiKey = apiKeyValue

	s.env.ExecSQL(`
		INSERT INTO end_users (id, external_id, project_id) VALUES ($1, $2, $3);
	`, s.endUserID, "test-user-123", s.projectID)

	s.env.ExecSQL(`
		INSERT INTO circuits (id, name, type, project_id, active) 
		VALUES ($1, $2, $3, $4, true);
	`, s.circuitID, "Circuit Objectives", "objective", s.projectID)

	s.env.ExecSQL(`
		INSERT INTO steps (id, circuit_id, name, step_number, completion_threshold, event_name) VALUES
		($1, $2, $3, 1, 1, 'objective_1'),
		($4, $2, $5, 2, 1, 'objective_2'),
		($6, $2, $7, 3, 1, 'objective_3');
	`,
		uuid.New().String(),
		s.circuitID,
		"step 1",
		uuid.New().String(),
		"step 2",
		uuid.New().String(),
		"step 3",
	)

	s.router = s.env.NewRouter()
}

func (s *TrackingEventObjectivesTestSuite) TearDownTest() {
	s.env.Cleanup()
}

func (s *TrackingEventObjectivesTestSuite) trackEvent(eventName string, statusCode int) map[string]interface{} {
	requestBody := map[string]interface{}{
		"user_id":    s.endUserID,
		"event_name": eventName,
	}
	jsonBody, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/api/tracking?api_key="+s.apiKey, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	s.Equal(statusCode, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		println("Unmarshal error:", err.Error())
	}

	return response
}

func (s *TrackingEventObjectivesTestSuite) TestFirstObjectiveCompletion() {
	response := s.trackEvent("objective_1", http.StatusOK)

	s.True(response["success"].(bool))
	s.True(response["updated"].(bool))
	s.True(response["stepCompleted"].(bool))
	s.Equal(float64(1), response["points"].(float64))
}

func (s *TrackingEventObjectivesTestSuite) TestObjectiveWithoutStepOrder() {
	response := s.trackEvent("objective_2", http.StatusOK)

	s.True(response["success"].(bool))
	s.True(response["stepCompleted"].(bool))
}

func (s *TrackingEventObjectivesTestSuite) TestSequentialObjectiveCompletion() {
	s.trackEvent("objective_1", http.StatusOK)
	response := s.trackEvent("objective_2", http.StatusOK)

	s.True(response["success"].(bool))
	s.True(response["updated"].(bool))
	s.True(response["stepCompleted"].(bool))
	s.Equal(float64(2), response["points"].(float64))
}

func (s *TrackingEventObjectivesTestSuite) TestCompleteCircuit() {
	s.trackEvent("objective_1", http.StatusOK)
	s.trackEvent("objective_2", http.StatusOK)
	response := s.trackEvent("objective_3", http.StatusOK)

	s.True(response["success"].(bool))
	s.True(response["updated"].(bool))
	s.True(response["stepCompleted"].(bool))
	s.True(response["circuitCompleted"].(bool))
	s.Equal(float64(3), response["points"].(float64))
}

func (s *TrackingEventObjectivesTestSuite) TestRepeatCompletedObjective() {
	s.trackEvent("objective_1", http.StatusOK)

	response := s.trackEvent("objective_1", http.StatusOK)

	s.True(response["success"].(bool))
	s.True(response["alreadyCompleted"].(bool))
}

func (s *TrackingEventObjectivesTestSuite) TestInvalidObjective() {
	response := s.trackEvent("invalid_objective", http.StatusBadRequest)

	s.False(response["success"].(bool))
}
