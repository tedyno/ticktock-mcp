package tools

import (
	"context"
	"fmt"
	"net/url"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/tedyno/ticktock-mcp/clockify"
)

func registerTimeEntryTools(s *server.MCPServer, r *registry) {
	s.AddTool(
		mcp.NewTool("clockify_time_entry_list",
			mcp.WithDescription("List time entries for the current user (paginated, default page 1, page_size 50)"),
			mcp.WithString("start", mcp.Description("Start date filter (ISO 8601, e.g. 2024-01-01T00:00:00Z)")),
			mcp.WithString("end", mcp.Description("End date filter (ISO 8601)")),
			mcp.WithString("project_id", mcp.Description("Filter by project ID")),
			mcp.WithNumber("page", mcp.Description("Page number (default 1)")),
			mcp.WithNumber("page_size", mcp.Description("Number of entries per page (default 50)")),
			mcp.WithString("workspace_id", mcp.Description("Workspace ID (uses default if not provided)")),
		),
		timeEntryListHandler(r),
	)

	s.AddTool(
		mcp.NewTool("clockify_time_entry_create",
			mcp.WithDescription("Create a manual time entry"),
			mcp.WithString("start", mcp.Required(), mcp.Description("Start time (ISO 8601)")),
			mcp.WithString("end", mcp.Required(), mcp.Description("End time (ISO 8601)")),
			mcp.WithString("description", mcp.Description("Entry description")),
			mcp.WithString("project_id", mcp.Description("Project ID")),
			mcp.WithString("task_id", mcp.Description("Task ID")),
			mcp.WithArray("tag_ids", mcp.Description("Tag IDs"), mcp.WithStringItems()),
			mcp.WithBoolean("billable", mcp.Description("Whether the entry is billable")),
			mcp.WithString("workspace_id", mcp.Description("Workspace ID (uses default if not provided)")),
		),
		timeEntryCreateHandler(r),
	)

	s.AddTool(
		mcp.NewTool("clockify_time_entry_update",
			mcp.WithDescription("Update an existing time entry"),
			mcp.WithString("entry_id", mcp.Required(), mcp.Description("Time entry ID to update")),
			mcp.WithString("start", mcp.Required(), mcp.Description("Start time (ISO 8601)")),
			mcp.WithString("end", mcp.Description("End time (ISO 8601)")),
			mcp.WithString("description", mcp.Description("Entry description")),
			mcp.WithString("project_id", mcp.Description("Project ID")),
			mcp.WithString("task_id", mcp.Description("Task ID")),
			mcp.WithArray("tag_ids", mcp.Description("Tag IDs"), mcp.WithStringItems()),
			mcp.WithBoolean("billable", mcp.Description("Whether the entry is billable")),
			mcp.WithString("workspace_id", mcp.Description("Workspace ID (uses default if not provided)")),
		),
		timeEntryUpdateHandler(r),
	)

	s.AddTool(
		mcp.NewTool("clockify_time_entry_delete",
			mcp.WithDescription("Delete a time entry"),
			mcp.WithString("entry_id", mcp.Required(), mcp.Description("Time entry ID to delete")),
			mcp.WithString("workspace_id", mcp.Description("Workspace ID (uses default if not provided)")),
		),
		timeEntryDeleteHandler(r),
	)
}

func timeEntryListHandler(r *registry) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		wsID := r.workspaceID(req.GetString("workspace_id", ""))
		if wsID == "" {
			return mcp.NewToolResultError("workspace_id is required"), nil
		}

		user, err := r.client.GetCurrentUser()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get current user: %v", err)), nil
		}

		params := url.Values{}
		if start := req.GetString("start", ""); start != "" {
			params.Set("start", start)
		}
		if end := req.GetString("end", ""); end != "" {
			params.Set("end", end)
		}
		if projectID := req.GetString("project_id", ""); projectID != "" {
			params.Set("project", projectID)
		}

		page := req.GetInt("page", 1)
		pageSize := req.GetInt("page_size", 50)

		entries, err := r.client.GetTimeEntries(wsID, user.ID, params, page, pageSize)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list time entries: %v", err)), nil
		}

		return resultJSON(map[string]any{"entries": entries})
	}
}

func timeEntryCreateHandler(r *registry) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		wsID := r.workspaceID(req.GetString("workspace_id", ""))
		if wsID == "" {
			return mcp.NewToolResultError("workspace_id is required"), nil
		}

		start, err := req.RequireString("start")
		if err != nil {
			return mcp.NewToolResultError("start is required"), nil
		}
		end, err := req.RequireString("end")
		if err != nil {
			return mcp.NewToolResultError("end is required"), nil
		}

		entry, err := r.client.CreateTimeEntry(wsID, clockify.CreateTimeEntryRequest{
			Start:       start,
			End:         end,
			Description: req.GetString("description", ""),
			ProjectID:   req.GetString("project_id", ""),
			TaskID:      req.GetString("task_id", ""),
			TagIDs:      req.GetStringSlice("tag_ids", nil),
			Billable:    req.GetBool("billable", false),
		})
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create time entry: %v", err)), nil
		}

		return resultJSON(entry)
	}
}

func timeEntryUpdateHandler(r *registry) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		wsID := r.workspaceID(req.GetString("workspace_id", ""))
		if wsID == "" {
			return mcp.NewToolResultError("workspace_id is required"), nil
		}

		entryID, err := req.RequireString("entry_id")
		if err != nil {
			return mcp.NewToolResultError("entry_id is required"), nil
		}
		start, err := req.RequireString("start")
		if err != nil {
			return mcp.NewToolResultError("start is required"), nil
		}

		entry, err := r.client.UpdateTimeEntry(wsID, entryID, clockify.UpdateTimeEntryRequest{
			Start:       start,
			End:         req.GetString("end", ""),
			Description: req.GetString("description", ""),
			ProjectID:   req.GetString("project_id", ""),
			TaskID:      req.GetString("task_id", ""),
			TagIDs:      req.GetStringSlice("tag_ids", nil),
			Billable:    req.GetBool("billable", false),
		})
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to update time entry: %v", err)), nil
		}

		return resultJSON(entry)
	}
}

func timeEntryDeleteHandler(r *registry) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		wsID := r.workspaceID(req.GetString("workspace_id", ""))
		if wsID == "" {
			return mcp.NewToolResultError("workspace_id is required"), nil
		}

		entryID, err := req.RequireString("entry_id")
		if err != nil {
			return mcp.NewToolResultError("entry_id is required"), nil
		}

		if err := r.client.DeleteTimeEntry(wsID, entryID); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to delete time entry: %v", err)), nil
		}

		return mcp.NewToolResultText("Time entry deleted successfully."), nil
	}
}
