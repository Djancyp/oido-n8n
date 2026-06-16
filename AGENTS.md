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

66 tools are wired via `mcp.AddTool()` in `RunMCPServer()`, grouped by resource:

- **Workflows:** `n8n_list_workflows`, `n8n_get_workflow`, `n8n_create_workflow`, `n8n_update_workflow`, `n8n_patch_workflow`, `n8n_delete_workflow`, `n8n_activate_workflow`, `n8n_deactivate_workflow`, `n8n_archive_workflow`, `n8n_unarchive_workflow`, `n8n_execute_workflow`, `n8n_validate_workflow`
- **Node knowledge:** `n8n_search_nodes`, `n8n_get_node_schema`
- **Executions:** `n8n_list_executions`, `n8n_get_execution`, `n8n_delete_execution`, `n8n_stop_execution`, `n8n_stop_executions`, `n8n_retry_execution`
- **Credentials:** `n8n_list_credentials`, `n8n_get_credential`, `n8n_create_credential`, `n8n_update_credential`, `n8n_delete_credential`, `n8n_get_credential_schema`, `n8n_test_credential`
- **Tags:** `n8n_list_tags`, `n8n_get_tag`, `n8n_create_tag`, `n8n_update_tag`, `n8n_delete_tag`, `n8n_list_execution_tags`, `n8n_update_execution_tags`
- **Data tables:** `n8n_list_data_tables`, `n8n_get_data_table`, `n8n_create_data_table`, `n8n_update_data_table`, `n8n_delete_data_table`, `n8n_list_columns`, `n8n_get_column`, `n8n_create_column`, `n8n_update_column`, `n8n_delete_column`, `n8n_list_rows`, `n8n_create_row`, `n8n_update_rows`, `n8n_upsert_rows`, `n8n_delete_rows`
- **Projects:** `n8n_list_projects`, `n8n_get_project`, `n8n_create_project`, `n8n_update_project`, `n8n_delete_project`
- **Variables:** `n8n_list_variables`, `n8n_get_variable`, `n8n_create_variable`, `n8n_update_variable`, `n8n_delete_variable`
- **Users:** `n8n_list_users`, `n8n_get_user`, `n8n_create_users`, `n8n_change_user_role`, `n8n_delete_user`
- **Audit / webhooks:** `n8n_generate_audit`, `n8n_trigger_webhook`

To add a tool: write `HandleX` + an `XArgs` struct, then add a `mcp.AddTool()` call in `RunMCPServer()`.

## Non-obvious facts for agents

1. **Authoring loop:** `n8n_search_nodes` → `n8n_get_node_schema` (per node) → write JSON → `n8n_validate_workflow` → `n8n_create_workflow`. The create/validate tools reject unknown node types via the embedded node DB.
2. **Workflow connection keys use node `name`, not `id`.** `n8n_validate_workflow` now enforces this — every connection source key and target `node` must match a defined node name, else `ERROR`.
   - **Node DB `properties` is minimal:** only `{name, type, default}` per param. No `required`, `options`, descriptions, or credential metadata. `n8n_get_node_schema` surfaces what's there; required fields / option values must come from n8n docs. Richer authoring guidance needs a regenerated `n8n-nodes.db` — that's the lever for closing the gap with full node-knowledge toolkits.
   - **`version` quality:** ~175/823 rows store `version` as `"[object Object]"`. `LookupNode` reports `VersionValid=false` for those, and the typeVersion check skips them.
3. **`doWebhook` vs `do`:** Webhooks hit `<N8N_API_URL>/webhook/<path>` directly (no `/api/v1` prefix, no API key). All other calls go through `<N8N_API_URL>/api/v1` with `X-N8N-API-KEY`.
4. **Plugin packaging:** `make dist` bundles `oido-extension.json`, `OIDO.md`, the binary, and contents of `commands/` and `skills/oido-n8n/`. If you add files to those dirs, rebuild dist.
5. **CGO_ENABLED=0 is required** (`go build` defaults to CGO enabled on some platforms, which breaks the pure-Go SQLite driver).
6. **Single `package main`** — no sub-packages, no tests currently. All Go source is flat in the root.
7. **The node DB is embedded at compile time** from `n8n-nodes.db`. It's written to a temp file at startup (read-only). To update it, replace the `.db` file and rebuild.
8. **Supported API resource groups:** workflows, executions, credentials, tags, webhooks, projects, variables, users, audit.
