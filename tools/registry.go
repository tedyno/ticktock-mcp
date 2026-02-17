package tools

import (
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/tedyno/ticktock-mcp/clockify"
)

// resultJSON marshals data to JSON and returns it as a text-only tool result.
// Unlike mcp.NewToolResultJSON, this does NOT set StructuredContent,
// avoiding Claude Code's Zod validation error on structuredContent field.
func resultJSON(data any) (*mcp.CallToolResult, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("marshal JSON result: %w", err)
	}
	return mcp.NewToolResultText(string(b)), nil
}

// RegisterAll registers all Clockify MCP tools on the given server.
func RegisterAll(s *server.MCPServer, client *clockify.Client, defaultWorkspaceID string) {
	r := &registry{client: client, defaultWorkspaceID: defaultWorkspaceID}

	registerTimerTools(s, r)
	registerTimeEntryTools(s, r)
	registerProjectTools(s, r)
	registerTaskTools(s, r)
	registerTagTools(s, r)
	registerClientTools(s, r)
	registerWorkspaceTools(s, r)
	registerUserTools(s, r)
	registerReportTools(s, r)
}

type registry struct {
	client             *clockify.Client
	defaultWorkspaceID string
}

// workspaceID returns the provided workspace ID or falls back to default.
func (r *registry) workspaceID(override string) string {
	if override != "" {
		return override
	}
	return r.defaultWorkspaceID
}
