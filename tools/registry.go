package tools

import (
	"github.com/mark3labs/mcp-go/server"
	"github.com/tedyno/ticktock-mcp/clockify"
)

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
