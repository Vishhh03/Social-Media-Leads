package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
	"github.com/social-media-lead/backend/internal/engine"
)

const (
	TaskResumeWorkflow = "workflow:resume"
)

// ResumeWorkflowPayload represents the data sent to the background job
type ResumeWorkflowPayload struct {
	ExecutionID int64 `json:"execution_id"`
}

// NewResumeWorkflowTask creates a new Asynq task for resuming a workflow
func NewResumeWorkflowTask(executionID int64) (*asynq.Task, error) {
	payload, err := json.Marshal(ResumeWorkflowPayload{ExecutionID: executionID})
	if err != nil {
		return nil, err
	}
	// no deadline specified here, using Asynq defaults
	return asynq.NewTask(TaskResumeWorkflow, payload), nil
}

// HandleResumeWorkflowTask processes the resume workflow job
func HandleResumeWorkflowTask(graphWalker *engine.GraphWalker) func(context.Context, *asynq.Task) error {
	return func(ctx context.Context, t *asynq.Task) error {
		var p ResumeWorkflowPayload
		if err := json.Unmarshal(t.Payload(), &p); err != nil {
			return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
		}

		log.Printf("[Worker] Starting to resume workflow execution %d", p.ExecutionID)

		err := graphWalker.ResumeExecution(ctx, p.ExecutionID)
		if err != nil {
			log.Printf("[Worker] Execution %d failed: %v", p.ExecutionID, err)
			return err // Return err to retry, or Wrap and return SkipRetry to drop
		}

		log.Printf("[Worker] Successfully resumed execution %d", p.ExecutionID)
		return nil
	}
}
