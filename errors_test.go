package llm

import (
	"errors"
	"testing"
)

func TestErrUnknownProvider_Is(t *testing.T) {
	err := &ErrUnknownProvider{Provider: "unknown"}
	if !errors.Is(err, &ErrUnknownProvider{Provider: "unknown"}) {
		t.Error("errors.Is should match same provider")
	}
	if !errors.Is(err, &ErrUnknownProvider{Provider: ""}) {
		t.Error("errors.Is should match empty provider target")
	}
	if errors.Is(err, &ErrUnknownProvider{Provider: "other"}) {
		t.Error("errors.Is should not match different provider")
	}
}
