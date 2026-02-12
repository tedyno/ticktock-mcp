package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerWorkspaceTools(s *server.MCPServer, r *registry) {
	s.AddTool(
		mcp.NewTool("clockify_workspace_list",
			mcp.WithDescription("List all workspaces available to the current user"),
		),
		workspaceListHandler(r),
	)
}

func workspaceListHandler(r *registry) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		workspaces, err := r.client.GetWorkspaces()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list workspaces: %v", err)), nil
		}

		return mcp.NewToolResultJSON(workspaces)
	}
}
