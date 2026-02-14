package llm

import (
	"context"
	"io"
)

// MockChatProvider is a ChatProvider for testing.
type MockChatProvider struct {
	CreateFunc       func(context.Context, *ChatRequest) (*ChatResponse, error)
	CreateStreamFunc func(context.Context, *ChatRequest) (StreamReader, error)
}

func (m *MockChatProvider) Create(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, req)
	}
	return &ChatResponse{Choices: []Choice{{Message: &Message{Content: "mock"}}}}, nil
}

func (m *MockChatProvider) CreateStream(ctx context.Context, req *ChatRequest) (StreamReader, error) {
	if m.CreateStreamFunc != nil {
		return m.CreateStreamFunc(ctx, req)
	}
	return &noopStreamReader{}, nil
}

// noopStreamReader is a StreamReader that immediately returns EOF.
type noopStreamReader struct{}

func (*noopStreamReader) Next() (*StreamChunk, error) { return nil, io.EOF }
func (*noopStreamReader) Close() error               { return nil }

// MockEmbeddingProvider is an EmbeddingProvider for testing.
type MockEmbeddingProvider struct {
	CreateFunc func(context.Context, *EmbeddingRequest) (*EmbeddingResponse, error)
}

func (m *MockEmbeddingProvider) Create(ctx context.Context, req *EmbeddingRequest) (*EmbeddingResponse, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, req)
	}
	return &EmbeddingResponse{Data: []EmbeddingData{{Embedding: []float64{0.1}}}}, nil
}
