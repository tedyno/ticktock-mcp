package tools

import (
	"context"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/tedyno/ticktock-mcp/clockify"
)

func registerTimerTools(s *server.MCPServer, r *registry) {
	s.AddTool(
		mcp.NewTool("clockify_timer_start",
			mcp.WithDescription("Start a new timer in Clockify"),
			mcp.WithString("description", mcp.Description("Timer description")),
			mcp.WithString("project_id", mcp.Description("Project ID")),
			mcp.WithString("task_id", mcp.Description("Task ID")),
			mcp.WithArray("tag_ids", mcp.Description("Tag IDs"), mcp.WithStringItems()),
			mcp.WithBoolean("billable", mcp.Description("Whether the entry is billable")),
			mcp.WithString("workspace_id", mcp.Description("Workspace ID (uses default if not provided)")),
		),
		timerStartHandler(r),
	)

	s.AddTool(
		mcp.NewTool("clockify_timer_stop",
			mcp.WithDescription("Stop the currently running timer"),
			mcp.WithString("workspace_id", mcp.Description("Workspace ID (uses default if not provided)")),
		),
		timerStopHandler(r),
	)

	s.AddTool(
		mcp.NewTool("clockify_timer_current",
			mcp.WithDescription("Get the currently running timer"),
			mcp.WithString("workspace_id", mcp.Description("Workspace ID (uses default if not provided)")),
		),
		timerCurrentHandler(r),
	)
}

func timerStartHandler(r *registry) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		wsID := r.workspaceID(req.GetString("workspace_id", ""))
		if wsID == "" {
			return mcp.NewToolResultError("workspace_id is required"), nil
		}

		entry, err := r.client.StartTimer(wsID, clockify.CreateTimeEntryRequest{
			Start:       time.Now().UTC().Format("2006-01-02T15:04:05Z"),
			Description: req.GetString("description", ""),
			ProjectID:   req.GetString("project_id", ""),
			TaskID:      req.GetString("task_id", ""),
			TagIDs:      req.GetStringSlice("tag_ids", nil),
			Billable:    req.GetBool("billable", false),
		})
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to start timer: %v", err)), nil
		}

		return resultJSON(entry)
	}
}

func timerStopHandler(r *registry) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		wsID := r.workspaceID(req.GetString("workspace_id", ""))
		if wsID == "" {
			return mcp.NewToolResultError("workspace_id is required"), nil
		}

		user, err := r.client.GetCurrentUser()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get current user: %v", err)), nil
		}

		entry, err := r.client.StopTimer(wsID, user.ID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to stop timer: %v", err)), nil
		}

		return resultJSON(entry)
	}
}

func timerCurrentHandler(r *registry) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		wsID := r.workspaceID(req.GetString("workspace_id", ""))
		if wsID == "" {
			return mcp.NewToolResultError("workspace_id is required"), nil
		}

		user, err := r.client.GetCurrentUser()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get current user: %v", err)), nil
		}

		entry, err := r.client.GetRunningTimer(wsID, user.ID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get running timer: %v", err)), nil
		}

		if entry == nil {
			return mcp.NewToolResultText("No timer is currently running."), nil
		}

		return resultJSON(entry)
	}
}
