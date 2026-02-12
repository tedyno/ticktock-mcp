package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerUserTools(s *server.MCPServer, r *registry) {
	s.AddTool(
		mcp.NewTool("clockify_user_current",
			mcp.WithDescription("Get the current authenticated user"),
		),
		userCurrentHandler(r),
	)

	s.AddTool(
		mcp.NewTool("clockify_user_list",
			mcp.WithDescription("List all users in a workspace"),
			mcp.WithString("workspace_id", mcp.Description("Workspace ID (uses default if not provided)")),
		),
		userListHandler(r),
	)
}

func userCurrentHandler(r *registry) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		user, err := r.client.GetCurrentUser()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get current user: %v", err)), nil
		}

		return mcp.NewToolResultJSON(user)
	}
}

func userListHandler(r *registry) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		wsID := r.workspaceID(req.GetString("workspace_id", ""))
		if wsID == "" {
			return mcp.NewToolResultError("workspace_id is required"), nil
		}

		users, err := r.client.GetWorkspaceUsers(wsID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list users: %v", err)), nil
		}

		return mcp.NewToolResultJSON(users)
	}
}
