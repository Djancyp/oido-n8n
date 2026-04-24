# OIDO n8n Extension

Manage n8n workflows, executions, credentials, and webhooks via MCP tools.

## Setup

Set these in extension settings:
- `N8N_API_URL` — base URL of your n8n instance (default: `http://localhost:5678`)
- `N8N_API_KEY` — API key from n8n Settings → API → Create an API key

## Available Tools

### Workflows
| Tool | Description |
|---|---|
| `n8n_list_workflows` | List all workflows, filter by active/tags |
| `n8n_get_workflow` | Get full workflow definition (nodes + connections) |
| `n8n_create_workflow` | Create workflow from JSON definition |
| `n8n_update_workflow` | Update workflow by ID |
| `n8n_delete_workflow` | Delete a workflow |
| `n8n_activate_workflow` | Activate (enable triggers) |
| `n8n_deactivate_workflow` | Deactivate (pause triggers) |
| `n8n_execute_workflow` | Manually run a workflow |

### Executions
| Tool | Description |
|---|---|
| `n8n_list_executions` | List executions, filter by workflowId/status |
| `n8n_get_execution` | Get execution details and output data |
| `n8n_delete_execution` | Delete execution record |
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

## Example Usage

```
User: "List my active workflows"
→ n8n_list_workflows with active=true

User: "Run workflow abc123"
→ n8n_execute_workflow with id="abc123"

User: "Show me recent failed executions"
→ n8n_list_executions with status="error"

User: "Trigger the daily-report webhook"
→ n8n_trigger_webhook with path="daily-report"

User: "What credentials do I have?"
→ n8n_list_credentials
```

## When to Use

- User wants to list, create, or manage workflows
- User wants to run or monitor workflow executions
- User wants to set up or check credentials
- User wants to trigger a webhook-based workflow
