# ticktock-mcp

MCP server for [Clockify](https://clockify.me) time tracking. Provides 27 tools for full Clockify management via the [Model Context Protocol](https://modelcontextprotocol.io).

## Features

- **Timer** — start, stop, get current running timer
- **Time entries** — create, list, update, delete
- **Projects** — CRUD operations
- **Tasks** — CRUD operations (per project)
- **Tags** — CRUD operations
- **Clients** — CRUD operations
- **Workspaces** — list available workspaces
- **Users** — current user, list workspace users
- **Reports** — summary and detailed reports with filters

Every tool supports an optional `workspace_id` parameter to override the default workspace.

## Installation

```bash
go install github.com/tedyno/ticktock-mcp@latest
```

Or build from source:

```bash
git clone https://github.com/tedyno/ticktock-mcp.git
cd ticktock-mcp
go build -o ticktock-mcp .
```

## Configuration

### API Key

Set your Clockify API key via environment variable:

```bash
export CLOCKIFY_API_KEY=your-api-key
```

Or create a config file at `~/.config/ticktock-mcp/config.json`:

```json
{
  "api_key": "your-api-key",
  "workspace_id": "optional-default-workspace-id"
}
```

Environment variables take priority over the config file.

You can get your API key from [Clockify Settings](https://app.clockify.me/user/preferences#advanced).

## Usage with Claude Code

Add to your Claude Code MCP configuration (`~/.claude/claude_desktop_config.json` or project `.mcp.json`):

```json
{
  "mcpServers": {
    "clockify": {
      "command": "ticktock-mcp",
      "env": {
        "CLOCKIFY_API_KEY": "your-api-key"
      }
    }
  }
}
```

Or if using a local build:

```json
{
  "mcpServers": {
    "clockify": {
      "command": "/path/to/ticktock-mcp",
      "env": {
        "CLOCKIFY_API_KEY": "your-api-key"
      }
    }
  }
}
```

You can also add it via CLI:

```bash
claude mcp add clockify -- /path/to/ticktock-mcp
```

## Available Tools

| Tool | Description |
|------|-------------|
| `clockify_timer_start` | Start a new timer |
| `clockify_timer_stop` | Stop the running timer |
| `clockify_timer_current` | Get the running timer |
| `clockify_time_entry_list` | List time entries |
| `clockify_time_entry_create` | Create a manual time entry |
| `clockify_time_entry_update` | Update a time entry |
| `clockify_time_entry_delete` | Delete a time entry |
| `clockify_project_list` | List projects |
| `clockify_project_create` | Create a project |
| `clockify_project_update` | Update a project |
| `clockify_project_delete` | Delete a project |
| `clockify_task_list` | List tasks in a project |
| `clockify_task_create` | Create a task |
| `clockify_task_update` | Update a task |
| `clockify_task_delete` | Delete a task |
| `clockify_tag_list` | List tags |
| `clockify_tag_create` | Create a tag |
| `clockify_tag_update` | Update a tag |
| `clockify_tag_delete` | Delete a tag |
| `clockify_client_list` | List clients |
| `clockify_client_create` | Create a client |
| `clockify_client_update` | Update a client |
| `clockify_client_delete` | Delete a client |
| `clockify_workspace_list` | List workspaces |
| `clockify_user_current` | Get current user |
| `clockify_user_list` | List workspace users |
| `clockify_report_summary` | Generate summary report |
| `clockify_report_detailed` | Generate detailed report |

## License

MIT
