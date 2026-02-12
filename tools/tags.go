package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/tedyno/ticktock-mcp/clockify"
)

func registerTagTools(s *server.MCPServer, r *registry) {
	s.AddTool(
		mcp.NewTool("clockify_tag_list",
			mcp.WithDescription("List all tags in a workspace"),
			mcp.WithString("workspace_id", mcp.Description("Workspace ID (uses default if not provided)")),
		),
		tagListHandler(r),
	)

	s.AddTool(
		mcp.NewTool("clockify_tag_create",
			mcp.WithDescription("Create a new tag"),
			mcp.WithString("name", mcp.Required(), mcp.Description("Tag name")),
			mcp.WithString("workspace_id", mcp.Description("Workspace ID (uses default if not provided)")),
		),
		tagCreateHandler(r),
	)

	s.AddTool(
		mcp.NewTool("clockify_tag_update",
			mcp.WithDescription("Update a tag"),
			mcp.WithString("tag_id", mcp.Required(), mcp.Description("Tag ID to update")),
			mcp.WithString("name", mcp.Description("New tag name")),
			mcp.WithBoolean("archived", mcp.Description("Whether the tag is archived")),
			mcp.WithString("workspace_id", mcp.Description("Workspace ID (uses default if not provided)")),
		),
		tagUpdateHandler(r),
	)

	s.AddTool(
		mcp.NewTool("clockify_tag_delete",
			mcp.WithDescription("Delete a tag"),
			mcp.WithString("tag_id", mcp.Required(), mcp.Description("Tag ID to delete")),
			mcp.WithString("workspace_id", mcp.Description("Workspace ID (uses default if not provided)")),
		),
		tagDeleteHandler(r),
	)
}

func tagListHandler(r *registry) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		wsID := r.workspaceID(req.GetString("workspace_id", ""))
		if wsID == "" {
			return mcp.NewToolResultError("workspace_id is required"), nil
		}

		tags, err := r.client.GetTags(wsID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list tags: %v", err)), nil
		}

		return mcp.NewToolResultJSON(tags)
	}
}

func tagCreateHandler(r *registry) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		wsID := r.workspaceID(req.GetString("workspace_id", ""))
		if wsID == "" {
			return mcp.NewToolResultError("workspace_id is required"), nil
		}

		name, err := req.RequireString("name")
		if err != nil {
			return mcp.NewToolResultError("name is required"), nil
		}

		tag, err := r.client.CreateTag(wsID, clockify.CreateTagRequest{Name: name})
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create tag: %v", err)), nil
		}

		return mcp.NewToolResultJSON(tag)
	}
}

func tagUpdateHandler(r *registry) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		wsID := r.workspaceID(req.GetString("workspace_id", ""))
		if wsID == "" {
			return mcp.NewToolResultError("workspace_id is required"), nil
		}

		tagID, err := req.RequireString("tag_id")
		if err != nil {
			return mcp.NewToolResultError("tag_id is required"), nil
		}

		updateReq := clockify.UpdateTagRequest{
			Name: req.GetString("name", ""),
		}

		args := req.GetArguments()
		if _, ok := args["archived"]; ok {
			a := req.GetBool("archived", false)
			updateReq.Archived = &a
		}

		tag, err := r.client.UpdateTag(wsID, tagID, updateReq)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to update tag: %v", err)), nil
		}

		return mcp.NewToolResultJSON(tag)
	}
}

func tagDeleteHandler(r *registry) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		wsID := r.workspaceID(req.GetString("workspace_id", ""))
		if wsID == "" {
			return mcp.NewToolResultError("workspace_id is required"), nil
		}

		tagID, err := req.RequireString("tag_id")
		if err != nil {
			return mcp.NewToolResultError("tag_id is required"), nil
		}

		if err := r.client.DeleteTag(wsID, tagID); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to delete tag: %v", err)), nil
		}

		return mcp.NewToolResultText("Tag deleted successfully."), nil
	}
}
