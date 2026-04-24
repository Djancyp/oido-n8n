package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

// N8nClient wraps the n8n REST API v1.
type N8nClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewN8nClient creates a client from environment variables.
func NewN8nClient() (*N8nClient, error) {
	apiURL := os.Getenv("N8N_API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:5678"
	}

	apiKey := os.Getenv("N8N_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("N8N_API_KEY environment variable is required")
	}

	return &N8nClient{
		baseURL:    apiURL + "/api/v1",
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}, nil
}

// --- Types ---

type Workflow struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Active      bool            `json:"active"`
	CreatedAt   string          `json:"createdAt"`
	UpdatedAt   string          `json:"updatedAt"`
	Nodes       json.RawMessage `json:"nodes,omitempty"`
	Connections json.RawMessage `json:"connections,omitempty"`
	Settings    json.RawMessage `json:"settings,omitempty"`
	Tags        []Tag           `json:"tags,omitempty"`
}

type Execution struct {
	ID         string `json:"id"`
	Finished   bool   `json:"finished"`
	Mode       string `json:"mode"`
	StartedAt  string `json:"startedAt"`
	StoppedAt  string `json:"stoppedAt,omitempty"`
	WorkflowID string `json:"workflowId"`
	Status     string `json:"status"`
}

type Credential struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type Tag struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ListResponse[T any] struct {
	Data       []T    `json:"data"`
	NextCursor string `json:"nextCursor,omitempty"`
}

type ExecuteResponse struct {
	ExecutionID int `json:"executionId"`
}

// --- HTTP helpers ---

func (c *N8nClient) do(method, path string, body interface{}) ([]byte, int, error) {
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, 0, fmt.Errorf("marshal request: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, c.baseURL+path, reqBody)
	if err != nil {
		return nil, 0, fmt.Errorf("build request: %w", err)
	}

	req.Header.Set("X-N8N-API-KEY", c.apiKey)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, resp.StatusCode, fmt.Errorf("n8n API error %d: %s", resp.StatusCode, string(data))
	}

	return data, resp.StatusCode, nil
}

func (c *N8nClient) doWebhook(baseURL, path, method string, body interface{}) ([]byte, int, error) {
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, 0, fmt.Errorf("marshal request: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}

	webhookURL := baseURL + "/webhook/" + path
	req, err := http.NewRequest(method, webhookURL, reqBody)
	if err != nil {
		return nil, 0, fmt.Errorf("build request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("webhook request failed: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	return data, resp.StatusCode, err
}

// --- Workflows ---

func (c *N8nClient) ListWorkflows(active *bool, tags []string, limit int, cursor string) (*ListResponse[Workflow], error) {
	params := url.Values{}
	if active != nil {
		params.Set("active", strconv.FormatBool(*active))
	}
	for _, t := range tags {
		params.Add("tags", t)
	}
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}
	if cursor != "" {
		params.Set("cursor", cursor)
	}

	path := "/workflows"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	data, _, err := c.do("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result ListResponse[Workflow]
	return &result, json.Unmarshal(data, &result)
}

func (c *N8nClient) GetWorkflow(id string) (*Workflow, error) {
	data, _, err := c.do("GET", "/workflows/"+id, nil)
	if err != nil {
		return nil, err
	}

	var result Workflow
	return &result, json.Unmarshal(data, &result)
}

func (c *N8nClient) CreateWorkflow(workflowJSON string) (*Workflow, error) {
	var body json.RawMessage
	if err := json.Unmarshal([]byte(workflowJSON), &body); err != nil {
		return nil, fmt.Errorf("invalid workflow JSON: %w", err)
	}

	data, _, err := c.do("POST", "/workflows", body)
	if err != nil {
		return nil, err
	}

	var result Workflow
	return &result, json.Unmarshal(data, &result)
}

func (c *N8nClient) UpdateWorkflow(id, workflowJSON string) (*Workflow, error) {
	var body json.RawMessage
	if err := json.Unmarshal([]byte(workflowJSON), &body); err != nil {
		return nil, fmt.Errorf("invalid workflow JSON: %w", err)
	}

	data, _, err := c.do("PATCH", "/workflows/"+id, body)
	if err != nil {
		return nil, err
	}

	var result Workflow
	return &result, json.Unmarshal(data, &result)
}

func (c *N8nClient) DeleteWorkflow(id string) error {
	_, _, err := c.do("DELETE", "/workflows/"+id, nil)
	return err
}

func (c *N8nClient) ActivateWorkflow(id string) (*Workflow, error) {
	data, _, err := c.do("POST", "/workflows/"+id+"/activate", nil)
	if err != nil {
		return nil, err
	}

	var result Workflow
	return &result, json.Unmarshal(data, &result)
}

func (c *N8nClient) DeactivateWorkflow(id string) (*Workflow, error) {
	data, _, err := c.do("POST", "/workflows/"+id+"/deactivate", nil)
	if err != nil {
		return nil, err
	}

	var result Workflow
	return &result, json.Unmarshal(data, &result)
}

func (c *N8nClient) ExecuteWorkflow(id string, dataJSON string) (*ExecuteResponse, error) {
	var body interface{}
	if dataJSON != "" {
		var raw json.RawMessage
		if err := json.Unmarshal([]byte(dataJSON), &raw); err != nil {
			return nil, fmt.Errorf("invalid data JSON: %w", err)
		}
		body = raw
	}

	data, _, err := c.do("POST", "/workflows/"+id+"/run", body)
	if err != nil {
		return nil, err
	}

	var result ExecuteResponse
	return &result, json.Unmarshal(data, &result)
}

// --- Executions ---

func (c *N8nClient) ListExecutions(workflowID, status string, limit int, cursor string) (*ListResponse[Execution], error) {
	params := url.Values{}
	if workflowID != "" {
		params.Set("workflowId", workflowID)
	}
	if status != "" {
		params.Set("status", status)
	}
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}
	if cursor != "" {
		params.Set("cursor", cursor)
	}

	path := "/executions"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	data, _, err := c.do("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result ListResponse[Execution]
	return &result, json.Unmarshal(data, &result)
}

func (c *N8nClient) GetExecution(id string) (*Execution, error) {
	data, _, err := c.do("GET", "/executions/"+id, nil)
	if err != nil {
		return nil, err
	}

	var result Execution
	return &result, json.Unmarshal(data, &result)
}

func (c *N8nClient) DeleteExecution(id string) error {
	_, _, err := c.do("DELETE", "/executions/"+id, nil)
	return err
}

func (c *N8nClient) StopExecution(id string) error {
	_, _, err := c.do("POST", "/executions/"+id+"/stop", nil)
	return err
}

// --- Credentials ---

func (c *N8nClient) ListCredentials(limit int) (*ListResponse[Credential], error) {
	path := "/credentials"
	if limit > 0 {
		path += "?limit=" + strconv.Itoa(limit)
	}

	data, _, err := c.do("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result ListResponse[Credential]
	return &result, json.Unmarshal(data, &result)
}

func (c *N8nClient) CreateCredential(name, credType, credDataJSON string) (*Credential, error) {
	var credData json.RawMessage
	if credDataJSON != "" {
		if err := json.Unmarshal([]byte(credDataJSON), &credData); err != nil {
			return nil, fmt.Errorf("invalid credential data JSON: %w", err)
		}
	}

	body := map[string]interface{}{
		"name": name,
		"type": credType,
		"data": credData,
	}

	data, _, err := c.do("POST", "/credentials", body)
	if err != nil {
		return nil, err
	}

	var result Credential
	return &result, json.Unmarshal(data, &result)
}

func (c *N8nClient) DeleteCredential(id string) error {
	_, _, err := c.do("DELETE", "/credentials/"+id, nil)
	return err
}

func (c *N8nClient) GetCredentialSchema(credentialType string) (json.RawMessage, error) {
	data, _, err := c.do("GET", "/credentials/schema/"+credentialType, nil)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// --- Tags ---

func (c *N8nClient) ListTags(limit int) (*ListResponse[Tag], error) {
	path := "/tags"
	if limit > 0 {
		path += "?limit=" + strconv.Itoa(limit)
	}

	data, _, err := c.do("GET", path, nil)
	if err != nil {
		return nil, err
	}

	var result ListResponse[Tag]
	return &result, json.Unmarshal(data, &result)
}

func (c *N8nClient) CreateTag(name string) (*Tag, error) {
	body := map[string]string{"name": name}

	data, _, err := c.do("POST", "/tags", body)
	if err != nil {
		return nil, err
	}

	var result Tag
	return &result, json.Unmarshal(data, &result)
}

// --- Webhooks ---

func (c *N8nClient) TriggerWebhook(path, method, bodyJSON string) (string, int, error) {
	n8nBaseURL := os.Getenv("N8N_API_URL")
	if n8nBaseURL == "" {
		n8nBaseURL = "http://localhost:5678"
	}

	var body interface{}
	if bodyJSON != "" {
		var raw json.RawMessage
		if err := json.Unmarshal([]byte(bodyJSON), &raw); err != nil {
			return "", 0, fmt.Errorf("invalid body JSON: %w", err)
		}
		body = raw
	}

	if method == "" {
		method = "POST"
	}

	data, status, err := c.doWebhook(n8nBaseURL, path, method, body)
	return string(data), status, err
}
