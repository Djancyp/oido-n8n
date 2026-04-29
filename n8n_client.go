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

type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Role      string `json:"role,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
	UpdatedAt string `json:"updatedAt,omitempty"`
}

type CreateUserRequest struct {
	Email     string `json:"email"`
	Role      string `json:"role,omitempty"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
}

type ListResponse[T any] struct {
	Data       []T    `json:"data"`
	NextCursor string `json:"nextCursor,omitempty"`
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

func (c *N8nClient) ListWorkflows(active *bool, tags []string, name string, limit int, cursor string) (*ListResponse[Workflow], error) {
	params := url.Values{}
	if active != nil {
		params.Set("active", strconv.FormatBool(*active))
	}
	for _, t := range tags {
		params.Add("tags", t)
	}
	if name != "" {
		params.Set("name", name)
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

	data, _, err := c.do("PUT", "/workflows/"+id, body)
	if err != nil {
		return nil, err
	}

	var result Workflow
	return &result, json.Unmarshal(data, &result)
}

func (c *N8nClient) UpdatePartialWorkflow(id string, patches map[string]json.RawMessage) (*Workflow, error) {
	existing, err := c.GetWorkflow(id)
	if err != nil {
		return nil, fmt.Errorf("fetch workflow: %w", err)
	}

	// Build a map from the existing workflow, then apply patches
	raw, err := json.Marshal(existing)
	if err != nil {
		return nil, fmt.Errorf("marshal workflow: %w", err)
	}
	var body map[string]json.RawMessage
	if err := json.Unmarshal(raw, &body); err != nil {
		return nil, fmt.Errorf("unmarshal workflow: %w", err)
	}
	for k, v := range patches {
		body[k] = v
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

// --- Executions ---

func (c *N8nClient) ListExecutions(workflowID, status string, limit int, cursor string, includeData bool) (*ListResponse[Execution], error) {
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
	if includeData {
		params.Set("includeData", "true")
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

func (c *N8nClient) GetExecution(id string, includeData bool) (*Execution, error) {
	path := "/executions/" + id
	if includeData {
		path += "?includeData=true"
	}
	data, _, err := c.do("GET", path, nil)
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

func (c *N8nClient) StopExecutions(workflowID string) (json.RawMessage, error) {
	body := map[string]string{}
	if workflowID != "" {
		body["workflowId"] = workflowID
	}
	data, _, err := c.do("POST", "/executions/stop", body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c *N8nClient) RetryExecution(id string) (json.RawMessage, error) {
	data, _, err := c.do("POST", "/executions/"+id+"/retry", nil)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c *N8nClient) ListExecutionTags(id string) ([]Tag, error) {
	data, _, err := c.do("GET", "/executions/"+id+"/tags", nil)
	if err != nil {
		return nil, err
	}
	var result []Tag
	return result, json.Unmarshal(data, &result)
}

func (c *N8nClient) UpdateExecutionTags(id string, tagIDs []string) ([]Tag, error) {
	body := make([]map[string]string, len(tagIDs))
	for i, tid := range tagIDs {
		body[i] = map[string]string{"id": tid}
	}
	data, _, err := c.do("PUT", "/executions/"+id+"/tags", body)
	if err != nil {
		return nil, err
	}
	var result []Tag
	return result, json.Unmarshal(data, &result)
}

// --- Credentials ---

func (c *N8nClient) ListCredentials(limit int, credentialID, credType string) (*ListResponse[Credential], error) {
	params := url.Values{}
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}
	if credentialID != "" {
		params.Set("credentialId", credentialID)
	}
	if credType != "" {
		params.Set("type", credType)
	}
	path := "/credentials"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}
	data, _, err := c.do("GET", path, nil)
	if err != nil {
		return nil, err
	}
	var result ListResponse[Credential]
	return &result, json.Unmarshal(data, &result)
}

func (c *N8nClient) GetCredential(id string) (*Credential, error) {
	data, _, err := c.do("GET", "/credentials/"+id, nil)
	if err != nil {
		return nil, err
	}
	var result Credential
	return &result, json.Unmarshal(data, &result)
}

func (c *N8nClient) UpdateCredential(id, name, credType, credDataJSON string) (*Credential, error) {
	var credData json.RawMessage
	if credDataJSON != "" {
		if err := json.Unmarshal([]byte(credDataJSON), &credData); err != nil {
			return nil, fmt.Errorf("invalid credential data JSON: %w", err)
		}
	}
	body := map[string]interface{}{"name": name, "type": credType, "data": credData}
	data, _, err := c.do("PUT", "/credentials/"+id, body)
	if err != nil {
		return nil, err
	}
	var result Credential
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

// --- Projects ---

type Project struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	CreatedAt string `json:"createdAt,omitempty"`
	UpdatedAt string `json:"updatedAt,omitempty"`
}

func (c *N8nClient) ListProjects(limit int, cursor string) (*ListResponse[Project], error) {
	params := url.Values{}
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}
	if cursor != "" {
		params.Set("cursor", cursor)
	}
	path := "/projects"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}
	data, _, err := c.do("GET", path, nil)
	if err != nil {
		return nil, err
	}
	var result ListResponse[Project]
	return &result, json.Unmarshal(data, &result)
}

func (c *N8nClient) GetProject(id string) (*Project, error) {
	data, _, err := c.do("GET", "/projects/"+id, nil)
	if err != nil {
		return nil, err
	}
	var result Project
	return &result, json.Unmarshal(data, &result)
}

func (c *N8nClient) CreateProject(name, projectType string) (*Project, error) {
	body := map[string]string{"name": name, "type": projectType}
	data, _, err := c.do("POST", "/projects", body)
	if err != nil {
		return nil, err
	}
	var result Project
	return &result, json.Unmarshal(data, &result)
}

func (c *N8nClient) UpdateProject(id, name string) (*Project, error) {
	body := map[string]string{"name": name}
	data, _, err := c.do("PUT", "/projects/"+id, body)
	if err != nil {
		return nil, err
	}
	var result Project
	return &result, json.Unmarshal(data, &result)
}

func (c *N8nClient) DeleteProject(id string) error {
	_, _, err := c.do("DELETE", "/projects/"+id, nil)
	return err
}

// --- Variables ---

type Variable struct {
	ID    string `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (c *N8nClient) ListVariables(limit int, cursor string) (*ListResponse[Variable], error) {
	params := url.Values{}
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}
	if cursor != "" {
		params.Set("cursor", cursor)
	}
	path := "/variables"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}
	data, _, err := c.do("GET", path, nil)
	if err != nil {
		return nil, err
	}
	var result ListResponse[Variable]
	return &result, json.Unmarshal(data, &result)
}

func (c *N8nClient) GetVariable(id string) (*Variable, error) {
	data, _, err := c.do("GET", "/variables/"+id, nil)
	if err != nil {
		return nil, err
	}
	var result Variable
	return &result, json.Unmarshal(data, &result)
}

func (c *N8nClient) CreateVariable(key, value string) (*Variable, error) {
	body := map[string]string{"key": key, "value": value}
	data, _, err := c.do("POST", "/variables", body)
	if err != nil {
		return nil, err
	}
	var result Variable
	return &result, json.Unmarshal(data, &result)
}

func (c *N8nClient) UpdateVariable(id, key, value string) (*Variable, error) {
	body := map[string]string{}
	if key != "" {
		body["key"] = key
	}
	if value != "" {
		body["value"] = value
	}
	data, _, err := c.do("PATCH", "/variables/"+id, body)
	if err != nil {
		return nil, err
	}
	var result Variable
	return &result, json.Unmarshal(data, &result)
}

func (c *N8nClient) DeleteVariable(id string) error {
	_, _, err := c.do("DELETE", "/variables/"+id, nil)
	return err
}

// --- Audit ---

func (c *N8nClient) GenerateAudit(additionalOptions json.RawMessage) (json.RawMessage, error) {
	var body interface{}
	if additionalOptions != nil {
		body = map[string]json.RawMessage{"additionalOptions": additionalOptions}
	}
	data, _, err := c.do("POST", "/audit", body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// --- Users ---

func (c *N8nClient) ListUsers(limit int, cursor string, includeRole bool) (*ListResponse[User], error) {
	params := url.Values{}
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}
	if cursor != "" {
		params.Set("cursor", cursor)
	}
	if includeRole {
		params.Set("includeRole", "true")
	}
	path := "/users"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}
	data, _, err := c.do("GET", path, nil)
	if err != nil {
		return nil, err
	}
	var result ListResponse[User]
	return &result, json.Unmarshal(data, &result)
}

func (c *N8nClient) GetUser(idOrEmail string, includeRole bool) (*User, error) {
	path := "/users/" + idOrEmail
	if includeRole {
		path += "?includeRole=true"
	}
	data, _, err := c.do("GET", path, nil)
	if err != nil {
		return nil, err
	}
	var result User
	return &result, json.Unmarshal(data, &result)
}

func (c *N8nClient) CreateUsers(users []CreateUserRequest) ([]User, error) {
	data, _, err := c.do("POST", "/users", users)
	if err != nil {
		return nil, err
	}
	var result []User
	return result, json.Unmarshal(data, &result)
}

func (c *N8nClient) ChangeUserRole(id, role string) (*User, error) {
	body := map[string]string{"newRoleName": role}
	data, _, err := c.do("PATCH", "/users/"+id+"/role", body)
	if err != nil {
		return nil, err
	}
	var result User
	return &result, json.Unmarshal(data, &result)
}

func (c *N8nClient) DeleteUser(id string) error {
	_, _, err := c.do("DELETE", "/users/"+id, nil)
	return err
}

// --- Tags ---

func (c *N8nClient) ListTags(limit int, cursor string) (*ListResponse[Tag], error) {
	params := url.Values{}
	if limit > 0 {
		params.Set("limit", strconv.Itoa(limit))
	}
	if cursor != "" {
		params.Set("cursor", cursor)
	}
	path := "/tags"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}
	data, _, err := c.do("GET", path, nil)
	if err != nil {
		return nil, err
	}
	var result ListResponse[Tag]
	return &result, json.Unmarshal(data, &result)
}

func (c *N8nClient) GetTag(id string) (*Tag, error) {
	data, _, err := c.do("GET", "/tags/"+id, nil)
	if err != nil {
		return nil, err
	}
	var result Tag
	return &result, json.Unmarshal(data, &result)
}

func (c *N8nClient) UpdateTag(id, name string) (*Tag, error) {
	body := map[string]string{"name": name}
	data, _, err := c.do("PUT", "/tags/"+id, body)
	if err != nil {
		return nil, err
	}
	var result Tag
	return &result, json.Unmarshal(data, &result)
}

func (c *N8nClient) DeleteTag(id string) error {
	_, _, err := c.do("DELETE", "/tags/"+id, nil)
	return err
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
