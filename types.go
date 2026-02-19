package llm

// Message represents a chat message.
// Content is either a string (text-only) or []ContentPart for multimodal (text + images).
type Message struct {
	Role       string     `json:"role"`
	Content    any        `json:"content"`
	Name       string     `json:"name,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
}

// ContentPart represents a part of multimodal content.
// Type can be "text", "image_url", "video_url", "input_audio", or "file".
type ContentPart struct {
	Type       string      `json:"type"`
	Text       string      `json:"text,omitempty"`
	ImageURL   *ImageURL   `json:"image_url,omitempty"`
	VideoURL   *ImageURL   `json:"video_url,omitempty"`
	InputAudio *InputAudio `json:"input_audio,omitempty"`
	File       *FileData   `json:"file,omitempty"`
}

// ImageURL represents an image or video URL for vision (supports https:// or data:...;base64,...).
type ImageURL struct {
	URL    string `json:"url"`
	Detail string `json:"detail,omitempty"`
}

// InputAudio represents audio input for speech-capable models (base64 data + format).
type InputAudio struct {
	Data   string `json:"data"`
	Format string `json:"format"` // e.g. "mp3", "wav"
}

// FileData represents a document file (PDF, etc.) for OpenRouter.
type FileData struct {
	Filename string `json:"filename"`
	FileData string `json:"file_data"` // URL or data:application/pdf;base64,...
}

// Tool represents a tool/function definition.
type Tool struct {
	Type     string      `json:"type"`
	Function FunctionDef `json:"function"`
}

// FunctionDef defines a function for tool calling.
type FunctionDef struct {
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Parameters  map[string]any `json:"parameters,omitempty"`
}

// ToolCall represents a tool call in the response.
// Index is used when accumulating streaming deltas (OpenAI format).
type ToolCall struct {
	Index    int          `json:"index,omitempty"`
	ID       string       `json:"id"`
	Type     string       `json:"type"`
	Function FunctionCall `json:"function"`
}

// FunctionCall represents a function call in the response.
type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// ResponseFormat specifies the output format.
type ResponseFormat struct {
	Type       string         `json:"type"`
	JSONSchema *JSONSchemaDef `json:"json_schema,omitempty"`
}

// JSONSchemaDef defines a JSON schema for structured output.
type JSONSchemaDef struct {
	Name   string         `json:"name"`
	Strict bool           `json:"strict,omitempty"`
	Schema map[string]any `json:"schema"`
}

// ChatRequest is the request for chat completions.
type ChatRequest struct {
	Model             string          `json:"model,omitempty"`
	Messages          []Message       `json:"messages"`
	Temperature       *float64        `json:"temperature,omitempty"`
	TopP              *float64        `json:"top_p,omitempty"`
	TopK              *int            `json:"top_k,omitempty"`
	MaxTokens         *int            `json:"max_tokens,omitempty"`
	Stop              any             `json:"stop,omitempty"`
	Stream            bool            `json:"stream,omitempty"`
	Seed              *int            `json:"seed,omitempty"`
	PresencePenalty   *float64        `json:"presence_penalty,omitempty"`
	FrequencyPenalty  *float64        `json:"frequency_penalty,omitempty"`
	ResponseFormat    *ResponseFormat `json:"response_format,omitempty"`
	Tools             []Tool          `json:"tools,omitempty"`
	ToolChoice        any             `json:"tool_choice,omitempty"`
	ParallelToolCalls *bool           `json:"parallel_tool_calls,omitempty"`
	User              string          `json:"user,omitempty"`
}

// ChatResponse is the response from chat completions.
type ChatResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   *Usage   `json:"usage,omitempty"`
}

// ChoiceError represents provider error details in a choice (e.g. streaming partial failure).
type ChoiceError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Choice represents a completion choice.
type Choice struct {
	Index        int          `json:"index"`
	Message      *Message     `json:"message,omitempty"`
	Delta        *Message     `json:"delta,omitempty"`
	FinishReason string       `json:"finish_reason,omitempty"`
	Error        *ChoiceError `json:"error,omitempty"`
}

// Usage represents token usage.
type Usage struct {
	PromptTokens     int     `json:"prompt_tokens"`
	CompletionTokens int     `json:"completion_tokens"`
	TotalTokens      int     `json:"total_tokens"`
	Cost             float64 `json:"cost,omitempty"`
}

// StreamChunk represents a single streaming chunk.
type StreamChunk struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   *Usage   `json:"usage,omitempty"`
}

// EmbeddingRequest is the request for creating embeddings.
// Input must be a non-empty string, []string, or []interface{} (with string elements).
type EmbeddingRequest struct {
	Model string `json:"model"`
	Input any    `json:"input"`
}

// EmbeddingResponse is the response from creating embeddings.
type EmbeddingResponse struct {
	Data  []EmbeddingData `json:"data"`
	Usage *EmbeddingUsage `json:"usage,omitempty"`
}

// EmbeddingData holds a single embedding.
type EmbeddingData struct {
	Object    string    `json:"object"`
	Embedding []float64 `json:"embedding"`
	Index     int       `json:"index"`
}

// EmbeddingUsage represents token usage for embeddings.
type EmbeddingUsage struct {
	PromptTokens int `json:"prompt_tokens"`
	TotalTokens  int `json:"total_tokens"`
}
