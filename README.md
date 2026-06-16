# oido-n8n

MCP plugin for n8n. Manage workflows, executions, credentials, tags, and webhooks directly from your AI assistant.

## Features

- **Workflow management** — list, create, update, delete, activate, deactivate, and execute workflows
- **Execution monitoring** — list, inspect, stop, and delete execution records
- **Credential management** — list, create, and delete credentials; inspect credential schemas
- **Tags** — list and create tags for workflow organization
- **Webhooks** — trigger webhook-based workflows without an API key
- **Static binary** — no runtime dependencies, CGO disabled

## Installation

```bash
make build
```

Then register `oido-n8n-mcp` as an MCP server or install via the Oido extension system.

## Configuration

| Variable | Required | Default | Description |
|---|---|---|---|
| `N8N_API_URL` | yes | `http://localhost:5678` | Base URL of your n8n instance |
| `N8N_API_KEY` | yes | — | API key from n8n Settings → API → Create an API key |

## Tools

### Workflows

| Tool | Description |
|---|---|
| `n8n_list_workflows` | List workflows, filter by active status, tags, or name |
| `n8n_get_workflow` | Get full workflow definition (nodes + connections) |
| `n8n_create_workflow` | Create a workflow from a JSON definition |
| `n8n_update_workflow` | Update an existing workflow by ID |
| `n8n_delete_workflow` | Delete a workflow |
| `n8n_activate_workflow` | Activate a workflow (enable triggers) |
| `n8n_deactivate_workflow` | Deactivate a workflow (pause triggers) |
| `n8n_execute_workflow` | Manually run a workflow |
| `n8n_validate_workflow` | Validate JSON: structure, node types, connection integrity, trigger presence |

### Node Knowledge

| Tool | Description |
|---|---|
| `n8n_search_nodes` | Find node types by keyword (comma-separated OR; filter by group t/i/o) |
| `n8n_get_node_schema` | Get a node's parameters, inputs, outputs, and version |

### Executions

| Tool | Description |
|---|---|
| `n8n_list_executions` | List executions, filter by workflow ID or status |
| `n8n_get_execution` | Get execution details and node output data |
| `n8n_delete_execution` | Delete an execution record |
| `n8n_stop_execution` | Stop a running execution |

### Credentials

| Tool | Description |
|---|---|
| `n8n_list_credentials` | List credential names and types |
| `n8n_create_credential` | Create a new credential |
| `n8n_delete_credential` | Delete a credential |
| `n8n_get_credential_schema` | Get required fields for a credential type |

### Tags

| Tool | Description |
|---|---|
| `n8n_list_tags` | List all tags |
| `n8n_create_tag` | Create a new tag |

### Webhooks

| Tool | Description |
|---|---|
| `n8n_trigger_webhook` | Trigger a workflow via webhook path (no API key needed) |

## Usage Examples

```
# List active workflows
"Show me all active workflows"

# Run a workflow
"Execute workflow abc123"

# Check failed executions
"Show me recent failed executions"

# Trigger a webhook
"Trigger the daily-report webhook"

# Inspect credentials
"What credentials do I have configured?"

# Create a tag
"Create a tag called production"
```

## Development

```bash
make build   # build binary
make dist    # build + zip for distribution
make clean   # remove binary and dist/
```

## Tech Stack

- Go 1.26+ with `CGO_ENABLED=0` (fully static)
- [`modelcontextprotocol/go-sdk`](https://github.com/modelcontextprotocol/go-sdk) — MCP server
- n8n REST API v1
