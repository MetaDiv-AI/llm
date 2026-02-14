package llm

import (
	"time"

	"github.com/MetaDiv-AI/logger"
)

// Provider identifies the LLM backend.
type Provider string

const (
	ProviderOpenRouter Provider = "openrouter"
)

// Client exposes Chat and Embeddings providers.
type Client struct {
	Chat      ChatProvider
	Embeddings EmbeddingProvider
}

// Option is a functional option for configuring a provider.
type Option func(*config)

type config struct {
	APIKey      string
	BaseURL     string
	Timeout     time.Duration
	MaxRetries  int
	Headers     map[string]string
	Debug       bool
	Logger      logger.Logger
	Referer     string
	Title       string
	ForwardedFor string
}

// WithAPIKey sets the API key.
func WithAPIKey(apiKey string) Option {
	return func(c *config) {
		c.APIKey = apiKey
	}
}

// WithBaseURL sets the base URL for API requests.
func WithBaseURL(baseURL string) Option {
	return func(c *config) {
		c.BaseURL = baseURL
	}
}

// WithTimeout sets the HTTP client timeout.
func WithTimeout(timeout time.Duration) Option {
	return func(c *config) {
		c.Timeout = timeout
	}
}

// WithMaxRetries sets the maximum number of retries.
func WithMaxRetries(n int) Option {
	return func(c *config) {
		c.MaxRetries = n
	}
}

// WithHeaders sets custom headers.
func WithHeaders(headers map[string]string) Option {
	return func(c *config) {
		if c.Headers == nil {
			c.Headers = make(map[string]string)
		}
		for k, v := range headers {
			c.Headers[k] = v
		}
	}
}

// WithReferer sets the HTTP-Referer header (OpenRouter app attribution).
func WithReferer(url string) Option {
	return func(c *config) {
		c.Referer = url
	}
}

// WithTitle sets the X-Title header (OpenRouter app title).
func WithTitle(title string) Option {
	return func(c *config) {
		c.Title = title
	}
}

// WithForwardedFor sets the X-Forwarded-For header.
func WithForwardedFor(ip string) Option {
	return func(c *config) {
		c.ForwardedFor = ip
	}
}

// WithDebug enables debug logging.
func WithDebug(debug bool) Option {
	return func(c *config) {
		c.Debug = debug
	}
}

// WithLogger sets a custom logger.
func WithLogger(log logger.Logger) Option {
	return func(c *config) {
		c.Logger = log
	}
}

// DefaultTimeout is the default HTTP client timeout (60s).
const DefaultTimeout = 60 * time.Second

// DefaultMaxRetries is the default number of retries for retryable errors (3).
const DefaultMaxRetries = 3

// NewClient creates a new client for the given provider with the specified options.
func NewClient(provider Provider, opts ...Option) (*Client, error) {
	cfg := &config{
		Headers:    make(map[string]string),
		Timeout:    DefaultTimeout,
		MaxRetries: DefaultMaxRetries,
	}
	for _, opt := range opts {
		opt(cfg)
	}

	switch provider {
	case ProviderOpenRouter:
		return newOpenRouterClient(cfg)
	default:
		return nil, &ErrUnknownProvider{Provider: string(provider)}
	}
}
