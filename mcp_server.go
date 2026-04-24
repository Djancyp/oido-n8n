package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type MCPHandler struct {
	client *N8nClient
}

func NewMCPHandler(client *N8nClient) *MCPHandler {
	return &MCPHandler{client: client}
}

// --- Arg types ---

type ListWorkflowsArgs struct {
	Active *bool    `json:"active,omitempty" jsonschema:"Filter by active status"`
	Tags   []string `json:"tags,omitempty"   jsonschema:"Filter by tag names"`
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
	WorkflowID string `json:"workflow_id,omitempty" jsonschema:"Filter by workflow ID"`
	Status     string `json:"status,omitempty"      jsonschema:"Filter: waiting, running, success, error, canceled"`
	Limit      int    `json:"limit,omitempty"       jsonschema:"Max results (default: 20)"`
	Cursor     string `json:"cursor,omitempty"      jsonschema:"Pagination cursor"`
}

type ExecutionIDArgs struct {
	ID string `json:"id" jsonschema:"Execution ID"`
}

type ListCredentialsArgs struct {
	Limit int `json:"limit,omitempty" jsonschema:"Max results (default: 20)"`
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
	Limit int `json:"limit,omitempty" jsonschema:"Max results (default: 100)"`
}

type CreateTagArgs struct {
	Name string `json:"name" jsonschema:"Tag name"`
}

type TriggerWebhookArgs struct {
	Path     string `json:"path"            jsonschema:"Webhook path (without /webhook/ prefix)"`
	Method   string `json:"method,omitempty" jsonschema:"HTTP method: GET, POST, PUT, PATCH, DELETE (default: POST)"`
	BodyJSON string `json:"body,omitempty"  jsonschema:"Request body as JSON string"`
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
	result, err := h.client.ListWorkflows(args.Active, args.Tags, args.Limit, args.Cursor)
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
	result, err := h.client.ListExecutions(args.WorkflowID, args.Status, args.Limit, args.Cursor)
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

func (h *MCPHandler) HandleGetExecution(_ context.Context, _ *mcp.CallToolRequest, args ExecutionIDArgs) (*mcp.CallToolResult, any, error) {
	if args.ID == "" {
		return errResult("id is required"), nil, nil
	}
	e, err := h.client.GetExecution(args.ID)
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

// --- Credential handlers ---

func (h *MCPHandler) HandleListCredentials(_ context.Context, _ *mcp.CallToolRequest, args ListCredentialsArgs) (*mcp.CallToolResult, any, error) {
	result, err := h.client.ListCredentials(args.Limit)
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

// --- Tag handlers ---

func (h *MCPHandler) HandleListTags(_ context.Context, _ *mcp.CallToolRequest, args ListTagsArgs) (*mcp.CallToolResult, any, error) {
	result, err := h.client.ListTags(args.Limit)
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

// --- Server bootstrap ---

func RunMCPServer() {
	n8nClient, err := NewN8nClient()
	if err != nil {
		log.Fatalf("Failed to create n8n client: %v", err)
	}

	handler := NewMCPHandler(n8nClient)

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "oido-n8n",
		Version: "1.0.0",
	}, nil)

	// Workflows
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_list_workflows",
		Description: "List all workflows. Optionally filter by active status or tags.",
	}, handler.HandleListWorkflows)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_get_workflow",
		Description: "Get full details of a workflow by ID including nodes and connections.",
	}, handler.HandleGetWorkflow)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_create_workflow",
		Description: "Create a new workflow from a JSON definition.",
	}, handler.HandleCreateWorkflow)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_update_workflow",
		Description: "Update an existing workflow by ID with a JSON patch.",
	}, handler.HandleUpdateWorkflow)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_delete_workflow",
		Description: "Permanently delete a workflow by ID.",
	}, handler.HandleDeleteWorkflow)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_activate_workflow",
		Description: "Activate a workflow so it responds to triggers.",
	}, handler.HandleActivateWorkflow)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_deactivate_workflow",
		Description: "Deactivate a workflow to pause trigger responses.",
	}, handler.HandleDeactivateWorkflow)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_execute_workflow",
		Description: "Manually execute a workflow, optionally passing input data.",
	}, handler.HandleExecuteWorkflow)

	// Executions
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_list_executions",
		Description: "List workflow executions. Filter by workflowId or status (waiting/running/success/error/canceled).",
	}, handler.HandleListExecutions)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_get_execution",
		Description: "Get full details of a specific execution by ID.",
	}, handler.HandleGetExecution)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_delete_execution",
		Description: "Delete an execution record by ID.",
	}, handler.HandleDeleteExecution)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_stop_execution",
		Description: "Stop a currently running execution.",
	}, handler.HandleStopExecution)

	// Credentials
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_list_credentials",
		Description: "List all stored credentials (names and types only, no secrets).",
	}, handler.HandleListCredentials)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_create_credential",
		Description: "Create a new credential. Use n8n_get_credential_schema first to know required fields.",
	}, handler.HandleCreateCredential)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_delete_credential",
		Description: "Delete a credential by ID.",
	}, handler.HandleDeleteCredential)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_get_credential_schema",
		Description: "Get the field schema for a credential type (e.g. httpBasicAuth, githubApi).",
	}, handler.HandleGetCredentialSchema)

	// Tags
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_list_tags",
		Description: "List all workflow tags.",
	}, handler.HandleListTags)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_create_tag",
		Description: "Create a new workflow tag.",
	}, handler.HandleCreateTag)

	// Webhooks
	mcp.AddTool(server, &mcp.Tool{
		Name:        "n8n_trigger_webhook",
		Description: "Trigger a workflow via its webhook path. Does not require API key.",
	}, handler.HandleTriggerWebhook)

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
