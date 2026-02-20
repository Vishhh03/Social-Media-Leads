package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/social-media-lead/backend/internal/api/handlers"
	"github.com/social-media-lead/backend/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a new mock store
	mockStore := NewMockStore()

	// Hash a real password for the test
	password := "securepassword123"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// Create our test user in the mock store
	testUser := &models.User{
		ID:           1,
		Email:        "test@example.com",
		PasswordHash: string(hash),
		IsActive:     true,
	}
	mockStore.UsersByEmail[testUser.Email] = testUser

	// Set up the handler
	handler := &handlers.AuthHandler{
		Store:     mockStore,
		JWTSecret: "super_secret_test_key",
	}

	r := gin.Default()
	r.POST("/api/v1/auth/login", handler.Login)

	tests := []struct {
		name           string
		payload        map[string]string
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Valid credentials",
			payload: map[string]string{
				"email":    "test@example.com",
				"password": "securepassword123",
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "token", // Just checking for token presence
		},
		{
			name: "Invalid password",
			payload: map[string]string{
				"email":    "test@example.com",
				"password": "wrongpassword",
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Invalid email or password",
		},
		{
			name: "User not found",
			payload: map[string]string{
				"email":    "nonexistent@example.com",
				"password": "securepassword123",
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Invalid email or password",
		},
		{
			name: "Missing fields",
			payload: map[string]string{
				"email": "test@example.com",
				// Missing password
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.payload)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != tc.expectedStatus {
				t.Errorf("expected status %v, got %v", tc.expectedStatus, w.Code)
			}

			responseBody := w.Body.String()
			if tc.expectedStatus == http.StatusOK {
				// Success format: {"token": "...", "user": {...}}
				var response map[string]interface{}
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Fatalf("Failed to parse response: %v", err)
				}
				if _, ok := response["token"]; !ok {
					t.Errorf("Expected token in response, got: %v", response)
				}
			} else {
				// Error format: {"error": "..."}
				if !bytes.Contains(w.Body.Bytes(), []byte(tc.expectedBody)) {
					t.Errorf("expected response to contain %q, but got %q", tc.expectedBody, responseBody)
				}
			}
		})
	}
}
