# Common Node Type Mistakes

## Nodes That Do NOT Exist

These are hallucinated or deprecated node types — never use them:

| Wrong | Correct | Notes |
|-------|---------|-------|
| `n8n-nodes-base.cron` | `n8n-nodes-base.scheduleTrigger` v1 | Old name, removed |
| `n8n-nodes-base.timeTrigger` | `n8n-nodes-base.scheduleTrigger` v1 | Never existed |
| `n8n-nodes-base.triggerCron` | `n8n-nodes-base.scheduleTrigger` v1 | Never existed |
| `n8n-nodes-base.start` | `n8n-nodes-base.manualTrigger` v1 | Never existed |
| `n8n-nodes-base.filter` | `n8n-nodes-base.if` v2 | Never existed; use IF node |
| `n8n-nodes-base.hackerNews` | `n8n-nodes-base.httpRequest` | HN has no native node |
| `n8n-nodes-base.sendEmail` | `n8n-nodes-base.emailSend` v2 | Transposed name |
| `n8n-nodes-base.errorHandler` | `n8n-nodes-base.stopAndError` | For explicit error stops |
| `n8n-nodes-base.functionItem` | `n8n-nodes-base.code` v2 | Deprecated in n8n ≥1.0 |
| `n8n-nodes-base.function` | `n8n-nodes-base.code` v2 | Deprecated in n8n ≥1.0 |
| `n8n-nodes-base.set` (v1) | `n8n-nodes-base.set` v3 | Use typeVersion 3 |

## Deprecated Nodes → Replacements

| Deprecated | Replacement | typeVersion |
|-----------|-------------|-------------|
| `cron` | `scheduleTrigger` | 1 |
| `functionItem` | `code` | 2 |
| `function` | `code` | 2 |
| `itemLists` | `splitInBatches` or `aggregate` | — |
| `moveBinaryData` | `extractFromFile` / `convertToFile` | — |

## Correct typeVersions for Common Nodes

Always use the latest stable version:

| Node | typeVersion |
|------|-------------|
| `manualTrigger` | 1 |
| `httpRequest` | 4 |
| `code` | 2 |
| `emailSend` | 2 |
| `set` | 3 |
| `if` | 2 |
| `splitInBatches` | 3 |
| `merge` | 3 |
| `stopAndError` | 1 |
| `noOp` | 1 |
| `scheduleTrigger` | 1 |
| `webhook` | 2 |
| `respondToWebhook` | 1 |

## Fetching External APIs Without a Native Node

Many services have no native node. Always fall back to `httpRequest`:

```json
{
  "id": "fetch",
  "name": "Fetch Data",
  "type": "n8n-nodes-base.httpRequest",
  "typeVersion": 4,
  "position": [460, 300],
  "parameters": {
    "url": "https://api.example.com/endpoint",
    "method": "GET",
    "responseFormat": "json"
  }
}
```

## Error Handling Patterns

### Stop execution and report error (inline)
Throw inside a Code node — workflow stops, error is shown in execution details:
```javascript
if (!data || data.length === 0) {
  throw new Error('No data returned from API');
}
```

### Explicit Stop and Error node
Use `n8n-nodes-base.stopAndError` on a false branch from an IF node:
```json
{
  "id": "err",
  "name": "Stop and Report Error",
  "type": "n8n-nodes-base.stopAndError",
  "typeVersion": 1,
  "position": [900, 420],
  "parameters": {
    "errorMessage": "Describe what went wrong here"
  }
}
```

### Error workflow (separate workflow triggered on failure)
Set in workflow settings → "Error Workflow". A separate workflow with an `errorTrigger` node receives the error data. Use this for alerts/logging across all workflows.
