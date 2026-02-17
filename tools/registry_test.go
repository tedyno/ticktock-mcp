package tools

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestResultJSON_ReturnsTextContentWithJSON(t *testing.T) {
	data := map[string]any{"items": []string{"a", "b"}}
	result, err := resultJSON(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Must have exactly one content item
	if len(result.Content) != 1 {
		t.Fatalf("expected 1 content item, got %d", len(result.Content))
	}

	// Content must be TextContent with valid JSON
	tc, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		t.Fatalf("expected TextContent, got %T", result.Content[0])
	}

	var parsed map[string]any
	if err := json.Unmarshal([]byte(tc.Text), &parsed); err != nil {
		t.Fatalf("content is not valid JSON: %v", err)
	}

	items, ok := parsed["items"].([]any)
	if !ok || len(items) != 2 {
		t.Fatalf("unexpected parsed content: %v", parsed)
	}
}

func TestResultJSON_DoesNotSetStructuredContent(t *testing.T) {
	data := map[string]any{"entries": []string{"x"}}
	result, err := resultJSON(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.StructuredContent != nil {
		t.Fatalf("expected StructuredContent to be nil, got %v", result.StructuredContent)
	}
}

func TestResultJSON_HandlesStruct(t *testing.T) {
	type item struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	data := item{ID: "1", Name: "test"}
	result, err := resultJSON(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tc := result.Content[0].(mcp.TextContent)
	var parsed map[string]any
	if err := json.Unmarshal([]byte(tc.Text), &parsed); err != nil {
		t.Fatalf("content is not valid JSON: %v", err)
	}
	if parsed["id"] != "1" || parsed["name"] != "test" {
		t.Fatalf("unexpected parsed content: %v", parsed)
	}
}

func TestNoToolUsesNewToolResultJSON(t *testing.T) {
	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatalf("failed to glob: %v", err)
	}
	for _, f := range files {
		if f == "registry_test.go" {
			continue
		}
		data, err := os.ReadFile(f)
		if err != nil {
			t.Fatalf("failed to read %s: %v", f, err)
		}
		if strings.Contains(string(data), "NewToolResultJSON") {
			t.Errorf("%s still uses mcp.NewToolResultJSON — use resultJSON() instead", f)
		}
	}
}
