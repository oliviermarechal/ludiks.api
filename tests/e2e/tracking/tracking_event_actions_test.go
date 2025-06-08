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

type TrackingEventActionsTestSuite struct {
	suite.Suite
	env       *testutil.TestEnv
	apiKey    string
	projectID string
	circuitID string
	endUserID string
	router    *gin.Engine
}

func TestTrackingEventActionsSuite(t *testing.T) {
	suite.Run(t, new(TrackingEventActionsTestSuite))
}

func (s *TrackingEventActionsTestSuite) SetupTest() {
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
    `, s.circuitID, "Circuit Actions", "actions", s.projectID)

	s.env.ExecSQL(`
        INSERT INTO steps (id, circuit_id, name, step_number, completion_threshold, event_name) VALUES
        ($1, $2, $3, 1, 10, 'test_event'),
        ($4, $2, $5, 2, 50, 'test_event'),
        ($6, $2, $7, 3, 100, 'test_event');
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

func (s *TrackingEventActionsTestSuite) TearDownTest() {
	s.env.Cleanup()
}

func (s *TrackingEventActionsTestSuite) trackEvent(value *int) map[string]interface{} {
	requestBody := map[string]interface{}{
		"user_id":    s.endUserID,
		"event_name": "test_event",
		"value":      value,
	}
	jsonBody, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/api/tracking?api_key="+s.apiKey, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	println("Status Code:", w.Code)
	println("Raw Response:", w.Body.String())

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		println("Unmarshal error:", err.Error())
	}

	return response
}

func (s *TrackingEventActionsTestSuite) TestFirstEventCreatesProgression() {
	response := s.trackEvent(nil)

	s.True(response["success"].(bool))
	s.True(response["updated"].(bool))
	s.Equal(float64(1), response["points"].(float64))
	s.False(response["stepCompleted"].(bool))
}

func (s *TrackingEventActionsTestSuite) TestCompleteFirstStep() {
	value := 10
	response := s.trackEvent(&value)

	s.True(response["success"].(bool))
	s.True(response["updated"].(bool))
	s.True(response["stepCompleted"].(bool))
	s.Equal(float64(10), response["points"].(float64))
}

func (s *TrackingEventActionsTestSuite) TestProgressiveCompletion() {
	// Premier event
	s.trackEvent(nil)

	// Deuxième event qui complète la première étape
	value := 9
	response := s.trackEvent(&value)
	s.True(response["stepCompleted"].(bool))
	s.Equal(float64(10), response["points"].(float64))

	// Event qui progresse vers la deuxième étape
	value = 20
	response = s.trackEvent(&value)
	s.False(response["stepCompleted"].(bool))
	s.Equal(float64(30), response["points"].(float64))
}

func (s *TrackingEventActionsTestSuite) TestCompleteMultipleStepsAtOnce() {
	value := 60
	response := s.trackEvent(&value)

	s.True(response["success"].(bool))
	s.True(response["updated"].(bool))
	s.True(response["stepCompleted"].(bool))
	s.Equal(float64(60), response["points"].(float64))
}

func (s *TrackingEventActionsTestSuite) TestCompleteCircuit() {
	value := 100
	response := s.trackEvent(&value)

	s.True(response["success"].(bool))
	s.True(response["updated"].(bool))
	s.True(response["stepCompleted"].(bool))
	s.True(response["circuitCompleted"].(bool))
}

func (s *TrackingEventActionsTestSuite) TestContinueProgressAfterCompletion() {
	// Compléter d'abord le circuit
	value := 100
	s.trackEvent(&value)

	// Ajouter plus de points
	value = 10
	response := s.trackEvent(&value)

	s.True(response["success"].(bool))
	s.Equal(float64(110), response["points"].(float64))
	s.False(response["stepCompleted"].(bool))
}
