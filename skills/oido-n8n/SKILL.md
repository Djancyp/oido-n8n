---
name: oido-n8n
description: >
  Complete n8n expert guide. Use for ANY n8n task: managing workflows, executions,
  credentials, webhooks via MCP tools; writing JavaScript or Python Code nodes;
  n8n expression syntax ($json, $node, $workflow); configuring nodes with
  operation-aware parameters; validating workflows and fixing errors; designing
  workflow patterns (webhook, HTTP API, database, AI agent, scheduled, batch).
  Always load this skill when the user mentions n8n, workflows, automation,
  Code nodes, expressions, or any n8n concept.
---

# n8n Expert Guide

## MCP Tools

| Task | Tool |
|------|------|
| List / find workflows | `n8n_list_workflows` |
| Inspect workflow | `n8n_get_workflow` — always before update |
| Create workflow | `n8n_validate_workflow` → `n8n_create_workflow` |
| Full replace | `n8n_update_workflow` |
| Surgical edit | `n8n_update_partial_workflow` (name/nodes/connections/settings) |
| Run | `n8n_execute_workflow` |
| List executions | `n8n_list_executions` — filter by `status` |
| Debug execution | `n8n_get_execution` with `include_data=true` |
| Retry failed | `n8n_retry_execution` |
| Bulk stop | `n8n_stop_executions` |
| Credential setup | `n8n_get_credential_schema` → `n8n_create_credential` |
| Rotate keys | `n8n_update_credential` |
| Webhook trigger | `n8n_trigger_webhook` (no API key needed) |
| Security check | `n8n_generate_audit` |
| Variables | `n8n_list/get/create/update/delete_variable` |
| Users | `n8n_list/get/create_users`, `n8n_change_user_role` |
| Projects | `n8n_list/get/create/update/delete_project` |
| Tags | `n8n_list/get/create/update/delete_tag` |

## Safe Operation Order

```
Create:  n8n_validate_workflow → n8n_create_workflow → n8n_activate_workflow
Update:  n8n_get_workflow → edit → n8n_validate_workflow → n8n_update_workflow
Cred:    n8n_get_credential_schema → n8n_create_credential
Delete:  confirm ID first → delete (irreversible)
```

**Workflow JSON — required shape:**
```json
{
  "nodes": [{ "id":"A","name":"Start","type":"n8n-nodes-base.manualTrigger","typeVersion":1,"position":[240,300],"parameters":{} }],
  "connections": { "Start": { "main": [[{ "node": "Next Node Name", "type": "main", "index": 0 }]] } }
}
```
Connection keys = node `name` (not `id`). Top-level `name` optional.
→ Full examples: [WORKFLOW_FORMAT.md](n8n-workflow-patterns/WORKFLOW_FORMAT.md)

---

## Domain Quick-Refs

**Workflow Patterns** — choose before building (webhook / HTTP API / DB / AI agent / scheduled / batch)
→ [n8n-workflow-patterns/SKILL_GUIDE.md](n8n-workflow-patterns/SKILL_GUIDE.md)

**JavaScript Code Node** — return `[{json:{}}]`, use `$input.all()`, webhook body = `$json.body`
→ [n8n-code-javascript/DATA_ACCESS.md](n8n-code-javascript/DATA_ACCESS.md)

**Python Code Node** — return `[{"json":{}}]`, use `_input.all()`, stdlib only
→ [n8n-code-python/DATA_ACCESS.md](n8n-code-python/DATA_ACCESS.md)

**Expression Syntax** — `{{ $json.field }}` in fields, NOT in code nodes. Webhook: `$json.body.x`
→ [n8n-expression-syntax/COMMON_MISTAKES.md](n8n-expression-syntax/COMMON_MISTAKES.md)

**Node Configuration** — fields depend on `resource`+`operation`. Always verify the node type exists before using it.
→ [n8n-node-configuration/NODE_REFERENCE.md](n8n-node-configuration/NODE_REFERENCE.md) ← look up real type strings by group
→ [n8n-node-configuration/COMMON_NODE_MISTAKES.md](n8n-node-configuration/COMMON_NODE_MISTAKES.md) ← hallucinated/deprecated types
→ [n8n-node-configuration/DEPENDENCIES.md](n8n-node-configuration/DEPENDENCIES.md)

**Validation** — ERROR blocks activation, WARN is advisory. Expect 2–3 fix cycles.
→ [n8n-validation-expert/ERROR_CATALOG.md](n8n-validation-expert/ERROR_CATALOG.md)

---

## Triggers

workflow, automation, n8n, execute, trigger, execution, failed, status,
Code node, JavaScript, Python, $input, $json, $node, pairedItem, SplitInBatches,
expression, `{{ }}`, $workflow, $execution,
node configuration, parameters, displayOptions, operation, resource,
validation error, validate, missing required, false positive,
credential, API key, OAuth, webhook, HTTP Request, database, AI agent, scheduled, batch
