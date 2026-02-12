package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/tedyno/ticktock-mcp/clockify"
)

func registerProjectTools(s *server.MCPServer, r *registry) {
	s.AddTool(
		mcp.NewTool("clockify_project_list",
			mcp.WithDescription("List all projects in a workspace"),
			mcp.WithBoolean("archived", mcp.Description("Include archived projects")),
			mcp.WithString("workspace_id", mcp.Description("Workspace ID (uses default if not provided)")),
		),
		projectListHandler(r),
	)

	s.AddTool(
		mcp.NewTool("clockify_project_create",
			mcp.WithDescription("Create a new project"),
			mcp.WithString("name", mcp.Required(), mcp.Description("Project name")),
			mcp.WithString("client_id", mcp.Description("Client ID")),
			mcp.WithBoolean("billable", mcp.Description("Whether the project is billable")),
			mcp.WithString("color", mcp.Description("Project color (hex, e.g. #FF0000)")),
			mcp.WithBoolean("is_public", mcp.Description("Whether the project is public")),
			mcp.WithString("workspace_id", mcp.Description("Workspace ID (uses default if not provided)")),
		),
		projectCreateHandler(r),
	)

	s.AddTool(
		mcp.NewTool("clockify_project_update",
			mcp.WithDescription("Update an existing project"),
			mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID to update")),
			mcp.WithString("name", mcp.Description("New project name")),
			mcp.WithString("client_id", mcp.Description("Client ID")),
			mcp.WithBoolean("billable", mcp.Description("Whether the project is billable")),
			mcp.WithString("color", mcp.Description("Project color (hex)")),
			mcp.WithBoolean("archived", mcp.Description("Whether the project is archived")),
			mcp.WithString("workspace_id", mcp.Description("Workspace ID (uses default if not provided)")),
		),
		projectUpdateHandler(r),
	)

	s.AddTool(
		mcp.NewTool("clockify_project_delete",
			mcp.WithDescription("Delete a project"),
			mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID to delete")),
			mcp.WithString("workspace_id", mcp.Description("Workspace ID (uses default if not provided)")),
		),
		projectDeleteHandler(r),
	)
}

func projectListHandler(r *registry) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		wsID := r.workspaceID(req.GetString("workspace_id", ""))
		if wsID == "" {
			return mcp.NewToolResultError("workspace_id is required"), nil
		}

		projects, err := r.client.GetProjects(wsID, req.GetBool("archived", false))
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list projects: %v", err)), nil
		}

		return mcp.NewToolResultJSON(projects)
	}
}

func projectCreateHandler(r *registry) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		wsID := r.workspaceID(req.GetString("workspace_id", ""))
		if wsID == "" {
			return mcp.NewToolResultError("workspace_id is required"), nil
		}

		name, err := req.RequireString("name")
		if err != nil {
			return mcp.NewToolResultError("name is required"), nil
		}

		project, err := r.client.CreateProject(wsID, clockify.CreateProjectRequest{
			Name:     name,
			ClientID: req.GetString("client_id", ""),
			Billable: req.GetBool("billable", false),
			Color:    req.GetString("color", ""),
			IsPublic: req.GetBool("is_public", true),
		})
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create project: %v", err)), nil
		}

		return mcp.NewToolResultJSON(project)
	}
}

func projectUpdateHandler(r *registry) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		wsID := r.workspaceID(req.GetString("workspace_id", ""))
		if wsID == "" {
			return mcp.NewToolResultError("workspace_id is required"), nil
		}

		projectID, err := req.RequireString("project_id")
		if err != nil {
			return mcp.NewToolResultError("project_id is required"), nil
		}

		updateReq := clockify.UpdateProjectRequest{
			Name:     req.GetString("name", ""),
			ClientID: req.GetString("client_id", ""),
			Color:    req.GetString("color", ""),
		}

		args := req.GetArguments()
		if _, ok := args["billable"]; ok {
			b := req.GetBool("billable", false)
			updateReq.Billable = &b
		}
		if _, ok := args["archived"]; ok {
			a := req.GetBool("archived", false)
			updateReq.Archived = &a
		}

		project, err := r.client.UpdateProject(wsID, projectID, updateReq)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to update project: %v", err)), nil
		}

		return mcp.NewToolResultJSON(project)
	}
}

func projectDeleteHandler(r *registry) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		wsID := r.workspaceID(req.GetString("workspace_id", ""))
		if wsID == "" {
			return mcp.NewToolResultError("workspace_id is required"), nil
		}

		projectID, err := req.RequireString("project_id")
		if err != nil {
			return mcp.NewToolResultError("project_id is required"), nil
		}

		if err := r.client.DeleteProject(wsID, projectID); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to delete project: %v", err)), nil
		}

		return mcp.NewToolResultText("Project deleted successfully."), nil
	}
}
