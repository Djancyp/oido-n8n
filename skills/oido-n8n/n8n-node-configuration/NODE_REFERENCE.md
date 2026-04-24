# Node Type Reference Index

Always use the exact type string. Wrong type = node not found in n8n.

## Files by Group

| Group | File | Contents |
|-------|------|----------|
| Core & Flow | [NODE_TYPES_CORE.md](NODE_TYPES_CORE.md) | Triggers, IF/Switch, Merge, Code, Set, HTTP Request, deprecated list |
| Communication | [NODE_TYPES_COMMUNICATION.md](NODE_TYPES_COMMUNICATION.md) | Email, Slack, Discord, Telegram, Teams, SMS, social |
| Data & Storage | [NODE_TYPES_DATA_STORAGE.md](NODE_TYPES_DATA_STORAGE.md) | Postgres, MySQL, MongoDB, Redis, Google Sheets, S3, Notion, queues |
| Services | [NODE_TYPES_SERVICES.md](NODE_TYPES_SERVICES.md) | CRM, DevOps, Google/Microsoft 365, payments, marketing, e-commerce |
| AI & LangChain | [NODE_TYPES_AI.md](NODE_TYPES_AI.md) | Agents, LLMs, memory, tools, vector stores, embeddings |
| Common Mistakes | [COMMON_NODE_MISTAKES.md](COMMON_NODE_MISTAKES.md) | Hallucinated types, deprecated nodes, correct typeVersions |

## Quick Lookup — Most Used Types

```
n8n-nodes-base.manualTrigger          typeVersion: 1
n8n-nodes-base.scheduleTrigger        typeVersion: 1
n8n-nodes-base.webhook                typeVersion: 2
n8n-nodes-base.httpRequest            typeVersion: 4
n8n-nodes-base.code                   typeVersion: 2
n8n-nodes-base.set                    typeVersion: 3
n8n-nodes-base.if                     typeVersion: 2
n8n-nodes-base.switch                 typeVersion: 3
n8n-nodes-base.merge                  typeVersion: 3
n8n-nodes-base.splitInBatches         typeVersion: 3
n8n-nodes-base.emailSend              typeVersion: 2
n8n-nodes-base.gmail                  typeVersion: 2
n8n-nodes-base.slack                  typeVersion: 2
n8n-nodes-base.googleSheets           typeVersion: 4
n8n-nodes-base.postgres               typeVersion: 2
n8n-nodes-base.respondToWebhook       typeVersion: 1
n8n-nodes-base.stopAndError           typeVersion: 1
n8n-nodes-base.executeWorkflow        typeVersion: 1
n8n-nodes-base.filter                 typeVersion: 1
n8n-nodes-base.aggregate              typeVersion: 1

@n8n/n8n-nodes-langchain.agent                         typeVersion: 1
@n8n/n8n-nodes-langchain.lmChatOpenAi                  typeVersion: 1
@n8n/n8n-nodes-langchain.lmChatAnthropic               typeVersion: 1
@n8n/n8n-nodes-langchain.memoryBufferWindow             typeVersion: 1
@n8n/n8n-nodes-langchain.toolHttpRequest                typeVersion: 1
@n8n/n8n-nodes-langchain.vectorStorePinecone            typeVersion: 1
```

## Rules

1. If no native node exists for a service → use `n8n-nodes-base.httpRequest` v4
2. AI sub-nodes (model, memory, tools) connect with `ai_languageModel` / `ai_memory` / `ai_tool` — not `main`
3. Never invent a node type. If unsure, search this index first.
