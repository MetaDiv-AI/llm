package llm

import (
	"context"
	"io"

	"github.com/MetaDiv-AI/openrouter"
	"github.com/MetaDiv-AI/openrouter/chat"
	"github.com/MetaDiv-AI/openrouter/embeddings"
)

type openRouterChat struct {
	or *openrouter.Client
}

type openRouterEmbedding struct {
	or *openrouter.Client
}

func newOpenRouterClient(cfg *config) (*Client, error) {
	opts := buildOpenRouterOptions(cfg)
	or, err := openrouter.NewClient(opts...)
	if err != nil {
		return nil, err
	}

	return &Client{
		Chat:      &openRouterChat{or: or},
		Embeddings: &openRouterEmbedding{or: or},
	}, nil
}

func buildOpenRouterOptions(cfg *config) []openrouter.Option {
	opts := []openrouter.Option{openrouter.WithAPIKey(cfg.APIKey)}
	if cfg.BaseURL != "" {
		opts = append(opts, openrouter.WithBaseURL(cfg.BaseURL))
	}
	if cfg.Timeout > 0 {
		opts = append(opts, openrouter.WithTimeout(cfg.Timeout))
	}
	if cfg.MaxRetries >= 0 {
		opts = append(opts, openrouter.WithMaxRetries(cfg.MaxRetries))
	}
	if len(cfg.Headers) > 0 {
		opts = append(opts, openrouter.WithHeaders(cfg.Headers))
	}
	if cfg.Debug {
		opts = append(opts, openrouter.WithDebug(true))
	}
	if cfg.Referer != "" {
		opts = append(opts, openrouter.WithReferer(cfg.Referer))
	}
	if cfg.Title != "" {
		opts = append(opts, openrouter.WithTitle(cfg.Title))
	}
	if cfg.ForwardedFor != "" {
		opts = append(opts, openrouter.WithForwardedFor(cfg.ForwardedFor))
	}
	if cfg.Logger != nil {
		opts = append(opts, openrouter.WithLogger(cfg.Logger))
	}
	return opts
}

func validateChatRequest(req *ChatRequest) error {
	if req == nil {
		return &ValidationError{Field: "request", Message: "cannot be nil"}
	}
	if req.Model == "" {
		return &ValidationError{Field: "model", Message: "cannot be empty"}
	}
	if len(req.Messages) == 0 {
		return &ValidationError{Field: "messages", Message: "cannot be empty"}
	}
	return nil
}

func validateEmbeddingRequest(req *EmbeddingRequest) error {
	if req == nil {
		return &ValidationError{Field: "request", Message: "cannot be nil"}
	}
	if req.Model == "" {
		return &ValidationError{Field: "model", Message: "cannot be empty"}
	}
	if req.Input == nil {
		return &ValidationError{Field: "input", Message: "cannot be nil"}
	}
	switch v := req.Input.(type) {
	case string:
		if v == "" {
			return &ValidationError{Field: "input", Message: "cannot be empty string"}
		}
	case []string:
		if len(v) == 0 {
			return &ValidationError{Field: "input", Message: "cannot be empty slice"}
		}
	case []interface{}:
		if len(v) == 0 {
			return &ValidationError{Field: "input", Message: "cannot be empty slice"}
		}
	}
	return nil
}

func (c *openRouterChat) Create(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	if err := validateChatRequest(req); err != nil {
		return nil, err
	}
	orReq := toORChatRequest(req)
	resp, err := c.or.Chat.Create(ctx, orReq)
	if err != nil {
		return nil, err
	}
	return toLLMChatResponse(resp), nil
}

func (c *openRouterChat) CreateStream(ctx context.Context, req *ChatRequest) (StreamReader, error) {
	if err := validateChatRequest(req); err != nil {
		return nil, err
	}
	orReq := toORChatRequest(req)
	orStream, err := c.or.Chat.CreateStream(ctx, orReq)
	if err != nil {
		return nil, err
	}
	return &orStreamReader{inner: orStream}, nil
}

func (c *openRouterEmbedding) Create(ctx context.Context, req *EmbeddingRequest) (*EmbeddingResponse, error) {
	if err := validateEmbeddingRequest(req); err != nil {
		return nil, err
	}
	embReq := &embeddings.CreateRequest{Model: req.Model, Input: req.Input}
	embResp, err := c.or.Embeddings.Create(ctx, embReq)
	if err != nil {
		return nil, err
	}
	return toLLMEmbeddingResponse(embResp), nil
}

type orStreamReader struct {
	inner *chat.StreamReader
}

func (s *orStreamReader) Next() (*StreamChunk, error) {
	chunk, err := s.inner.Next()
	if err == io.EOF {
		if chunk != nil {
			return toLLMStreamChunk(chunk), io.EOF
		}
		return nil, io.EOF
	}
	if err != nil {
		return nil, err
	}
	return toLLMStreamChunk(chunk), nil
}

func (s *orStreamReader) Close() error {
	s.inner.Close()
	return nil
}

func toORChatRequest(req *ChatRequest) *chat.ChatRequest {
	if req == nil {
		return nil
	}
	orReq := &chat.ChatRequest{
		Model:             req.Model,
		Temperature:       req.Temperature,
		TopP:              req.TopP,
		TopK:              req.TopK,
		MaxTokens:         req.MaxTokens,
		Stop:              req.Stop,
		Stream:            req.Stream,
		Seed:              req.Seed,
		PresencePenalty:   req.PresencePenalty,
		FrequencyPenalty:  req.FrequencyPenalty,
		ResponseFormat:    toORResponseFormat(req.ResponseFormat),
		Tools:             toORTools(req.Tools),
		ToolChoice:        req.ToolChoice,
		ParallelToolCalls: req.ParallelToolCalls,
		User:              req.User,
	}
	orReq.Messages = make([]chat.Message, len(req.Messages))
	for i, m := range req.Messages {
		orReq.Messages[i] = toORMessage(m)
	}
	return orReq
}

func toORMessage(m Message) chat.Message {
	or := chat.Message{
		Role:      m.Role,
		Content:   m.Content,
		Name:      m.Name,
		ToolCallID: m.ToolCallID,
	}
	if len(m.ToolCalls) > 0 {
		or.ToolCalls = make([]chat.ToolCall, len(m.ToolCalls))
		for i, tc := range m.ToolCalls {
			or.ToolCalls[i] = chat.ToolCall{
				ID:   tc.ID,
				Type: tc.Type,
				Function: chat.FunctionCall{
					Name:      tc.Function.Name,
					Arguments: tc.Function.Arguments,
				},
			}
		}
	}
	return or
}

func toORResponseFormat(r *ResponseFormat) *chat.ResponseFormat {
	if r == nil {
		return nil
	}
	or := &chat.ResponseFormat{Type: r.Type}
	if r.JSONSchema != nil {
		or.JSONSchema = &chat.JSONSchemaDef{
			Name:   r.JSONSchema.Name,
			Strict: r.JSONSchema.Strict,
			Schema: r.JSONSchema.Schema,
		}
	}
	return or
}

func toORTools(tools []Tool) []chat.Tool {
	if len(tools) == 0 {
		return nil
	}
	or := make([]chat.Tool, len(tools))
	for i, t := range tools {
		or[i] = chat.Tool{
			Type: t.Type,
			Function: chat.FunctionDef{
				Name:        t.Function.Name,
				Description: t.Function.Description,
				Parameters:  t.Function.Parameters,
			},
		}
	}
	return or
}

func toLLMChatResponse(resp *chat.ChatResponse) *ChatResponse {
	if resp == nil {
		return nil
	}
	out := &ChatResponse{
		ID:      resp.ID,
		Object:  resp.Object,
		Created: resp.Created,
		Model:   resp.Model,
		Choices: make([]Choice, len(resp.Choices)),
	}
	if resp.Usage != nil {
		out.Usage = &Usage{
			PromptTokens:     resp.Usage.PromptTokens,
			CompletionTokens: resp.Usage.CompletionTokens,
			TotalTokens:      resp.Usage.TotalTokens,
			Cost:             resp.Usage.Cost,
		}
	}
	for i, ch := range resp.Choices {
		out.Choices[i] = toLLMChoice(ch)
	}
	return out
}

func toLLMChoice(ch chat.Choice) Choice {
	out := Choice{
		Index:        ch.Index,
		FinishReason: ch.FinishReason,
	}
	if ch.Message != nil {
		out.Message = toLLMMessage(ch.Message)
	}
	if ch.Delta != nil {
		out.Delta = toLLMMessage(ch.Delta)
	}
	if ch.Error != nil {
		out.Error = &ChoiceError{
			Code:    ch.Error.Code,
			Message: ch.Error.Message,
		}
	}
	return out
}

func toLLMMessage(m *chat.Message) *Message {
	if m == nil {
		return nil
	}
	out := &Message{
		Role:       m.Role,
		Content:    m.Content,
		Name:       m.Name,
		ToolCallID: m.ToolCallID,
	}
	if len(m.ToolCalls) > 0 {
		out.ToolCalls = make([]ToolCall, len(m.ToolCalls))
		for i, tc := range m.ToolCalls {
			out.ToolCalls[i] = ToolCall{
				Index: tc.Index,
				ID:    tc.ID,
				Type:  tc.Type,
				Function: FunctionCall{
					Name:      tc.Function.Name,
					Arguments: tc.Function.Arguments,
				},
			}
		}
	}
	return out
}

func toLLMStreamChunk(c *chat.StreamChunk) *StreamChunk {
	if c == nil {
		return nil
	}
	out := &StreamChunk{
		ID:      c.ID,
		Object:  c.Object,
		Created: c.Created,
		Model:   c.Model,
		Choices: make([]Choice, len(c.Choices)),
	}
	if c.Usage != nil {
		out.Usage = &Usage{
			PromptTokens:     c.Usage.PromptTokens,
			CompletionTokens: c.Usage.CompletionTokens,
			TotalTokens:      c.Usage.TotalTokens,
			Cost:             c.Usage.Cost,
		}
	}
	for i, ch := range c.Choices {
		out.Choices[i] = toLLMChoice(ch)
	}
	return out
}

func toLLMEmbeddingResponse(resp *embeddings.CreateResponse) *EmbeddingResponse {
	if resp == nil {
		return nil
	}
	out := &EmbeddingResponse{
		Data: make([]EmbeddingData, len(resp.Data)),
	}
	for i, d := range resp.Data {
		out.Data[i] = EmbeddingData{
			Object:    d.Object,
			Embedding: d.Embedding,
			Index:     d.Index,
		}
	}
	if resp.Usage != nil {
		out.Usage = &EmbeddingUsage{
			PromptTokens: resp.Usage.PromptTokens,
			TotalTokens:  resp.Usage.TotalTokens,
		}
	}
	return out
}
