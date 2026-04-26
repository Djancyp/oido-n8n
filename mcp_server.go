package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
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

type ExecuteWorkflowArgs struct {
	ID       string `json:"id"              jsonschema:"Workflow ID"`
	DataJSON string `json:"data,omitempty"  jsonschema:"Optional input data as JSON string"`
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
	Keyword string `json:"keyword"         jsonschema:"Search term matched against node name and display name (partial match)"`
	Group   string `json:"group,omitempty" jsonschema:"Filter by group: t=trigger, i=action, o=output"`
	Limit   int    `json:"limit,omitempty" jsonschema:"Max results (default: 20)"`
}

type GetNodeSchemaArgs struct {
	Name string `json:"name" jsonschema:"Exact node type name e.g. n8n-nodes-base.httpRequest"`
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

func (h *MCPHandler) HandleExecuteWorkflow(_ context.Context, _ *mcp.CallToolRequest, args ExecuteWorkflowArgs) (*mcp.CallToolResult, any, error) {
	if args.ID == "" {
		return errResult("id is required"), nil, nil
	}
	resp, err := h.client.ExecuteWorkflow(args.ID, args.DataJSON)
	if err != nil {
		return errResult(err.Error()), nil, nil
	}
	return textResult(fmt.Sprintf("Workflow execution started. executionId=%d", resp.ExecutionID)), nil, nil
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
	if args.Keyword == "" {
		return errResult("keyword is required"), nil, nil
	}
	results, err := SearchNodes(h.nodeDB, args.Keyword, args.Group, args.Limit)
	if err != nil {
		return errResult(err.Error()), nil, nil
	}
	if len(results) == 0 {
		return textResult("No node types found matching: " + args.Keyword), nil, nil
	}
	return jsonResult(results), nil, nil
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

// --- Server bootstrap ---

func RunMCPServer() {
	n8nClient, err := NewN8nClient()
	if err != nil {
		log.Fatalf("Failed to create n8n client: %v", err)
	}

	nodeDB, dbCleanup, err := InitNodeDB()
	if err != nil {
		log.Fatalf("Failed to load node DB: %v", err)
	}
	defer dbCleanup()

	handler := NewMCPHandler(n8nClient, nodeDB)

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "oido-n8n",
		Version: "1.0.0",
	}, nil)

	mcp.AddTool(server, &mcp.Tool{
		Name: "n8n_search_nodes",
		Description: "Search available n8n node types by keyword (partial match on name and display name). " +
			"Always call this first to get the correct type string and version before building a workflow. " +
			"Optional group filter: 't'=triggers, 'i'=actions, 'o'=outputs.",
	}, handler.HandleSearchNodes)

	mcp.AddTool(server, &mcp.Tool{
		Name: "n8n_create_workflow",
		Description: `Create a new n8n workflow.

BEFORE CALLING: use n8n_search_nodes to find correct type strings — never guess node type names.
This tool validates the workflow internally and will REJECT it if node types are wrong or required fields are missing.

Every node requires: id, name, type, typeVersion, position, parameters.
Connection keys use node "name" (not "id").

Minimal example:
{
  "name": "My Workflow",
  "nodes": [{"id":"t1","name":"Schedule","type":"n8n-nodes-base.scheduleTrigger","typeVersion":1,"position":[240,300],"parameters":{}}],
  "connections": {}
}`,
	}, handler.HandleCreateWorkflow)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_update_workflow",
		Description: "Update an existing workflow by ID. Provide the full updated workflow JSON. Use n8n_search_nodes to verify node types before updating.",
	}, handler.HandleUpdateWorkflow)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_delete_workflow",
		Description: "Permanently delete a workflow by ID. Irreversible.",
	}, handler.HandleDeleteWorkflow)

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
