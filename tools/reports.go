package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/tedyno/ticktock-mcp/clockify"
)

func registerReportTools(s *server.MCPServer, r *registry) {
	s.AddTool(
		mcp.NewTool("clockify_report_summary",
			mcp.WithDescription("Generate a summary report for a workspace"),
			mcp.WithString("start", mcp.Required(), mcp.Description("Report start date (ISO 8601, e.g. 2024-01-01T00:00:00Z)")),
			mcp.WithString("end", mcp.Required(), mcp.Description("Report end date (ISO 8601)")),
			mcp.WithString("project_id", mcp.Description("Filter by project ID")),
			mcp.WithString("user_id", mcp.Description("Filter by user ID")),
			mcp.WithString("workspace_id", mcp.Description("Workspace ID (uses default if not provided)")),
		),
		reportSummaryHandler(r),
	)

	s.AddTool(
		mcp.NewTool("clockify_report_detailed",
			mcp.WithDescription("Generate a detailed report for a workspace"),
			mcp.WithString("start", mcp.Required(), mcp.Description("Report start date (ISO 8601, e.g. 2024-01-01T00:00:00Z)")),
			mcp.WithString("end", mcp.Required(), mcp.Description("Report end date (ISO 8601)")),
			mcp.WithString("project_id", mcp.Description("Filter by project ID")),
			mcp.WithString("user_id", mcp.Description("Filter by user ID")),
			mcp.WithNumber("page", mcp.Description("Page number (default 1)")),
			mcp.WithNumber("page_size", mcp.Description("Page size (default 50)")),
			mcp.WithString("workspace_id", mcp.Description("Workspace ID (uses default if not provided)")),
		),
		reportDetailedHandler(r),
	)
}

func reportSummaryHandler(r *registry) server.ToolHandlerFunc {
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

		reportReq := clockify.SummaryReportRequest{
			DateRangeStart: start,
			DateRangeEnd:   end,
			SummaryFilter:  &clockify.SummaryFilter{Groups: []string{"PROJECT", "TIMEENTRY"}},
		}

		if projectID := req.GetString("project_id", ""); projectID != "" {
			reportReq.Projects = &clockify.ReportProjectFilter{IDs: []string{projectID}}
		}
		if userID := req.GetString("user_id", ""); userID != "" {
			reportReq.Users = &clockify.ReportUsersFilter{IDs: []string{userID}}
		}

		report, err := r.client.GetSummaryReport(wsID, reportReq)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get summary report: %v", err)), nil
		}

		return mcp.NewToolResultJSON(report)
	}
}

func reportDetailedHandler(r *registry) server.ToolHandlerFunc {
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

		page := req.GetInt("page", 1)
		pageSize := req.GetInt("page_size", 50)

		reportReq := clockify.DetailedReportRequest{
			DateRangeStart: start,
			DateRangeEnd:   end,
			Page:           page,
			PageSize:       pageSize,
			DetailedFilter: &clockify.DetailedFilter{Page: page, PageSize: pageSize},
		}

		if projectID := req.GetString("project_id", ""); projectID != "" {
			reportReq.Projects = &clockify.ReportProjectFilter{IDs: []string{projectID}}
		}
		if userID := req.GetString("user_id", ""); userID != "" {
			reportReq.Users = &clockify.ReportUsersFilter{IDs: []string{userID}}
		}

		report, err := r.client.GetDetailedReport(wsID, reportReq)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get detailed report: %v", err)), nil
		}

		return mcp.NewToolResultJSON(report)
	}
}
