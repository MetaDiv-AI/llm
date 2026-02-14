package llm

import (
	"context"
	"errors"
	"io"
	"testing"
)

func TestValidateChatRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     *ChatRequest
		wantErr bool
	}{
		{"nil", nil, true},
		{"empty model", &ChatRequest{Model: "", Messages: []Message{{Role: "user", Content: "hi"}}}, true},
		{"empty messages", &ChatRequest{Model: "m", Messages: nil}, true},
		{"empty messages slice", &ChatRequest{Model: "m", Messages: []Message{}}, true},
		{"valid", &ChatRequest{Model: "m", Messages: []Message{{Role: "user", Content: "hi"}}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateChatRequest(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateChatRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && !errors.Is(err, ErrInvalidRequest) {
				t.Errorf("expected errors.Is(err, ErrInvalidRequest)")
			}
		})
	}
}

func TestValidateEmbeddingRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     *EmbeddingRequest
		wantErr bool
	}{
		{"nil", nil, true},
		{"empty model", &EmbeddingRequest{Model: "", Input: "x"}, true},
		{"nil input", &EmbeddingRequest{Model: "m", Input: nil}, true},
		{"empty string", &EmbeddingRequest{Model: "m", Input: ""}, true},
		{"empty []string", &EmbeddingRequest{Model: "m", Input: []string{}}, true},
		{"empty []interface{}", &EmbeddingRequest{Model: "m", Input: []interface{}{}}, true},
		{"valid string", &EmbeddingRequest{Model: "m", Input: "x"}, false},
		{"valid []string", &EmbeddingRequest{Model: "m", Input: []string{"x"}}, false},
		{"valid []interface{}", &EmbeddingRequest{Model: "m", Input: []interface{}{"x"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateEmbeddingRequest(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateEmbeddingRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && !errors.Is(err, ErrInvalidRequest) {
				t.Errorf("expected errors.Is(err, ErrInvalidRequest)")
			}
		})
	}
}

func TestMultimodalMessage_NilParts(t *testing.T) {
	m := MultimodalMessage("user", nil)
	if m.Content == nil {
		t.Error("Content should not be nil")
	}
	parts, ok := m.Content.([]ContentPart)
	if !ok || len(parts) != 0 {
		t.Errorf("Content = %v, want []ContentPart{}", m.Content)
	}
}

func TestMockChatProvider_CreateStreamNoOp(t *testing.T) {
	mock := &MockChatProvider{}
	stream, err := mock.CreateStream(context.Background(), &ChatRequest{Model: "m", Messages: []Message{{Role: "user", Content: "hi"}}})
	if err != nil {
		t.Fatalf("CreateStream: %v", err)
	}
	if stream == nil {
		t.Fatal("stream should not be nil")
	}
	defer stream.Close()
	chunk, err := stream.Next()
	if chunk != nil || err == nil {
		t.Errorf("Next() = %v, %v; want nil, io.EOF", chunk, err)
	}
	if !errors.Is(err, io.EOF) {
		t.Errorf("Next() err = %v, want io.EOF", err)
	}
}
