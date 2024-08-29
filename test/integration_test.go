package testing

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"test-task/internal/config"
	db "test-task/internal/database"
	"test-task/internal/modules/auth"
	"test-task/internal/routes"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func initializeApp() (*gin.Engine, *config.Config, func()) {
	cfg := &config.Config{
		Port:         "8080",
		DBSource:     "postgresql://postgres:1111@localhost:5432/test-task?sslmode=disable",
		JWTSecretKey: "testtest",
	}

	dbHandler := db.InitDB(cfg.DBSource)
	service, err := auth.InitAuthService(dbHandler, cfg)
	if err != nil {
		panic("mock module initialization error")
	}

	app := gin.New()
	router := routes.NewAppRouter(app, "/api", "/v1")
	router.RegisterAuthRoutes(auth.NewHandler(service, cfg))

	testServer := httptest.NewServer(app)

	cleanup := func() {
		dbHandler.DB.Exec("TRUNCATE TABLE users, tokens RESTART IDENTITY CASCADE")
		testServer.Close()
	}

	return app, cfg, cleanup
}

func TestAuthFlow(t *testing.T) {
	app, cfg, cleanup := initializeApp()
	defer cleanup()

	userPayload := map[string]string{
		"email":    "testuser@example.com",
		"password": "password",
	}

	// Step 1: User signs up
	signUpResp, signUpBody, err := sendRequest(http.MethodPost, "http://localhost:"+cfg.Port+"/api/v1/auth/signup", userPayload, app)
	if err != nil {
		t.Fatalf("Failed to sign up user: %v", err)
	}
	assert.Equal(t, http.StatusOK, signUpResp.StatusCode)

	t.Logf("Sign Up Response Body: %s", signUpBody)

	var signUpResponse map[string]interface{}
	if err := json.Unmarshal(signUpBody, &signUpResponse); err != nil {
		t.Fatalf("Failed to decode sign up response: %v", err)
	}
	// Step 2: User logs in
	loginPayload := userPayload
	loginResp, loginBody, err := sendRequest(http.MethodPost, "http://localhost:"+cfg.Port+"/api/v1/auth/login", loginPayload, app)
	if err != nil {
		t.Fatalf("Failed to log in user: %v", err)
	}
	assert.Equal(t, http.StatusOK, loginResp.StatusCode)

	var loginResponse map[string]interface{}
	if err := json.Unmarshal(loginBody, &loginResponse); err != nil {
		t.Fatalf("Failed to decode login response: %v", err)
	}

	accessToken, ok := loginResponse["access_token"].(string)
	if !ok {
		t.Fatalf("Access token not found in login response")
	}

	refreshToken, ok := loginResponse["refresh_token"].(string)
	if !ok {
		t.Fatalf("Refresh token not found in login response")
	}

	// Step 3: User refreshes token
	refreshPayload := map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}
	refreshResp, refreshBody, err := sendRequest(http.MethodPost, "http://localhost:"+cfg.Port+"/api/v1/auth/refresh-tokens", refreshPayload, app)
	if err != nil {
		t.Fatalf("Failed to refresh tokens: %v", err)
	}
	assert.Equal(t, http.StatusOK, refreshResp.StatusCode)

	var refreshResponse map[string]interface{}
	if err := json.Unmarshal(refreshBody, &refreshResponse); err != nil {
		t.Fatalf("Failed to decode refresh tokens response: %v", err)
	}

	assert.Contains(t, refreshResponse, "access_token")
	assert.Contains(t, refreshResponse, "refresh_token")
}

func sendRequest(method, url string, payload interface{}, app *gin.Engine) (*http.Response, []byte, error) {
	var body *bytes.Reader
	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			return nil, nil, err
		}
		body = bytes.NewReader(data)
	} else {
		body = bytes.NewReader([]byte{})
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	app.ServeHTTP(recorder, req)

	responseBody := recorder.Body.Bytes()

	return recorder.Result(), responseBody, nil
}
