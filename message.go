package llm

// TextMessage returns a text-only message for the given role and content.
func TextMessage(role, text string) Message {
	return Message{Role: role, Content: text}
}

// MultimodalMessage returns a message with multimodal content (text + images).
// If parts is nil, an empty slice is used.
func MultimodalMessage(role string, parts []ContentPart) Message {
	if parts == nil {
		parts = []ContentPart{}
	}
	return Message{Role: role, Content: parts}
}
