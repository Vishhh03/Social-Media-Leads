package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/social-media-lead/backend/internal/cache"
	"github.com/social-media-lead/backend/internal/config"
	"github.com/social-media-lead/backend/internal/engine"
	"github.com/social-media-lead/backend/internal/meta"
	"github.com/social-media-lead/backend/internal/models"
	"github.com/social-media-lead/backend/internal/store"
)

// WebhookHandler handles Meta platform webhook events.
type WebhookHandler struct {
	Store       store.Store
	Config      *config.Config
	MetaClient  *meta.Client
	GraphWalker *engine.GraphWalker
	Cache       *cache.RedisClient
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

	// 0. Idempotency Check
	if h.Cache != nil && platformMsgID != "" {
		isNew, err := h.Cache.MarkWebhookProcessed(ctx, platformMsgID)
		if err != nil {
			log.Printf("[Webhook] Redis error checking idempotency: %v", err)
		} else if !isNew {
			log.Printf("[Webhook] Ignored duplicate message: %s", platformMsgID)
			return
		}
	}

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

	// 4. Handle Property Visit Q&A Flow
	if h.processVisitBookingFlow(ctx, channel, contact, content) {
		// If flow handled it, skip generic workflow orchestrator
		return
	}

	// 5. Trigger the new Workflow DAG Orchestrator
	h.triggerWorkflows(ctx, channel, contact, content)

	// Legacy automation triggers
	h.checkAutomationTriggers(ctx, channel, contact, content)
}

// processVisitBookingFlow runs the Property Visit state machine logic.
// Returns true if the flow sent an automated reply.
func (h *WebhookHandler) processVisitBookingFlow(ctx context.Context, channel *models.Channel, contact *models.Contact, content string) bool {
	// Escape Hatch: If agent replied manually recently, bot is paused.
	if contact.BotPaused {
		log.Printf("[BookingFlow] Interaction skipped for contact %d (Bot Paused by Agent)", contact.ID)
		return false
	}

	// ---- Load tenant config (Redis → Postgres fallback) ----
	cfg := h.loadTenantConfig(ctx, channel.UserID)
	if cfg == nil {
		// No wizard set up yet for this user — skip booking flow entirely
		log.Printf("[BookingFlow] No config for user %d — skipping booking flow", channel.UserID)
		return false
	}

	contentLower := strings.ToLower(strings.TrimSpace(content))
	var reply string

	switch contact.BookingState {
	case "new", "":
		reply = fmt.Sprintf("Hi 👋 Thanks for your interest in %s!\n\nAre you looking for:\n1. Self-use\n2. Investment", cfg.ProjectName)
		contact.BookingState = "qualified"

	case "qualified":
		if cfg.BrochureURL != "" {
			reply = fmt.Sprintf("Great choice! Here is the %s brochure: %s\n\nWould you like to schedule a site visit? We have slots tomorrow at 10 AM, 2 PM, or 4 PM. Reply with your preferred time.", cfg.ProjectName, cfg.BrochureURL)
		} else {
			reply = fmt.Sprintf("Great! Would you like to schedule a site visit for %s? We have slots tomorrow at 10 AM, 2 PM, or 4 PM. Reply with your preferred time.", cfg.ProjectName)
		}
		contact.BookingState = "offered_slots"

	case "offered_slots":
		slot := ""
		visitTime := time.Now().AddDate(0, 0, 1)
		if strings.Contains(contentLower, "10") {
			slot = "10:00 AM"
			visitTime = time.Date(visitTime.Year(), visitTime.Month(), visitTime.Day(), 10, 0, 0, 0, visitTime.Location())
		} else if strings.Contains(contentLower, "2") {
			slot = "2:00 PM"
			visitTime = time.Date(visitTime.Year(), visitTime.Month(), visitTime.Day(), 14, 0, 0, 0, visitTime.Location())
		} else if strings.Contains(contentLower, "4") {
			slot = "4:00 PM"
			visitTime = time.Date(visitTime.Year(), visitTime.Month(), visitTime.Day(), 16, 0, 0, 0, visitTime.Location())
		} else {
			reply = "I can answer more questions, but to ensure you get the best experience, would you like to book a site visit? We have slots tomorrow at 10 AM, 2 PM, or 4 PM."
			break
		}

		// Prevent double booking via Redis TTL lock
		if h.Cache != nil {
			locked, err := h.Cache.ReserveSlot(ctx, cfg.ProjectName, visitTime, contact.ID, 5*time.Minute)
			if err != nil || !locked {
				reply = fmt.Sprintf("I'm sorry, the %s slot just got taken! Please choose another time: 10 AM, 2 PM, or 4 PM.", slot)
				break
			}
		}

		visit := &models.Visit{
			UserID:            channel.UserID,
			ContactID:         contact.ID,
			ProjectName:       cfg.ProjectName,
			VisitTime:         visitTime,
			Status:            "confirmed",
			LeadSourceChannel: contact.Platform,
		}

		if err := h.Store.CreateVisit(ctx, visit); err != nil {
			log.Printf("[BookingFlow] Failed to save visit: %v", err)
			reply = "There was an error booking your visit. Please hold on, our agent will contact you."
			contact.BotPaused = true
		} else {
			reply = fmt.Sprintf("Perfect! Your visit to %s is confirmed for tomorrow at %s. Our agent will be in touch shortly to confirm details.", cfg.ProjectName, slot)
			contact.BookingState = "booked"
			go h.sendAgentNotification(channel, contact, visitTime)
		}

	case "booked":
		reply = fmt.Sprintf("Your visit to %s is already confirmed! Our team will be in touch. 🏡", cfg.ProjectName)
	}

	_ = h.Store.UpdateContactState(ctx, contact.ID, contact.BookingState, contact.BotPaused)

	if reply != "" {
		h.sendAutoReply(ctx, channel, contact, reply)
		return true
	}
	return false
}

// loadTenantConfig fetches the wizard config from Redis (10 min TTL) or falls back to Postgres.
func (h *WebhookHandler) loadTenantConfig(ctx context.Context, userID int64) *models.PropertyVisitConfig {
	if h.Cache != nil {
		if data, _ := h.Cache.GetCachedVisitConfig(ctx, userID); len(data) > 0 {
			var cfg models.PropertyVisitConfig
			if err := json.Unmarshal(data, &cfg); err == nil {
				return &cfg
			}
		}
	}

	cfg, err := h.Store.GetPropertyVisitConfig(ctx, userID)
	if err != nil || cfg == nil {
		return nil
	}

	// Populate cache for subsequent webhooks
	if h.Cache != nil {
		if data, err := json.Marshal(cfg); err == nil {
			_ = h.Cache.CacheVisitConfig(ctx, userID, data)
		}
	}
	return cfg
}


func (h *WebhookHandler) sendAutoReply(ctx context.Context, channel *models.Channel, contact *models.Contact, text string) {
	if h.MetaClient == nil {
		log.Println("[BookingFlow] MetaClient is nil (test mode), skipping real API call.")
		return
	}
	result, err := h.MetaClient.SendMessage(
		contact.Platform,
		channel.AccountID,
		contact.PlatformUserID,
		text,
		channel.AccessToken,
	)
	if err != nil || !result.Success {
		log.Printf("[BookingFlow] Failed to send auto-reply: %v", err)
		return
	}

	autoMsg := &models.Message{
		UserID:        channel.UserID,
		ChannelID:     channel.ID,
		ContactID:     contact.ID,
		Platform:      contact.Platform,
		Direction:     "outbound",
		Content:       text,
		MessageType:   "text",
		PlatformMsgID: result.MessageID,
		Status:        "sent",
		IsAutomated:   true,
	}
	_ = h.Store.CreateMessage(ctx, autoMsg)
}

func (h *WebhookHandler) sendAgentNotification(channel *models.Channel, contact *models.Contact, visitTime time.Time) {
	// For MVP, structured structured console log format
	msg := fmt.Sprintf(`
=========================================
🔥 New Site Visit Booked
Project: Project Alpha
Lead: %s
Visit: %s
Open Chat: https://wa.me/%s
Summary: Intent (Qualified), Time Booked.
=========================================
`, contact.Name, visitTime.Format("Jan 02, 3:04 PM"), contact.Phone)
	
	log.Println(msg)
}


// triggerWorkflows executes any active DAG workflow matching the meta_dm_received trigger
func (h *WebhookHandler) triggerWorkflows(ctx context.Context, channel *models.Channel, contact *models.Contact, content string) {
	workflows, err := h.Store.GetActiveWorkflowsByTrigger(ctx, channel.UserID, "trigger_meta_dm")
	if err != nil {
		log.Printf("[Webhook] Failed to fetch active workflows: %v", err)
		return
	}

	for _, w := range workflows {
		log.Printf("[Webhook] Execution Engine starting Workflow %d: '%s'", w.ID, w.Name)
		
		initialState := map[string]interface{}{
			"received_message": content,
			"platform":         contact.Platform,
			"contact_name":     contact.Name,
		}
		
		// Run GraphWalker in a separate goroutine so it doesn't block the webhook response
		go func(workflowID, contactID int64, state map[string]interface{}) {
			err := h.GraphWalker.StartWorkflow(context.Background(), workflowID, contactID, state)
			if err != nil {
				log.Printf("[Engine] Workflow %d execution failed for contact %d: %v", workflowID, contactID, err)
			}
		}(w.ID, contact.ID, initialState)
	}
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
