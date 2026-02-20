package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/social-media-lead/backend/internal/ai"
)

type AIHandler struct {
	LLMClient ai.LLMClient
}

type GenerateWorkflowRequest struct {
	Prompt string `json:"prompt" binding:"required"`
}

// schema struct exactly matching OpenAI JSON Schema requirements
var workflowSchema = map[string]interface{}{
	"type": "object",
	"properties": map[string]interface{}{
		"nodes": map[string]interface{}{
			"type": "array",
			"items": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"id":   map[string]interface{}{"type": "string"},
					"type": map[string]interface{}{
						"type": "string",
						"enum": []string{
							"trigger_meta_dm", "trigger_keyword",
							"action_send_message", "action_delay",
							"action_ai_reply", "logic_ai_router",
						},
					},
					"position": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"x": map[string]interface{}{"type": "number"},
							"y": map[string]interface{}{"type": "number"},
						},
						"required": []string{"x", "y"},
						"additionalProperties": false,
					},
					"data": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"label":       map[string]interface{}{"type": "string"},
							"description": map[string]interface{}{"type": "string"},
							"message":     map[string]interface{}{"type": "string"},
							"prompt":      map[string]interface{}{"type": "string"},
							"isRouter":    map[string]interface{}{"type": "boolean"},
						},
						"additionalProperties": false,
					},
				},
				"required": []string{"id", "type", "position", "data"},
				"additionalProperties": false,
			},
		},
		"edges": map[string]interface{}{
			"type": "array",
			"items": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"id":           map[string]interface{}{"type": "string"},
					"source":       map[string]interface{}{"type": "string"},
					"target":       map[string]interface{}{"type": "string"},
					"sourceHandle": map[string]interface{}{"type": "string"},
				},
				"required": []string{"id", "source", "target"},
				"additionalProperties": false,
			},
		},
	},
	"required": []string{"nodes", "edges"},
	"additionalProperties": false,
}

func (h *AIHandler) GenerateWorkflow(c *gin.Context) {
	var req GenerateWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	systemPrompt := `You are an expert AI architect generating automation workflows for a Social Media Lead SaaS app.
The user will provide a desired behavior (e.g., "reply to DM pricing inquiries, wait a day, follow up").
Your goal is to output a strictly formatted Directed Acyclic Graph containing nodes and edges.

Valid Node Types:
- trigger_meta_dm: A new inbound Instagram/Messenger DM arrives.
- trigger_keyword: Fires if the message contains specific words.
- action_send_message: Sends a static text reply (put in data.message).
- action_delay: Pauses the workflow.
- action_ai_reply: Uses the Knowledge Base to answer a question (put instructions in data.prompt).
- logic_ai_router: Branches based on intent (Outputs: sourceHandle="hot" or "cold").

Requirements:
- Always start with exactly 1 trigger node (ID: "1", positioned at x: 250, y: 50).
- Sequence all subsequent nodes cleanly, spacing them vertically (y + 150 each).
- Ensure edges connect source node IDs to target node IDs sequentially.
- Use visually descriptive text for data.label and data.description.

User Prompt: ` + req.Prompt

	rawJSON, err := h.LLMClient.GenerateStructuredJSON(context.Background(), systemPrompt, workflowSchema)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate workflow via AI: " + err.Error()})
		return
	}

	// Just proxy the raw JSON string back to the client, it is guaranteed to match the layout schema
	var finalResponse map[string]interface{}
	if err := json.Unmarshal([]byte(rawJSON), &finalResponse); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse AI response"})
		return
	}

	c.JSON(http.StatusOK, finalResponse)
}
