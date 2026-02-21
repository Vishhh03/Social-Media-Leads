package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/social-media-lead/backend/internal/api/handlers"
	"github.com/social-media-lead/backend/internal/meta"
)

func TestInboxHandlers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockStore := NewMockStore()
	mockMetaClient := meta.NewClient()

	handler := &handlers.InboxHandler{
		Store:      mockStore,
		MetaClient: mockMetaClient,
	}

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", int64(1))
	})
	r.GET("/inbox/conversations", handler.GetConversations)
	r.GET("/inbox/messages/:contact_id", handler.GetMessages)
	r.GET("/inbox/contacts", handler.GetContacts)

	tests := []struct {
		name           string
		method         string
		url            string
		expectedStatus int
	}{
		{"Get Conversations", http.MethodGet, "/inbox/conversations", http.StatusOK},
		{"Get Messages", http.MethodGet, "/inbox/messages/1", http.StatusOK},
		{"Get Messages Invalid ID", http.MethodGet, "/inbox/messages/abc", http.StatusBadRequest},
		{"Get Contacts", http.MethodGet, "/inbox/contacts", http.StatusOK},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.url, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			if w.Code != tc.expectedStatus {
				t.Errorf("expected status %v, got %v", tc.expectedStatus, w.Code)
			}
		})
	}
}
