package llm

import (
	"testing"

	"github.com/MetaDiv-AI/openrouter/chat"
	"github.com/MetaDiv-AI/openrouter/embeddings"
)

func TestToORChatRequest(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		got := toORChatRequest(nil)
		if got != nil {
			t.Errorf("toORChatRequest(nil) = %v, want nil", got)
		}
	})

	t.Run("simple", func(t *testing.T) {
		req := &ChatRequest{
			Model: "test-model",
			Messages: []Message{
				{Role: "user", Content: "hello"},
			},
		}
		got := toORChatRequest(req)
		if got == nil {
			t.Fatal("toORChatRequest returned nil")
		}
		if got.Model != "test-model" {
			t.Errorf("Model = %q, want test-model", got.Model)
		}
		if len(got.Messages) != 1 {
			t.Errorf("len(Messages) = %d, want 1", len(got.Messages))
		}
		if got.Messages[0].Role != "user" || got.Messages[0].Content != "hello" {
			t.Errorf("Messages[0] = %+v", got.Messages[0])
		}
	})

	t.Run("multimodal", func(t *testing.T) {
		req := &ChatRequest{
			Model: "vision-model",
			Messages: []Message{{
				Role: "user",
				Content: []ContentPart{
					{Type: "text", Text: "describe"},
					{Type: "image_url", ImageURL: &ImageURL{URL: "https://example.com/img.png"}},
				},
			}},
		}
		got := toORChatRequest(req)
		if got == nil {
			t.Fatal("toORChatRequest returned nil")
		}
		content, ok := got.Messages[0].Content.([]ContentPart)
		if !ok {
			t.Fatalf("Content type = %T, want []ContentPart", got.Messages[0].Content)
		}
		if len(content) != 2 {
			t.Errorf("len(Content) = %d, want 2", len(content))
		}
		if content[1].ImageURL == nil || content[1].ImageURL.URL != "https://example.com/img.png" {
			t.Errorf("ImageURL = %+v", content[1].ImageURL)
		}
	})
}

func TestToLLMChatResponse(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		got := toLLMChatResponse(nil)
		if got != nil {
			t.Errorf("toLLMChatResponse(nil) = %v, want nil", got)
		}
	})

	t.Run("with usage", func(t *testing.T) {
		resp := &chat.ChatResponse{
			ID: "id1", Model: "m1",
			Choices: []chat.Choice{{
				Index: 0,
				Message: &chat.Message{Role: "assistant", Content: "hi"},
				FinishReason: "stop",
			}},
			Usage: &chat.Usage{PromptTokens: 10, CompletionTokens: 5, TotalTokens: 15},
		}
		got := toLLMChatResponse(resp)
		if got == nil {
			t.Fatal("toLLMChatResponse returned nil")
		}
		if got.Usage == nil || got.Usage.PromptTokens != 10 || got.Usage.TotalTokens != 15 {
			t.Errorf("Usage = %+v", got.Usage)
		}
		if len(got.Choices) != 1 || got.Choices[0].Message.Content != "hi" {
			t.Errorf("Choices = %+v", got.Choices)
		}
	})

	t.Run("with choice error", func(t *testing.T) {
		resp := &chat.ChatResponse{
			Choices: []chat.Choice{{
				Index: 0,
				Error: &chat.ChoiceError{Code: 500, Message: "provider error"},
			}},
		}
		got := toLLMChatResponse(resp)
		if got == nil || len(got.Choices) != 1 {
			t.Fatal("expected one choice")
		}
		if got.Choices[0].Error == nil {
			t.Fatal("expected Choice.Error to be set")
		}
		if got.Choices[0].Error.Code != 500 || got.Choices[0].Error.Message != "provider error" {
			t.Errorf("Choice.Error = %+v", got.Choices[0].Error)
		}
	})
}

func TestToLLMEmbeddingResponse(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		got := toLLMEmbeddingResponse(nil)
		if got != nil {
			t.Errorf("toLLMEmbeddingResponse(nil) = %v, want nil", got)
		}
	})

	t.Run("with data", func(t *testing.T) {
		resp := &embeddings.CreateResponse{
			Data: []embeddings.EmbeddingData{
				{Object: "embedding", Embedding: []float64{0.1, 0.2}, Index: 0},
			},
			Usage: &embeddings.Usage{PromptTokens: 3, TotalTokens: 3},
		}
		got := toLLMEmbeddingResponse(resp)
		if got == nil {
			t.Fatal("toLLMEmbeddingResponse returned nil")
		}
		if len(got.Data) != 1 || len(got.Data[0].Embedding) != 2 {
			t.Errorf("Data = %+v", got.Data)
		}
		if got.Usage == nil || got.Usage.TotalTokens != 3 {
			t.Errorf("Usage = %+v", got.Usage)
		}
	})
}
