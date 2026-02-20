package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/social-media-lead/backend/internal/ai"
	"github.com/social-media-lead/backend/internal/models"
	"github.com/social-media-lead/backend/internal/store"
)

// GraphWalker is responsible for traversing a Workflow DAG and executing node logic
type GraphWalker struct {
	Store     store.Store
	LLMClient ai.LLMClient
}

func NewGraphWalker(store store.Store, llmClient ai.LLMClient) *GraphWalker {
	return &GraphWalker{
		Store:     store,
		LLMClient: llmClient,
	}
}

// StartWorkflow initiates a new workflow execution for a given contact
func (gw *GraphWalker) StartWorkflow(ctx context.Context, workflowID, contactID int64, initialState map[string]interface{}) error {
	w, err := gw.Store.GetWorkflowByID(ctx, workflowID)
	if err != nil {
		return fmt.Errorf("failed to get workflow: %w", err)
	}

	graph, err := models.ParseWorkflowGraph(w.Nodes, w.Edges)
	if err != nil {
		return fmt.Errorf("failed to parse workflow graph: %w", err)
	}

	// Find the trigger node
	var startNode *models.ReactFlowNode
	for _, n := range graph.Nodes {
		if n.Type == models.NodeTypeTriggerDM || n.Type == models.NodeTypeTriggerKeyword {
			// Re-assign explicitly because implicit memory address of loop var is bad conceptually
			nCopy := n
			startNode = &nCopy
			break
		}
	}

	if startNode == nil {
		return fmt.Errorf("no trigger node found in workflow %d", workflowID)
	}

	stateBytes, _ := json.Marshal(initialState)

	exec := &models.WorkflowExecution{
		WorkflowID:    workflowID,
		ContactID:     contactID,
		CurrentNodeID: startNode.ID,
		Status:        "running",
		StateData:     stateBytes,
	}

	if err := gw.Store.CreateWorkflowExecution(ctx, exec); err != nil {
		return fmt.Errorf("failed to create execution: %w", err)
	}

	// Begin execution loop
	return gw.ResumeExecution(ctx, exec.ID)
}

// ResumeExecution picks up an execution from its CurrentNodeID and walks the DAG
func (gw *GraphWalker) ResumeExecution(ctx context.Context, executionID int64) error {
	exec, err := gw.Store.GetWorkflowExecutionByID(ctx, executionID)
	if err != nil {
		return err
	}

	w, err := gw.Store.GetWorkflowByID(ctx, exec.WorkflowID)
	if err != nil {
		return err
	}

	graph, err := models.ParseWorkflowGraph(w.Nodes, w.Edges)
	if err != nil {
		return err
	}

	var stateData map[string]interface{}
	if err := json.Unmarshal(exec.StateData, &stateData); err != nil {
		stateData = make(map[string]interface{})
	}

	currentNodeID := exec.CurrentNodeID

	for {
		// Find current node
		node := findNode(graph.Nodes, currentNodeID)
		if node == nil {
			exec.Status = "completed"
			gw.Store.UpdateWorkflowExecution(ctx, exec)
			log.Printf("Execution %d completed. Node %s not found (end of flow).", executionID, currentNodeID)
			return nil
		}

		// Execute node logic
		log.Printf("Executing Node %s (%s) for Execution %d", node.ID, node.Type, executionID)
		
		nextNodeID, err := gw.processNode(ctx, node, graph, exec, stateData)
		if err != nil {
			exec.Status = "failed"
			gw.Store.UpdateWorkflowExecution(ctx, exec)
			return fmt.Errorf("node %s failed: %w", node.ID, err)
		}

		// Update state in DB
		newStateBytes, _ := json.Marshal(stateData)
		exec.StateData = newStateBytes
		
		if nextNodeID == "" {
			// Flow finished
			exec.Status = "completed"
			exec.CurrentNodeID = node.ID // Keep last valid node
			gw.Store.UpdateWorkflowExecution(ctx, exec)
			log.Printf("Execution %d completed successfully.", executionID)
			return nil
		}

		// If it's a delay node, we would pause here and rely on Asynq to resume later
		if node.Type == models.NodeTypeActionDelay {
			exec.Status = "waiting"
			exec.CurrentNodeID = nextNodeID
			gw.Store.UpdateWorkflowExecution(ctx, exec)
			
			// TODO: Phase 4 - Schedule Asynq worker using the configured delay
			log.Printf("Execution %d paused at Delay node %s. Will resume at %s", executionID, node.ID, nextNodeID)
			return nil
		}

		// Move to next node
		currentNodeID = nextNodeID
		exec.CurrentNodeID = currentNodeID
		gw.Store.UpdateWorkflowExecution(ctx, exec)
	}
}

func (gw *GraphWalker) processNode(ctx context.Context, node *models.ReactFlowNode, graph *models.WorkflowGraph, exec *models.WorkflowExecution, stateData map[string]interface{}) (string, error) {
	// Execute specific behaviors
	switch node.Type {
	case models.NodeTypeTriggerDM, models.NodeTypeTriggerKeyword:
		// Triggers just pass through, state is already populated by StartWorkflow
		log.Printf("Processing Trigger: %v", node.Data["label"])
		return gw.findNextNode(graph.Edges, node.ID, ""), nil

	case models.NodeTypeActionSendMessage:
		// Send a message using Meta API (Mocked for now)
		msg := "Hello!"
		if val, ok := node.Data["message"]; ok {
			msg = val.(string)
		}
		log.Printf("[Meta API] Sending message to Contact %d: %s", exec.ContactID, msg)
		return gw.findNextNode(graph.Edges, node.ID, ""), nil
		
	case models.NodeTypeActionAIReply:
		// Read prompt from UI
		prompt := "Reply to the user's message."
		if val, ok := node.Data["prompt"]; ok {
			prompt = val.(string)
		}
		
		// If knowledge base RAG was executed before this, context would be in stateData["kb_context"]
		userMsg := ""
		if val, ok := stateData["received_message"]; ok {
			userMsg = val.(string)
		}
		
		fullPrompt := fmt.Sprintf("%s\n\nUser Message: %s", prompt, userMsg)
		
		// Call LLM
		reply, err := gw.LLMClient.GenerateText(ctx, fullPrompt)
		if err != nil {
			return "", err
		}
		
		log.Printf("[Meta API -> AI] Sending generated reply to Contact %d: %s", exec.ContactID, reply)
		return gw.findNextNode(graph.Edges, node.ID, ""), nil

	case models.NodeTypeActionDelay:
		// For delay, we just return the next node to schedule
		log.Printf("Delay node executed")
		return gw.findNextNode(graph.Edges, node.ID, ""), nil

	default:
		log.Printf("Unknown node type: %s", node.Type)
		return gw.findNextNode(graph.Edges, node.ID, ""), nil
	}
}

func (gw *GraphWalker) findNextNode(edges []models.ReactFlowEdge, sourceNodeID, sourceHandle string) string {
	for _, edge := range edges {
		if edge.Source == sourceNodeID {
			// If a specific source handle is requested (e.g. AI intent routing), match it
			if sourceHandle != "" && edge.SourceHandle != sourceHandle {
				continue
			}
			return edge.Target
		}
	}
	return ""
}

func findNode(nodes []models.ReactFlowNode, nodeID string) *models.ReactFlowNode {
	for _, n := range nodes {
		if n.ID == nodeID {
			nCopy := n
			return &nCopy
		}
	}
	return nil
}
