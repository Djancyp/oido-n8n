---
name: oido-n8n
description: Manage n8n workflows, executions, credentials, and webhooks via MCP
---

# Oido n8n Extension

## Overview

Full control over a self-hosted n8n instance via its REST API v1. Use these tools to list, create, run, and monitor workflows, manage credentials, and trigger webhooks.

## Available Tools

### `n8n_list_workflows`
List all workflows. Filter by active status or tags.
- `active` (boolean, optional): Filter to active/inactive only
- `tags` (string[], optional): Filter by tag names
- `limit` (number, optional): Max results
- `cursor` (string, optional): Pagination cursor

**When to use:** User wants to see their workflows, check what's running, find a workflow by name.

### `n8n_get_workflow`
Get full workflow definition including nodes and connections.
- `id` (string, required): Workflow ID

**When to use:** User wants to inspect a workflow, see its nodes, or before updating it.

### `n8n_create_workflow`
Create a new workflow from a JSON definition.
- `workflow_json` (string, required): Full n8n workflow JSON

**When to use:** User wants to build a new workflow programmatically.

### `n8n_update_workflow`
Update an existing workflow.
- `id` (string, required): Workflow ID
- `workflow_json` (string, required): Updated workflow JSON

### `n8n_delete_workflow`
Delete a workflow permanently.
- `id` (string, required): Workflow ID

### `n8n_activate_workflow` / `n8n_deactivate_workflow`
Toggle a workflow's trigger responses.
- `id` (string, required): Workflow ID

**When to use:** User wants to pause or resume a workflow without deleting it.

### `n8n_execute_workflow`
Manually trigger a workflow execution.
- `id` (string, required): Workflow ID
- `data` (string, optional): Input data as JSON string

**When to use:** User wants to run a workflow now, test it, or trigger it with specific data.

### `n8n_list_executions`
List workflow execution history.
- `workflow_id` (string, optional): Filter by workflow
- `status` (string, optional): waiting | running | success | error | canceled
- `limit` (number, optional): Max results
- `cursor` (string, optional): Pagination cursor

**When to use:** User wants to check if workflows ran successfully, debug failures.

### `n8n_get_execution`
Get full details of a specific execution including output data.
- `id` (string, required): Execution ID

### `n8n_delete_execution` / `n8n_stop_execution`
Clean up or halt executions.
- `id` (string, required): Execution ID

### `n8n_list_credentials`
List credential names and types (no secret values returned).
- `limit` (number, optional)

### `n8n_create_credential`
Create a new credential.
- `name` (string, required): Credential name
- `type` (string, required): Credential type (e.g. `httpBasicAuth`, `githubApi`)
- `data` (string, optional): Credential fields as JSON

**Tip:** Use `n8n_get_credential_schema` first to know required fields.

### `n8n_delete_credential`
Delete a credential by ID.

### `n8n_get_credential_schema`
Get the field schema for a credential type.
- `credential_type` (string, required): e.g. `httpBasicAuth`, `githubApi`, `slackApi`

### `n8n_list_tags` / `n8n_create_tag`
Manage workflow tags.

### `n8n_trigger_webhook`
Trigger a webhook-enabled workflow directly (bypasses API key).
- `path` (string, required): Webhook path (without `/webhook/` prefix)
- `method` (string, optional): GET, POST, PUT, PATCH, DELETE (default: POST)
- `body` (string, optional): Request body as JSON string

## Example Interactions

```
User: "What workflows are active?"
→ n8n_list_workflows with active=true

User: "Run the data-sync workflow"
→ n8n_list_workflows to find the ID, then n8n_execute_workflow

User: "Did my nightly backup workflow succeed?"
→ n8n_list_executions with status="error" + workflow_id

User: "Create a Slack credential"
→ n8n_get_credential_schema with credential_type="slackApi"
→ n8n_create_credential with the required fields

User: "Trigger my order-processing webhook"
→ n8n_trigger_webhook with path="order-processing"
```

## Triggers

Use these tools when you see:
- "workflow", "automation", "n8n"
- "run", "execute", "trigger"
- "execution", "failed", "status"
- "credential", "API key", "connect"
- "webhook", "endpoint"

## Related Commands
- `/workflow-list` — List all workflows
- `/workflow-run` — Execute a workflow by ID
- `/workflow-status` — Check recent execution status
- `/execution-status` — Get details of a specific execution
