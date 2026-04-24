# Workflow JSON Format

## Required Structure

```json
{
  "nodes": [...],
  "connections": {...}
}
```

Top-level `name` is optional — n8n auto-assigns one if omitted.

## Node Object

Every node requires these fields:

```json
{
  "id": "UniqueNodeId",
  "name": "Human Readable Name",
  "type": "n8n-nodes-base.nodeType",
  "typeVersion": 1,
  "position": [x, y],
  "parameters": {}
}
```

| Field | Notes |
|-------|-------|
| `id` | Unique within workflow. Can be any string. |
| `name` | Shown in canvas. Used as connection key. |
| `type` | Full node type: `n8n-nodes-base.<type>` |
| `typeVersion` | Check node docs. Usually `1`, some nodes at `2`–`4`. |
| `position` | `[x, y]` canvas coordinates. Increment x by ~220 per node. |
| `parameters` | Node-specific config. Can be `{}` for defaults. |

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
