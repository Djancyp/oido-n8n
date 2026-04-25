---
name: oido-n8n
description: >
  n8n workflow creation expert. Use when building, designing, or creating n8n workflows:
  picking the right node types, writing JavaScript or Python Code nodes, using
  expression syntax ($json, $node), configuring node parameters, validating workflow
  JSON, and choosing workflow patterns (webhook, HTTP API, database, AI agent,
  scheduled, batch). Always load when user asks to create, build, or design an n8n workflow.
---

# n8n Workflow Creation Guide

## Create Workflow — Step by Step

```
1. n8n_validate_workflow   ← check JSON before sending
2. n8n_create_workflow     ← POST to n8n
3. n8n_activate_workflow   ← enable triggers (optional)
```

## Workflow JSON — Minimum Required

```json
{
  "name": "My Workflow",
  "nodes": [
    {
      "id": "trigger",
      "name": "Manual Trigger",
      "type": "n8n-nodes-base.manualTrigger",
      "typeVersion": 1,
      "position": [240, 300],
      "parameters": {}
    }
  ],
  "connections": {
    "Manual Trigger": {
      "main": [[{ "node": "Next Node Name", "type": "main", "index": 0 }]]
    }
  }
}
```

**Rules:**
- Connection keys = node `name` (NOT `id`)
- Last node in chain has no `connections` entry
- `name` at top level is optional — n8n auto-assigns
- AI sub-nodes use `ai_languageModel` / `ai_memory` / `ai_tool` instead of `main`

→ Full schema (all node fields + settings): [WORKFLOW_FORMAT.md](n8n-workflow-patterns/WORKFLOW_FORMAT.md)

## Node Type — Look Up Before Using

Never guess a node type. Look it up:
→ [NODE_REFERENCE.md](n8n-node-configuration/NODE_REFERENCE.md) — quick-lookup + all groups
→ [COMMON_NODE_MISTAKES.md](n8n-node-configuration/COMMON_NODE_MISTAKES.md) — hallucinated/deprecated types

**Most used:**
```
n8n-nodes-base.manualTrigger       v1
n8n-nodes-base.scheduleTrigger     v1
n8n-nodes-base.webhook             v2
n8n-nodes-base.httpRequest         v4   ← use for any API without a native node
n8n-nodes-base.code                v2   ← replaces function / functionItem
n8n-nodes-base.set                 v3
n8n-nodes-base.if                  v2
n8n-nodes-base.emailSend           v2   ← NOT sendEmail
n8n-nodes-base.slack               v2
n8n-nodes-base.googleSheets        v4
n8n-nodes-base.postgres            v2
n8n-nodes-base.stopAndError        v1   ← NOT errorHandler
n8n-nodes-base.respondToWebhook    v1
@n8n/n8n-nodes-langchain.agent     v1
```

## Choose a Pattern First

| Pattern | Trigger | Use when |
|---------|---------|----------|
| Webhook | `webhook` v2 | Receiving external events |
| HTTP API | `scheduleTrigger` or `webhook` | Fetching from REST APIs |
| Database | `scheduleTrigger` | Read/write/sync DB records |
| AI Agent | `chatTrigger` or `webhook` | LLM + tools + memory |
| Scheduled | `scheduleTrigger` | Reports, recurring jobs |
| Batch | `scheduleTrigger` + `splitInBatches` | Large datasets, rate limits |

→ Pattern skeletons: [SKILL_GUIDE.md](n8n-workflow-patterns/SKILL_GUIDE.md)

## Code Node (JavaScript)

```javascript
// typeVersion: 2  |  mode: "runOnceForAllItems"
const items = $input.all();

if (!items[0].json.hits) throw new Error('No data from API');

return items.map(item => ({
  json: { ...item.json, processed: true }
}));
```

**Rules:** return `[{json:{}}]` · webhook body = `$json.body` · no `{{ }}` syntax inside code
→ [DATA_ACCESS.md](n8n-code-javascript/DATA_ACCESS.md) · [ERROR_PATTERNS.md](n8n-code-javascript/ERROR_PATTERNS.md)

## Expression Syntax (in node parameter fields)

```
{{ $json.fieldName }}              current item
{{ $node["Node Name"].json.field }} previous node
{{ $json.body.data }}              webhook body field
```

NOT inside Code nodes — use `$input.first().json.field` there instead.
→ [COMMON_MISTAKES.md](n8n-expression-syntax/COMMON_MISTAKES.md)

## Validation

Run `n8n_validate_workflow` before every create. Expect 2–3 fix cycles.
- `ERROR` — blocks activation, must fix
- `WARN` — advisory, fix before production

→ [ERROR_CATALOG.md](n8n-validation-expert/ERROR_CATALOG.md)

---

## Triggers

create workflow, build workflow, design workflow, new workflow, automate,
Code node, JavaScript, Python, $input, $json, expression, `{{ }}`,
node type, httpRequest, webhook, trigger, schedule, AI agent,
validate workflow, workflow JSON, connections, nodes
