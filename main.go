package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/server"
	"github.com/tedyno/ticktock-mcp/clockify"
	"github.com/tedyno/ticktock-mcp/config"
	"github.com/tedyno/ticktock-mcp/tools"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	client := clockify.NewClient(cfg.APIKey)

	// Resolve default workspace ID
	workspaceID := cfg.WorkspaceID
	if workspaceID == "" {
		workspaces, err := client.GetWorkspaces()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching workspaces: %v\n", err)
			os.Exit(1)
		}
		if len(workspaces) == 0 {
			fmt.Fprintf(os.Stderr, "Error: no workspaces found for this API key\n")
			os.Exit(1)
		}
		workspaceID = workspaces[0].ID
	}

	s := server.NewMCPServer(
		"ticktock-mcp",
		"0.1.0",
		server.WithToolCapabilities(false),
	)

	tools.RegisterAll(s, client, workspaceID)

	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
