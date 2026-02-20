package store

import (
	"context"

	"github.com/social-media-lead/backend/internal/models"
)

// CreateKnowledgeBaseEntry adds a new RAG document chunk and its embedding to the DB.
func (s *Storage) CreateKnowledgeBaseEntry(ctx context.Context, entry *models.KnowledgeBase, embedding []float32) error {
	query := `
		INSERT INTO knowledge_base (user_id, title, content, embedding)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`
	return s.DB.QueryRow(ctx, query,
		entry.UserID,
		entry.Title,
		entry.Content,
		embedding, // pgvector handles []float32 mapping automatically with pgx
	).Scan(&entry.ID, &entry.CreatedAt, &entry.UpdatedAt)
}

// GetKnowledgeBaseEntriesByUser retrieves all KB entries for a user (without embeddings).
func (s *Storage) GetKnowledgeBaseEntriesByUser(ctx context.Context, userID int64) ([]models.KnowledgeBase, error) {
	query := `
		SELECT id, user_id, title, content, created_at, updated_at
		FROM knowledge_base
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	rows, err := s.DB.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []models.KnowledgeBase
	for rows.Next() {
		var kb models.KnowledgeBase
		if err := rows.Scan(
			&kb.ID, &kb.UserID, &kb.Title, &kb.Content,
			&kb.CreatedAt, &kb.UpdatedAt,
		); err != nil {
			return nil, err
		}
		entries = append(entries, kb)
	}
	return entries, nil
}

// SearchKnowledgeBase finds top `limit` relevant documents based on cosine similarity `<=>`.
func (s *Storage) SearchKnowledgeBase(ctx context.Context, userID int64, queryEmbedding []float32, limit int) ([]models.KnowledgeBase, error) {
	query := `
		SELECT id, user_id, title, content, created_at, updated_at
		FROM knowledge_base
		WHERE user_id = $1
		ORDER BY embedding <=> $2
		LIMIT $3
	`
	rows, err := s.DB.Query(ctx, query, userID, queryEmbedding, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []models.KnowledgeBase
	for rows.Next() {
		var kb models.KnowledgeBase
		if err := rows.Scan(
			&kb.ID, &kb.UserID, &kb.Title, &kb.Content,
			&kb.CreatedAt, &kb.UpdatedAt,
		); err != nil {
			return nil, err
		}
		entries = append(entries, kb)
	}
	return entries, nil
}

// DeleteKnowledgeBaseEntry removes a document chunk.
func (s *Storage) DeleteKnowledgeBaseEntry(ctx context.Context, entryID, userID int64) error {
	query := `DELETE FROM knowledge_base WHERE id = $1 AND user_id = $2`
	_, err := s.DB.Exec(ctx, query, entryID, userID)
	return err
}

// ============================================
// Workflow Blueprints
// ============================================

func (s *Storage) CreateWorkflow(ctx context.Context, w *models.Workflow) error {
	query := `
		INSERT INTO workflows (user_id, name, trigger_type, status, prompt, nodes, edges)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`
	return s.DB.QueryRow(ctx, query,
		w.UserID, w.Name, w.TriggerType, w.Status, w.Prompt, w.Nodes, w.Edges,
	).Scan(&w.ID, &w.CreatedAt, &w.UpdatedAt)
}

func (s *Storage) GetWorkflowByID(ctx context.Context, workflowID int64) (*models.Workflow, error) {
	query := `
		SELECT id, user_id, name, trigger_type, status, prompt, nodes, edges, created_at, updated_at
		FROM workflows WHERE id = $1
	`
	var w models.Workflow
	err := s.DB.QueryRow(ctx, query, workflowID).Scan(
		&w.ID, &w.UserID, &w.Name, &w.TriggerType, &w.Status, &w.Prompt,
		&w.Nodes, &w.Edges, &w.CreatedAt, &w.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &w, nil
}

func (s *Storage) GetWorkflowsByUser(ctx context.Context, userID int64) ([]models.Workflow, error) {
	query := `
		SELECT id, user_id, name, trigger_type, status, prompt, nodes, edges, created_at, updated_at
		FROM workflows WHERE user_id = $1 ORDER BY created_at DESC
	`
	rows, err := s.DB.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var flows []models.Workflow
	for rows.Next() {
		var w models.Workflow
		if err := rows.Scan(
			&w.ID, &w.UserID, &w.Name, &w.TriggerType, &w.Status, &w.Prompt,
			&w.Nodes, &w.Edges, &w.CreatedAt, &w.UpdatedAt,
		); err != nil {
			return nil, err
		}
		flows = append(flows, w)
	}
	return flows, nil
}

func (s *Storage) GetActiveWorkflowsByTrigger(ctx context.Context, userID int64, triggerType string) ([]models.Workflow, error) {
	query := `
		SELECT id, user_id, name, trigger_type, status, prompt, nodes, edges, created_at, updated_at
		FROM workflows 
		WHERE user_id = $1 AND trigger_type = $2 AND status = 'published'
	`
	rows, err := s.DB.Query(ctx, query, userID, triggerType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var flows []models.Workflow
	for rows.Next() {
		var w models.Workflow
		if err := rows.Scan(
			&w.ID, &w.UserID, &w.Name, &w.TriggerType, &w.Status, &w.Prompt,
			&w.Nodes, &w.Edges, &w.CreatedAt, &w.UpdatedAt,
		); err != nil {
			return nil, err
		}
		flows = append(flows, w)
	}
	return flows, nil
}

func (s *Storage) UpdateWorkflow(ctx context.Context, w *models.Workflow) error {
	query := `
		UPDATE workflows 
		SET name = $1, status = $2, prompt = $3, nodes = $4, edges = $5, updated_at = NOW()
		WHERE id = $6 AND user_id = $7
		RETURNING updated_at
	`
	return s.DB.QueryRow(ctx, query,
		w.Name, w.Status, w.Prompt, w.Nodes, w.Edges, w.ID, w.UserID,
	).Scan(&w.UpdatedAt)
}

func (s *Storage) DeleteWorkflow(ctx context.Context, workflowID, userID int64) error {
	query := `DELETE FROM workflows WHERE id = $1 AND user_id = $2`
	_, err := s.DB.Exec(ctx, query, workflowID, userID)
	return err
}

// ============================================
// Workflow Executions (Running State)
// ============================================

func (s *Storage) CreateWorkflowExecution(ctx context.Context, exec *models.WorkflowExecution) error {
	query := `
		INSERT INTO workflow_executions (workflow_id, contact_id, current_node_id, status, state_data)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	return s.DB.QueryRow(ctx, query,
		exec.WorkflowID, exec.ContactID, exec.CurrentNodeID, exec.Status, exec.StateData,
	).Scan(&exec.ID, &exec.CreatedAt, &exec.UpdatedAt)
}

func (s *Storage) GetWorkflowExecutionByID(ctx context.Context, executionID int64) (*models.WorkflowExecution, error) {
	query := `
		SELECT id, workflow_id, contact_id, current_node_id, status, state_data, created_at, updated_at
		FROM workflow_executions WHERE id = $1
	`
	var exec models.WorkflowExecution
	err := s.DB.QueryRow(ctx, query, executionID).Scan(
		&exec.ID, &exec.WorkflowID, &exec.ContactID, &exec.CurrentNodeID,
		&exec.Status, &exec.StateData, &exec.CreatedAt, &exec.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &exec, nil
}

func (s *Storage) UpdateWorkflowExecution(ctx context.Context, exec *models.WorkflowExecution) error {
	query := `
		UPDATE workflow_executions 
		SET current_node_id = $1, status = $2, state_data = $3, updated_at = NOW()
		WHERE id = $4
		RETURNING updated_at
	`
	return s.DB.QueryRow(ctx, query,
		exec.CurrentNodeID, exec.Status, exec.StateData, exec.ID,
	).Scan(&exec.UpdatedAt)
}
