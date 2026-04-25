# Workflow JSON Format

## Top-Level Structure

```json
{
  "name": "My Workflow",
  "nodes": [...],
  "connections": {...},
  "settings": {...}
}
```

| Field | Required | Notes |
|-------|----------|-------|
| `nodes` | Yes | Array of node objects |
| `connections` | Yes | Wiring between nodes |
| `name` | No | n8n auto-assigns if omitted |
| `settings` | No | Execution behaviour, error workflow, timezone |

---

## Node Object — Full Schema

```json
{
  "id": "0f5532f9-36ba-4bef-86c7-30d607400b15",
  "name": "Jira",
  "type": "n8n-nodes-base.jira",
  "typeVersion": 1,
  "position": [-100, 80],
  "parameters": {
    "additionalProperties": {}
  },
  "credentials": {
    "jiraSoftwareCloudApi": {
      "id": "35",
      "name": "jiraApi"
    }
  },
  "disabled": false,
  "executeOnce": false,
  "alwaysOutputData": false,
  "retryOnFail": false,
  "maxTries": 1,
  "waitBetweenTries": 1,
  "onError": "stopWorkflow",
  "notesInFlow": false,
  "notes": "",
  "webhookId": ""
}
```

### Node Field Reference

| Field | Required | Default | Notes |
|-------|----------|---------|-------|
| `id` | Yes | — | UUID or any unique string |
| `name` | Yes | — | Canvas label. **Used as connection key.** |
| `type` | Yes | — | Full type string e.g. `n8n-nodes-base.httpRequest` |
| `typeVersion` | Yes | — | See [NODE_REFERENCE.md](../n8n-node-configuration/NODE_REFERENCE.md) |
| `position` | Yes | — | `[x, y]`. Increment x by ~220 per node |
| `parameters` | Yes | `{}` | Node-specific config. Use `{}` for defaults |
| `credentials` | No | — | Required when node needs auth. See below. |
| `disabled` | No | `false` | Skip node without removing it |
| `executeOnce` | No | `false` | Run only for first item, ignore rest |
| `alwaysOutputData` | No | `false` | Output empty item even if no data |
| `retryOnFail` | No | `false` | Auto-retry on error |
| `maxTries` | No | `1` | Max retry attempts (requires `retryOnFail: true`) |
| `waitBetweenTries` | No | `1000` | ms between retries |
| `onError` | No | `"stopWorkflow"` | `"stopWorkflow"` or `"continueRegularOutput"` or `"continueErrorOutput"` |
| `notesInFlow` | No | `false` | Show `notes` as sticky label on canvas |
| `notes` | No | `""` | Internal documentation for the node |
| `webhookId` | No | `""` | Set by n8n for webhook nodes |

### Credentials Format

```json
"credentials": {
  "<credentialType>": {
    "id": "<credential-id-from-n8n>",
    "name": "<display-name>"
  }
}
```

Get credential IDs with `n8n_list_credentials`. Credential type keys match the node's expected auth type (e.g. `jiraSoftwareCloudApi`, `slackApi`, `googleSheetsOAuth2Api`).

---

## Settings Object — Full Schema

```json
"settings": {
  "saveExecutionProgress": true,
  "saveManualExecutions": true,
  "saveDataErrorExecution": "all",
  "saveDataSuccessExecution": "all",
  "executionTimeout": 3600,
  "errorWorkflow": "VzqKEW0ShTXA5vPj",
  "timezone": "America/New_York",
  "executionOrder": "v1",
  "callerPolicy": "workflowsFromSameOwner",
  "callerIds": "14, 18, 23",
  "timeSavedPerExecution": 1,
  "availableInMCP": false
}
```

| Field | Values | Notes |
|-------|--------|-------|
| `saveExecutionProgress` | `true` / `false` | Save node-by-node progress (higher DB load) |
| `saveManualExecutions` | `true` / `false` | Save runs triggered manually |
| `saveDataErrorExecution` | `"all"` / `"none"` | Persist data for failed runs |
| `saveDataSuccessExecution` | `"all"` / `"none"` | Persist data for successful runs |
| `executionTimeout` | seconds (`-1` = no limit) | Kill execution after N seconds |
| `errorWorkflow` | workflow ID string | Workflow to trigger on failure |
| `timezone` | IANA tz string | e.g. `"America/New_York"`, `"UTC"`, `"Europe/London"` |
| `executionOrder` | `"v1"` | Execution engine version. Use `"v1"`. |
| `callerPolicy` | `"workflowsFromSameOwner"` / `"any"` / `"workflowsFromAList"` | Who can call this as sub-workflow |
| `callerIds` | comma-separated IDs | Allowed caller workflow IDs (when `callerPolicy` = `"workflowsFromAList"`) |
| `timeSavedPerExecution` | number (minutes) | For ROI tracking in n8n |
| `availableInMCP` | `true` / `false` | Expose workflow as an MCP tool |

## Connections Object

Keys are node **names** (not IDs). Each output port is an array of target arrays.

```json
"connections": {
  "Node Name": {
    "main": [
      [
        { "node": "Next Node Name", "type": "main", "index": 0 }
      ]
    ]
  }
}
```

- `main` = the standard output port
- Outer array = output ports (index 0 = first output, index 1 = second, etc.)
- Inner array = all connections from that port (fan-out)
- Last node in the chain has **no entry** in connections

## Full Example

```json
{
  "nodes": [
    {
      "parameters": {},
      "id": "ManualTrigger",
      "name": "Manual Trigger",
      "type": "n8n-nodes-base.manualTrigger",
      "typeVersion": 1,
      "position": [240, 300]
    },
    {
      "parameters": {
        "url": "https://hn.algolia.com/api/v1/search_by_date?tags=story&hitsPerPage=10",
        "responseFormat": "json",
        "method": "GET"
      },
      "id": "HTTP Request",
      "name": "Get HN Latest",
      "type": "n8n-nodes-base.httpRequest",
      "typeVersion": 4,
      "position": [460, 300]
    },
    {
      "parameters": {
        "jsCode": "const hits = items[0].json.hits;\nconst top10 = hits.slice(0,10).map((item,i) => `${i+1}. ${item.title}\\n${item.url||'No URL'}\\n`).join('\\n');\nreturn [{ json: { text: top10 } }];"
      },
      "id": "Code",
      "name": "Format Top 10",
      "type": "n8n-nodes-base.code",
      "typeVersion": 2,
      "position": [680, 300]
    },
    {
      "parameters": {
        "fromEmail": "sender@example.com",
        "toEmail": "recipient@example.com",
        "subject": "Top 10 HN Stories",
        "text": "={{$json.text}}"
      },
      "id": "Email",
      "name": "Send Email",
      "type": "n8n-nodes-base.emailSend",
      "typeVersion": 2,
      "position": [900, 300]
    }
  ],
  "connections": {
    "Manual Trigger": {
      "main": [[{ "node": "Get HN Latest", "type": "main", "index": 0 }]]
    },
    "Get HN Latest": {
      "main": [[{ "node": "Format Top 10", "type": "main", "index": 0 }]]
    },
    "Format Top 10": {
      "main": [[{ "node": "Send Email", "type": "main", "index": 0 }]]
    }
  }
}
```

## Branching (IF node)

IF node has two outputs: index 0 = true branch, index 1 = false branch.

```json
"connections": {
  "IF": {
    "main": [
      [{ "node": "On True", "type": "main", "index": 0 }],
      [{ "node": "On False", "type": "main", "index": 0 }]
    ]
  }
}
```

## Fan-out (one node → multiple)

```json
"connections": {
  "Trigger": {
    "main": [
      [
        { "node": "Branch A", "type": "main", "index": 0 },
        { "node": "Branch B", "type": "main", "index": 0 }
      ]
    ]
  }
}
```
