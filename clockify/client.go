package clockify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const baseURL = "https://api.clockify.me/api/v1"
const reportsURL = "https://reports.api.clockify.me/v1"

type Client struct {
	apiKey     string
	httpClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) do(method, endpoint string, body any, result any) error {
	return c.doWithBase(baseURL, method, endpoint, body, result)
}

func (c *Client) doReports(method, endpoint string, body any, result any) error {
	return c.doWithBase(reportsURL, method, endpoint, body, result)
}

func (c *Client) doWithBase(base, method, endpoint string, body any, result any) error {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, base+endpoint, bodyReader)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("X-Api-Key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode == 429 {
		return fmt.Errorf("clockify rate limit exceeded, try again later")
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("clockify API error (%d): %s", resp.StatusCode, string(respBody))
	}

	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("unmarshal response: %w", err)
		}
	}

	return nil
}

// --- Workspace ---

type Workspace struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (c *Client) GetWorkspaces() ([]Workspace, error) {
	var result []Workspace
	err := c.do("GET", "/workspaces", nil, &result)
	return result, err
}

// --- User ---

type User struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	ActiveWorkspace string `json:"activeWorkspace"`
}

func (c *Client) GetCurrentUser() (*User, error) {
	var result User
	err := c.do("GET", "/user", nil, &result)
	return &result, err
}

func (c *Client) GetWorkspaceUsers(workspaceID string) ([]User, error) {
	var result []User
	err := c.do("GET", fmt.Sprintf("/workspaces/%s/users", workspaceID), nil, &result)
	return result, err
}

// --- Project ---

type Project struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	ClientID string `json:"clientId,omitempty"`
	Billable bool   `json:"billable"`
	Color    string `json:"color,omitempty"`
	Archived bool   `json:"archived"`
}

type CreateProjectRequest struct {
	Name     string `json:"name"`
	ClientID string `json:"clientId,omitempty"`
	Billable bool   `json:"billable"`
	Color    string `json:"color,omitempty"`
	IsPublic bool   `json:"isPublic"`
}

type UpdateProjectRequest struct {
	Name     string `json:"name,omitempty"`
	ClientID string `json:"clientId,omitempty"`
	Billable *bool  `json:"billable,omitempty"`
	Color    string `json:"color,omitempty"`
	Archived *bool  `json:"archived,omitempty"`
}

func (c *Client) GetProjects(workspaceID string, archived bool, page, pageSize int) ([]Project, error) {
	var result []Project
	q := url.Values{}
	if archived {
		q.Set("archived", "true")
	}
	q.Set("page", fmt.Sprintf("%d", page))
	q.Set("page-size", fmt.Sprintf("%d", pageSize))
	endpoint := fmt.Sprintf("/workspaces/%s/projects?%s", workspaceID, q.Encode())
	err := c.do("GET", endpoint, nil, &result)
	return result, err
}

func (c *Client) CreateProject(workspaceID string, req CreateProjectRequest) (*Project, error) {
	var result Project
	err := c.do("POST", fmt.Sprintf("/workspaces/%s/projects", workspaceID), req, &result)
	return &result, err
}

func (c *Client) UpdateProject(workspaceID, projectID string, req UpdateProjectRequest) (*Project, error) {
	var result Project
	err := c.do("PUT", fmt.Sprintf("/workspaces/%s/projects/%s", workspaceID, projectID), req, &result)
	return &result, err
}

func (c *Client) DeleteProject(workspaceID, projectID string) error {
	return c.do("DELETE", fmt.Sprintf("/workspaces/%s/projects/%s", workspaceID, projectID), nil, nil)
}

// --- Task ---

type Task struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	ProjectID string `json:"projectId"`
	Billable  bool   `json:"billable"`
	Status    string `json:"status"`
}

type CreateTaskRequest struct {
	Name     string `json:"name"`
	Billable bool   `json:"billable"`
}

type UpdateTaskRequest struct {
	Name     string `json:"name,omitempty"`
	Billable *bool  `json:"billable,omitempty"`
	Status   string `json:"status,omitempty"`
}

func (c *Client) GetTasks(workspaceID, projectID string) ([]Task, error) {
	var result []Task
	err := c.do("GET", fmt.Sprintf("/workspaces/%s/projects/%s/tasks", workspaceID, projectID), nil, &result)
	return result, err
}

func (c *Client) CreateTask(workspaceID, projectID string, req CreateTaskRequest) (*Task, error) {
	var result Task
	err := c.do("POST", fmt.Sprintf("/workspaces/%s/projects/%s/tasks", workspaceID, projectID), req, &result)
	return &result, err
}

func (c *Client) UpdateTask(workspaceID, projectID, taskID string, req UpdateTaskRequest) (*Task, error) {
	var result Task
	err := c.do("PUT", fmt.Sprintf("/workspaces/%s/projects/%s/tasks/%s", workspaceID, projectID, taskID), req, &result)
	return &result, err
}

func (c *Client) DeleteTask(workspaceID, projectID, taskID string) error {
	return c.do("DELETE", fmt.Sprintf("/workspaces/%s/projects/%s/tasks/%s", workspaceID, projectID, taskID), nil, nil)
}

// --- Tag ---

type Tag struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	WorkspaceID string `json:"workspaceId"`
	Archived    bool   `json:"archived"`
}

type CreateTagRequest struct {
	Name string `json:"name"`
}

type UpdateTagRequest struct {
	Name     string `json:"name,omitempty"`
	Archived *bool  `json:"archived,omitempty"`
}

func (c *Client) GetTags(workspaceID string) ([]Tag, error) {
	var result []Tag
	err := c.do("GET", fmt.Sprintf("/workspaces/%s/tags", workspaceID), nil, &result)
	return result, err
}

func (c *Client) CreateTag(workspaceID string, req CreateTagRequest) (*Tag, error) {
	var result Tag
	err := c.do("POST", fmt.Sprintf("/workspaces/%s/tags", workspaceID), req, &result)
	return &result, err
}

func (c *Client) UpdateTag(workspaceID, tagID string, req UpdateTagRequest) (*Tag, error) {
	var result Tag
	err := c.do("PUT", fmt.Sprintf("/workspaces/%s/tags/%s", workspaceID, tagID), req, &result)
	return &result, err
}

func (c *Client) DeleteTag(workspaceID, tagID string) error {
	return c.do("DELETE", fmt.Sprintf("/workspaces/%s/tags/%s", workspaceID, tagID), nil, nil)
}

// --- Client (Clockify client entity) ---

type ClockifyClient struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	WorkspaceID string `json:"workspaceId"`
	Archived    bool   `json:"archived"`
}

type CreateClientRequest struct {
	Name string `json:"name"`
}

type UpdateClientRequest struct {
	Name     string `json:"name,omitempty"`
	Archived *bool  `json:"archived,omitempty"`
}

func (c *Client) GetClients(workspaceID string) ([]ClockifyClient, error) {
	var result []ClockifyClient
	err := c.do("GET", fmt.Sprintf("/workspaces/%s/clients", workspaceID), nil, &result)
	return result, err
}

func (c *Client) CreateClient(workspaceID string, req CreateClientRequest) (*ClockifyClient, error) {
	var result ClockifyClient
	err := c.do("POST", fmt.Sprintf("/workspaces/%s/clients", workspaceID), req, &result)
	return &result, err
}

func (c *Client) UpdateClient(workspaceID, clientID string, req UpdateClientRequest) (*ClockifyClient, error) {
	var result ClockifyClient
	err := c.do("PUT", fmt.Sprintf("/workspaces/%s/clients/%s", workspaceID, clientID), req, &result)
	return &result, err
}

func (c *Client) DeleteClient(workspaceID, clientID string) error {
	return c.do("DELETE", fmt.Sprintf("/workspaces/%s/clients/%s", workspaceID, clientID), nil, nil)
}

// --- Time Entry ---

type TimeInterval struct {
	Start    string `json:"start"`
	End      string `json:"end,omitempty"`
	Duration string `json:"duration,omitempty"`
}

type TimeEntry struct {
	ID           string       `json:"id"`
	Description  string       `json:"description"`
	ProjectID    string       `json:"projectId,omitempty"`
	TaskID       string       `json:"taskId,omitempty"`
	TagIDs       []string     `json:"tagIds,omitempty"`
	Billable     bool         `json:"billable"`
	TimeInterval TimeInterval `json:"timeInterval"`
	UserID       string       `json:"userId,omitempty"`
}

type CreateTimeEntryRequest struct {
	Start       string   `json:"start"`
	End         string   `json:"end,omitempty"`
	Description string   `json:"description,omitempty"`
	ProjectID   string   `json:"projectId,omitempty"`
	TaskID      string   `json:"taskId,omitempty"`
	TagIDs      []string `json:"tagIds,omitempty"`
	Billable    bool     `json:"billable"`
}

type UpdateTimeEntryRequest struct {
	Start       string   `json:"start"`
	End         string   `json:"end,omitempty"`
	Description string   `json:"description,omitempty"`
	ProjectID   string   `json:"projectId,omitempty"`
	TaskID      string   `json:"taskId,omitempty"`
	TagIDs      []string `json:"tagIds,omitempty"`
	Billable    bool     `json:"billable"`
}

func (c *Client) GetTimeEntries(workspaceID, userID string, params url.Values) ([]TimeEntry, error) {
	var result []TimeEntry
	endpoint := fmt.Sprintf("/workspaces/%s/user/%s/time-entries", workspaceID, userID)
	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}
	err := c.do("GET", endpoint, nil, &result)
	return result, err
}

func (c *Client) CreateTimeEntry(workspaceID string, req CreateTimeEntryRequest) (*TimeEntry, error) {
	var result TimeEntry
	err := c.do("POST", fmt.Sprintf("/workspaces/%s/time-entries", workspaceID), req, &result)
	return &result, err
}

func (c *Client) UpdateTimeEntry(workspaceID, entryID string, req UpdateTimeEntryRequest) (*TimeEntry, error) {
	var result TimeEntry
	err := c.do("PUT", fmt.Sprintf("/workspaces/%s/time-entries/%s", workspaceID, entryID), req, &result)
	return &result, err
}

func (c *Client) DeleteTimeEntry(workspaceID, entryID string) error {
	return c.do("DELETE", fmt.Sprintf("/workspaces/%s/time-entries/%s", workspaceID, entryID), nil, nil)
}

// --- Timer ---

func (c *Client) StartTimer(workspaceID string, req CreateTimeEntryRequest) (*TimeEntry, error) {
	// A timer is just a time entry without an end time
	req.End = ""
	return c.CreateTimeEntry(workspaceID, req)
}

func (c *Client) StopTimer(workspaceID, userID string) (*TimeEntry, error) {
	var result TimeEntry
	body := map[string]string{"end": time.Now().UTC().Format("2006-01-02T15:04:05Z")}
	err := c.do("PATCH", fmt.Sprintf("/workspaces/%s/user/%s/time-entries", workspaceID, userID), body, &result)
	return &result, err
}

func (c *Client) GetRunningTimer(workspaceID, userID string) (*TimeEntry, error) {
	var result []TimeEntry
	params := url.Values{"in-progress": {"true"}}
	endpoint := fmt.Sprintf("/workspaces/%s/user/%s/time-entries?%s", workspaceID, userID, params.Encode())
	err := c.do("GET", endpoint, nil, &result)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, nil
	}
	return &result[0], nil
}

// --- Reports ---

type SummaryReportRequest struct {
	DateRangeStart string              `json:"dateRangeStart"`
	DateRangeEnd   string              `json:"dateRangeEnd"`
	SummaryFilter  *SummaryFilter      `json:"summaryFilter,omitempty"`
	Users          *ReportUsersFilter  `json:"users,omitempty"`
	Projects       *ReportProjectFilter `json:"projects,omitempty"`
}

type SummaryFilter struct {
	Groups []string `json:"groups,omitempty"`
}

type ReportUsersFilter struct {
	IDs      []string `json:"ids,omitempty"`
	Contains string   `json:"contains,omitempty"`
	Status   string   `json:"status,omitempty"`
}

type ReportProjectFilter struct {
	IDs      []string `json:"ids,omitempty"`
	Contains string   `json:"contains,omitempty"`
	Status   string   `json:"status,omitempty"`
}

type SummaryReport struct {
	Totals  []ReportTotal  `json:"totals,omitempty"`
	GroupOne []ReportGroup `json:"groupOne,omitempty"`
}

type ReportTotal struct {
	TotalTime     int64   `json:"totalTime"`
	TotalBillable int64   `json:"totalBillableTime"`
	TotalAmount   float64 `json:"totalAmount"`
}

type ReportGroup struct {
	Name     string `json:"name"`
	Duration int64  `json:"duration"`
}

type DetailedReportRequest struct {
	DateRangeStart string              `json:"dateRangeStart"`
	DateRangeEnd   string              `json:"dateRangeEnd"`
	DetailedFilter *DetailedFilter     `json:"detailedFilter,omitempty"`
	Users          *ReportUsersFilter  `json:"users,omitempty"`
	Projects       *ReportProjectFilter `json:"projects,omitempty"`
	SortColumn     string              `json:"sortColumn,omitempty"`
	SortOrder      string              `json:"sortOrder,omitempty"`
	Page           int                 `json:"page,omitempty"`
	PageSize       int                 `json:"pageSize,omitempty"`
}

type DetailedFilter struct {
	Page     int `json:"page,omitempty"`
	PageSize int `json:"pageSize,omitempty"`
}

type DetailedReport struct {
	TimeEntries []DetailedReportEntry `json:"timeentries,omitempty"`
	TotalCount  int                   `json:"totalsCount,omitempty"`
}

type DetailedReportEntry struct {
	Description  string `json:"description"`
	ProjectName  string `json:"projectName"`
	UserName     string `json:"userName"`
	TimeInterval TimeInterval `json:"timeInterval"`
	Duration     int64  `json:"duration"`
}

func (c *Client) GetSummaryReport(workspaceID string, req SummaryReportRequest) (*SummaryReport, error) {
	var result SummaryReport
	err := c.doReports("POST", fmt.Sprintf("/workspaces/%s/reports/summary", workspaceID), req, &result)
	return &result, err
}

func (c *Client) GetDetailedReport(workspaceID string, req DetailedReportRequest) (*DetailedReport, error) {
	var result DetailedReport
	err := c.doReports("POST", fmt.Sprintf("/workspaces/%s/reports/detailed", workspaceID), req, &result)
	return &result, err
}
