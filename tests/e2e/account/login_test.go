package account_test

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"ludiks/tests/e2e/testutil"

	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/argon2"
)

type LoginTestSuite struct {
	suite.Suite
	env *testutil.TestEnv
}

func TestLoginTestSuite(t *testing.T) {
	suite.Run(t, new(LoginTestSuite))
}

func (s *LoginTestSuite) SetupSuite() {
	s.env = testutil.NewTestEnv()
	encryptedPassword := s.hashPassword("password123")
	s.env.ResetDatabase()
	s.env.ExecSQL(`
		INSERT INTO users (email, password) VALUES ($1, $2);
	`, "test@example.com", encryptedPassword)
}

func (s *LoginTestSuite) TearDownSuite() {
	s.env.Cleanup()
}

func (s *LoginTestSuite) hashPassword(password string) string {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return ""
	}

	time := uint32(1)
	memory := uint32(64 * 1024)
	threads := uint8(4)
	keyLen := uint32(32)

	hash := argon2.IDKey([]byte(password), salt, time, memory, threads, keyLen)

	b64Salt := base64.StdEncoding.EncodeToString(salt)
	b64Hash := base64.StdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s", memory, time, threads, b64Salt, b64Hash)

	return encodedHash
}

func (s *LoginTestSuite) TestLogin() {
	router := s.env.NewRouter()

	requestBody := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	jsonBody, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/api/accounts/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)

	responseBody := w.Body.String()
	s.T().Logf("Response body: %s", responseBody)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	s.Equal("test@example.com", response["user"].(map[string]interface{})["email"])
	s.NotNil(response["token"])
}

func (s *LoginTestSuite) TestLoginBadCredentials() {
	router := s.env.NewRouter()

	requestBody := map[string]string{
		"email":    "test@example.com",
		"password": "mauvais_mot_de_passe",
	}
	jsonBody, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/api/accounts/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	s.Equal(http.StatusBadRequest, w.Code)

	var response struct {
		Error string `json:"error"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	s.NoError(err, "Should be able to parse error response")
	s.Equal("Invalid credentials", response.Error)
}
