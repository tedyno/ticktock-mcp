package tools

import (
	"encoding/json"
	"testing"
)

func TestResultJSON_WireFormat(t *testing.T) {
	tests := []struct {
		name string
		data any
	}{
		{"map with slice", map[string]any{"entries": []map[string]any{{"id": "1"}, {"id": "2"}}}},
		{"map with empty slice", map[string]any{"entries": []any{}}},
		{"struct pointer", &struct {
			Name string `json:"name"`
		}{Name: "test"}},
		{"nested map", map[string]any{"totals": map[string]any{"total": 123}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := resultJSON(tt.data)
			if err != nil {
				t.Fatalf("resultJSON failed: %v", err)
			}

			// Serialize the CallToolResult to JSON (as the MCP server would)
			wireBytes, err := json.Marshal(result)
			if err != nil {
				t.Fatalf("failed to marshal result: %v", err)
			}

			// Parse the wire format and verify no structuredContent
			var wire map[string]any
			if err := json.Unmarshal(wireBytes, &wire); err != nil {
				t.Fatalf("failed to unmarshal wire: %v", err)
			}

			if _, exists := wire["structuredContent"]; exists {
				t.Errorf("wire format contains structuredContent — this will cause Claude Code validation errors.\nWire: %s", string(wireBytes))
			}

			// Verify content field has text with valid JSON
			content, ok := wire["content"].([]any)
			if !ok || len(content) == 0 {
				t.Fatalf("expected content array, got: %v", wire["content"])
			}

			firstContent, ok := content[0].(map[string]any)
			if !ok {
				t.Fatalf("expected content[0] to be object, got: %T", content[0])
			}

			text, ok := firstContent["text"].(string)
			if !ok || text == "" {
				t.Fatalf("expected non-empty text in content[0], got: %v", firstContent)
			}

			// Verify the text is valid JSON
			var parsed any
			if err := json.Unmarshal([]byte(text), &parsed); err != nil {
				t.Errorf("content text is not valid JSON: %v\nText: %s", err, text)
			}
		})
	}
}

func TestResultJSON_ErrorOnMarshalFailure(t *testing.T) {
	// Channels can't be marshaled to JSON
	data := map[string]any{"bad": make(chan int)}
	_, err := resultJSON(data)
	if err == nil {
		t.Fatal("expected error for unmarshalable data, got nil")
	}
}
