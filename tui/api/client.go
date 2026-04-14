package api

import (
	"bytes"
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

func (c *Client) doRequest(method, path string, body interface{}) ([]byte, error) {
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

	req, err := http.NewRequest(method, c.BaseURL+path, reqBody)
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

func (c *Client) GetUser() (*UserProfile, error) {
	resp, err := c.doRequest("GET", "/user/", nil)
	if err != nil {
		return nil, err
	}
	var profile UserProfile
	if err := json.Unmarshal(resp, &profile); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	return &profile, nil
}

func (c *Client) CreateUser(profile *UserProfile) (*UserProfile, error) {
	resp, err := c.doRequest("POST", "/user/", profile)
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

func (c *Client) GetTaskCategories() ([]TaskCategory, error) {
	resp, err := c.doRequest("GET", "/task-categories/", nil)
	if err != nil {
		return nil, err
	}
	var categories []TaskCategory
	if err := json.Unmarshal(resp, &categories); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	return categories, nil
}

func (c *Client) CreateTaskCategory(category *TaskCategory) (*TaskCategory, error) {
	resp, err := c.doRequest("POST", "/task-categories/", category)
	if err != nil {
		return nil, err
	}
	var created TaskCategory
	if err := json.Unmarshal(resp, &created); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	return &created, nil
}

func (c *Client) UpdateTaskCategory(id int, category *TaskCategory) (*TaskCategory, error) {
	resp, err := c.doRequest("PATCH", fmt.Sprintf("/task-categories/%d", id), category)
	if err != nil {
		return nil, err
	}
	var updated TaskCategory
	if err := json.Unmarshal(resp, &updated); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	return &updated, nil
}

func (c *Client) DeleteTaskCategory(id int) error {
	_, err := c.doRequest("DELETE", fmt.Sprintf("/task-categories/%d", id), nil)
	return err
}

// Hard Routines

func (c *Client) GetHardRoutines() ([]HardRoutine, error) {
	resp, err := c.doRequest("GET", "/hard-routines/", nil)
	if err != nil {
		return nil, err
	}
	var routines []HardRoutine
	if err := json.Unmarshal(resp, &routines); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	return routines, nil
}

func (c *Client) CreateHardRoutine(routine *HardRoutine) (*HardRoutine, error) {
	resp, err := c.doRequest("POST", "/hard-routines/", routine)
	if err != nil {
		return nil, err
	}
	var created HardRoutine
	if err := json.Unmarshal(resp, &created); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	return &created, nil
}

func (c *Client) UpdateHardRoutine(id int, routine *HardRoutine) (*HardRoutine, error) {
	resp, err := c.doRequest("PATCH", fmt.Sprintf("/hard-routines/%d", id), routine)
	if err != nil {
		return nil, err
	}
	var updated HardRoutine
	if err := json.Unmarshal(resp, &updated); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	return &updated, nil
}

func (c *Client) DeleteHardRoutine(id int) error {
	_, err := c.doRequest("DELETE", fmt.Sprintf("/hard-routines/%d", id), nil)
	return err
}

// Tasks

func (c *Client) GetTasks() ([]Task, error) {
	resp, err := c.doRequest("GET", "/tasks/", nil)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	if err := json.Unmarshal(resp, &tasks); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	return tasks, nil
}

func (c *Client) CreateTask(task *Task) (*Task, error) {
	resp, err := c.doRequest("POST", "/tasks/", task)
	if err != nil {
		return nil, err
	}
	var created Task
	if err := json.Unmarshal(resp, &created); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	return &created, nil
}

func (c *Client) UpdateTask(id int, task *Task) (*Task, error) {
	resp, err := c.doRequest("PATCH", fmt.Sprintf("/tasks/%d", id), task)
	if err != nil {
		return nil, err
	}
	var updated Task
	if err := json.Unmarshal(resp, &updated); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	return &updated, nil
}

func (c *Client) DeleteTask(id int) error {
	_, err := c.doRequest("DELETE", fmt.Sprintf("/tasks/%d", id), nil)
	return err
}

// Wrap (auto-schedule)

func (c *Client) WrapTask(req *WrapTaskRequest) ([]Task, error) {
	resp, err := c.doRequest("POST", "/tasks/wrap", req)
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

func (c *Client) ShiftTasks(req *ShiftTasksRequest) ([]Task, error) {
	resp, err := c.doRequest("POST", "/tasks/shift", req)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	if err := json.Unmarshal(resp, &tasks); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	return tasks, nil
}
