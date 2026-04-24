# Workflow Pattern Selection Guide

## The 6 Core Patterns

| Pattern | Use when | Trigger node |
|---------|----------|--------------|
| **Webhook** | Receiving external events, need instant response | Webhook |
| **HTTP API** | Fetching/syncing with external APIs | Schedule or Webhook |
| **Database** | Read/write/sync database records | Schedule |
| **AI Agent** | LLM reasoning with tools and memory | Webhook or Manual |
| **Scheduled** | Recurring reports, jobs, notifications | Schedule |
| **Batch** | Large datasets, API rate limits | Schedule |

## Pattern Skeletons

### Webhook Processing
```
Webhook → [Validate] → [Transform] → Respond to Webhook
                                   ↘ Notify (async)
```
→ Full template: [webhook_processing.md](webhook_processing.md)

### HTTP API Integration
```
Trigger → HTTP Request → IF (ok?) → Transform → Action
                                  ↘ Error Handler
```
→ Full template: [http_api_integration.md](http_api_integration.md)

### Database Operations
```
Schedule → DB Query → [Transform] → DB Write → [Verify]
```
→ Full template: [database_operations.md](database_operations.md)

### AI Agent
```
Trigger → AI Agent [Model + Tools + Memory] → Format → Respond
```
→ Full template: [ai_agent_workflow.md](ai_agent_workflow.md)

### Scheduled Task
```
Schedule → Fetch → Process → Deliver → Log
```
→ Full template: [scheduled_tasks.md](scheduled_tasks.md)

### Batch Processing
```
Schedule → Prepare → SplitInBatches → Process → Accumulate → Aggregate
```
Key: set batch size to stay under API rate limits. Use a counter node or $vars to track progress.

## Node JSON format

→ [WORKFLOW_FORMAT.md](WORKFLOW_FORMAT.md)
