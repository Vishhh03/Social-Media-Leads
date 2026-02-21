package handlers_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

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

func TestHandleWebhook(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockStore := NewMockStore()

	handler := &handlers.WebhookHandler{
		Store: mockStore,
	}

	r := gin.Default()
	r.POST("/webhooks/meta", handler.HandleWebhook)

	t.Run("Valid WhatsApp Webhook", func(t *testing.T) {
		payload := `{
			"object": "whatsapp_business_account",
			"entry": [
				{
					"id": "12345",
					"changes": [
						{
							"value": {
								"messaging_product": "whatsapp",
								"metadata": {
									"display_phone_number": "1234567890",
									"phone_number_id": "12345"
								},
								"contacts": [
									{
										"profile": {
											"name": "Jane Doe"
										},
										"wa_id": "0987654321"
									}
								],
								"messages": [
									{
										"from": "0987654321",
										"id": "wamid.123",
										"timestamp": "1618214300",
										"text": {
											"body": "Hello world"
										},
										"type": "text"
									}
								]
							},
							"field": "messages"
						}
					]
				}
			]
		}`
		
		req := httptest.NewRequest(http.MethodPost, "/webhooks/meta", strings.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status OK, got %v", w.Code)
		}
		
		// Allow the async goroutine to process the payload
		time.Sleep(50 * time.Millisecond)
	})
	
	t.Run("Valid Instagram Webhook", func(t *testing.T) {
		payload := `{
			"object": "instagram",
			"entry": [
				{
					"id": "ig_page_123",
					"time": 1618214300,
					"messaging": [
						{
							"sender": {
								"id": "ig_user_456"
							},
							"recipient": {
								"id": "ig_page_123"
							},
							"timestamp": 1618214300,
							"message": {
								"mid": "mid.123",
								"text": "Hi from IG"
							}
						}
					]
				}
			]
		}`
		
		req := httptest.NewRequest(http.MethodPost, "/webhooks/meta", strings.NewReader(payload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status OK, got %v", w.Code)
		}
		
		// Allow asynchronous goroutine to process the payload
		time.Sleep(50 * time.Millisecond)
	})

	t.Run("Invalid Payload", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/webhooks/meta", strings.NewReader("not json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status Bad Request, got %v", w.Code)
		}
	})
}

