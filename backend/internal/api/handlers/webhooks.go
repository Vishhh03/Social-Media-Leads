package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/social-media-lead/backend/internal/config"
	"github.com/social-media-lead/backend/internal/meta"
	"github.com/social-media-lead/backend/internal/models"
	"github.com/social-media-lead/backend/internal/store"
)

// WebhookHandler handles Meta platform webhook events.
type WebhookHandler struct {
	Store      *store.Storage
	Config     *config.Config
	MetaClient *meta.Client
}

// VerifyWebhook handles the GET request from Meta to verify the webhook URL.
// Meta sends: hub.mode, hub.verify_token, hub.challenge
func (h *WebhookHandler) VerifyWebhook(c *gin.Context) {
	mode := c.Query("hub.mode")
	token := c.Query("hub.verify_token")
	challenge := c.Query("hub.challenge")

	if mode == "subscribe" && token == h.Config.Meta.VerifyToken {
		log.Printf("[Webhook] Verification successful")
		c.String(http.StatusOK, challenge)
		return
	}

	log.Printf("[Webhook] Verification failed: mode=%s token=%s", mode, token)
	c.JSON(http.StatusForbidden, gin.H{"error": "Verification failed"})
}

// HandleWebhook processes incoming webhook events from Meta (WhatsApp, IG, FB).
func (h *WebhookHandler) HandleWebhook(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	// Meta expects a 200 OK response immediately
	c.JSON(http.StatusOK, gin.H{"status": "received"})

	// Process the webhook asynchronously
	go h.processWebhookPayload(payload)
}

// processWebhookPayload parses Meta webhook events and stores the message.
func (h *WebhookHandler) processWebhookPayload(payload map[string]interface{}) {
	object, _ := payload["object"].(string)

	entries, ok := payload["entry"].([]interface{})
	if !ok {
		log.Println("[Webhook] No entries found in payload")
		return
	}

	for _, entry := range entries {
		entryMap, ok := entry.(map[string]interface{})
		if !ok {
			continue
		}

		switch object {
		case "whatsapp_business_account":
			h.processWhatsAppEntry(entryMap)
		case "instagram":
			h.processInstagramEntry(entryMap)
		case "page":
			h.processFacebookEntry(entryMap)
		default:
			log.Printf("[Webhook] Unknown object type: %s", object)
		}
	}
}

// processWhatsAppEntry processes a WhatsApp Business webhook entry.
func (h *WebhookHandler) processWhatsAppEntry(entry map[string]interface{}) {
	changes, ok := entry["changes"].([]interface{})
	if !ok {
		return
	}

	for _, change := range changes {
		changeMap, ok := change.(map[string]interface{})
		if !ok {
			continue
		}

		value, ok := changeMap["value"].(map[string]interface{})
		if !ok {
			continue
		}

		// Extract the phone_number_id (this is the account_id for channel lookup)
		phoneNumberID := ""
		if metadata, ok := value["metadata"].(map[string]interface{}); ok {
			phoneNumberID, _ = metadata["phone_number_id"].(string)
		}

		messages, ok := value["messages"].([]interface{})
		if !ok {
			continue
		}

		contacts, _ := value["contacts"].([]interface{})
		senderName := ""
		if len(contacts) > 0 {
			if contact, ok := contacts[0].(map[string]interface{}); ok {
				if profile, ok := contact["profile"].(map[string]interface{}); ok {
					senderName, _ = profile["name"].(string)
				}
			}
		}

		for _, msg := range messages {
			msgMap, ok := msg.(map[string]interface{})
			if !ok {
				continue
			}

			senderID, _ := msgMap["from"].(string)
			msgID, _ := msgMap["id"].(string)
			msgType, _ := msgMap["type"].(string)

			content := ""
			if msgType == "text" {
				if textObj, ok := msgMap["text"].(map[string]interface{}); ok {
					content, _ = textObj["body"].(string)
				}
			}

			log.Printf("[WhatsApp] Message from %s (%s): %s", senderName, senderID, content)

			h.storeIncomingMessage("whatsapp", phoneNumberID, senderID, senderName, msgID, content, msgType)
		}
	}
}

// processInstagramEntry processes an Instagram webhook entry.
func (h *WebhookHandler) processInstagramEntry(entry map[string]interface{}) {
	// The entry ID is the Instagram page/account ID
	pageID := fmt.Sprintf("%v", entry["id"])

	messaging, ok := entry["messaging"].([]interface{})
	if !ok {
		return
	}

	for _, event := range messaging {
		eventMap, ok := event.(map[string]interface{})
		if !ok {
			continue
		}

		sender, ok := eventMap["sender"].(map[string]interface{})
		if !ok {
			continue
		}
		senderID := fmt.Sprintf("%v", sender["id"])

		message, ok := eventMap["message"].(map[string]interface{})
		if !ok {
			continue
		}

		content, _ := message["text"].(string)
		msgID, _ := message["mid"].(string)

		log.Printf("[Instagram] Message from %s: %s", senderID, content)

		h.storeIncomingMessage("instagram", pageID, senderID, "", msgID, content, "text")
	}
}

// processFacebookEntry processes a Facebook Page webhook entry.
func (h *WebhookHandler) processFacebookEntry(entry map[string]interface{}) {
	// The entry ID is the Facebook page ID
	pageID := fmt.Sprintf("%v", entry["id"])

	messaging, ok := entry["messaging"].([]interface{})
	if !ok {
		return
	}

	for _, event := range messaging {
		eventMap, ok := event.(map[string]interface{})
		if !ok {
			continue
		}

		sender, ok := eventMap["sender"].(map[string]interface{})
		if !ok {
			continue
		}
		senderID := fmt.Sprintf("%v", sender["id"])

		message, ok := eventMap["message"].(map[string]interface{})
		if !ok {
			continue
		}

		content, _ := message["text"].(string)
		msgID, _ := message["mid"].(string)

		log.Printf("[Facebook] Message from %s: %s", senderID, content)

		h.storeIncomingMessage("facebook", pageID, senderID, "", msgID, content, "text")
	}
}

// storeIncomingMessage resolves the user from the channel, upserts the contact,
// saves the message to the DB, and checks automation triggers.
func (h *WebhookHandler) storeIncomingMessage(platform, accountID, senderID, senderName, platformMsgID, content, msgType string) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// 1. Resolve which user owns this account by looking up the channel
	channel, err := h.Store.GetChannelByAccountID(ctx, platform, accountID)
	if err != nil {
		log.Printf("[Webhook] No channel found for %s account %s: %v", platform, accountID, err)
		return
	}

	// 2. Upsert the contact (find or create)
	contact := &models.Contact{
		UserID:         channel.UserID,
		ChannelID:      channel.ID,
		Platform:       platform,
		PlatformUserID: senderID,
		Name:           senderName,
	}
	if err := h.Store.GetOrCreateContact(ctx, contact); err != nil {
		log.Printf("[Webhook] Failed to upsert contact %s: %v", senderID, err)
		return
	}

	// Update contact name if it was empty and we now have one
	if senderName != "" && contact.Name == "" {
		contact.Name = senderName
	}

	// 3. Save the inbound message
	msg := &models.Message{
		UserID:        channel.UserID,
		ChannelID:     channel.ID,
		ContactID:     contact.ID,
		Platform:      platform,
		Direction:     "inbound",
		Content:       content,
		MessageType:   msgType,
		PlatformMsgID: platformMsgID,
		Status:        "received",
		IsAutomated:   false,
	}

	if err := h.Store.CreateMessage(ctx, msg); err != nil {
		log.Printf("[Webhook] Failed to store message: %v", err)
		return
	}

	log.Printf("[Webhook] ✅ Stored message #%d from contact #%d (user #%d)", msg.ID, contact.ID, channel.UserID)

	// 4. Check for automation triggers
	h.checkAutomationTriggers(ctx, channel, contact, content)
}

// checkAutomationTriggers checks if incoming message matches any automation rules
// and sends auto-replies via the Meta API.
func (h *WebhookHandler) checkAutomationTriggers(ctx context.Context, channel *models.Channel, contact *models.Contact, content string) {
	contentLower := strings.ToLower(strings.TrimSpace(content))
	if contentLower == "" {
		return
	}

	// Fetch all active automations for this user
	automations, err := h.Store.GetAutomationsByUser(ctx, channel.UserID)
	if err != nil {
		log.Printf("[Automation] Failed to fetch automations for user #%d: %v", channel.UserID, err)
		return
	}

	for _, automation := range automations {
		matched := false

		switch automation.TriggerType {
		case "keyword":
			for _, keyword := range automation.Keywords {
				if strings.Contains(contentLower, strings.ToLower(keyword)) {
					matched = true
					break
				}
			}
		case "first_message":
			// TODO: Track if this is the first message from this contact
			// For now, skip first_message triggers
		}

		if !matched {
			continue
		}

		log.Printf("[Automation] ✅ Matched automation '%s' for keyword in '%s'", automation.Name, content)

		// Apply delay if configured
		if automation.DelayMs > 0 {
			time.Sleep(time.Duration(automation.DelayMs) * time.Millisecond)
		}

		// Send the auto-reply via Meta API
		result, err := h.MetaClient.SendMessage(
			contact.Platform,
			channel.AccountID,
			contact.PlatformUserID,
			automation.ReplyText,
			channel.AccessToken,
		)
		if err != nil {
			log.Printf("[Automation] Failed to send reply: %v", err)
			continue
		}

		if !result.Success {
			log.Printf("[Automation] API error sending reply: %s", result.Error)
			continue
		}

		// Store the outbound automated message in DB
		autoMsg := &models.Message{
			UserID:        channel.UserID,
			ChannelID:     channel.ID,
			ContactID:     contact.ID,
			Platform:      contact.Platform,
			Direction:     "outbound",
			Content:       automation.ReplyText,
			MessageType:   "text",
			PlatformMsgID: result.MessageID,
			Status:        "sent",
			IsAutomated:   true,
		}

		if err := h.Store.CreateMessage(ctx, autoMsg); err != nil {
			log.Printf("[Automation] Failed to store auto-reply message: %v", err)
		} else {
			log.Printf("[Automation] ✅ Sent and stored auto-reply #%d", autoMsg.ID)
		}
	}
}
