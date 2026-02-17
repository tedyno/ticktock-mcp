package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/tedyno/ticktock-mcp/clockify"
)

func registerClientTools(s *server.MCPServer, r *registry) {
	s.AddTool(
		mcp.NewTool("clockify_client_list",
			mcp.WithDescription("List clients in a workspace (paginated, default page 1, page_size 50)"),
			mcp.WithNumber("page", mcp.Description("Page number (default 1)")),
			mcp.WithNumber("page_size", mcp.Description("Number of clients per page (default 50)")),
			mcp.WithString("workspace_id", mcp.Description("Workspace ID (uses default if not provided)")),
		),
		clientListHandler(r),
	)

	s.AddTool(
		mcp.NewTool("clockify_client_create",
			mcp.WithDescription("Create a new client"),
			mcp.WithString("name", mcp.Required(), mcp.Description("Client name")),
			mcp.WithString("workspace_id", mcp.Description("Workspace ID (uses default if not provided)")),
		),
		clientCreateHandler(r),
	)

	s.AddTool(
		mcp.NewTool("clockify_client_update",
			mcp.WithDescription("Update a client"),
			mcp.WithString("client_id", mcp.Required(), mcp.Description("Client ID to update")),
			mcp.WithString("name", mcp.Description("New client name")),
			mcp.WithBoolean("archived", mcp.Description("Whether the client is archived")),
			mcp.WithString("workspace_id", mcp.Description("Workspace ID (uses default if not provided)")),
		),
		clientUpdateHandler(r),
	)

	s.AddTool(
		mcp.NewTool("clockify_client_delete",
			mcp.WithDescription("Delete a client"),
			mcp.WithString("client_id", mcp.Required(), mcp.Description("Client ID to delete")),
			mcp.WithString("workspace_id", mcp.Description("Workspace ID (uses default if not provided)")),
		),
		clientDeleteHandler(r),
	)
}

func clientListHandler(r *registry) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		wsID := r.workspaceID(req.GetString("workspace_id", ""))
		if wsID == "" {
			return mcp.NewToolResultError("workspace_id is required"), nil
		}

		page := req.GetInt("page", 1)
		pageSize := req.GetInt("page_size", 50)

		clients, err := r.client.GetClients(wsID, page, pageSize)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list clients: %v", err)), nil
		}

		return resultJSON(map[string]any{"clients": clients})
	}
}

func clientCreateHandler(r *registry) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		wsID := r.workspaceID(req.GetString("workspace_id", ""))
		if wsID == "" {
			return mcp.NewToolResultError("workspace_id is required"), nil
		}

		name, err := req.RequireString("name")
		if err != nil {
			return mcp.NewToolResultError("name is required"), nil
		}

		client, err := r.client.CreateClient(wsID, clockify.CreateClientRequest{Name: name})
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create client: %v", err)), nil
		}

		return resultJSON(client)
	}
}

func clientUpdateHandler(r *registry) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		wsID := r.workspaceID(req.GetString("workspace_id", ""))
		if wsID == "" {
			return mcp.NewToolResultError("workspace_id is required"), nil
		}

		clientID, err := req.RequireString("client_id")
		if err != nil {
			return mcp.NewToolResultError("client_id is required"), nil
		}

		updateReq := clockify.UpdateClientRequest{
			Name: req.GetString("name", ""),
		}

		args := req.GetArguments()
		if _, ok := args["archived"]; ok {
			a := req.GetBool("archived", false)
			updateReq.Archived = &a
		}

		client, err := r.client.UpdateClient(wsID, clientID, updateReq)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to update client: %v", err)), nil
		}

		return resultJSON(client)
	}
}

func clientDeleteHandler(r *registry) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		wsID := r.workspaceID(req.GetString("workspace_id", ""))
		if wsID == "" {
			return mcp.NewToolResultError("workspace_id is required"), nil
		}

		clientID, err := req.RequireString("client_id")
		if err != nil {
			return mcp.NewToolResultError("client_id is required"), nil
		}

		if err := r.client.DeleteClient(wsID, clientID); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to delete client: %v", err)), nil
		}

		return mcp.NewToolResultText("Client deleted successfully."), nil
	}
}
