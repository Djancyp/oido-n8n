package main

import (
	"context"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// newTestHandler spins up the embedded node DB so validation exercises the
// real registry path. Skips if the DB cannot be initialized.
func newTestHandler(t *testing.T) (*MCPHandler, func()) {
	t.Helper()
	db, cleanup, err := InitNodeDB()
	if err != nil {
		t.Fatalf("InitNodeDB: %v", err)
	}
	return &MCPHandler{nodeDB: db}, cleanup
}

func validate(t *testing.T, h *MCPHandler, wf string) string {
	t.Helper()
	res, _, err := h.HandleValidateWorkflow(context.Background(), nil, ValidateWorkflowArgs{WorkflowJSON: wf})
	if err != nil {
		t.Fatalf("HandleValidateWorkflow returned err: %v", err)
	}
	if len(res.Content) == 0 {
		t.Fatal("empty validation result")
	}
	tc, ok := res.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatalf("expected TextContent, got %T", res.Content[0])
	}
	return tc.Text
}

func TestValidate_DetectsBrokenConnectionTarget(t *testing.T) {
	h, cleanup := newTestHandler(t)
	defer cleanup()

	// Connection points at "Slak" but the node is named "Slack".
	wf := `{
	  "name": "t",
	  "nodes": [
	    {"id":"1","name":"Trigger","type":"n8n-nodes-base.manualTrigger","typeVersion":1,"position":[0,0]},
	    {"id":"2","name":"Slack","type":"n8n-nodes-base.slack","typeVersion":1,"position":[1,0]}
	  ],
	  "connections": {"Trigger": {"main": [[{"node":"Slak","type":"main","index":0}]]}}
	}`
	out := validate(t, h, wf)
	if !strings.Contains(out, `"Slak"`) || !strings.Contains(out, "not defined") {
		t.Fatalf("expected broken-target error, got:\n%s", out)
	}
}

func TestValidate_DetectsUnknownSourceNode(t *testing.T) {
	h, cleanup := newTestHandler(t)
	defer cleanup()

	wf := `{
	  "name": "t",
	  "nodes": [
	    {"id":"1","name":"Trigger","type":"n8n-nodes-base.manualTrigger","typeVersion":1,"position":[0,0]}
	  ],
	  "connections": {"Ghost": {"main": [[]]}}
	}`
	out := validate(t, h, wf)
	if !strings.Contains(out, "Ghost") {
		t.Fatalf("expected unknown-source error, got:\n%s", out)
	}
}

func TestValidate_UnknownNodeType(t *testing.T) {
	h, cleanup := newTestHandler(t)
	defer cleanup()

	wf := `{
	  "name": "t",
	  "nodes": [
	    {"id":"1","name":"X","type":"n8n-nodes-base.totallyFakeNode","typeVersion":1,"position":[0,0]}
	  ],
	  "connections": {}
	}`
	out := validate(t, h, wf)
	if !strings.Contains(out, "not found in node registry") {
		t.Fatalf("expected unknown-type error, got:\n%s", out)
	}
}

func TestValidate_CleanWorkflowPasses(t *testing.T) {
	h, cleanup := newTestHandler(t)
	defer cleanup()

	wf := `{
	  "name": "ok",
	  "nodes": [
	    {"id":"1","name":"Trigger","type":"n8n-nodes-base.manualTrigger","typeVersion":1,"position":[0,0]},
	    {"id":"2","name":"NoOp","type":"n8n-nodes-base.noOp","typeVersion":1,"position":[1,0]}
	  ],
	  "connections": {"Trigger": {"main": [[{"node":"NoOp","type":"main","index":0}]]}}
	}`
	out := validate(t, h, wf)
	if strings.Contains(out, "ERROR") {
		t.Fatalf("expected no errors, got:\n%s", out)
	}
}
