package models

import "encoding/json"

// ============================================
// Workflow Graph Definitions (React Flow Compat)
// ============================================

// NodeType defines the specific capability of the node
type NodeType string

const (
	NodeTypeTriggerDM       NodeType = "trigger_meta_dm"
	NodeTypeTriggerKeyword  NodeType = "trigger_keyword"
	
	// Native Actions
	NodeTypeActionSendMessage NodeType = "action_send_message"
	NodeTypeActionDelay       NodeType = "action_delay"
	NodeTypeActionAddTag      NodeType = "action_add_tag"
	
	// AI Powered Actions
	NodeTypeActionAIReply     NodeType = "action_ai_reply" // Generates a response and sends it
	NodeTypeActionRAGSearch   NodeType = "action_rag_search" // Queries knowledge base
	NodeTypeLogicAIRouter      NodeType = "logic_ai_router" // Classifies intent to branch path
)

// ReactFlowNode represents a single block on the visual builder canvas
type ReactFlowNode struct {
	ID       string                 `json:"id"`
	Type     NodeType               `json:"type"`
	Position map[string]float64     `json:"position"`
	Data     map[string]interface{} `json:"data"` // Configuration specific to the node type
}

// ReactFlowEdge represents a connection between two nodes
type ReactFlowEdge struct {
	ID           string `json:"id"`
	Source       string `json:"source"`
	SourceHandle string `json:"sourceHandle,omitempty"` // For nodes with multiple outputs (like AI Router)
	Target       string `json:"target"`
	TargetHandle string `json:"targetHandle,omitempty"`
}

// WorkflowGraph is the unmarshaled version of the `nodes` and `edges` JSONB columns
type WorkflowGraph struct {
	Nodes []ReactFlowNode `json:"nodes"`
	Edges []ReactFlowEdge `json:"edges"`
}

// ParseWorkflowGraph converts the raw JSONB bytes into the typed Go structs
func ParseWorkflowGraph(nodesBytes, edgesBytes []byte) (*WorkflowGraph, error) {
	var nodes []ReactFlowNode
	var edges []ReactFlowEdge

	if err := json.Unmarshal(nodesBytes, &nodes); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(edgesBytes, &edges); err != nil {
		return nil, err
	}

	return &WorkflowGraph{
		Nodes: nodes,
		Edges: edges,
	}, nil
}
