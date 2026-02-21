package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/social-media-lead/backend/internal/models"
	"github.com/social-media-lead/backend/internal/store"
)

type WorkflowHandler struct {
	Store store.Store
}

// CreateWorkflowRequest holds data for creating/updating a workflow
type CreateWorkflowRequest struct {
	Name        string          `json:"name" binding:"required"`
	TriggerType string          `json:"trigger_type" binding:"required"`
	Status      string          `json:"status" binding:"required"`
	Prompt      string          `json:"prompt"`
	Nodes       json.RawMessage `json:"nodes" binding:"required"`
	Edges       json.RawMessage `json:"edges" binding:"required"`
}

func (h *WorkflowHandler) ListWorkflows(c *gin.Context) {
	userID := c.GetInt64("user_id")

	workflows, err := h.Store.GetWorkflowsByUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch workflows"})
		return
	}

	c.JSON(http.StatusOK, workflows)
}

func (h *WorkflowHandler) GetWorkflow(c *gin.Context) {
	userID := c.GetInt64("user_id")
	workflowID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workflow ID"})
		return
	}

	workflow, err := h.Store.GetWorkflowByID(c.Request.Context(), workflowID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workflow not found"})
		return
	}

	if workflow.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.JSON(http.StatusOK, workflow)
}

func (h *WorkflowHandler) CreateWorkflow(c *gin.Context) {
	userID := c.GetInt64("user_id")

	var req CreateWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	w := &models.Workflow{
		UserID:      userID,
		Name:        req.Name,
		TriggerType: req.TriggerType,
		Status:      req.Status,
		Prompt:      req.Prompt,
		Nodes:       []byte(req.Nodes),
		Edges:       []byte(req.Edges),
	}

	if err := h.Store.CreateWorkflow(c.Request.Context(), w); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create workflow"})
		return
	}

	c.JSON(http.StatusCreated, w)
}

func (h *WorkflowHandler) UpdateWorkflow(c *gin.Context) {
	userID := c.GetInt64("user_id")
	workflowID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workflow ID"})
		return
	}

	var req CreateWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify ownership before update
	existing, err := h.Store.GetWorkflowByID(c.Request.Context(), workflowID)
	if err != nil || existing.UserID != userID {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workflow not found"})
		return
	}

	existing.Name = req.Name
	existing.Status = req.Status
	existing.Prompt = req.Prompt
	existing.Nodes = []byte(req.Nodes)
	existing.Edges = []byte(req.Edges)

	if err := h.Store.UpdateWorkflow(c.Request.Context(), existing); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update workflow"})
		return
	}

	c.JSON(http.StatusOK, existing)
}

func (h *WorkflowHandler) DeleteWorkflow(c *gin.Context) {
	userID := c.GetInt64("user_id")
	workflowID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workflow ID"})
		return
	}

	if err := h.Store.DeleteWorkflow(c.Request.Context(), workflowID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete workflow"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Workflow deleted successfully"})
}
