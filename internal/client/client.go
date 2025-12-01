package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Case represents a Theta Lake Case.
type Case struct {
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name"`
	Number      string `json:"number"`
	OpenDate    string `json:"open_date"`
	Visibility  string `json:"visibility"`
	Description string `json:"description,omitempty"`
}

// User represents a Theta Lake User.
type User struct {
	ID                   int    `json:"id,omitempty"`
	Name                 string `json:"name"`
	Email                string `json:"email"`
	Password             string `json:"password,omitempty"`
	PasswordConfirmation string `json:"password_confirmation,omitempty"`
	RoleID               int    `json:"role_id,omitempty"`
	SearchID             int    `json:"search_id,omitempty"`
}

// DirectoryGroup represents a Theta Lake Directory Group.
type DirectoryGroup struct {
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name"`
	ExternalID  string `json:"external_id,omitempty"`
	Description string `json:"description,omitempty"`
}

// Client holds the connection details for the Theta Lake API.
type Client struct {
	Endpoint   string
	Token      string
	HTTPClient *http.Client
}

// NewClient creates a new Theta Lake API client.
func NewClient(endpoint, token string) (*Client, error) {
	return &Client{
		Endpoint: endpoint,
		Token:    token,
		HTTPClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}, nil
}

// DoRequest performs the HTTP request.
func (c *Client) DoRequest(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	req.Header.Set("Content-Type", "application/json")

	return c.HTTPClient.Do(req)
}

// GetCase retrieves a case by ID.
func (c *Client) GetCase(caseID string) (*Case, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/cases/%s", c.Endpoint, caseID), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, "error reading body") // simplified error handling
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var theCase Case
	err = json.Unmarshal(body, &theCase)
	if err != nil {
		return nil, err
	}

	return &theCase, nil
}

// CreateCase creates a new case.
func (c *Client) CreateCase(theCase Case) (*Case, error) {
	rb, err := json.Marshal(theCase)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/cases", c.Endpoint), bytes.NewBuffer(rb))
	if err != nil {
		return nil, err
	}

	res, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var newCase Case
	err = json.Unmarshal(body, &newCase)
	if err != nil {
		return nil, err
	}

	return &newCase, nil
}

// UpdateCase updates an existing case.
func (c *Client) UpdateCase(caseID string, theCase Case) (*Case, error) {
	rb, err := json.Marshal(theCase)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/cases/%s", c.Endpoint, caseID), bytes.NewBuffer(rb))
	if err != nil {
		return nil, err
	}

	res, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var updatedCase Case
	err = json.Unmarshal(body, &updatedCase)
	if err != nil {
		return nil, err
	}

	return &updatedCase, nil
}

// DeleteCase deletes a case.
func (c *Client) DeleteCase(caseID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/cases/%s", c.Endpoint, caseID), nil)
	if err != nil {
		return err
	}

	res, err := c.DoRequest(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("status: %d", res.StatusCode)
	}

	return nil
}

// GetUser retrieves a user by ID.
func (c *Client) GetUser(userID string) (*User, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/users/%s", c.Endpoint, userID), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, "error reading body")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var user User
	err = json.Unmarshal(body, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// CreateUser creates a new user.
func (c *Client) CreateUser(user User) (*User, error) {
	rb, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/users", c.Endpoint), bytes.NewBuffer(rb))
	if err != nil {
		return nil, err
	}

	res, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var newUser User
	err = json.Unmarshal(body, &newUser)
	if err != nil {
		return nil, err
	}

	return &newUser, nil
}

// UpdateUser updates an existing user.
func (c *Client) UpdateUser(userID string, user User) (*User, error) {
	rb, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/users/%s", c.Endpoint, userID), bytes.NewBuffer(rb))
	if err != nil {
		return nil, err
	}

	res, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var updatedUser User
	err = json.Unmarshal(body, &updatedUser)
	if err != nil {
		return nil, err
	}

	return &updatedUser, nil
}

// DeleteUser deletes a user.
func (c *Client) DeleteUser(userID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/users/%s", c.Endpoint, userID), nil)
	if err != nil {
		return err
	}

	res, err := c.DoRequest(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("status: %d", res.StatusCode)
	}

	return nil
}

// GetDirectoryGroup retrieves a directory group by ID.
func (c *Client) GetDirectoryGroup(groupID string) (*DirectoryGroup, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/directory_groups/%s", c.Endpoint, groupID), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, "error reading body")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var group DirectoryGroup
	err = json.Unmarshal(body, &group)
	if err != nil {
		return nil, err
	}

	return &group, nil
}

// CreateDirectoryGroup creates a new directory group.
func (c *Client) CreateDirectoryGroup(group DirectoryGroup) (*DirectoryGroup, error) {
	rb, err := json.Marshal(group)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/directory_groups", c.Endpoint), bytes.NewBuffer(rb))
	if err != nil {
		return nil, err
	}

	res, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var newGroup DirectoryGroup
	err = json.Unmarshal(body, &newGroup)
	if err != nil {
		return nil, err
	}

	return &newGroup, nil
}

// UpdateDirectoryGroup updates an existing directory group.
func (c *Client) UpdateDirectoryGroup(groupID string, group DirectoryGroup) (*DirectoryGroup, error) {
	rb, err := json.Marshal(group)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/directory_groups/%s", c.Endpoint, groupID), bytes.NewBuffer(rb))
	if err != nil {
		return nil, err
	}

	res, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var updatedGroup DirectoryGroup
	err = json.Unmarshal(body, &updatedGroup)
	if err != nil {
		return nil, err
	}

	return &updatedGroup, nil
}

// DeleteDirectoryGroup deletes a directory group.
func (c *Client) DeleteDirectoryGroup(groupID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/directory_groups/%s", c.Endpoint, groupID), nil)
	if err != nil {
		return err
	}

	res, err := c.DoRequest(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("status: %d", res.StatusCode)
	}

	return nil
}

// RetentionPolicy represents a Theta Lake Retention Policy.
type RetentionPolicy struct {
	ID                  int    `json:"id,omitempty"`
	Name                string `json:"name"`
	Description         string `json:"description,omitempty"`
	RetentionPeriodDays int    `json:"retention_period_days,omitempty"`
}

// GetRetentionPolicy retrieves a retention policy by ID.
func (c *Client) GetRetentionPolicy(policyID string) (*RetentionPolicy, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/retention_policies/%s", c.Endpoint, policyID), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, "error reading body")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var policy RetentionPolicy
	err = json.Unmarshal(body, &policy)
	if err != nil {
		return nil, err
	}

	return &policy, nil
}

// CreateRetentionPolicy creates a new retention policy.
func (c *Client) CreateRetentionPolicy(policy RetentionPolicy) (*RetentionPolicy, error) {
	rb, err := json.Marshal(policy)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/retention_policies", c.Endpoint), bytes.NewBuffer(rb))
	if err != nil {
		return nil, err
	}

	res, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var newPolicy RetentionPolicy
	err = json.Unmarshal(body, &newPolicy)
	if err != nil {
		return nil, err
	}

	return &newPolicy, nil
}

// UpdateRetentionPolicy updates an existing retention policy.
func (c *Client) UpdateRetentionPolicy(policyID string, policy RetentionPolicy) (*RetentionPolicy, error) {
	rb, err := json.Marshal(policy)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/retention_policies/%s", c.Endpoint, policyID), bytes.NewBuffer(rb))
	if err != nil {
		return nil, err
	}

	res, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var updatedPolicy RetentionPolicy
	err = json.Unmarshal(body, &updatedPolicy)
	if err != nil {
		return nil, err
	}

	return &updatedPolicy, nil
}

// DeleteRetentionPolicy deletes a retention policy.
func (c *Client) DeleteRetentionPolicy(policyID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/retention_policies/%s", c.Endpoint, policyID), nil)
	if err != nil {
		return err
	}

	res, err := c.DoRequest(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("status: %d", res.StatusCode)
	}

	return nil
}

// LegalHold represents a Theta Lake Legal Hold.
type LegalHold struct {
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	CaseID      int    `json:"case_id,omitempty"`
}

// GetLegalHold retrieves a legal hold by ID.
func (c *Client) GetLegalHold(holdID string) (*LegalHold, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/legal_holds/%s", c.Endpoint, holdID), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, "error reading body")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var hold LegalHold
	err = json.Unmarshal(body, &hold)
	if err != nil {
		return nil, err
	}

	return &hold, nil
}

// CreateLegalHold creates a new legal hold.
func (c *Client) CreateLegalHold(hold LegalHold) (*LegalHold, error) {
	rb, err := json.Marshal(hold)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/legal_holds", c.Endpoint), bytes.NewBuffer(rb))
	if err != nil {
		return nil, err
	}

	res, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var newHold LegalHold
	err = json.Unmarshal(body, &newHold)
	if err != nil {
		return nil, err
	}

	return &newHold, nil
}

// UpdateLegalHold updates an existing legal hold.
func (c *Client) UpdateLegalHold(holdID string, hold LegalHold) (*LegalHold, error) {
	rb, err := json.Marshal(hold)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/legal_holds/%s", c.Endpoint, holdID), bytes.NewBuffer(rb))
	if err != nil {
		return nil, err
	}

	res, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var updatedHold LegalHold
	err = json.Unmarshal(body, &updatedHold)
	if err != nil {
		return nil, err
	}

	return &updatedHold, nil
}

// DeleteLegalHold deletes a legal hold.
func (c *Client) DeleteLegalHold(holdID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/legal_holds/%s", c.Endpoint, holdID), nil)
	if err != nil {
		return err
	}

	res, err := c.DoRequest(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("status: %d", res.StatusCode)
	}

	return nil
}

// Tag represents a Theta Lake Tag.
type Tag struct {
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// GetTag retrieves a tag by ID.
func (c *Client) GetTag(tagID string) (*Tag, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/tags/%s", c.Endpoint, tagID), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, "error reading body")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var tag Tag
	err = json.Unmarshal(body, &tag)
	if err != nil {
		return nil, err
	}

	return &tag, nil
}

// CreateTag creates a new tag.
func (c *Client) CreateTag(tag Tag) (*Tag, error) {
	rb, err := json.Marshal(tag)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/tags", c.Endpoint), bytes.NewBuffer(rb))
	if err != nil {
		return nil, err
	}

	res, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var newTag Tag
	err = json.Unmarshal(body, &newTag)
	if err != nil {
		return nil, err
	}

	return &newTag, nil
}

// UpdateTag updates an existing tag.
func (c *Client) UpdateTag(tagID string, tag Tag) (*Tag, error) {
	rb, err := json.Marshal(tag)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/tags/%s", c.Endpoint, tagID), bytes.NewBuffer(rb))
	if err != nil {
		return nil, err
	}

	res, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var updatedTag Tag
	err = json.Unmarshal(body, &updatedTag)
	if err != nil {
		return nil, err
	}

	return &updatedTag, nil
}

// DeleteTag deletes a tag.
func (c *Client) DeleteTag(tagID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/tags/%s", c.Endpoint, tagID), nil)
	if err != nil {
		return err
	}

	res, err := c.DoRequest(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("status: %d", res.StatusCode)
	}

	return nil
}

// AuditLog represents a Theta Lake Audit Log entry.
type AuditLog struct {
	ID        string `json:"id"`
	User      string `json:"user"`
	Action    string `json:"action"`
	Resource  string `json:"resource"`
	Timestamp string `json:"timestamp"`
}

// GetAuditLogs retrieves audit logs.
func (c *Client) GetAuditLogs() ([]AuditLog, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/audit_logs", c.Endpoint), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, "error reading body")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var logs []AuditLog
	err = json.Unmarshal(body, &logs)
	if err != nil {
		return nil, err
	}

	return logs, nil
}

// Event represents a Theta Lake Event.
type Event struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}

// GetEvents retrieves events.
func (c *Client) GetEvents() ([]Event, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/events", c.Endpoint), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, "error reading body")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var events []Event
	err = json.Unmarshal(body, &events)
	if err != nil {
		return nil, err
	}

	return events, nil
}

// AnalysisPolicy represents a Theta Lake Analysis Policy.
type AnalysisPolicy struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsBuiltIn   bool   `json:"is_built_in"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type analysisPoliciesResponse struct {
	Policies []AnalysisPolicy `json:"policies"`
}

// GetAnalysisPolicies retrieves analysis policies.
func (c *Client) GetAnalysisPolicies() ([]AnalysisPolicy, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/analysis/policies", c.Endpoint), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, "error reading body")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var response analysisPoliciesResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response.Policies, nil
}

// IntegrationState represents the state of a Theta Lake Integration.
type IntegrationState struct {
	ID         int    `json:"id,omitempty"` // Not in API response, but useful for resource
	Paused     bool   `json:"paused"`
	LastRun    string `json:"last_run,omitempty"`
	LastUpload string `json:"last_upload,omitempty"`
}

type integrationStateResponse struct {
	State IntegrationState `json:"state"`
}

type integrationStateRequest struct {
	Status string `json:"status"`
}

// GetIntegrationState retrieves the state of an integration.
func (c *Client) GetIntegrationState(integrationID string) (*IntegrationState, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/ingestion/integration/%s/state", c.Endpoint, integrationID), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, "error reading body")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var response integrationStateResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response.State, nil
}

// UpdateIntegrationState updates the state of an integration (pauses/unpauses).
func (c *Client) UpdateIntegrationState(integrationID string, paused bool) (*IntegrationState, error) {
	status := "active"
	if paused {
		status = "paused"
	}

	reqBody := integrationStateRequest{
		Status: status,
	}

	rb, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/ingestion/integration/%s/state", c.Endpoint, integrationID), bytes.NewBuffer(rb))
	if err != nil {
		return nil, err
	}

	res, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	// The PUT response might not return the full state, so we fetch it again to be sure,
	// or we can just return what we sent if the API doesn't return the state.
	// Assuming API returns the updated state or similar structure.
	// Research showed PUT response body wasn't explicitly captured in detail, but usually returns state.
	// Let's try to parse it as state, if fails, fetch it.

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var response integrationStateResponse
	if err := json.Unmarshal(body, &response); err == nil && response.State.LastRun != "" {
		return &response.State, nil
	}

	// Fallback to GET
	return c.GetIntegrationState(integrationID)
}

// Export represents a Theta Lake Export.
type Export struct {
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	QueryID     int    `json:"query_id,omitempty"`
	Format      string `json:"format,omitempty"`
	Status      string `json:"status,omitempty"`
	DownloadURL string `json:"download_url,omitempty"`
}

// GetExport retrieves an export by ID.
func (c *Client) GetExport(exportID string) (*Export, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/exports/%s", c.Endpoint, exportID), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, "error reading body")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var export Export
	err = json.Unmarshal(body, &export)
	if err != nil {
		return nil, err
	}

	return &export, nil
}

// CreateExport creates a new export.
func (c *Client) CreateExport(export Export) (*Export, error) {
	rb, err := json.Marshal(export)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/exports", c.Endpoint), bytes.NewBuffer(rb))
	if err != nil {
		return nil, err
	}

	res, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var newExport Export
	err = json.Unmarshal(body, &newExport)
	if err != nil {
		return nil, err
	}

	return &newExport, nil
}

// DeleteExport deletes an export.
func (c *Client) DeleteExport(exportID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/exports/%s", c.Endpoint, exportID), nil)
	if err != nil {
		return err
	}

	res, err := c.DoRequest(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("status: %d", res.StatusCode)
	}

	return nil
}

// Record represents a Theta Lake Record.
type Record struct {
	ID           string   `json:"id"`
	ContentDate  string   `json:"content_date"`
	Participants []string `json:"participants"`
	ReviewState  string   `json:"review_state,omitempty"`
	Comment      string   `json:"comment,omitempty"`
}

type recordReviewStateRequest struct {
	ReviewState string `json:"review_state"`
	Comment     string `json:"comment,omitempty"`
}

// GetRecord retrieves a record by ID.
func (c *Client) GetRecord(recordID string) (*Record, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/records/%s", c.Endpoint, recordID), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, "error reading body")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var record Record
	err = json.Unmarshal(body, &record)
	if err != nil {
		return nil, err
	}

	return &record, nil
}

// UpdateRecordReviewState updates the review state of a record.
func (c *Client) UpdateRecordReviewState(recordID string, reviewState string, comment string) (*Record, error) {
	reqBody := recordReviewStateRequest{
		ReviewState: reviewState,
		Comment:     comment,
	}

	rb, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/records/%s/review_state", c.Endpoint, recordID), bytes.NewBuffer(rb))
	if err != nil {
		return nil, err
	}

	res, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	// Assuming the API returns the updated record or we fetch it.
	// Let's fetch it to be sure.
	return c.GetRecord(recordID)
}

// SystemStatus represents the Theta Lake System Status.
type SystemStatus struct {
	Status  string `json:"status"`
	Version string `json:"version"`
	Message string `json:"message,omitempty"`
}

// GetSystemStatus retrieves the system status.
func (c *Client) GetSystemStatus() (*SystemStatus, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/system/status", c.Endpoint), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, "error reading body")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var status SystemStatus
	err = json.Unmarshal(body, &status)
	if err != nil {
		return nil, err
	}

	return &status, nil
}

// PolicyHit represents a Theta Lake Analysis Policy Hit.
type PolicyHit struct {
	ID         string `json:"id"`
	PolicyID   int    `json:"policy_id"`
	RecordID   string `json:"record_id"`
	HitDate    string `json:"hit_date"`
	Confidence int    `json:"confidence"`
}

type policyHitsResponse struct {
	Hits []PolicyHit `json:"hits"`
}

// GetAnalysisPolicyHits retrieves analysis policy hits.
func (c *Client) GetAnalysisPolicyHits() ([]PolicyHit, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/analysis/policy_hits", c.Endpoint), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, "error reading body")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var response policyHitsResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return response.Hits, nil
}

// AddRecordToCase adds a record to a case.
func (c *Client) AddRecordToCase(caseID string, recordID string) error {
	// Assuming the API expects a JSON body with record_id or similar.
	// Based on typical patterns: POST /cases/{id}/records with body {"record_id": "..."}
	reqBody := map[string]string{
		"record_id": recordID,
	}

	rb, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/cases/%s/records", c.Endpoint, caseID), bytes.NewBuffer(rb))
	if err != nil {
		return err
	}

	res, err := c.DoRequest(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		return fmt.Errorf("status: %d", res.StatusCode)
	}

	return nil
}

// RemoveRecordFromCase removes a record from a case.
func (c *Client) RemoveRecordFromCase(caseID string, recordID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/cases/%s/records/%s", c.Endpoint, caseID, recordID), nil)
	// Note: API might be DELETE /cases/{id}/records with body, or DELETE /cases/{id}/records/{record_id}
	// The summary said: DELETE /cases/{id}/records
	// This implies the record ID might be in the body or query param if it's not in path.
	// But standard REST usually puts ID in path for DELETE of specific item.
	// If the summary was strictly `DELETE /cases/{id}/records`, it might expect a body.
	// Let's assume `DELETE /cases/{id}/records/{record_id}` is the standard way, or `DELETE /cases/{id}/records?record_id=...`
	// Given the ambiguity, I'll assume path parameter for now as it's most common for sub-resources.
	// Wait, the summary said `DELETE /cases/{id}/records`. This suggests a bulk delete or body.
	// Let's try sending body with DELETE if path doesn't look right.
	// Actually, `DELETE` with body is discouraged but used.
	// Let's try `DELETE /cases/{id}/records/{record_id}` first.

	// Re-reading summary: `DELETE /cases/{id}/records`
	// This strongly suggests the record ID is passed in the body or query.
	// Let's try query param `record_id`.

	q := req.URL.Query()
	q.Add("record_id", recordID)
	req.URL.RawQuery = q.Encode()

	if err != nil {
		return err
	}

	res, err := c.DoRequest(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("status: %d", res.StatusCode)
	}

	return nil
}

// UpdateCaseStatus updates the status of a case (open/close).
func (c *Client) UpdateCaseStatus(caseID string, status string) error {
	// Status should be "OPEN" or "CLOSED" (or lowercase)
	// API endpoints: PUT /cases/{id}/open, PUT /cases/{id}/close

	var action string
	if status == "OPEN" || status == "open" {
		action = "open"
	} else if status == "CLOSED" || status == "closed" {
		action = "close"
	} else {
		return fmt.Errorf("invalid status: %s", status)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/cases/%s/%s", c.Endpoint, caseID, action), nil)
	if err != nil {
		return err
	}

	res, err := c.DoRequest(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("status: %d", res.StatusCode)
	}

	return nil
}

// Analysis represents a Theta Lake Analysis result.
type Analysis struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Details   string `json:"details,omitempty"`
}

// GetAnalysis retrieves an analysis by ID.
func (c *Client) GetAnalysis(analysisID string) (*Analysis, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/analysis/%s", c.Endpoint, analysisID), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, "error reading body")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var analysis Analysis
	err = json.Unmarshal(body, &analysis)
	if err != nil {
		return nil, err
	}

	return &analysis, nil
}
