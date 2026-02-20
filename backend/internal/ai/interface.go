package ai

import (
	"context"
)

// LLMClient defines the interface for interacting with Language Models
// This allows us to hot-swap between OpenAI, Anthropic, Gemini, or Local LLMs
type LLMClient interface {
	// GenerateText sends a prompt and receives a string response
	// Useful for AI replies in workflows
	GenerateText(ctx context.Context, prompt string) (string, error)

	// GenerateEmbedding takes a string and converts it to a vector array
	// Useful for RAG operations with pgvector
	GenerateEmbedding(ctx context.Context, text string) ([]float32, error)

	// GenerateStructuredJSON sends a prompt and forces the LLM to reply matching
	// a specific JSON schema.
	// Useful for the "Magic Prompt" visual flow builder
	GenerateStructuredJSON(ctx context.Context, prompt string, schema any) (string, error)
}
