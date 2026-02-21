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
)

func TestWorkflowHandlers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockStore := NewMockStore()

	handler := &handlers.WorkflowHandler{
		Store: mockStore,
	}

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", int64(1))
	})

	r.GET("/api/v1/workflows", handler.ListWorkflows)
	r.GET("/api/v1/workflows/:id", handler.GetWorkflow)
	r.POST("/api/v1/workflows", handler.CreateWorkflow)
	r.PUT("/api/v1/workflows/:id", handler.UpdateWorkflow)
	r.DELETE("/api/v1/workflows/:id", handler.DeleteWorkflow)

	t.Run("Create Workflow", func(t *testing.T) {
		payload := map[string]interface{}{
			"name":         "New AI Workflow",
			"trigger_type": "trigger_meta_dm",
			"status":       "published",
			"prompt":       "Respond nicely",
			"nodes":        []interface{}{},
			"edges":        []interface{}{},
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/workflows", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected 201 Created, got %v", w.Code)
		}
	})

	t.Run("List Workflows", func(t *testing.T) {
		// Populate mock
		mockStore.Workflows[2] = &models.Workflow{ID: 2, UserID: 1, Name: "Test Flow"}

		req := httptest.NewRequest(http.MethodGet, "/api/v1/workflows", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200 OK, got %v", w.Code)
		}
	})

	t.Run("Get Workflow", func(t *testing.T) {
		mockStore.Workflows[3] = &models.Workflow{ID: 3, UserID: 1, Name: "Specific Flow"}

		req := httptest.NewRequest(http.MethodGet, "/api/v1/workflows/3", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200 OK, got %v", w.Code)
		}
	})

	t.Run("Update Workflow", func(t *testing.T) {
		mockStore.Workflows[4] = &models.Workflow{ID: 4, UserID: 1, Name: "Old Flow"}

		payload := map[string]interface{}{
			"name":         "Updated Flow",
			"trigger_type": "trigger_meta_dm",
			"status":       "published",
			"prompt":       "Respond nicely",
			"nodes":        []interface{}{},
			"edges":        []interface{}{},
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPut, "/api/v1/workflows/4", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200 OK, got %v", w.Code)
		}
	})

	t.Run("Delete Workflow", func(t *testing.T) {
		mockStore.Workflows[5] = &models.Workflow{ID: 5, UserID: 1, Name: "To Delete"}

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/workflows/5", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200 OK, got %v", w.Code)
		}

		if _, exists := mockStore.Workflows[5]; exists {
			t.Errorf("expected workflow to be deleted")
		}
	})
}
