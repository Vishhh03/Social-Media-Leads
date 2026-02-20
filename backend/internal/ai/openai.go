package ai

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

// OpenAIClient implements the LLMClient interface for OpenAI's models
type OpenAIClient struct {
	client *openai.Client
	model  string
}

// NewOpenAIClient initializes a new OpenAI provider
func NewOpenAIClient(apiKey string, model string) *OpenAIClient {
	if model == "" {
		model = openai.GPT4oMini // Default recommended for speed/cost
	}
	return &OpenAIClient{
		client: openai.NewClient(apiKey),
		model:  model,
	}
}

// GenerateText sends a simple prompt to the OpenAI Chat Completion API
func (c *OpenAIClient) GenerateText(ctx context.Context, prompt string) (string, error) {
	resp, err := c.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: c.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)
	if err != nil {
		return "", fmt.Errorf("openai error: %v", err)
	}
	return resp.Choices[0].Message.Content, nil
}

// GenerateEmbedding calls the OpenAI Embeddings API (typically used with text-embedding-3-small)
func (c *OpenAIClient) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	req := openai.EmbeddingRequest{
		Input: []string{text},
		Model: openai.SmallEmbedding3, // 1536 dimensions
	}

	resp, err := c.client.CreateEmbeddings(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("openai embedding error: %v", err)
	}

	return resp.Data[0].Embedding, nil
}

// GenerateStructuredJSON forces the LLM to reply with a JSON structure matching the provided `schema`
// Requires a model that supports JSON formats like JSON Schema (e.g. gpt-4o)
func (c *OpenAIClient) GenerateStructuredJSON(ctx context.Context, prompt string, schema any) (string, error) {
	schemaBytes, err := json.Marshal(schema)
	if err != nil {
		return "", fmt.Errorf("failed to marshal schema: %v", err)
	}

	rawSchema := json.RawMessage(schemaBytes)

	resp, err := c.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: c.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			ResponseFormat: &openai.ChatCompletionResponseFormat{
				Type: openai.ChatCompletionResponseFormatTypeJSONSchema,
				JSONSchema: &openai.ChatCompletionResponseFormatJSONSchema{
					Name:   "workflow_graph",
					Strict: true,
					Schema: &rawSchema,
				},
			},
		},
	)

	if err != nil {
		return "", fmt.Errorf("openai structured generation error: %v", err)
	}

	return resp.Choices[0].Message.Content, nil
}
