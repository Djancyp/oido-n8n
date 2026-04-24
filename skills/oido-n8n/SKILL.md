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

One skill covering all n8n domains. Jump to the relevant section.

---

## MCP Tools

### Tool Selection

| Task | Tool | Notes |
|------|------|-------|
| List workflows | `n8n_list_workflows` | Filter: `active=true/false`, `tags` |
| Inspect workflow | `n8n_get_workflow` | Always call before update |
| Create workflow | `n8n_create_workflow` | Run `n8n_validate_workflow` first |
| Update full workflow | `n8n_update_workflow` | Full JSON replacement |
| Update specific fields | `n8n_update_partial_workflow` | Surgical: name/nodes/connections/settings only |
| Validate before save | `n8n_validate_workflow` | Client-side structural check |
| Run workflow | `n8n_execute_workflow` | Pass optional `data` JSON |
| List executions | `n8n_list_executions` | Filter: `status=error/success/running` |
| Debug execution | `n8n_get_execution` | Returns full output data |
| Stop execution | `n8n_stop_execution` | Running executions only |
| List credentials | `n8n_list_credentials` | Names/types only, no secrets |
| Create credential | `n8n_create_credential` | Get schema first |
| Get field schema | `n8n_get_credential_schema` | Always call before create |
| Trigger webhook | `n8n_trigger_webhook` | No API key required |

### Safe Operation Order

**Creating a workflow:**
```
n8n_validate_workflow → n8n_create_workflow → n8n_activate_workflow
```

**Updating a workflow:**
```
n8n_get_workflow → [edit] → n8n_validate_workflow → n8n_update_workflow
```
Or for small changes: `n8n_update_partial_workflow` directly.

**Creating a credential:**
```
n8n_get_credential_schema → n8n_create_credential
```

**Deleting a workflow:** Irreversible. Confirm ID with `n8n_get_workflow` first.

---

## Workflow Patterns

Choose the right pattern before building:

| Pattern | Use when | Key nodes |
|---------|----------|-----------|
| **Webhook** | Receiving external events, instant response needed | Webhook, Respond to Webhook |
| **HTTP API** | Fetching/syncing with external APIs | Schedule/Webhook, HTTP Request, Set |
| **Database** | Read/write/sync database data | Schedule, DB node, IF, Set |
| **AI Agent** | LLM reasoning with tools | Trigger, AI Agent, Tool nodes |
| **Scheduled** | Recurring reports/jobs | Schedule, HTTP/DB, Send nodes |
| **Batch** | Large datasets, API rate limits | Schedule, SplitInBatches, Merge |

### Pattern Structure

**Webhook Processing** (most common):
```
Webhook → Validate input → Transform → Respond to Webhook
                                    → Notify (async branch)
```

**HTTP API Integration:**
```
Trigger → HTTP Request → IF (success?) → Transform → Action → Error Handler
```

**AI Agent:**
```
Trigger → AI Agent [Model + Tools + Memory] → Output formatter → Respond
```

For full pattern templates see:
- [webhook_processing.md](n8n-workflow-patterns/webhook_processing.md)
- [http_api_integration.md](n8n-workflow-patterns/http_api_integration.md)
- [database_operations.md](n8n-workflow-patterns/database_operations.md)
- [ai_agent_workflow.md](n8n-workflow-patterns/ai_agent_workflow.md)
- [scheduled_tasks.md](n8n-workflow-patterns/scheduled_tasks.md)

---

## Code Nodes — JavaScript

### Essential Rules

1. **Default mode: "Run Once for All Items"** — use for 95% of cases
2. **Access items:** `$input.all()` (all) or `$input.first()` (one)
3. **Return format is mandatory:**
   ```javascript
   return [{ json: { key: "value" } }]
   ```
4. **Webhook body is nested:**
   ```javascript
   const body = $input.first().json.body  // NOT $json directly
   ```
5. **Built-ins:** `$helpers.httpRequest()`, `DateTime` (Luxon), `$jmespath()`

### Common Pattern

```javascript
// Run Once for All Items
const items = $input.all();

const result = items.map(item => ({
  json: {
    ...item.json,
    processed: true,
    timestamp: DateTime.now().toISO()
  }
}));

return result;
```

### Mode: Run Once for Each Item

Use only when: per-item side effects with different outputs, or `pairedItem` tracking needed.
```javascript
// Run Once for Each Item
const item = $input.item;
return [{ json: { value: item.json.value * 2 }, pairedItem: 0 }];
```

### Deep dives:
- [DATA_ACCESS.md](n8n-code-javascript/DATA_ACCESS.md)
- [COMMON_PATTERNS.md](n8n-code-javascript/COMMON_PATTERNS.md)
- [ERROR_PATTERNS.md](n8n-code-javascript/ERROR_PATTERNS.md)
- [BUILTIN_FUNCTIONS.md](n8n-code-javascript/BUILTIN_FUNCTIONS.md)

---

## Code Nodes — Python

### Essential Rules

1. **Access items:** `_input.all()` or `items` variable
2. **Return format is mandatory:**
   ```python
   return [{"json": {"key": "value"}}]
   ```
3. **Webhook body:**
   ```python
   body = _input.first().json["body"]  # NOT json directly
   ```
4. **Standard library only** — no pip installs. Available: `json`, `datetime`, `re`, `math`, `hashlib`, `base64`, `urllib`

### Common Pattern

```python
items = _input.all()

result = []
for item in items:
    result.append({
        "json": {
            **item.json,
            "processed": True
        }
    })

return result
```

### Deep dives:
- [DATA_ACCESS.md](n8n-code-python/DATA_ACCESS.md)
- [COMMON_PATTERNS.md](n8n-code-python/COMMON_PATTERNS.md)
- [ERROR_PATTERNS.md](n8n-code-python/ERROR_PATTERNS.md)
- [STANDARD_LIBRARY.md](n8n-code-python/STANDARD_LIBRARY.md)

---

## Expression Syntax

### Core Rules

- All expressions: `{{ expression }}`
- Current item data: `{{ $json.fieldName }}`
- Previous node data: `{{ $node["Node Name"].json.field }}`
- Workflow metadata: `{{ $workflow.name }}`, `{{ $execution.id }}`
- String interpolation: `Hello {{ $json.name }}!`

### CRITICAL: Expressions vs Code Nodes

| Context | Syntax | Example |
|---------|--------|---------|
| Expression field | `{{ $json.name }}` | Works in node parameter fields |
| Code node (JS) | `$input.first().json.name` | No `{{ }}`, no `$json` shorthand |
| Code node (Python) | `_input.first().json["name"]` | No `{{ }}` |

### Common Mistakes

```
❌ {{ $json.body.data }}      in webhook — body is $json.body, not $json
✓  {{ $json.body.data }}      only after a Set node that extracts body

❌ $json.name                 in expression field (missing {{ }})
✓  {{ $json.name }}

❌ {{ $input.first().json }}  in expression field (Code node syntax)
✓  {{ $json }}
```

### Deep dives:
- [EXAMPLES.md](n8n-expression-syntax/EXAMPLES.md)
- [COMMON_MISTAKES.md](n8n-expression-syntax/COMMON_MISTAKES.md)

---

## Node Configuration

### Operation-Aware Config

Fields change based on `resource` + `operation`. Always check which fields apply:

```javascript
// Slack: operation=post needs channel + text
// Slack: operation=update needs messageId + text (NOT channel)
```

**Pattern:** `get_node(detail="standard")` → check `displayOptions` → configure only visible fields.

### Property Dependencies

Fields appear/disappear based on other field values. Never assume a field is always required — it depends on the current operation.

```javascript
// HTTP Request: method=GET → no body fields visible
// HTTP Request: method=POST → body fields appear
```

### Surgical vs Full Update

| Change | Use |
|--------|-----|
| Change one node's params | `n8n_update_partial_workflow` with `nodes` |
| Rename workflow | `n8n_update_partial_workflow` with `name` |
| Restructure entire workflow | `n8n_update_workflow` with full JSON |

### Deep dives:
- [DEPENDENCIES.md](n8n-node-configuration/DEPENDENCIES.md)
- [OPERATION_PATTERNS.md](n8n-node-configuration/OPERATION_PATTERNS.md)

---

## Validation

### Error Severity

| Level | Meaning | Action |
|-------|---------|--------|
| `ERROR` | Blocks execution | Must fix before activate |
| `WARN` | Works but fragile | Fix before production |
| `INFO` | Best practice | Optional improvement |

### Common Error Types

- `missing_required` — required field absent for this operation
- `invalid_value` — value not in allowed enum
- `type_mismatch` — wrong type (string vs number)
- `invalid_reference` — referenced node doesn't exist
- `invalid_expression` — `{{ }}` syntax error

### Iterative Fix Loop

Validation is never one-shot. Expect 2–3 cycles:
```
validate → fix errors → validate → fix warnings → validate → ✓
```

### False Positives

Some warnings can be ignored:
- `best_practice` warnings on dynamically-built node params
- `performance` warnings when dataset is known to be small

### Deep dives:
- [ERROR_CATALOG.md](n8n-validation-expert/ERROR_CATALOG.md)
- [FALSE_POSITIVES.md](n8n-validation-expert/FALSE_POSITIVES.md)

---

## Triggers

Load this skill when you see ANY of:

**Workflow management:** workflow, automation, n8n, activate, deactivate, run, execute, trigger, execution, failed, status, history

**Code nodes:** Code node, JavaScript, Python, $input, $json, $node, pairedItem, Run Once, SplitInBatches, $helpers, DateTime, Luxon

**Expressions:** expression, `{{ }}`, $json, $workflow, $execution, $node reference

**Node config:** node configuration, parameters, displayOptions, operation, resource, required fields, property dependencies

**Validation:** validation error, validate workflow, missing required, invalid value, false positive, auto-fix

**Credentials:** credential, API key, httpBasicAuth, connect service, OAuth

**Patterns:** webhook, HTTP Request, database sync, AI agent, scheduled task, batch processing, SplitInBatches
