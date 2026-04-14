package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	BaseURL string
	client  *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		client: &http.Client{
			Timeout: 20 * time.Second, // this might be an issue when adding the AI layer
		},
	}
}

func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}) ([]byte, error) {
	var reqBody *bytes.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request: %w", err)
		}
		reqBody = bytes.NewReader(b)
	} else {
		reqBody = bytes.NewReader(nil)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.BaseURL+path, reqBody)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		apiErr := ParseError(respBody)
		return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, apiErr.Error())
	}

	return respBody, nil
}

// User Profile

func (c *Client) GetUser(ctx context.Context) (*UserProfile, error) {
	resp, err := c.doRequest(ctx, "GET", "/user/", nil)
	if err != nil {
		return nil, err
	}
	var profile UserProfile
	if err := json.Unmarshal(resp, &profile); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	return &profile, nil
}

func (c *Client) CreateUser(ctx context.Context, profile *UserProfile) (*UserProfile, error) {
	resp, err := c.doRequest(ctx, "POST", "/user/", profile)
	if err != nil {
		return nil, err
	}
	var created UserProfile
	if err := json.Unmarshal(resp, &created); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	return &created, nil
}

// Task Categories

func (c *Client) GetTaskCategories(ctx context.Context) ([]TaskCategory, error) {
	resp, err := c.doRequest(ctx, "GET", "/task-categories/", nil)
	if err != nil {
		return nil, err
	}
	var categories []TaskCategory
	if err := json.Unmarshal(resp, &categories); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	return categories, nil
}

func (c *Client) CreateTaskCategory(ctx context.Context, category *TaskCategory) (*TaskCategory, error) {
	resp, err := c.doRequest(ctx, "POST", "/task-categories/", category)
	if err != nil {
		return nil, err
	}
	var created TaskCategory
	if err := json.Unmarshal(resp, &created); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	return &created, nil
}

func (c *Client) UpdateTaskCategory(ctx context.Context, id int, category *TaskCategory) (*TaskCategory, error) {
	resp, err := c.doRequest(ctx, "PATCH", fmt.Sprintf("/task-categories/%d", id), category)
	if err != nil {
		return nil, err
	}
	var updated TaskCategory
	if err := json.Unmarshal(resp, &updated); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	return &updated, nil
}

func (c *Client) DeleteTaskCategory(ctx context.Context, id int) error {
	_, err := c.doRequest(ctx, "DELETE", fmt.Sprintf("/task-categories/%d", id), nil)
	return err
}

// Hard Routines

func (c *Client) GetHardRoutines(ctx context.Context) ([]HardRoutine, error) {
	resp, err := c.doRequest(ctx, "GET", "/hard-routines/", nil)
	if err != nil {
		return nil, err
	}
	var routines []HardRoutine
	if err := json.Unmarshal(resp, &routines); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	return routines, nil
}

func (c *Client) CreateHardRoutine(ctx context.Context, routine *HardRoutine) (*HardRoutine, error) {
	resp, err := c.doRequest(ctx, "POST", "/hard-routines/", routine)
	if err != nil {
		return nil, err
	}
	var created HardRoutine
	if err := json.Unmarshal(resp, &created); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	return &created, nil
}

func (c *Client) UpdateHardRoutine(ctx context.Context, id int, routine *HardRoutine) (*HardRoutine, error) {
	resp, err := c.doRequest(ctx, "PATCH", fmt.Sprintf("/hard-routines/%d", id), routine)
	if err != nil {
		return nil, err
	}
	var updated HardRoutine
	if err := json.Unmarshal(resp, &updated); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	return &updated, nil
}

func (c *Client) DeleteHardRoutine(ctx context.Context, id int) error {
	_, err := c.doRequest(ctx, "DELETE", fmt.Sprintf("/hard-routines/%d", id), nil)
	return err
}

// Tasks

func (c *Client) GetTasks(ctx context.Context) ([]Task, error) {
	resp, err := c.doRequest(ctx, "GET", "/tasks/", nil)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	if err := json.Unmarshal(resp, &tasks); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	return tasks, nil
}

func (c *Client) CreateTask(ctx context.Context, task *Task) (*Task, error) {
	resp, err := c.doRequest(ctx, "POST", "/tasks/", task)
	if err != nil {
		return nil, err
	}
	var created Task
	if err := json.Unmarshal(resp, &created); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	return &created, nil
}

func (c *Client) UpdateTask(ctx context.Context, id int, task *Task) (*Task, error) {
	resp, err := c.doRequest(ctx, "PATCH", fmt.Sprintf("/tasks/%d", id), task)
	if err != nil {
		return nil, err
	}
	var updated Task
	if err := json.Unmarshal(resp, &updated); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	return &updated, nil
}

func (c *Client) DeleteTask(ctx context.Context, id int) error {
	_, err := c.doRequest(ctx, "DELETE", fmt.Sprintf("/tasks/%d", id), nil)
	return err
}

// Wrap

func (c *Client) WrapTask(ctx context.Context, req *WrapTaskRequest) ([]Task, error) {
	resp, err := c.doRequest(ctx, "POST", "/tasks/wrap", req)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	if err := json.Unmarshal(resp, &tasks); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	return tasks, nil
}

// Shift

func (c *Client) ShiftTasks(ctx context.Context, req *ShiftTasksRequest) ([]Task, error) {
	resp, err := c.doRequest(ctx, "POST", "/tasks/shift", req)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	if err := json.Unmarshal(resp, &tasks); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	return tasks, nil
}
