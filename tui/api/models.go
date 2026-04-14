package api

import (
	"encoding/json"
	"time"
)

type Flexibility string
type Energy string
type Status string
type Weekday string

const (
	FlexibilityLow    Flexibility = "L"
	FlexibilityMedium Flexibility = "M"
	FlexibilityHigh   Flexibility = "H"

	EnergyLow    Energy = "L"
	EnergyMedium Energy = "M"
	EnergyHigh   Energy = "H"

	StatusPending    Status = "pending"
	StatusInProgress Status = "in_progress"
	StatusCompleted  Status = "completed"
	StatusSkipped    Status = "skipped"

	WeekdayMonday    Weekday = "monday"
	WeekdayTuesday   Weekday = "tuesday"
	WeekdayWednesday Weekday = "wednesday"
	WeekdayThursday  Weekday = "thursday"
	WeekdayFriday    Weekday = "friday"
	WeekdaySaturday  Weekday = "saturday"
	WeekdaySunday    Weekday = "sunday"
)

type TaskCategory struct {
	ID                    int         `json:"id,omitempty"`
	Name                  string      `json:"name"`
	SchedulingFlexibility Flexibility `json:"scheduling_flexibility"`
	EnergyRequired        Energy      `json:"energy_required"`
	NeedsFocusBlock       bool        `json:"needs_focus_block"`
}

type HardRoutine struct {
	ID        int       `json:"id,omitempty"`
	Name      string    `json:"name"`
	Weekdays  []Weekday `json:"weekdays"`
	StartTime string    `json:"start_time"` // HH:MM format
	Duration  int       `json:"duration"`   // minutes
	IsActive  bool      `json:"is_active"`
}

type Task struct {
	ID                int         `json:"id,omitempty"`
	Title             string      `json:"title"`
	Description       string      `json:"description,omitempty"`
	CategoryID        int         `json:"category_id,omitempty"`
	Flexibility       Flexibility `json:"flexibility,omitempty"`
	EnergyRequired    Energy      `json:"energy_required,omitempty"`
	ScheduledStart    time.Time   `json:"scheduled_start"`
	ScheduledDate     string      `json:"scheduled_date"`
	EstimatedDuration int         `json:"estimated_duration"`
	Priority          int         `json:"priority"`
	Status            Status      `json:"status"`
	ParentTaskID      *int        `json:"parent_task_id,omitempty"`
	ActualStart       *time.Time  `json:"actual_start,omitempty"`
	ActualEnd         *time.Time  `json:"actual_end,omitempty"`
	ActualDuration    *int        `json:"actual_duration,omitempty"`
	ActualDate        *string     `json:"actual_date,omitempty"`
	LastStartedAt     *time.Time  `json:"last_started_at,omitempty"`
}

type UserProfile struct {
	ID                int    `json:"id,omitempty"`
	DefaultWorkStart  string `json:"default_work_start"`  // HH:MM format
	DefaultSleepStart string `json:"default_sleep_start"` // HH:MM format
	Timezone          string `json:"timezone"`
}

type ShiftTasksRequest struct {
	ShiftFromTime      time.Time `json:"shift_from_time"`
	ShiftAmountMinutes int       `json:"shift_amount_minutes"`
}

type WrapTaskRequest struct {
	TaskID int    `json:"task_id"`
	Date   string `json:"date"` // YYYY-MM-DD format
}

type APIError struct {
	Detail string `json:"detail"`
}

func (e *APIError) Error() string {
	return e.Detail
}

// Helper to parse error response
func ParseError(body []byte) *APIError {
	var err APIError
	if json.Unmarshal(body, &err) == nil && err.Detail != "" {
		return &err
	}
	return &APIError{Detail: string(body)}
}
