---
name: OpenAI Structured Outputs in Go
description: Enforcing strict JSON formatting in Go for robust agentic interactions.
---
# OpenAI Structured Outputs in Go

## Key Findings & Caveats
1. **Determinism is the Goal:** When building complex UIs like a React Flow DAG from AI prompts, standard JSON prompting is unreliable. The LLM might hallucinate keys, fail to close brackets, or create invalid node connections.
2. **Response format JSON Schema:** Use `openai.ChatCompletionResponseFormatTypeJSONSchema` with `Strict: true` when calling the OpenAI API `CreateChatCompletion` endpoint.
3. **Defining the Schema:** The schema must be defined meticulously to represent exactly what React or the Backend expects (e.g., `id`, `type`, `position.x`, `position.y`, `data`). 
4. **Enforcing `additionalProperties: false`:** OpenAI's structured outputs dictate that ALL objects and child objects inside the JSON schema must explicitly define `"additionalProperties": false` to guarantee the model does not inject unhandled data. If any nest misses this, OpenAI will throw a 400 validation error.
5. **Enums are Powerful:** Limit hallucination of categorical states (like the specific Node Types available in the frontend canvas) by defining `"enum": ["trigger_meta_dm", "action_send_message"]` in the schema type definition.
