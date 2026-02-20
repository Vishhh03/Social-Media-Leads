package handlers_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/social-media-lead/backend/internal/api/handlers"
	"github.com/social-media-lead/backend/internal/config"
)

func TestVerifyWebhook(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Setup Config
	cfg := &config.Config{
		Meta: config.MetaConfig{
			VerifyToken: "secret_token",
		},
	}

	handler := &handlers.WebhookHandler{
		Config: cfg,
	}

	r := gin.Default()
	r.GET("/webhooks/meta", handler.VerifyWebhook)

	t.Run("Valid token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/webhooks/meta?hub.mode=subscribe&hub.verify_token=secret_token&hub.challenge=12345", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status OK, got %v", w.Code)
		}

		body, _ := io.ReadAll(w.Body)
		if string(body) != "12345" {
			t.Errorf("expected body 12345, got %v", string(body))
		}
	})

	t.Run("Invalid token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/webhooks/meta?hub.mode=subscribe&hub.verify_token=wrong_token&hub.challenge=12345", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Errorf("expected status Forbidden, got %v", w.Code)
		}
		
		body, _ := io.ReadAll(w.Body)
		if !strings.Contains(string(body), "Verification failed") {
			t.Errorf("expected error message in body, got %v", string(body))
		}
	})

	t.Run("Invalid mode", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/webhooks/meta?hub.mode=unsubscribe&hub.verify_token=secret_token&hub.challenge=12345", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Errorf("expected status Forbidden, got %v", w.Code)
		}
	})
}
