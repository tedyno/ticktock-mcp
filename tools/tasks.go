package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/tedyno/ticktock-mcp/clockify"
)

func registerTaskTools(s *server.MCPServer, r *registry) {
	s.AddTool(
		mcp.NewTool("clockify_task_list",
			mcp.WithDescription("List tasks for a project (paginated, default page 1, page_size 50)"),
			mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID")),
			mcp.WithNumber("page", mcp.Description("Page number (default 1)")),
			mcp.WithNumber("page_size", mcp.Description("Number of tasks per page (default 50)")),
			mcp.WithString("workspace_id", mcp.Description("Workspace ID (uses default if not provided)")),
		),
		taskListHandler(r),
	)

	s.AddTool(
		mcp.NewTool("clockify_task_create",
			mcp.WithDescription("Create a new task in a project"),
			mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID")),
			mcp.WithString("name", mcp.Required(), mcp.Description("Task name")),
			mcp.WithBoolean("billable", mcp.Description("Whether the task is billable")),
			mcp.WithString("workspace_id", mcp.Description("Workspace ID (uses default if not provided)")),
		),
		taskCreateHandler(r),
	)

	s.AddTool(
		mcp.NewTool("clockify_task_update",
			mcp.WithDescription("Update a task"),
			mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID")),
			mcp.WithString("task_id", mcp.Required(), mcp.Description("Task ID to update")),
			mcp.WithString("name", mcp.Description("New task name")),
			mcp.WithBoolean("billable", mcp.Description("Whether the task is billable")),
			mcp.WithString("status", mcp.Description("Task status (ACTIVE or DONE)")),
			mcp.WithString("workspace_id", mcp.Description("Workspace ID (uses default if not provided)")),
		),
		taskUpdateHandler(r),
	)

	s.AddTool(
		mcp.NewTool("clockify_task_delete",
			mcp.WithDescription("Delete a task"),
			mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID")),
			mcp.WithString("task_id", mcp.Required(), mcp.Description("Task ID to delete")),
			mcp.WithString("workspace_id", mcp.Description("Workspace ID (uses default if not provided)")),
		),
		taskDeleteHandler(r),
	)
}

func taskListHandler(r *registry) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		wsID := r.workspaceID(req.GetString("workspace_id", ""))
		if wsID == "" {
			return mcp.NewToolResultError("workspace_id is required"), nil
		}

		projectID, err := req.RequireString("project_id")
		if err != nil {
			return mcp.NewToolResultError("project_id is required"), nil
		}

		page := req.GetInt("page", 1)
		pageSize := req.GetInt("page_size", 50)

		tasks, err := r.client.GetTasks(wsID, projectID, page, pageSize)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list tasks: %v", err)), nil
		}

		return mcp.NewToolResultJSON(tasks)
	}
}

func taskCreateHandler(r *registry) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		wsID := r.workspaceID(req.GetString("workspace_id", ""))
		if wsID == "" {
			return mcp.NewToolResultError("workspace_id is required"), nil
		}

		projectID, err := req.RequireString("project_id")
		if err != nil {
			return mcp.NewToolResultError("project_id is required"), nil
		}
		name, err := req.RequireString("name")
		if err != nil {
			return mcp.NewToolResultError("name is required"), nil
		}

		task, err := r.client.CreateTask(wsID, projectID, clockify.CreateTaskRequest{
			Name:     name,
			Billable: req.GetBool("billable", false),
		})
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create task: %v", err)), nil
		}

		return mcp.NewToolResultJSON(task)
	}
}

func taskUpdateHandler(r *registry) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		wsID := r.workspaceID(req.GetString("workspace_id", ""))
		if wsID == "" {
			return mcp.NewToolResultError("workspace_id is required"), nil
		}

		projectID, err := req.RequireString("project_id")
		if err != nil {
			return mcp.NewToolResultError("project_id is required"), nil
		}
		taskID, err := req.RequireString("task_id")
		if err != nil {
			return mcp.NewToolResultError("task_id is required"), nil
		}

		updateReq := clockify.UpdateTaskRequest{
			Name:   req.GetString("name", ""),
			Status: req.GetString("status", ""),
		}

		args := req.GetArguments()
		if _, ok := args["billable"]; ok {
			b := req.GetBool("billable", false)
			updateReq.Billable = &b
		}

		task, err := r.client.UpdateTask(wsID, projectID, taskID, updateReq)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to update task: %v", err)), nil
		}

		return mcp.NewToolResultJSON(task)
	}
}

func taskDeleteHandler(r *registry) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		wsID := r.workspaceID(req.GetString("workspace_id", ""))
		if wsID == "" {
			return mcp.NewToolResultError("workspace_id is required"), nil
		}

		projectID, err := req.RequireString("project_id")
		if err != nil {
			return mcp.NewToolResultError("project_id is required"), nil
		}
		taskID, err := req.RequireString("task_id")
		if err != nil {
			return mcp.NewToolResultError("task_id is required"), nil
		}

		if err := r.client.DeleteTask(wsID, projectID, taskID); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to delete task: %v", err)), nil
		}

		return mcp.NewToolResultText("Task deleted successfully."), nil
	}
}
