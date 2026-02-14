// Package llm provides a provider-agnostic Go client for LLM chat (with multimodal support) and embeddings.
//
// It defines interfaces (ChatProvider, EmbeddingProvider) that abstract over different backends,
// with OpenRouter as the first supported provider. Use NewClient to create a client for a given provider.
package llm
