# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.0] - 2025-02-14

### Added

- **ToolCall.Index** field for accumulating streaming deltas (OpenAI format)

### Changed

- OpenRouter provider now maps `Index` when converting tool calls to LLM messages

## [1.0.0] - 2025-02-14

### Added

- Provider-agnostic Go client for LLM chat and embeddings
- **ChatProvider** interface with `Create` (sync) and `CreateStream` (streaming)
- **EmbeddingProvider** interface with `Create`
- **StreamReader** interface for streaming chat completions
- **OpenRouter** as the first supported provider (`ProviderOpenRouter`)
- Multimodal chat support (text + images) via `ContentPart` and `ImageURL`
- Tool calling support (`Tool`, `ToolCall`, `FunctionDef`, `FunctionCall`)
- Structured output via `ResponseFormat` and `JSONSchemaDef`
- Functional options for client configuration:
  - `WithAPIKey`, `WithBaseURL`, `WithTimeout`, `WithMaxRetries`
  - `WithHeaders`, `WithReferer`, `WithTitle`, `WithForwardedFor`
  - `WithDebug`, `WithLogger`
- Default timeout (60s) and max retries (3)
- Request validation for `ChatRequest` and `EmbeddingRequest`
- Typed errors: `ErrUnknownProvider`, `ErrInvalidRequest`, `ValidationError`
- `errors.Is` support for `ErrUnknownProvider` and `ValidationError`
- Helper constructors: `TextMessage`, `MultimodalMessage`
- Mock providers for testing: `MockChatProvider`, `MockEmbeddingProvider`
- `ChoiceError` propagation for provider-level errors in streaming responses

### Dependencies

- `github.com/MetaDiv-AI/logger` v1.0.0
- `github.com/MetaDiv-AI/openrouter` v1.0.0
