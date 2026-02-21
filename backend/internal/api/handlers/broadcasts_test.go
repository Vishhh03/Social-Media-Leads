package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/social-media-lead/backend/internal/api/handlers"
)

func TestBroadcastHandlers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockStore := NewMockStore()

	handler := &handlers.BroadcastHandler{
		Store: mockStore,
	}

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", int64(1))
	})
	r.POST("/broadcasts", handler.CreateBroadcast)
	r.GET("/broadcasts", handler.ListBroadcasts)

	t.Run("Create Broadcast Valid", func(t *testing.T) {
		payload := map[string]interface{}{
			"name":    "Promo",
			"content": "50% off!",
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/broadcasts", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusCreated {
			t.Errorf("expected 201, got %v", w.Code)
		}
	})

	t.Run("Create Broadcast Invalid", func(t *testing.T) {
		payload := map[string]interface{}{
			"name": "Promo",
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/broadcasts", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %v", w.Code)
		}
	})

	t.Run("List Broadcasts", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/broadcasts", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %v", w.Code)
		}
	})
}
