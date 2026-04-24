# Core & Flow Nodes

## Triggers

| Display Name | Type | typeVersion | Notes |
|---|---|---|---|
| Manual Trigger | `n8n-nodes-base.manualTrigger` | 1 | Run by hand |
| Schedule Trigger | `n8n-nodes-base.scheduleTrigger` | 1 | Cron / interval |
| Webhook | `n8n-nodes-base.webhook` | 2 | Receive HTTP |
| Form Trigger | `n8n-nodes-base.formTrigger` | 2 | HTML form submissions |
| Chat Trigger | `n8n-nodes-base.chatTrigger` | 1 | AI chat interface |
| Error Trigger | `n8n-nodes-base.errorTrigger` | 1 | Catches failed workflow — use in error workflow |
| Execute Workflow Trigger | `n8n-nodes-base.executeWorkflowTrigger` | 1 | Called by another workflow |
| Email Trigger (IMAP) | `n8n-nodes-base.emailReadImap` | 2 | Trigger on new email |

## Flow Control

| Display Name | Type | typeVersion | Notes |
|---|---|---|---|
| IF | `n8n-nodes-base.if` | 2 | Two outputs: true / false |
| Switch | `n8n-nodes-base.switch` | 3 | Multiple output branches |
| Merge | `n8n-nodes-base.merge` | 3 | Combine branches |
| Split in Batches | `n8n-nodes-base.splitInBatches` | 3 | Loop over chunks |
| Wait | `n8n-nodes-base.wait` | 1 | Pause execution |
| Stop and Error | `n8n-nodes-base.stopAndError` | 1 | Halt + report error |
| No Operation | `n8n-nodes-base.noOp` | 1 | Pass-through placeholder |
| Execute Workflow | `n8n-nodes-base.executeWorkflow` | 1 | Call another workflow |
| Respond to Webhook | `n8n-nodes-base.respondToWebhook` | 1 | Send HTTP response |

## Data Transformation

| Display Name | Type | typeVersion | Notes |
|---|---|---|---|
| Code | `n8n-nodes-base.code` | 2 | JS or Python. Replaces deprecated `function`/`functionItem` |
| Set | `n8n-nodes-base.set` | 3 | Add / overwrite fields |
| Edit Fields (Set) | `n8n-nodes-base.set` | 3 | Same node, newer UI name |
| Filter | `n8n-nodes-base.filter` | 1 | Keep items matching condition |
| Remove Duplicates | `n8n-nodes-base.removeDuplicates` | 1 | Deduplicate by field |
| Sort | `n8n-nodes-base.sort` | 1 | Sort items by field |
| Limit | `n8n-nodes-base.limit` | 1 | Keep first N items |
| Aggregate | `n8n-nodes-base.aggregate` | 1 | Merge fields from all items into one |
| Split Out | `n8n-nodes-base.splitOut` | 1 | Array field → individual items |
| Summarize | `n8n-nodes-base.summarize` | 1 | Group by + aggregate (count/sum/avg) |
| Rename Keys | `n8n-nodes-base.renameKeys` | 1 | Rename item fields |
| Compare Datasets | `n8n-nodes-base.compareDatasets` | 3 | Diff two inputs |

## HTTP & Utilities

| Display Name | Type | typeVersion | Notes |
|---|---|---|---|
| HTTP Request | `n8n-nodes-base.httpRequest` | 4 | Call any REST API |
| GraphQL | `n8n-nodes-base.graphql` | 1 | GraphQL queries |
| HTML | `n8n-nodes-base.html` | 1 | Extract / generate HTML |
| Markdown | `n8n-nodes-base.markdown` | 1 | Convert markdown ↔ HTML |
| XML | `n8n-nodes-base.xml` | 1 | Parse / generate XML |
| JSON | `n8n-nodes-base.set` | 3 | Use Set node for JSON manipulation |
| Crypto | `n8n-nodes-base.crypto` | 1 | Hash / sign / encrypt |
| Date & Time | `n8n-nodes-base.dateTime` | 2 | Parse / format / add dates |
| JWT | `n8n-nodes-base.jwt` | 1 | Sign / verify JWT |
| Compression | `n8n-nodes-base.compression` | 1 | Zip / gzip |
| SSH | `n8n-nodes-base.ssh` | 1 | Run remote shell commands |
| FTP | `n8n-nodes-base.ftp` | 1 | FTP / SFTP file transfer |
| Read/Write File | `n8n-nodes-base.readWriteFile` | 1 | Local filesystem |

## Deprecated → Do Not Use

| Old Type | Use Instead |
|---|---|
| `n8n-nodes-base.function` | `n8n-nodes-base.code` v2 |
| `n8n-nodes-base.functionItem` | `n8n-nodes-base.code` v2 |
| `n8n-nodes-base.itemLists` | `n8n-nodes-base.splitOut` or `aggregate` |
| `n8n-nodes-base.moveBinaryData` | `n8n-nodes-base.extractFromFile` |
| `n8n-nodes-base.readBinaryFile` | `n8n-nodes-base.readWriteFile` |
| `n8n-nodes-base.writeBinaryFile` | `n8n-nodes-base.readWriteFile` |
