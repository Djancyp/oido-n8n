# oido-n8n

## What

Go MCP server wrapping the n8n REST API v1. An Oido Studio plugin — packaged as a zip via `make dist`, uploaded through the Plugins UI.

## Build & dist

```bash
make build            # CGO_ENABLED=0 go build -o oido-n8n-mcp .
make dist             # build + zip into dist/ with oido-extension.json, OIDO.md, commands/, skills/
make clean            # rm binary + dist/
```

Output: `oido-n8n-mcp` (binary), `oido-n8n.zip` (plugin archive).

## Architecture

- `main.go` → `mcp_server.go` (tool registration + handlers)
- `n8n_client.go` — HTTP client to n8n's `/api/v1/*`, auth via `X-N8N-API-KEY` header
- `node_db.go` — embedded SQLite (`n8n-nodes.db` via `//go:embed`) for node type validation at create/validate time
- `commands/*.toml` — Oido command definitions (workflow-list, workflow-run, execution-status, workflow-status)
- `skills/oido-n8n/*.md` — skill files for agent skill routing

## Env vars

| Var | Default | Required |
|---|---|---|
| `N8N_API_URL` | `http://localhost:5678` | no |
| `N8N_API_KEY` | — | **yes** (server fails to start without it) |

## Registered MCP tools

Only 5 tools are wired up via `mcp.AddTool()` in `RunMCPServer()` — not the full handler set:

| Tool | Description |
|---|---|
| `n8n_search_nodes` | Search node types by keyword before building workflows |
| `n8n_create_workflow` | Create workflow from JSON (validates internally) |
| `n8n_update_workflow` | Update workflow by ID (full JSON) |
| `n8n_delete_workflow` | Delete workflow by ID |
| `n8n_get_workflow` | Get full workflow definition (nodes + connections) by ID |
| `n8n_validate_workflow` | Validate workflow JSON without creating |
| `n8n_list_workflows` | List workflows, filter by active/tags |
| `n8n_list_credentials` | List credentials (secrets excluded) |
| `n8n_get_credential` | Get credential by ID |
| `n8n_create_credential` | Create credential — use `n8n_get_credential_schema` first |
| `n8n_update_credential` | Update credential by ID |
| `n8n_delete_credential` | Delete credential by ID |
| `n8n_get_credential_schema` | Get required fields for a credential type |

All other handler functions (`HandleListExecutions`, `HandleCreateCredential`, etc.) are **not registered** — adding them requires a `mcp.AddTool()` call in `RunMCPServer()`.

## Non-obvious facts for agents

1. **Always `n8n_search_nodes` before creating/updating workflows.** The `n8n_create_workflow` and `n8n_validate_workflow` tools reject unknown node types by checking the embedded node DB. Don't guess node type names.
2. **Workflow connection keys use node `name`, not `id`.** This is a common gotcha when authoring workflow JSON.
3. **`doWebhook` vs `do`:** Webhooks hit `<N8N_API_URL>/webhook/<path>` directly (no `/api/v1` prefix, no API key). All other calls go through `<N8N_API_URL>/api/v1` with `X-N8N-API-KEY`.
4. **Plugin packaging:** `make dist` bundles `oido-extension.json`, `OIDO.md`, the binary, and contents of `commands/` and `skills/oido-n8n/`. If you add files to those dirs, rebuild dist.
5. **CGO_ENABLED=0 is required** (`go build` defaults to CGO enabled on some platforms, which breaks the pure-Go SQLite driver).
6. **Single `package main`** — no sub-packages, no tests currently. All Go source is flat in the root.
7. **The node DB is embedded at compile time** from `n8n-nodes.db`. It's written to a temp file at startup (read-only). To update it, replace the `.db` file and rebuild.
8. **Supported API resource groups:** workflows, executions, credentials, tags, webhooks, projects, variables, users, audit.
