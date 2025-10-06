package notify

import (
	"context"
	"fmt"
)

// Notifier defines the interface that all notification providers must implement
type Notifier interface {
	// Send sends a simple text message
	Send(ctx context.Context, message string) error

	// SendWithOptions sends a message with additional options
	SendWithOptions(ctx context.Context, msg *Message) error

	// Name returns the name of the notification provider
	Name() string
}

// Message represents a notification message with options
type Message struct {
	// Text is the main message content
	Text string

	// Title is an optional title for the message
	Title string

	// Priority defines the message priority (high, normal, low)
	Priority string

	// Channel defines the target channel/chat (provider-specific)
	Channel string

	// Attachments for rich messages (provider-specific)
	Attachments []Attachment

	// Metadata for additional provider-specific data
	Metadata map[string]interface{}
}

// Attachment represents a message attachment
type Attachment struct {
	Title      string
	Text       string
	ImageURL   string
	Color      string
	Fields     []Field
	Footer     string
	FooterIcon string
}

// Field represents a key-value field in an attachment
type Field struct {
	Title string
	Value string
	Short bool
}

// Priority constants
const (
	PriorityHigh   = "high"
	PriorityNormal = "normal"
	PriorityLow    = "low"
)

// NotificationError represents an error that occurred during notification
type NotificationError struct {
	Provider string
	Message  string
	Err      error
}

func (e *NotificationError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s notification error: %s - %v", e.Provider, e.Message, e.Err)
	}
	return fmt.Sprintf("%s notification error: %s", e.Provider, e.Message)
}

func (e *NotificationError) Unwrap() error {
	return e.Err
}
