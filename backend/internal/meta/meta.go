package meta

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const graphAPIBase = "https://graph.facebook.com/v21.0"

// Client handles outbound messaging via Meta's Graph API.
type Client struct {
	HTTPClient *http.Client
}

// NewClient creates a Meta API client with sensible defaults.
func NewClient() *Client {
	return &Client{
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SendResult contains the API response after sending a message.
type SendResult struct {
	MessageID string `json:"message_id,omitempty"`
	Success   bool   `json:"success"`
	Error     string `json:"error,omitempty"`
}

// SendWhatsAppMessage sends a text message via WhatsApp Cloud API.
// phoneNumberID is the business phone number ID (from the channel).
// recipientPhone is the end-user's phone number (e.g., "15551234567").
func (c *Client) SendWhatsAppMessage(phoneNumberID, recipientPhone, text, accessToken string) (*SendResult, error) {
	url := fmt.Sprintf("%s/%s/messages", graphAPIBase, phoneNumberID)

	payload := map[string]interface{}{
		"messaging_product": "whatsapp",
		"to":                recipientPhone,
		"type":              "text",
		"text": map[string]string{
			"body": text,
		},
	}

	return c.sendRequest(url, payload, accessToken)
}

// SendInstagramMessage sends a text message via Instagram Messaging API.
// recipientID is the Instagram-scoped user ID.
func (c *Client) SendInstagramMessage(recipientID, text, accessToken string) (*SendResult, error) {
	url := fmt.Sprintf("%s/me/messages", graphAPIBase)

	payload := map[string]interface{}{
		"recipient": map[string]string{
			"id": recipientID,
		},
		"message": map[string]string{
			"text": text,
		},
	}

	return c.sendRequest(url, payload, accessToken)
}

// SendFacebookMessage sends a text message via Facebook Messenger Platform.
// recipientID is the page-scoped user ID.
func (c *Client) SendFacebookMessage(recipientID, text, accessToken string) (*SendResult, error) {
	url := fmt.Sprintf("%s/me/messages", graphAPIBase)

	payload := map[string]interface{}{
		"recipient": map[string]string{
			"id": recipientID,
		},
		"message": map[string]string{
			"text": text,
		},
	}

	return c.sendRequest(url, payload, accessToken)
}

// SendMessage dispatches a message to the correct platform.
func (c *Client) SendMessage(platform, accountID, recipientID, text, accessToken string) (*SendResult, error) {
	switch platform {
	case "whatsapp":
		return c.SendWhatsAppMessage(accountID, recipientID, text, accessToken)
	case "instagram":
		return c.SendInstagramMessage(recipientID, text, accessToken)
	case "facebook":
		return c.SendFacebookMessage(recipientID, text, accessToken)
	default:
		return nil, fmt.Errorf("unsupported platform: %s", platform)
	}
}

// sendRequest makes a POST request to the Meta Graph API.
func (c *Client) sendRequest(url string, payload interface{}, accessToken string) (*SendResult, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		log.Printf("[Meta API] Error %d: %s", resp.StatusCode, string(respBody))
		return &SendResult{
			Success: false,
			Error:   fmt.Sprintf("API error %d: %s", resp.StatusCode, string(respBody)),
		}, nil
	}

	// Parse the response for message ID
	var apiResp map[string]interface{}
	if err := json.Unmarshal(respBody, &apiResp); err == nil {
		// WhatsApp returns messages[0].id, IG/FB returns message_id
		if messages, ok := apiResp["messages"].([]interface{}); ok && len(messages) > 0 {
			if msg, ok := messages[0].(map[string]interface{}); ok {
				if id, ok := msg["id"].(string); ok {
					return &SendResult{Success: true, MessageID: id}, nil
				}
			}
		}
		if msgID, ok := apiResp["message_id"].(string); ok {
			return &SendResult{Success: true, MessageID: msgID}, nil
		}
	}

	log.Printf("[Meta API] Sent successfully: %s", string(respBody))
	return &SendResult{Success: true}, nil
}
