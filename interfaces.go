package llm

import "context"

// ChatProvider provides chat completion operations.
type ChatProvider interface {
	Create(ctx context.Context, req *ChatRequest) (*ChatResponse, error)
	CreateStream(ctx context.Context, req *ChatRequest) (StreamReader, error)
}

// EmbeddingProvider provides embedding operations.
type EmbeddingProvider interface {
	Create(ctx context.Context, req *EmbeddingRequest) (*EmbeddingResponse, error)
}

// StreamReader reads a streaming chat completion.
// Next returns (*StreamChunk, nil) for data, (nil, io.EOF) when done, or (nil, err) on error.
// Callers should use errors.Is(err, io.EOF) for EOF detection.
type StreamReader interface {
	Next() (*StreamChunk, error)
	Close() error
}
