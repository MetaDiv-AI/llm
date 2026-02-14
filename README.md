# llm

Provider-agnostic Go client for LLM chat (with multimodal support) and embeddings.

## Installation

```bash
go get github.com/MetaDiv-AI/llm
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/MetaDiv-AI/llm"
	"github.com/MetaDiv-AI/logger"
)

func main() {
	log := logger.New().Development().Build()
	defer log.Sync()

	// OPENROUTER_API_KEY env var is used when WithAPIKey is omitted
	client, err := llm.NewClient(llm.ProviderOpenRouter,
		llm.WithAPIKey(os.Getenv("OPENROUTER_API_KEY")),
		llm.WithLogger(log),
	)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	// Chat completion
	resp, err := client.Chat.Create(ctx, &llm.ChatRequest{
		Model: "anthropic/claude-sonnet-4",
		Messages: []llm.Message{{Role: "user", Content: "Hello"}},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Choices[0].Message.Content)
}
```

## Multimodal Chat

Chat supports vision (images) via `ContentPart`:

```go
resp, err := client.Chat.Create(ctx, &llm.ChatRequest{
	Model: "anthropic/claude-sonnet-4",
	Messages: []llm.Message{{
		Role: "user",
		Content: []llm.ContentPart{
			{Type: "text", Text: "What's in this image?"},
			{Type: "image_url", ImageURL: &llm.ImageURL{URL: "https://example.com/image.jpg"}},
		},
	}},
})
```

## Streaming

```go
import (
	"errors"
	"io"
)

stream, err := client.Chat.CreateStream(ctx, &llm.ChatRequest{
	Model:    "anthropic/claude-sonnet-4",
	Messages: []llm.Message{{Role: "user", Content: "Hello"}},
})
if err != nil {
	panic(err)
}
defer stream.Close()

for {
	chunk, err := stream.Next()
	if err != nil {
		if errors.Is(err, io.EOF) {
			break
		}
		panic(err)
	}
	if chunk != nil && len(chunk.Choices) > 0 && chunk.Choices[0].Delta != nil && chunk.Choices[0].Delta.Content != nil {
		fmt.Print(chunk.Choices[0].Delta.Content)
	}
}
```

## Embeddings

```go
emb, err := client.Embeddings.Create(ctx, &llm.EmbeddingRequest{
	Model: "openai/text-embedding-3-small",
	Input: "The quick brown fox",
})
if err != nil {
	panic(err)
}
fmt.Println(emb.Data[0].Embedding)
```

## Supported Providers

- **OpenRouter** (`llm.ProviderOpenRouter`) - Access to multiple models via OpenRouter API. When using OpenRouter, the API key can be set via `OPENROUTER_API_KEY` env var if `WithAPIKey` is omitted.

## Error Handling

- `ErrUnknownProvider` is returned when the provider is not supported. Use `errors.Is(err, &llm.ErrUnknownProvider{Provider: "openrouter"})` or `errors.As` to check.
- `ErrInvalidRequest` and `ValidationError` are returned when a request fails validation (e.g. empty model, empty messages). Use `errors.Is(err, llm.ErrInvalidRequest)` to detect validation errors.
- For streaming, `StreamReader.Next()` returns `io.EOF` when done. Use `errors.Is(err, io.EOF)` for EOF detection.

## License

MIT
