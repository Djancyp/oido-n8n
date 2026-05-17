package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type MCPHandler struct {
	client *N8nClient
	nodeDB *sql.DB
}

func NewMCPHandler(client *N8nClient, nodeDB *sql.DB) *MCPHandler {
	return &MCPHandler{client: client, nodeDB: nodeDB}
}

// --- Arg types ---

type ListWorkflowsArgs struct {
	Active *bool    `json:"active,omitempty" jsonschema:"Filter by active status"`
	Tags   []string `json:"tags,omitempty"   jsonschema:"Filter by tag names"`
	Name   string   `json:"name,omitempty"   jsonschema:"Filter by workflow name (partial match)"`
	Limit  int      `json:"limit,omitempty"  jsonschema:"Max results (default: 20)"`
	Cursor string   `json:"cursor,omitempty" jsonschema:"Pagination cursor"`
}

type WorkflowIDArgs struct {
	ID string `json:"id" jsonschema:"Workflow ID"`
}

type CreateWorkflowArgs struct {
	WorkflowJSON string `json:"workflow_json" jsonschema:"Full workflow definition as JSON string"`
}

type UpdateWorkflowArgs struct {
	ID           string `json:"id"            jsonschema:"Workflow ID"`
	WorkflowJSON string `json:"workflow_json" jsonschema:"Partial or full workflow update as JSON string"`
}

type ListExecutionsArgs struct {
	WorkflowID  string `json:"workflow_id,omitempty"  jsonschema:"Filter by workflow ID"`
	Status      string `json:"status,omitempty"       jsonschema:"Filter: waiting, running, success, error, canceled"`
	Limit       int    `json:"limit,omitempty"        jsonschema:"Max results (default: 20)"`
	Cursor      string `json:"cursor,omitempty"       jsonschema:"Pagination cursor"`
	IncludeData bool   `json:"include_data,omitempty" jsonschema:"Include node input/output data in results"`
}

type ExecutionIDArgs struct {
	ID string `json:"id" jsonschema:"Execution ID"`
}

type GetExecutionArgs struct {
	ID          string `json:"id"                    jsonschema:"Execution ID"`
	IncludeData bool   `json:"include_data,omitempty" jsonschema:"Include node input/output data"`
}

type StopExecutionsArgs struct {
	WorkflowID string `json:"workflow_id,omitempty" jsonschema:"Stop all running executions for this workflow ID. Omit to stop all running executions."`
}

type RetryExecutionArgs struct {
	ID string `json:"id" jsonschema:"Execution ID of the failed execution to retry"`
}

type ExecutionTagsArgs struct {
	ID string `json:"id" jsonschema:"Execution ID"`
}

type UpdateExecutionTagsArgs struct {
	ID     string   `json:"id"      jsonschema:"Execution ID"`
	TagIDs []string `json:"tag_ids" jsonschema:"List of tag IDs to set on the execution (replaces existing tags)"`
}

type ListCredentialsArgs struct {
	Limit        int    `json:"limit,omitempty"         jsonschema:"Max results (default: 20)"`
	CredentialID string `json:"credential_id,omitempty" jsonschema:"Filter by a specific credential ID"`
	Type         string `json:"type,omitempty"          jsonschema:"Filter by credential type (e.g. n8n-nodes-base.aws, githubApi)"`
}

type GetCredentialArgs struct {
	ID string `json:"id" jsonschema:"Credential ID"`
}

type UpdateCredentialArgs struct {
	ID       string `json:"id"             jsonschema:"Credential ID"`
	Name     string `json:"name"           jsonschema:"Credential name"`
	Type     string `json:"type"           jsonschema:"Credential type (e.g. httpBasicAuth, githubApi)"`
	DataJSON string `json:"data,omitempty" jsonschema:"Updated credential fields as JSON string"`
}

type CreateCredentialArgs struct {
	Name     string `json:"name"           jsonschema:"Credential name"`
	Type     string `json:"type"           jsonschema:"Credential type (e.g. httpBasicAuth, gitlabApi)"`
	DataJSON string `json:"data,omitempty" jsonschema:"Credential fields as JSON string"`
}

type CredentialIDArgs struct {
	ID string `json:"id" jsonschema:"Credential ID"`
}

type GetCredentialSchemaArgs struct {
	CredentialType string `json:"credential_type" jsonschema:"Credential type name (e.g. httpBasicAuth)"`
}

type ListTagsArgs struct {
	Limit  int    `json:"limit,omitempty"  jsonschema:"Max results (default: 100)"`
	Cursor string `json:"cursor,omitempty" jsonschema:"Pagination cursor"`
}

type CreateTagArgs struct {
	Name string `json:"name" jsonschema:"Tag name"`
}

type TagIDArgs struct {
	ID string `json:"id" jsonschema:"Tag ID"`
}

type UpdateTagArgs struct {
	ID   string `json:"id"   jsonschema:"Tag ID"`
	Name string `json:"name" jsonschema:"New tag name"`
}

type TriggerWebhookArgs struct {
	Path     string `json:"path"            jsonschema:"Webhook path (without /webhook/ prefix)"`
	Method   string `json:"method,omitempty" jsonschema:"HTTP method: GET, POST, PUT, PATCH, DELETE (default: POST)"`
	BodyJSON string `json:"body,omitempty"  jsonschema:"Request body as JSON string"`
}

type UpdatePartialWorkflowArgs struct {
	ID              string `json:"id"                    jsonschema:"Workflow ID"`
	Name            string `json:"name,omitempty"        jsonschema:"New workflow name"`
	NodesJSON       string `json:"nodes,omitempty"       jsonschema:"Replacement nodes array as JSON string"`
	ConnectionsJSON string `json:"connections,omitempty" jsonschema:"Replacement connections map as JSON string"`
	SettingsJSON    string `json:"settings,omitempty"   jsonschema:"Replacement settings object as JSON string"`
}

type ValidateWorkflowArgs struct {
	WorkflowJSON string `json:"workflow_json" jsonschema:"Workflow definition as JSON string"`
}

type ListUsersArgs struct {
	Limit       int    `json:"limit,omitempty"       jsonschema:"Max results (default: 20)"`
	Cursor      string `json:"cursor,omitempty"      jsonschema:"Pagination cursor"`
	IncludeRole bool   `json:"include_role,omitempty" jsonschema:"Include role field in response"`
}

type GetUserArgs struct {
	ID          string `json:"id"                    jsonschema:"User ID or email address"`
	IncludeRole bool   `json:"include_role,omitempty" jsonschema:"Include role field in response"`
}

type CreateUsersArgs struct {
	UsersJSON string `json:"users_json" jsonschema:"JSON array of user objects with email (required), role, firstName, lastName"`
}

type ChangeUserRoleArgs struct {
	ID   string `json:"id"   jsonschema:"User ID"`
	Role string `json:"role" jsonschema:"New global role (e.g. global:admin, global:member)"`
}

type DeleteUserArgs struct {
	ID string `json:"id" jsonschema:"User ID"`
}

type ListProjectsArgs struct {
	Limit  int    `json:"limit,omitempty"  jsonschema:"Max results (default: 20)"`
	Cursor string `json:"cursor,omitempty" jsonschema:"Pagination cursor"`
}

type ProjectIDArgs struct {
	ID string `json:"id" jsonschema:"Project ID"`
}

type CreateProjectArgs struct {
	Name string `json:"name"           jsonschema:"Project name"`
	Type string `json:"type,omitempty" jsonschema:"Project type: team (default) or enterprise"`
}

type UpdateProjectArgs struct {
	ID   string `json:"id"   jsonschema:"Project ID"`
	Name string `json:"name" jsonschema:"New project name"`
}

type ListVariablesArgs struct {
	Limit  int    `json:"limit,omitempty"  jsonschema:"Max results (default: 20)"`
	Cursor string `json:"cursor,omitempty" jsonschema:"Pagination cursor"`
}

type VariableIDArgs struct {
	ID string `json:"id" jsonschema:"Variable ID"`
}

type CreateVariableArgs struct {
	Key   string `json:"key"   jsonschema:"Variable key name"`
	Value string `json:"value" jsonschema:"Variable value"`
}

type UpdateVariableArgs struct {
	ID    string `json:"id"             jsonschema:"Variable ID"`
	Key   string `json:"key,omitempty"   jsonschema:"New key name"`
	Value string `json:"value,omitempty" jsonschema:"New value"`
}

type GenerateAuditArgs struct {
	OptionsJSON string `json:"options,omitempty" jsonschema:"Optional audit configuration as JSON (e.g. {\"categories\":[\"credentials\",\"workflows\"]})"`
}

type SearchNodesArgs struct {
	Keywords string `json:"keywords"       jsonschema:"Comma-separated search terms matched against node name and display name (partial match, OR logic)"`
	Group    string `json:"group,omitempty" jsonschema:"Filter by group: t=trigger, i=action, o=output"`
	Limit    int    `json:"limit,omitempty" jsonschema:"Max results (default: 20)"`
}

type GetNodeSchemaArgs struct {
	Name string `json:"name" jsonschema:"Exact node type name e.g. n8n-nodes-base.httpRequest"`
}

type ExecuteWorkflowArgs struct {
	WorkflowID  string `json:"workflow_id"        jsonschema:"Workflow ID to execute"`
	RunDataJSON string `json:"run_data,omitempty" jsonschema:"Optional input data as JSON object"`
}

type ArchiveWorkflowArgs struct {
	ID string `json:"id" jsonschema:"Workflow ID"`
}

type TestCredentialArgs struct {
	ID string `json:"id" jsonschema:"Credential ID to test"`
}

type ListDataTablesArgs struct {
	Limit  int    `json:"limit,omitempty"  jsonschema:"Max results (default: 20)"`
	Cursor string `json:"cursor,omitempty" jsonschema:"Pagination cursor"`
}

type DataTableIDArgs struct {
	ID string `json:"id" jsonschema:"Data table ID"`
}

type CreateDataTableArgs struct {
	Name        string `json:"name"                   jsonschema:"Table name"`
	ColumnsJSON string `json:"columns_json,omitempty" jsonschema:"Optional JSON array of column objects with name and dataType fields"`
}

type UpdateDataTableArgs struct {
	ID   string `json:"id"   jsonschema:"Data table ID"`
	Name string `json:"name" jsonschema:"New table name"`
}

type ListRowsArgs struct {
	TableID string `json:"table_id"         jsonschema:"Data table ID"`
	Limit   int    `json:"limit,omitempty"  jsonschema:"Max results (default: 20)"`
	Cursor  string `json:"cursor,omitempty" jsonschema:"Pagination cursor"`
}

type CreateRowArgs struct {
	TableID     string `json:"table_id"  jsonschema:"Data table ID"`
	RowDataJSON string `json:"row_data"  jsonschema:"Row data as JSON object with column names as keys"`
}

type UpdateRowsArgs struct {
	TableID    string `json:"table_id"           jsonschema:"Data table ID"`
	FilterJSON string `json:"filter,omitempty"   jsonschema:"Filter to select rows to update as JSON object"`
	UpdateJSON string `json:"update"             jsonschema:"Fields to update as JSON object"`
}

type UpsertRowsArgs struct {
	TableID  string `json:"table_id"  jsonschema:"Data table ID"`
	RowsJSON string `json:"rows"      jsonschema:"Rows to upsert as JSON array of objects"`
}

type DeleteRowsArgs struct {
	TableID    string `json:"table_id"         jsonschema:"Data table ID"`
	FilterJSON string `json:"filter,omitempty" jsonschema:"Filter to select rows to delete as JSON object. Omit to delete all rows."`
}

type ListColumnsArgs struct {
	TableID string `json:"table_id" jsonschema:"Data table ID"`
}

type CreateColumnArgs struct {
	TableID  string `json:"table_id"           jsonschema:"Data table ID"`
	Name     string `json:"name"               jsonschema:"Column name"`
	DataType string `json:"data_type,omitempty" jsonschema:"Column data type (e.g. string, number, boolean, date)"`
}

type ColumnIDArgs struct {
	TableID  string `json:"table_id"  jsonschema:"Data table ID"`
	ColumnID string `json:"column_id" jsonschema:"Column ID"`
}

type UpdateColumnArgs struct {
	TableID  string `json:"table_id"  jsonschema:"Data table ID"`
	ColumnID string `json:"column_id" jsonschema:"Column ID"`
	Name     string `json:"name"      jsonschema:"New column name"`
}

// --- Helpers ---

func errResult(msg string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: "Error: " + msg}},
		IsError: true,
	}
}

func textResult(text string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: text}},
	}
}

func jsonResult(v interface{}) *mcp.CallToolResult {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return errResult(err.Error())
	}
	return textResult(string(data))
}

// --- Workflow handlers ---

func (h *MCPHandler) HandleListWorkflows(_ context.Context, _ *mcp.CallToolRequest, args ListWorkflowsArgs) (*mcp.CallToolResult, any, error) {
	result, err := h.client.ListWorkflows(args.Active, args.Tags, args.Name, args.Limit, args.Cursor)
	if err != nil {
		return errResult(err.Error()), nil, nil
	}

	if len(result.Data) == 0 {
		return textResult("No workflows found."), nil, nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Workflows (%d):\n\n", len(result.Data)))
	sb.WriteString(fmt.Sprintf("%-24s  %-6s  %-20s  %s\n", "ID", "Active", "Updated", "Name"))
	sb.WriteString(strings.Repeat("-", 80) + "\n")
	for _, w := range result.Data {
		active := "false"
		if w.Active {
			active = "true"
		}
		updated := truncate(w.UpdatedAt, 20)
		sb.WriteString(fmt.Sprintf("%-24s  %-6s  %-20s  %s\n", w.ID, active, updated, w.Name))
	}
	if result.NextCursor != "" {
		sb.WriteString(fmt.Sprintf("\nNext cursor: %s", result.NextCursor))
	}
	return textResult(sb.String()), nil, nil
}

func (h *MCPHandler) HandleGetWorkflow(_ context.Context, _ *mcp.CallToolRequest, args WorkflowIDArgs) (*mcp.CallToolResult, any, error) {
	if args.ID == "" {
		return errResult("id is required"), nil, nil
	}
	w, err := h.client.GetWorkflow(args.ID)
	if err != nil {
		return errResult(err.Error()), nil, nil
	}
	return jsonResult(w), nil, nil
}

func (h *MCPHandler) HandleCreateWorkflow(_ context.Context, _ *mcp.CallToolRequest, args CreateWorkflowArgs) (*mcp.CallToolResult, any, error) {
	if args.WorkflowJSON == "" {
		return errResult("workflow_json is required"), nil, nil
	}

	// Pre-flight: parse + validate before hitting the API
	var wf map[string]json.RawMessage
	if err := json.Unmarshal([]byte(args.WorkflowJSON), &wf); err != nil {
		return textResult("ERROR: invalid JSON — " + err.Error() + "\nFix the JSON then call n8n_create_workflow again."), nil, nil
	}
	if _, ok := wf["nodes"]; !ok {
		return textResult("ERROR: missing required field 'nodes'"), nil, nil
	}
	if _, ok := wf["connections"]; !ok {
		return textResult("ERROR: missing required field 'connections'"), nil, nil
	}
	var preErrs []string
	if nodesRaw, ok := wf["nodes"]; ok {
		var nodes []map[string]json.RawMessage
		if err := json.Unmarshal(nodesRaw, &nodes); err == nil {
			seenIDs := map[string]bool{}
			for i, node := range nodes {
				prefix := fmt.Sprintf("nodes[%d]", i)
				for _, req := range []string{"id", "name", "type", "typeVersion", "position"} {
					if _, ok := node[req]; !ok {
						preErrs = append(preErrs, fmt.Sprintf("ERROR: %s missing required field '%s'", prefix, req))
					}
				}
				if idRaw, ok := node["id"]; ok {
					var id string
					if json.Unmarshal(idRaw, &id) == nil {
						if seenIDs[id] {
							preErrs = append(preErrs, fmt.Sprintf("ERROR: duplicate node id %q", id))
						}
						seenIDs[id] = true
					}
				}
				if h.nodeDB != nil {
					if typeRaw, ok := node["type"]; ok {
						var nodeType string
						if json.Unmarshal(typeRaw, &nodeType) == nil && nodeType != "" {
							var count int
							_ = h.nodeDB.QueryRow(`SELECT COUNT(*) FROM node_types WHERE name = ?`, nodeType).Scan(&count)
							if count == 0 {
								parts := strings.Split(nodeType, ".")
								keyword := parts[len(parts)-1]
								preErrs = append(preErrs, fmt.Sprintf(
									"ERROR: %s type %q not found — run n8n_search_nodes keyword=%q to find the correct type",
									prefix, nodeType, keyword,
								))
							}
						}
					}
				}
			}
		}
	}
	if len(preErrs) > 0 {
		return textResult("Workflow has errors — fix before creating:\n\n" + strings.Join(preErrs, "\n")), nil, nil
	}

	w, err := h.client.CreateWorkflow(args.WorkflowJSON)
	if err != nil {
		return errResult(err.Error()), nil, nil
	}
	return textResult(fmt.Sprintf("Workflow created: id=%s name=%q", w.ID, w.Name)), nil, nil
}

func (h *MCPHandler) HandleUpdateWorkflow(_ context.Context, _ *mcp.CallToolRequest, args UpdateWorkflowArgs) (*mcp.CallToolResult, any, error) {
	if args.ID == "" {
		return errResult("id is required"), nil, nil
	}
	if args.WorkflowJSON == "" {
		return errResult("workflow_json is required"), nil, nil
	}
	w, err := h.client.UpdateWorkflow(args.ID, args.WorkflowJSON)
	if err != nil {
		return errResult(err.Error()), nil, nil
	}
	return textResult(fmt.Sprintf("Workflow updated: id=%s name=%q active=%v", w.ID, w.Name, w.Active)), nil, nil
}

func (h *MCPHandler) HandleDeleteWorkflow(_ context.Context, _ *mcp.CallToolRequest, args WorkflowIDArgs) (*mcp.CallToolResult, any, error) {
	if args.ID == "" {
		return errResult("id is required"), nil, nil
	}
	if err := h.client.DeleteWorkflow(args.ID); err != nil {
		return errResult(err.Error()), nil, nil
	}
	return textResult(fmt.Sprintf("Workflow %s deleted.", args.ID)), nil, nil
}

func (h *MCPHandler) HandleActivateWorkflow(_ context.Context, _ *mcp.CallToolRequest, args WorkflowIDArgs) (*mcp.CallToolResult, any, error) {
	if args.ID == "" {
		return errResult("id is required"), nil, nil
	}
	w, err := h.client.ActivateWorkflow(args.ID)
	if err != nil {
		return errResult(err.Error()), nil, nil
	}
	return textResult(fmt.Sprintf("Workflow %q (id=%s) activated.", w.Name, w.ID)), nil, nil
}

func (h *MCPHandler) HandleDeactivateWorkflow(_ context.Context, _ *mcp.CallToolRequest, args WorkflowIDArgs) (*mcp.CallToolResult, any, error) {
	if args.ID == "" {
		return errResult("id is required"), nil, nil
	}
	w, err := h.client.DeactivateWorkflow(args.ID)
	if err != nil {
		return errResult(err.Error()), nil, nil
	}
	return textResult(fmt.Sprintf("Workflow %q (id=%s) deactivated.", w.Name, w.ID)), nil, nil
}

// --- Execution handlers ---

func (h *MCPHandler) HandleListExecutions(_ context.Context, _ *mcp.CallToolRequest, args ListExecutionsArgs) (*mcp.CallToolResult, any, error) {
	result, err := h.client.ListExecutions(args.WorkflowID, args.Status, args.Limit, args.Cursor, args.IncludeData)
	if err != nil {
		return errResult(err.Error()), nil, nil
	}

	if len(result.Data) == 0 {
		return textResult("No executions found."), nil, nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Executions (%d):\n\n", len(result.Data)))
	sb.WriteString(fmt.Sprintf("%-20s  %-12s  %-10s  %-22s  %s\n", "ID", "WorkflowID", "Status", "StartedAt", "Finished"))
	sb.WriteString(strings.Repeat("-", 90) + "\n")
	for _, e := range result.Data {
		sb.WriteString(fmt.Sprintf("%-20s  %-12s  %-10s  %-22s  %v\n",
			e.ID, truncate(e.WorkflowID, 12), e.Status, truncate(e.StartedAt, 22), e.Finished))
	}
	if result.NextCursor != "" {
		sb.WriteString(fmt.Sprintf("\nNext cursor: %s", result.NextCursor))
	}
	return textResult(sb.String()), nil, nil
}

func (h *MCPHandler) HandleGetExecution(_ context.Context, _ *mcp.CallToolRequest, args GetExecutionArgs) (*mcp.CallToolResult, any, error) {
	if args.ID == "" {
		return errResult("id is required"), nil, nil
	}
	e, err := h.client.GetExecution(args.ID, args.IncludeData)
	if err != nil {
		return errResult(err.Error()), nil, nil
	}
	return jsonResult(e), nil, nil
}

func (h *MCPHandler) HandleDeleteExecution(_ context.Context, _ *mcp.CallToolRequest, args ExecutionIDArgs) (*mcp.CallToolResult, any, error) {
	if args.ID == "" {
		return errResult("id is required"), nil, nil
	}
	if err := h.client.DeleteExecution(args.ID); err != nil {
		return errResult(err.Error()), nil, nil
	}
	return textResult(fmt.Sprintf("Execution %s deleted.", args.ID)), nil, nil
}

func (h *MCPHandler) HandleStopExecution(_ context.Context, _ *mcp.CallToolRequest, args ExecutionIDArgs) (*mcp.CallToolResult, any, error) {
	if args.ID == "" {
		return errResult("id is required"), nil, nil
	}
	if err := h.client.StopExecution(args.ID); err != nil {
		return errResult(err.Error()), nil, nil
	}
	return textResult(fmt.Sprintf("Execution %s stopped.", args.ID)), nil, nil
}

func (h *MCPHandler) HandleStopExecutions(_ context.Context, _ *mcp.CallToolRequest, args StopExecutionsArgs) (*mcp.CallToolResult, any, error) {
	result, err := h.client.StopExecutions(args.WorkflowID)
	if err != nil {
		return errResult(err.Error()), nil, nil
	}
	pretty, _ := json.MarshalIndent(json.RawMessage(result), "", "  ")
	return textResult(string(pretty)), nil, nil
}

func (h *MCPHandler) HandleRetryExecution(_ context.Context, _ *mcp.CallToolRequest, args RetryExecutionArgs) (*mcp.CallToolResult, any, error) {
	if args.ID == "" {
		return errResult("id is required"), nil, nil
	}
	result, err := h.client.RetryExecution(args.ID)
	if err != nil {
		return errResult(err.Error()), nil, nil
	}
	pretty, _ := json.MarshalIndent(json.RawMessage(result), "", "  ")
	return textResult(string(pretty)), nil, nil
}

// --- Credential handlers ---

func (h *MCPHandler) HandleListCredentials(_ context.Context, _ *mcp.CallToolRequest, args ListCredentialsArgs) (*mcp.CallToolResult, any, error) {
	result, err := h.client.ListCredentials(args.Limit, args.CredentialID, args.Type)
	if err != nil {
		return errResult(err.Error()), nil, nil
	}

	if len(result.Data) == 0 {
		return textResult("No credentials found."), nil, nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Credentials (%d):\n\n", len(result.Data)))
	sb.WriteString(fmt.Sprintf("%-24s  %-24s  %s\n", "ID", "Type", "Name"))
	sb.WriteString(strings.Repeat("-", 70) + "\n")
	for _, cred := range result.Data {
		sb.WriteString(fmt.Sprintf("%-24s  %-24s  %s\n", cred.ID, cred.Type, cred.Name))
	}
	return textResult(sb.String()), nil, nil
}

func (h *MCPHandler) HandleCreateCredential(_ context.Context, _ *mcp.CallToolRequest, args CreateCredentialArgs) (*mcp.CallToolResult, any, error) {
	if args.Name == "" {
		return errResult("name is required"), nil, nil
	}
	if args.Type == "" {
		return errResult("type is required"), nil, nil
	}
	cred, err := h.client.CreateCredential(args.Name, args.Type, args.DataJSON)
	if err != nil {
		return errResult(err.Error()), nil, nil
	}
	return textResult(fmt.Sprintf("Credential created: id=%s name=%q type=%s", cred.ID, cred.Name, cred.Type)), nil, nil
}

func (h *MCPHandler) HandleDeleteCredential(_ context.Context, _ *mcp.CallToolRequest, args CredentialIDArgs) (*mcp.CallToolResult, any, error) {
	if args.ID == "" {
		return errResult("id is required"), nil, nil
	}
	if err := h.client.DeleteCredential(args.ID); err != nil {
		return errResult(err.Error()), nil, nil
	}
	return textResult(fmt.Sprintf("Credential %s deleted.", args.ID)), nil, nil
}

func (h *MCPHandler) HandleGetCredentialSchema(_ context.Context, _ *mcp.CallToolRequest, args GetCredentialSchemaArgs) (*mcp.CallToolResult, any, error) {
	if args.CredentialType == "" {
		return errResult("credential_type is required"), nil, nil
	}
	schema, err := h.client.GetCredentialSchema(args.CredentialType)
	if err != nil {
		return errResult(err.Error()), nil, nil
	}

	pretty, _ := json.MarshalIndent(schema, "", "  ")
	return textResult(string(pretty)), nil, nil
}

// --- Project handlers ---

func (h *MCPHandler) HandleListProjects(_ context.Context, _ *mcp.CallToolRequest, args ListProjectsArgs) (*mcp.CallToolResult, any, error) {
	result, err := h.client.ListProjects(args.Limit, args.Cursor)
	if err != nil {
		return errResult(err.Error()), nil, nil
	}
	if len(result.Data) == 0 {
		return textResult("No projects found."), nil, nil
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Projects (%d):\n\n", len(result.Data)))
	sb.WriteString(fmt.Sprintf("%-24s  %-12s  %s\n", "ID", "Type", "Name"))
	sb.WriteString(strings.Repeat("-", 60) + "\n")
	for _, p := range result.Data {
		sb.WriteString(fmt.Sprintf("%-24s  %-12s  %s\n", p.ID, p.Type, p.Name))
	}
	if result.NextCursor != "" {
		sb.WriteString(fmt.Sprintf("\nNext cursor: %s", result.NextCursor))
	}
	return textResult(sb.String()), nil, nil
}

func (h *MCPHandler) HandleGetProject(_ context.Context, _ *mcp.CallToolRequest, args ProjectIDArgs) (*mcp.CallToolResult, any, error) {
	if args.ID == "" {
		return errResult("id is required"), nil, nil
	}
	p, err := h.client.GetProject(args.ID)
	if err != nil {
		return errResult(err.Error()), nil, nil
	}
	return jsonResult(p), nil, nil
}

func (h *MCPHandler) HandleCreateProject(_ context.Context, _ *mcp.CallToolRequest, args CreateProjectArgs) (*mcp.CallToolResult, any, error) {
	if args.Name == "" {
		return errResult("name is required"), nil, nil
	}
	if args.Type == "" {
		args.Type = "team"
	}
	p, err := h.client.CreateProject(args.Name, args.Type)
	if err != nil {
		return errResult(err.Error()), nil, nil
	}
	return textResult(fmt.Sprintf("Project created: id=%s name=%q type=%s", p.ID, p.Name, p.Type)), nil, nil
}

func (h *MCPHandler) HandleUpdateProject(_ context.Context, _ *mcp.CallToolRequest, args UpdateProjectArgs) (*mcp.CallToolResult, any, error) {
	if args.ID == "" {
		return errResult("id is required"), nil, nil
	}
	if args.Name == "" {
		return errResult("name is required"), nil, nil
	}
	p, err := h.client.UpdateProject(args.ID, args.Name)
	if err != nil {
		return errResult(err.Error()), nil, nil
	}
	return textResult(fmt.Sprintf("Project updated: id=%s name=%q", p.ID, p.Name)), nil, nil
}

func (h *MCPHandler) HandleDeleteProject(_ context.Context, _ *mcp.CallToolRequest, args ProjectIDArgs) (*mcp.CallToolResult, any, error) {
	if args.ID == "" {
		return errResult("id is required"), nil, nil
	}
	if err := h.client.DeleteProject(args.ID); err != nil {
		return errResult(err.Error()), nil, nil
	}
	return textResult(fmt.Sprintf("Project %s deleted.", args.ID)), nil, nil
}

// --- Variable handlers ---

func (h *MCPHandler) HandleListVariables(_ context.Context, _ *mcp.CallToolRequest, args ListVariablesArgs) (*mcp.CallToolResult, any, error) {
	result, err := h.client.ListVariables(args.Limit, args.Cursor)
	if err != nil {
		return errResult(err.Error()), nil, nil
	}
	if len(result.Data) == 0 {
		return textResult("No variables found."), nil, nil
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Variables (%d):\n\n", len(result.Data)))
	sb.WriteString(fmt.Sprintf("%-24s  %-24s  %s\n", "ID", "Key", "Value"))
	sb.WriteString(strings.Repeat("-", 72) + "\n")
	for _, v := range result.Data {
		sb.WriteString(fmt.Sprintf("%-24s  %-24s  %s\n", v.ID, v.Key, v.Value))
	}
	if result.NextCursor != "" {
		sb.WriteString(fmt.Sprintf("\nNext cursor: %s", result.NextCursor))
	}
	return textResult(sb.String()), nil, nil
}

func (h *MCPHandler) HandleGetVariable(_ context.Context, _ *mcp.CallToolRequest, args VariableIDArgs) (*mcp.CallToolResult, any, error) {
	if args.ID == "" {
		return errResult("id is required"), nil, nil
	}
	v, err := h.client.GetVariable(args.ID)
	if err != nil {
		return errResult(err.Error()), nil, nil
	}
	return textResult(fmt.Sprintf("id=%s key=%q value=%q", v.ID, v.Key, v.Value)), nil, nil
}

func (h *MCPHandler) HandleCreateVariable(_ context.Context, _ *mcp.CallToolRequest, args CreateVariableArgs) (*mcp.CallToolResult, any, error) {
	if args.Key == "" {
		return errResult("key is required"), nil, nil
	}
	v, err := h.client.CreateVariable(args.Key, args.Value)
	if err != nil {
		return errResult(err.Error()), nil, nil
	}
	return textResult(fmt.Sprintf("Variable created: id=%s key=%q", v.ID, v.Key)), nil, nil
}

func (h *MCPHandler) HandleUpdateVariable(_ context.Context, _ *mcp.CallToolRequest, args UpdateVariableArgs) (*mcp.CallToolResult, any, error) {
	if args.ID == "" {
		return errResult("id is required"), nil, nil
	}
	if args.Key == "" && args.Value == "" {
		return errResult("at least one of key or value is required"), nil, nil
	}
	v, err := h.client.UpdateVariable(args.ID, args.Key, args.Value)
	if err != nil {
		return errResult(err.Error()), nil, nil
	}
	return textResult(fmt.Sprintf("Variable updated: id=%s key=%q value=%q", v.ID, v.Key, v.Value)), nil, nil
}

func (h *MCPHandler) HandleDeleteVariable(_ context.Context, _ *mcp.CallToolRequest, args VariableIDArgs) (*mcp.CallToolResult, any, error) {
	if args.ID == "" {
		return errResult("id is required"), nil, nil
	}
	if err := h.client.DeleteVariable(args.ID); err != nil {
		return errResult(err.Error()), nil, nil
	}
	return textResult(fmt.Sprintf("Variable %s deleted.", args.ID)), nil, nil
}

// --- Tag get/update/delete handlers ---

func (h *MCPHandler) HandleGetTag(_ context.Context, _ *mcp.CallToolRequest, args TagIDArgs) (*mcp.CallToolResult, any, error) {
	if args.ID == "" {
		return errResult("id is required"), nil, nil
	}
	tag, err := h.client.GetTag(args.ID)
	if err != nil {
		return errResult(err.Error()), nil, nil
	}
	return textResult(fmt.Sprintf("id=%s name=%q", tag.ID, tag.Name)), nil, nil
}

func (h *MCPHandler) HandleUpdateTag(_ context.Context, _ *mcp.CallToolRequest, args UpdateTagArgs) (*mcp.CallToolResult, any, error) {
	if args.ID == "" {
		return errResult("id is required"), nil, nil
	}
	if args.Name == "" {
		return errResult("name is required"), nil, nil
	}
	tag, err := h.client.UpdateTag(args.ID, args.Name)
	if err != nil {
		return errResult(err.Error()), nil, nil
	}
	return textResult(fmt.Sprintf("Tag updated: id=%s name=%q", tag.ID, tag.Name)), nil, nil
}

func (h *MCPHandler) HandleDeleteTag(_ context.Context, _ *mcp.CallToolRequest, args TagIDArgs) (*mcp.CallToolResult, any, error) {
	if args.ID == "" {
		return errResult("id is required"), nil, nil
	}
	if err := h.client.DeleteTag(args.ID); err != nil {
		return errResult(err.Error()), nil, nil
	}
	return textResult(fmt.Sprintf("Tag %s deleted and unlinked from all workflows.", args.ID)), nil, nil
}

// --- Credential get/update handlers ---

func (h *MCPHandler) HandleGetCredential(_ context.Context, _ *mcp.CallToolRequest, args GetCredentialArgs) (*mcp.CallToolResult, any, error) {
	if args.ID == "" {
		return errResult("id is required"), nil, nil
	}
	cred, err := h.client.GetCredential(args.ID)
	if err != nil {
		return errResult(err.Error()), nil, nil
	}
	return jsonResult(cred), nil, nil
}

func (h *MCPHandler) HandleUpdateCredential(_ context.Context, _ *mcp.CallToolRequest, args UpdateCredentialArgs) (*mcp.CallToolResult, any, error) {
	if args.ID == "" {
		return errResult("id is required"), nil, nil
	}
	if args.Name == "" {
		return errResult("name is required"), nil, nil
	}
	if args.Type == "" {
		return errResult("type is required"), nil, nil
	}
	cred, err := h.client.UpdateCredential(args.ID, args.Name, args.Type, args.DataJSON)
	if err != nil {
		return errResult(err.Error()), nil, nil
	}
	return textResult(fmt.Sprintf("Credential updated: id=%s name=%q type=%s", cred.ID, cred.Name, cred.Type)), nil, nil
}

// --- Tag handlers ---

func (h *MCPHandler) HandleListTags(_ context.Context, _ *mcp.CallToolRequest, args ListTagsArgs) (*mcp.CallToolResult, any, error) {
	result, err := h.client.ListTags(args.Limit, args.Cursor)
	if err != nil {
		return errResult(err.Error()), nil, nil
	}

	if len(result.Data) == 0 {
		return textResult("No tags found."), nil, nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Tags (%d):\n\n", len(result.Data)))
	for _, t := range result.Data {
		sb.WriteString(fmt.Sprintf("  %s  %s\n", t.ID, t.Name))
	}
	return textResult(sb.String()), nil, nil
}

func (h *MCPHandler) HandleCreateTag(_ context.Context, _ *mcp.CallToolRequest, args CreateTagArgs) (*mcp.CallToolResult, any, error) {
	if args.Name == "" {
		return errResult("name is required"), nil, nil
	}
	tag, err := h.client.CreateTag(args.Name)
	if err != nil {
		return errResult(err.Error()), nil, nil
	}
	return textResult(fmt.Sprintf("Tag created: id=%s name=%q", tag.ID, tag.Name)), nil, nil
}

// --- Webhook handler ---

func (h *MCPHandler) HandleTriggerWebhook(_ context.Context, _ *mcp.CallToolRequest, args TriggerWebhookArgs) (*mcp.CallToolResult, any, error) {
	if args.Path == "" {
		return errResult("path is required"), nil, nil
	}
	resp, status, err := h.client.TriggerWebhook(args.Path, args.Method, args.BodyJSON)
	if err != nil {
		return errResult(err.Error()), nil, nil
	}
	return textResult(fmt.Sprintf("HTTP %d\n%s", status, resp)), nil, nil
}

// --- Execution tag handlers ---

func (h *MCPHandler) HandleListExecutionTags(_ context.Context, _ *mcp.CallToolRequest, args ExecutionTagsArgs) (*mcp.CallToolResult, any, error) {
	if args.ID == "" {
		return errResult("id is required"), nil, nil
	}
	tags, err := h.client.ListExecutionTags(args.ID)
	if err != nil {
		return errResult(err.Error()), nil, nil
	}
	if len(tags) == 0 {
		return textResult("No tags on this execution."), nil, nil
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Tags (%d):\n", len(tags)))
	for _, t := range tags {
		sb.WriteString(fmt.Sprintf("  %s  %s\n", t.ID, t.Name))
	}
	return textResult(sb.String()), nil, nil
}

func (h *MCPHandler) HandleUpdateExecutionTags(_ context.Context, _ *mcp.CallToolRequest, args UpdateExecutionTagsArgs) (*mcp.CallToolResult, any, error) {
	if args.ID == "" {
		return errResult("id is required"), nil, nil
	}
	tags, err := h.client.UpdateExecutionTags(args.ID, args.TagIDs)
	if err != nil {
		return errResult(err.Error()), nil, nil
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Execution %s tags updated (%d):\n", args.ID, len(tags)))
	for _, t := range tags {
		sb.WriteString(fmt.Sprintf("  %s  %s\n", t.ID, t.Name))
	}
	return textResult(sb.String()), nil, nil
}

// --- Audit handler ---

func (h *MCPHandler) HandleGenerateAudit(_ context.Context, _ *mcp.CallToolRequest, args GenerateAuditArgs) (*mcp.CallToolResult, any, error) {
	var opts json.RawMessage
	if args.OptionsJSON != "" {
		if err := json.Unmarshal([]byte(args.OptionsJSON), &opts); err != nil {
			return errResult("invalid options JSON: " + err.Error()), nil, nil
		}
	}
	result, err := h.client.GenerateAudit(opts)
	if err != nil {
		return errResult(err.Error()), nil, nil
	}
	pretty, _ := json.MarshalIndent(json.RawMessage(result), "", "  ")
	return textResult(string(pretty)), nil, nil
}

// --- Node DB handlers ---

func (h *MCPHandler) HandleSearchNodes(_ context.Context, _ *mcp.CallToolRequest, args SearchNodesArgs) (*mcp.CallToolResult, any, error) {
	if args.Keywords == "" {
		return errResult("keywords is required"), nil, nil
	}
	results, err := SearchNodes(h.nodeDB, args.Keywords, args.Group, args.Limit)
	if err != nil {
		return errResult(err.Error()), nil, nil
	}
	if len(results) == 0 {
		return textResult("No node types found matching: " + args.Keywords), nil, nil
	}
	type slim struct {
		Name        string `json:"name"`
		DisplayName string `json:"display_name"`
	}
	out := make([]slim, len(results))
	for i, r := range results {
		out[i] = slim{Name: r.Name, DisplayName: r.DisplayName}
	}
	return jsonResult(out), nil, nil
}

func (h *MCPHandler) HandleGetNodeSchema(_ context.Context, _ *mcp.CallToolRequest, args GetNodeSchemaArgs) (*mcp.CallToolResult, any, error) {
	if args.Name == "" {
		return errResult("name is required"), nil, nil
	}
	schema, err := GetNodeSchema(h.nodeDB, args.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return errResult("node type not found: " + args.Name), nil, nil
		}
		return errResult(err.Error()), nil, nil
	}
	return jsonResult(schema), nil, nil
}

// --- User handlers ---

func (h *MCPHandler) HandleListUsers(_ context.Context, _ *mcp.CallToolRequest, args ListUsersArgs) (*mcp.CallToolResult, any, error) {
	result, err := h.client.ListUsers(args.Limit, args.Cursor, args.IncludeRole)
	if err != nil {
		return errResult(err.Error()), nil, nil
	}
	if len(result.Data) == 0 {
		return textResult("No users found."), nil, nil
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Users (%d):\n\n", len(result.Data)))
	sb.WriteString(fmt.Sprintf("%-36s  %-30s  %-20s  %s\n", "ID", "Email", "Name", "Role"))
	sb.WriteString(strings.Repeat("-", 100) + "\n")
	for _, u := range result.Data {
		name := strings.TrimSpace(u.FirstName + " " + u.LastName)
		sb.WriteString(fmt.Sprintf("%-36s  %-30s  %-20s  %s\n", u.ID, u.Email, name, u.Role))
	}
	if result.NextCursor != "" {
		sb.WriteString(fmt.Sprintf("\nNext cursor: %s", result.NextCursor))
	}
	return textResult(sb.String()), nil, nil
}

func (h *MCPHandler) HandleGetUser(_ context.Context, _ *mcp.CallToolRequest, args GetUserArgs) (*mcp.CallToolResult, any, error) {
	if args.ID == "" {
		return errResult("id is required"), nil, nil
	}
	u, err := h.client.GetUser(args.ID, args.IncludeRole)
	if err != nil {
		return errResult(err.Error()), nil, nil
	}
	return jsonResult(u), nil, nil
}

func (h *MCPHandler) HandleCreateUsers(_ context.Context, _ *mcp.CallToolRequest, args CreateUsersArgs) (*mcp.CallToolResult, any, error) {
	if args.UsersJSON == "" {
		return errResult("users_json is required"), nil, nil
	}
	var users []CreateUserRequest
	if err := json.Unmarshal([]byte(args.UsersJSON), &users); err != nil {
		return errResult("invalid users_json: " + err.Error()), nil, nil
	}
	if len(users) == 0 {
		return errResult("users_json array must not be empty"), nil, nil
	}
	created, err := h.client.CreateUsers(users)
	if err != nil {
		return errResult(err.Error()), nil, nil
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Created %d user(s):\n", len(created)))
	for _, u := range created {
		sb.WriteString(fmt.Sprintf("  id=%s email=%s\n", u.ID, u.Email))
	}
	return textResult(sb.String()), nil, nil
}

func (h *MCPHandler) HandleChangeUserRole(_ context.Context, _ *mcp.CallToolRequest, args ChangeUserRoleArgs) (*mcp.CallToolResult, any, error) {
	if args.ID == "" {
		return errResult("id is required"), nil, nil
	}
	if args.Role == "" {
		return errResult("role is required"), nil, nil
	}
	u, err := h.client.ChangeUserRole(args.ID, args.Role)
	if err != nil {
		return errResult(err.Error()), nil, nil
	}
	return textResult(fmt.Sprintf("User %s role updated to %q.", u.ID, u.Role)), nil, nil
}

func (h *MCPHandler) HandleDeleteUser(_ context.Context, _ *mcp.CallToolRequest, args DeleteUserArgs) (*mcp.CallToolResult, any, error) {
	if args.ID == "" {
		return errResult("id is required"), nil, nil
	}
	if err := h.client.DeleteUser(args.ID); err != nil {
		return errResult(err.Error()), nil, nil
	}
	return textResult(fmt.Sprintf("User %s deleted.", args.ID)), nil, nil
}

// --- Partial update handler ---

func (h *MCPHandler) HandleUpdatePartialWorkflow(_ context.Context, _ *mcp.CallToolRequest, args UpdatePartialWorkflowArgs) (*mcp.CallToolResult, any, error) {
	if args.ID == "" {
		return errResult("id is required"), nil, nil
	}
	hasChange := args.Name != "" || args.NodesJSON != "" || args.ConnectionsJSON != "" || args.SettingsJSON != ""
	if !hasChange {
		return errResult("at least one of name, nodes, connections, or settings is required"), nil, nil
	}

	patches := make(map[string]json.RawMessage)

	if args.Name != "" {
		nameJSON, _ := json.Marshal(args.Name)
		patches["name"] = nameJSON
	}
	for field, raw := range map[string]string{
		"nodes":       args.NodesJSON,
		"connections": args.ConnectionsJSON,
		"settings":    args.SettingsJSON,
	} {
		if raw == "" {
			continue
		}
		var v json.RawMessage
		if err := json.Unmarshal([]byte(raw), &v); err != nil {
			return errResult(fmt.Sprintf("invalid %s JSON: %s", field, err)), nil, nil
		}
		patches[field] = v
	}

	w, err := h.client.UpdatePartialWorkflow(args.ID, patches)
	if err != nil {
		return errResult(err.Error()), nil, nil
	}
	return textResult(fmt.Sprintf("Workflow updated: id=%s name=%q active=%v", w.ID, w.Name, w.Active)), nil, nil
}

// --- Validation handler ---

func (h *MCPHandler) HandleValidateWorkflow(_ context.Context, _ *mcp.CallToolRequest, args ValidateWorkflowArgs) (*mcp.CallToolResult, any, error) {
	if args.WorkflowJSON == "" {
		return errResult("workflow_json is required"), nil, nil
	}

	var wf map[string]json.RawMessage
	if err := json.Unmarshal([]byte(args.WorkflowJSON), &wf); err != nil {
		return textResult("ERROR: invalid JSON — " + err.Error()), nil, nil
	}

	var errs []string

	// Top-level field checks
	nameRaw, hasName := wf["name"]
	if !hasName || string(nameRaw) == `""` || string(nameRaw) == "null" {
		errs = append(errs, "WARN: no 'name' field — n8n will auto-assign one")
	}
	if _, ok := wf["nodes"]; !ok {
		errs = append(errs, "ERROR: missing required field 'nodes'")
	}
	if _, ok := wf["connections"]; !ok {
		errs = append(errs, "ERROR: missing required field 'connections'")
	}

	// Validate nodes array
	if nodesRaw, ok := wf["nodes"]; ok {
		var nodes []map[string]json.RawMessage
		if err := json.Unmarshal(nodesRaw, &nodes); err != nil {
			errs = append(errs, "ERROR: 'nodes' must be a JSON array")
		} else {
			seenIDs := map[string]bool{}
			for i, node := range nodes {
				prefix := fmt.Sprintf("nodes[%d]", i)
				for _, req := range []string{"id", "name", "type", "typeVersion", "position"} {
					if _, ok := node[req]; !ok {
						errs = append(errs, fmt.Sprintf("ERROR: %s missing required field '%s'", prefix, req))
					}
				}
				if idRaw, ok := node["id"]; ok {
					var id string
					if json.Unmarshal(idRaw, &id) == nil {
						if seenIDs[id] {
							errs = append(errs, fmt.Sprintf("ERROR: duplicate node id %q", id))
						}
						seenIDs[id] = true
					}
				}
			}
		}
	}

	// Validate connections is an object
	if connRaw, ok := wf["connections"]; ok {
		var conn map[string]json.RawMessage
		if err := json.Unmarshal(connRaw, &conn); err != nil {
			errs = append(errs, "ERROR: 'connections' must be a JSON object")
		}
	}

	// Validate node types exist in DB
	if h.nodeDB != nil {
		if nodesRaw, ok := wf["nodes"]; ok {
			var nodes []map[string]json.RawMessage
			if err := json.Unmarshal(nodesRaw, &nodes); err == nil {
				for i, node := range nodes {
					typeRaw, ok := node["type"]
					if !ok {
						continue
					}
					var nodeType string
					if err := json.Unmarshal(typeRaw, &nodeType); err != nil || nodeType == "" {
						continue
					}
					var count int
					_ = h.nodeDB.QueryRow(`SELECT COUNT(*) FROM node_types WHERE name = ?`, nodeType).Scan(&count)
					if count == 0 {
						// Derive a search keyword from the type string for a helpful suggestion
						parts := strings.Split(nodeType, ".")
						keyword := parts[len(parts)-1]
						errs = append(errs, fmt.Sprintf(
							"ERROR: nodes[%d] type %q not found in node registry — run n8n_search_nodes keyword=%q to find the correct type",
							i, nodeType, keyword,
						))
					}
				}
			}
		}
	}

	if len(errs) == 0 {
		return textResult("✓ valid — workflow structure looks correct"), nil, nil
	}
	return textResult(strings.Join(errs, "\n")), nil, nil
}

// --- Execute / archive handlers ---

func (h *MCPHandler) HandleExecuteWorkflow(_ context.Context, _ *mcp.CallToolRequest, args ExecuteWorkflowArgs) (*mcp.CallToolResult, any, error) {
	if args.WorkflowID == "" {
		return errResult("workflow_id is required"), nil, nil
	}
	result, err := h.client.ExecuteWorkflow(args.WorkflowID, args.RunDataJSON)
	if err != nil {
		return errResult(err.Error()), nil, nil
	}
	pretty, _ := json.MarshalIndent(json.RawMessage(result), "", "  ")
	return textResult(string(pretty)), nil, nil
}

func (h *MCPHandler) HandleArchiveWorkflow(_ context.Context, _ *mcp.CallToolRequest, args ArchiveWorkflowArgs) (*mcp.CallToolResult, any, error) {
	if args.ID == "" {
		return errResult("id is required"), nil, nil
	}
	w, err := h.client.ArchiveWorkflow(args.ID)
	if err != nil {
		return errResult(err.Error()), nil, nil
	}
	return textResult(fmt.Sprintf("Workflow %q (id=%s) archived.", w.Name, w.ID)), nil, nil
}

func (h *MCPHandler) HandleUnarchiveWorkflow(_ context.Context, _ *mcp.CallToolRequest, args ArchiveWorkflowArgs) (*mcp.CallToolResult, any, error) {
	if args.ID == "" {
		return errResult("id is required"), nil, nil
	}
	w, err := h.client.UnarchiveWorkflow(args.ID)
	if err != nil {
		return errResult(err.Error()), nil, nil
	}
	return textResult(fmt.Sprintf("Workflow %q (id=%s) unarchived.", w.Name, w.ID)), nil, nil
}

func (h *MCPHandler) HandleTestCredential(_ context.Context, _ *mcp.CallToolRequest, args TestCredentialArgs) (*mcp.CallToolResult, any, error) {
	if args.ID == "" {
		return errResult("id is required"), nil, nil
	}
	result, err := h.client.TestCredential(args.ID)
	if err != nil {
		return errResult(err.Error()), nil, nil
	}
	pretty, _ := json.MarshalIndent(json.RawMessage(result), "", "  ")
	return textResult(string(pretty)), nil, nil
}

// --- DataTable handlers ---

func (h *MCPHandler) HandleListDataTables(_ context.Context, _ *mcp.CallToolRequest, args ListDataTablesArgs) (*mcp.CallToolResult, any, error) {
	result, err := h.client.ListDataTables(args.Limit, args.Cursor)
	if err != nil {
		return errResult(err.Error()), nil, nil
	}
	if len(result.Data) == 0 {
		return textResult("No data tables found."), nil, nil
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Data Tables (%d):\n\n", len(result.Data)))
	sb.WriteString(fmt.Sprintf("%-24s  %s\n", "ID", "Name"))
	sb.WriteString(strings.Repeat("-", 50) + "\n")
	for _, t := range result.Data {
		sb.WriteString(fmt.Sprintf("%-24s  %s\n", t.ID, t.Name))
	}
	if result.NextCursor != "" {
		sb.WriteString(fmt.Sprintf("\nNext cursor: %s", result.NextCursor))
	}
	return textResult(sb.String()), nil, nil
}

func (h *MCPHandler) HandleCreateDataTable(_ context.Context, _ *mcp.CallToolRequest, args CreateDataTableArgs) (*mcp.CallToolResult, any, error) {
	if args.Name == "" {
		return errResult("name is required"), nil, nil
	}
	var cols []DataTableColumn
	if args.ColumnsJSON != "" {
		if err := json.Unmarshal([]byte(args.ColumnsJSON), &cols); err != nil {
			return errResult("invalid columns_json: " + err.Error()), nil, nil
		}
	}
	t, err := h.client.CreateDataTable(args.Name, cols)
	if err != nil {
		return errResult(err.Error()), nil, nil
	}
	return textResult(fmt.Sprintf("Data table created: id=%s name=%q", t.ID, t.Name)), nil, nil
}

func (h *MCPHandler) HandleGetDataTable(_ context.Context, _ *mcp.CallToolRequest, args DataTableIDArgs) (*mcp.CallToolResult, any, error) {
	if args.ID == "" {
		return errResult("id is required"), nil, nil
	}
	t, err := h.client.GetDataTable(args.ID)
	if err != nil {
		return errResult(err.Error()), nil, nil
	}
	return jsonResult(t), nil, nil
}

func (h *MCPHandler) HandleUpdateDataTable(_ context.Context, _ *mcp.CallToolRequest, args UpdateDataTableArgs) (*mcp.CallToolResult, any, error) {
	if args.ID == "" {
		return errResult("id is required"), nil, nil
	}
	if args.Name == "" {
		return errResult("name is required"), nil, nil
	}
	t, err := h.client.UpdateDataTable(args.ID, args.Name)
	if err != nil {
		return errResult(err.Error()), nil, nil
	}
	return textResult(fmt.Sprintf("Data table updated: id=%s name=%q", t.ID, t.Name)), nil, nil
}

func (h *MCPHandler) HandleDeleteDataTable(_ context.Context, _ *mcp.CallToolRequest, args DataTableIDArgs) (*mcp.CallToolResult, any, error) {
	if args.ID == "" {
		return errResult("id is required"), nil, nil
	}
	if err := h.client.DeleteDataTable(args.ID); err != nil {
		return errResult(err.Error()), nil, nil
	}
	return textResult(fmt.Sprintf("Data table %s deleted.", args.ID)), nil, nil
}

func (h *MCPHandler) HandleListRows(_ context.Context, _ *mcp.CallToolRequest, args ListRowsArgs) (*mcp.CallToolResult, any, error) {
	if args.TableID == "" {
		return errResult("table_id is required"), nil, nil
	}
	data, err := h.client.ListRows(args.TableID, args.Limit, args.Cursor)
	if err != nil {
		return errResult(err.Error()), nil, nil
	}
	pretty, _ := json.MarshalIndent(json.RawMessage(data), "", "  ")
	return textResult(string(pretty)), nil, nil
}

func (h *MCPHandler) HandleCreateRow(_ context.Context, _ *mcp.CallToolRequest, args CreateRowArgs) (*mcp.CallToolResult, any, error) {
	if args.TableID == "" {
		return errResult("table_id is required"), nil, nil
	}
	if args.RowDataJSON == "" {
		return errResult("row_data is required"), nil, nil
	}
	data, err := h.client.CreateRow(args.TableID, args.RowDataJSON)
	if err != nil {
		return errResult(err.Error()), nil, nil
	}
	pretty, _ := json.MarshalIndent(json.RawMessage(data), "", "  ")
	return textResult(string(pretty)), nil, nil
}

func (h *MCPHandler) HandleUpdateRows(_ context.Context, _ *mcp.CallToolRequest, args UpdateRowsArgs) (*mcp.CallToolResult, any, error) {
	if args.TableID == "" {
		return errResult("table_id is required"), nil, nil
	}
	if args.UpdateJSON == "" {
		return errResult("update is required"), nil, nil
	}
	data, err := h.client.UpdateRows(args.TableID, args.FilterJSON, args.UpdateJSON)
	if err != nil {
		return errResult(err.Error()), nil, nil
	}
	pretty, _ := json.MarshalIndent(json.RawMessage(data), "", "  ")
	return textResult(string(pretty)), nil, nil
}

func (h *MCPHandler) HandleUpsertRows(_ context.Context, _ *mcp.CallToolRequest, args UpsertRowsArgs) (*mcp.CallToolResult, any, error) {
	if args.TableID == "" {
		return errResult("table_id is required"), nil, nil
	}
	if args.RowsJSON == "" {
		return errResult("rows is required"), nil, nil
	}
	data, err := h.client.UpsertRows(args.TableID, args.RowsJSON)
	if err != nil {
		return errResult(err.Error()), nil, nil
	}
	pretty, _ := json.MarshalIndent(json.RawMessage(data), "", "  ")
	return textResult(string(pretty)), nil, nil
}

func (h *MCPHandler) HandleDeleteRows(_ context.Context, _ *mcp.CallToolRequest, args DeleteRowsArgs) (*mcp.CallToolResult, any, error) {
	if args.TableID == "" {
		return errResult("table_id is required"), nil, nil
	}
	data, err := h.client.DeleteRows(args.TableID, args.FilterJSON)
	if err != nil {
		return errResult(err.Error()), nil, nil
	}
	pretty, _ := json.MarshalIndent(json.RawMessage(data), "", "  ")
	return textResult(string(pretty)), nil, nil
}

func (h *MCPHandler) HandleListColumns(_ context.Context, _ *mcp.CallToolRequest, args ListColumnsArgs) (*mcp.CallToolResult, any, error) {
	if args.TableID == "" {
		return errResult("table_id is required"), nil, nil
	}
	cols, err := h.client.ListColumns(args.TableID)
	if err != nil {
		return errResult(err.Error()), nil, nil
	}
	if len(cols) == 0 {
		return textResult("No columns found."), nil, nil
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Columns (%d):\n\n", len(cols)))
	sb.WriteString(fmt.Sprintf("%-24s  %-20s  %s\n", "ID", "Name", "DataType"))
	sb.WriteString(strings.Repeat("-", 60) + "\n")
	for _, c := range cols {
		sb.WriteString(fmt.Sprintf("%-24s  %-20s  %s\n", c.ID, c.Name, c.DataType))
	}
	return textResult(sb.String()), nil, nil
}

func (h *MCPHandler) HandleCreateColumn(_ context.Context, _ *mcp.CallToolRequest, args CreateColumnArgs) (*mcp.CallToolResult, any, error) {
	if args.TableID == "" {
		return errResult("table_id is required"), nil, nil
	}
	if args.Name == "" {
		return errResult("name is required"), nil, nil
	}
	col, err := h.client.CreateColumn(args.TableID, args.Name, args.DataType)
	if err != nil {
		return errResult(err.Error()), nil, nil
	}
	return textResult(fmt.Sprintf("Column created: id=%s name=%q type=%s", col.ID, col.Name, col.DataType)), nil, nil
}

func (h *MCPHandler) HandleGetColumn(_ context.Context, _ *mcp.CallToolRequest, args ColumnIDArgs) (*mcp.CallToolResult, any, error) {
	if args.TableID == "" {
		return errResult("table_id is required"), nil, nil
	}
	if args.ColumnID == "" {
		return errResult("column_id is required"), nil, nil
	}
	col, err := h.client.GetColumn(args.TableID, args.ColumnID)
	if err != nil {
		return errResult(err.Error()), nil, nil
	}
	return textResult(fmt.Sprintf("id=%s name=%q type=%s", col.ID, col.Name, col.DataType)), nil, nil
}

func (h *MCPHandler) HandleUpdateColumn(_ context.Context, _ *mcp.CallToolRequest, args UpdateColumnArgs) (*mcp.CallToolResult, any, error) {
	if args.TableID == "" {
		return errResult("table_id is required"), nil, nil
	}
	if args.ColumnID == "" {
		return errResult("column_id is required"), nil, nil
	}
	if args.Name == "" {
		return errResult("name is required"), nil, nil
	}
	col, err := h.client.UpdateColumn(args.TableID, args.ColumnID, args.Name)
	if err != nil {
		return errResult(err.Error()), nil, nil
	}
	return textResult(fmt.Sprintf("Column updated: id=%s name=%q", col.ID, col.Name)), nil, nil
}

func (h *MCPHandler) HandleDeleteColumn(_ context.Context, _ *mcp.CallToolRequest, args ColumnIDArgs) (*mcp.CallToolResult, any, error) {
	if args.TableID == "" {
		return errResult("table_id is required"), nil, nil
	}
	if args.ColumnID == "" {
		return errResult("column_id is required"), nil, nil
	}
	if err := h.client.DeleteColumn(args.TableID, args.ColumnID); err != nil {
		return errResult(err.Error()), nil, nil
	}
	return textResult(fmt.Sprintf("Column %s deleted from table %s.", args.ColumnID, args.TableID)), nil, nil
}

// --- Server bootstrap ---

func RunMCPServer() {
	n8nClient, err := NewN8nClient()
	if err != nil {
		log.Printf("Warning: n8n client init failed (tools will return errors): %v", err)
		n8nClient = &N8nClient{httpClient: &http.Client{}}
	}

	nodeDB, dbCleanup, err := InitNodeDB()
	if err != nil {
		log.Printf("Warning: node DB init failed (node validation disabled): %v", err)
		nodeDB = nil
		dbCleanup = func() {}
	}
	defer dbCleanup()

	handler := NewMCPHandler(n8nClient, nodeDB)

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "oido-n8n",
		Version: "1.0.0",
	}, nil)

	// --- Workflow tools ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_list_workflows",
		Description: "List workflows with ID, name, and active status. Filter by active=true/false, tag names, or name.",
	}, handler.HandleListWorkflows)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_get_workflow",
		Description: "Get full workflow definition (nodes + connections) by ID.",
	}, handler.HandleGetWorkflow)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_create_workflow",
		Description: "Create a new n8n workflow from JSON definition. Validates node types and required fields internally.",
	}, handler.HandleCreateWorkflow)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_update_workflow",
		Description: "Replace a workflow by ID with a full updated JSON definition.",
	}, handler.HandleUpdateWorkflow)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_patch_workflow",
		Description: "Partially update a workflow: change name, nodes, connections, or settings without replacing the whole definition.",
	}, handler.HandleUpdatePartialWorkflow)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_delete_workflow",
		Description: "Permanently delete a workflow by ID. Irreversible.",
	}, handler.HandleDeleteWorkflow)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_activate_workflow",
		Description: "Activate a workflow (enable its triggers so it runs automatically).",
	}, handler.HandleActivateWorkflow)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_deactivate_workflow",
		Description: "Deactivate a workflow (pause its triggers).",
	}, handler.HandleDeactivateWorkflow)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_execute_workflow",
		Description: "Manually execute a workflow by ID. Optionally provide input run_data as JSON.",
	}, handler.HandleExecuteWorkflow)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_archive_workflow",
		Description: "Archive a workflow (hides it from the main list without deleting).",
	}, handler.HandleArchiveWorkflow)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_unarchive_workflow",
		Description: "Unarchive a previously archived workflow.",
	}, handler.HandleUnarchiveWorkflow)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_validate_workflow",
		Description: "Validate workflow JSON structure and check node types against the node registry.",
	}, handler.HandleValidateWorkflow)
	mcp.AddTool(server, &mcp.Tool{
		Name: "n8n_search_nodes",
		Description: "Search available n8n node types by keyword (partial match on name and display name). " +
			"Supports comma-separated keywords (OR logic). Optional group filter: 't'=triggers, 'i'=actions, 'o'=outputs.",
	}, handler.HandleSearchNodes)

	// --- Execution tools ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_list_executions",
		Description: "List workflow executions. Filter by workflow ID or status (waiting/running/success/error/canceled).",
	}, handler.HandleListExecutions)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_get_execution",
		Description: "Get execution details by ID. Set include_data=true to include node input/output data.",
	}, handler.HandleGetExecution)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_delete_execution",
		Description: "Delete an execution record by ID.",
	}, handler.HandleDeleteExecution)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_stop_execution",
		Description: "Stop a specific running execution by ID.",
	}, handler.HandleStopExecution)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_stop_executions",
		Description: "Stop all running executions, optionally filtered by workflow ID.",
	}, handler.HandleStopExecutions)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_retry_execution",
		Description: "Retry a failed execution by ID.",
	}, handler.HandleRetryExecution)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_list_execution_tags",
		Description: "List tags attached to an execution.",
	}, handler.HandleListExecutionTags)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_update_execution_tags",
		Description: "Replace all tags on an execution with a new set of tag IDs.",
	}, handler.HandleUpdateExecutionTags)

	// --- Credential tools ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_list_credentials",
		Description: "List credentials (secrets excluded). Filter by type or credential ID.",
	}, handler.HandleListCredentials)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_get_credential",
		Description: "Get credential metadata by ID (secrets excluded).",
	}, handler.HandleGetCredential)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_create_credential",
		Description: "Create a credential. Use n8n_get_credential_schema first to discover required fields.",
	}, handler.HandleCreateCredential)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_update_credential",
		Description: "Update a credential by ID (name, type, data).",
	}, handler.HandleUpdateCredential)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_delete_credential",
		Description: "Delete a credential by ID. Irreversible.",
	}, handler.HandleDeleteCredential)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_get_credential_schema",
		Description: "Get required fields and schema for a credential type before creating one.",
	}, handler.HandleGetCredentialSchema)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_test_credential",
		Description: "Test whether a credential is valid and can connect to its service.",
	}, handler.HandleTestCredential)

	// --- Tag tools ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_list_tags",
		Description: "List all workflow tags.",
	}, handler.HandleListTags)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_create_tag",
		Description: "Create a new workflow tag.",
	}, handler.HandleCreateTag)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_get_tag",
		Description: "Get a tag by ID.",
	}, handler.HandleGetTag)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_update_tag",
		Description: "Rename a tag by ID.",
	}, handler.HandleUpdateTag)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_delete_tag",
		Description: "Delete a tag and unlink it from all workflows.",
	}, handler.HandleDeleteTag)

	// --- Webhook tools ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_trigger_webhook",
		Description: "Trigger a workflow via its webhook path. No API key needed — uses the n8n webhook endpoint directly.",
	}, handler.HandleTriggerWebhook)

	// --- Project tools ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_list_projects",
		Description: "List n8n projects.",
	}, handler.HandleListProjects)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_get_project",
		Description: "Get a project by ID.",
	}, handler.HandleGetProject)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_create_project",
		Description: "Create a new project.",
	}, handler.HandleCreateProject)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_update_project",
		Description: "Rename a project by ID.",
	}, handler.HandleUpdateProject)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_delete_project",
		Description: "Delete a project by ID.",
	}, handler.HandleDeleteProject)

	// --- Variable tools ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_list_variables",
		Description: "List n8n environment variables.",
	}, handler.HandleListVariables)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_get_variable",
		Description: "Get a variable by ID.",
	}, handler.HandleGetVariable)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_create_variable",
		Description: "Create a new environment variable.",
	}, handler.HandleCreateVariable)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_update_variable",
		Description: "Update a variable's key or value by ID.",
	}, handler.HandleUpdateVariable)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_delete_variable",
		Description: "Delete a variable by ID.",
	}, handler.HandleDeleteVariable)

	// --- User tools ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_list_users",
		Description: "List n8n users. Set include_role=true to include global roles.",
	}, handler.HandleListUsers)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_get_user",
		Description: "Get a user by ID or email address.",
	}, handler.HandleGetUser)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_create_users",
		Description: "Invite users by providing a JSON array of user objects (email required, plus optional role, firstName, lastName).",
	}, handler.HandleCreateUsers)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_change_user_role",
		Description: "Change a user's global role (e.g. global:admin, global:member).",
	}, handler.HandleChangeUserRole)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_delete_user",
		Description: "Delete a user by ID.",
	}, handler.HandleDeleteUser)

	// --- Audit tool ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_generate_audit",
		Description: "Generate a security audit report for the n8n instance. Optionally filter by categories (credentials, workflows, etc.).",
	}, handler.HandleGenerateAudit)

	// --- DataTable tools ---
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_list_data_tables",
		Description: "List all data tables in the n8n instance.",
	}, handler.HandleListDataTables)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_create_data_table",
		Description: "Create a new data table with an optional set of columns.",
	}, handler.HandleCreateDataTable)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_get_data_table",
		Description: "Get data table metadata by ID.",
	}, handler.HandleGetDataTable)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_update_data_table",
		Description: "Rename a data table by ID.",
	}, handler.HandleUpdateDataTable)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_delete_data_table",
		Description: "Delete a data table and all its data. Irreversible.",
	}, handler.HandleDeleteDataTable)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_list_rows",
		Description: "List rows in a data table.",
	}, handler.HandleListRows)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_create_row",
		Description: "Insert a new row into a data table. Provide row_data as a JSON object with column names as keys.",
	}, handler.HandleCreateRow)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_update_rows",
		Description: "Update rows in a data table matching an optional filter. Provide update as a JSON object of fields to set.",
	}, handler.HandleUpdateRows)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_upsert_rows",
		Description: "Insert or update rows in a data table. Provide rows as a JSON array.",
	}, handler.HandleUpsertRows)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_delete_rows",
		Description: "Delete rows from a data table matching an optional filter. Omit filter to delete all rows.",
	}, handler.HandleDeleteRows)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_list_columns",
		Description: "List all columns in a data table.",
	}, handler.HandleListColumns)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_create_column",
		Description: "Add a column to a data table. Specify data_type (string, number, boolean, date).",
	}, handler.HandleCreateColumn)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_get_column",
		Description: "Get a data table column by ID.",
	}, handler.HandleGetColumn)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_update_column",
		Description: "Rename a column in a data table.",
	}, handler.HandleUpdateColumn)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_delete_column",
		Description: "Delete a column from a data table. Irreversible.",
	}, handler.HandleDeleteColumn)

	ctx := context.Background()
	log.Println("Oido n8n MCP Server starting on stdio...")
	if err := server.Run(ctx, &mcp.StdioTransport{}); err != nil {
		log.Fatalf("MCP server error: %v", err)
	}
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
