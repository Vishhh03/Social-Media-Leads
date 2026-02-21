package workers

import (
	"log"

	"github.com/hibiken/asynq"
	"github.com/social-media-lead/backend/internal/engine"
)

// StartServer starts the Asynq worker server to process background jobs
func StartServer(redisAddr string, graphWalker *engine.GraphWalker) *asynq.Server {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: 10,
			// Optionally specify multiple queues with different priority.
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	// mux maps a type to a handler
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskResumeWorkflow, HandleResumeWorkflowTask(graphWalker))

	// start the background server process
	go func() {
		if err := srv.Run(mux); err != nil {
			log.Fatalf("could not run asynq server: %v", err)
		}
	}()
	
	return srv
}
