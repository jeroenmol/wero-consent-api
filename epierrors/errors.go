package epierrors

import (
	"fmt"
	"strings"
)

// Problem represents a standard RFC 7807 problem detail returned by the API.
type Problem struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`
	CallID   string `json:"callId,omitempty"`
}

func (p *Problem) Error() string {
	if p.Detail != "" {
		return fmt.Sprintf("%s (status %d): %s", p.Title, p.Status, p.Detail)
	}
	return fmt.Sprintf("%s (status %d)", p.Title, p.Status)
}

// ConstraintViolation is a single field-level validation failure.
type ConstraintViolation struct {
	Field            string `json:"field"`
	Message          string `json:"message"`
	UserMessageKey   string `json:"userMessageKey,omitempty"`
	UserMessageValue string `json:"userMessageValue,omitempty"`
}

// ConstraintViolationProblem is returned when request validation fails (HTTP 400).
type ConstraintViolationProblem struct {
	Type       string                `json:"type"`
	Title      string                `json:"title"`
	Status     int                   `json:"status"`
	CallID     string                `json:"callId,omitempty"`
	Violations []ConstraintViolation `json:"violations"`
}

func (p *ConstraintViolationProblem) Error() string {
	msgs := make([]string, len(p.Violations))
	for i, v := range p.Violations {
		msgs[i] = fmt.Sprintf("%s: %s", v.Field, v.Message)
	}
	return fmt.Sprintf("constraint violation: %s", strings.Join(msgs, "; "))
}
