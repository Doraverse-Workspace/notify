package notify

import (
	"context"
	"testing"
)

// MockNotifier is a mock implementation of the Notifier interface for testing
type MockNotifier struct {
	name        string
	sendCalled  bool
	lastMessage string
	shouldFail  bool
}

func NewMockNotifier(name string) *MockNotifier {
	return &MockNotifier{name: name}
}

func (m *MockNotifier) Name() string {
	return m.name
}

func (m *MockNotifier) Send(ctx context.Context, message string) error {
	m.sendCalled = true
	m.lastMessage = message
	if m.shouldFail {
		return &NotificationError{
			Provider: m.name,
			Message:  "mock error",
		}
	}
	return nil
}

func (m *MockNotifier) SendWithOptions(ctx context.Context, msg *Message) error {
	m.sendCalled = true
	m.lastMessage = msg.Text
	if m.shouldFail {
		return &NotificationError{
			Provider: m.name,
			Message:  "mock error",
		}
	}
	return nil
}

func TestMessage(t *testing.T) {
	msg := &Message{
		Text:     "Test message",
		Title:    "Test Title",
		Priority: PriorityHigh,
		Channel:  "#test",
	}

	if msg.Text != "Test message" {
		t.Errorf("Expected text 'Test message', got '%s'", msg.Text)
	}

	if msg.Priority != PriorityHigh {
		t.Errorf("Expected priority '%s', got '%s'", PriorityHigh, msg.Priority)
	}
}

func TestAttachment(t *testing.T) {
	att := Attachment{
		Title: "Test Attachment",
		Text:  "Attachment text",
		Color: "good",
		Fields: []Field{
			{Title: "Field1", Value: "Value1", Short: true},
		},
	}

	if att.Title != "Test Attachment" {
		t.Errorf("Expected title 'Test Attachment', got '%s'", att.Title)
	}

	if len(att.Fields) != 1 {
		t.Errorf("Expected 1 field, got %d", len(att.Fields))
	}
}

func TestNotificationError(t *testing.T) {
	err := &NotificationError{
		Provider: "test",
		Message:  "test error",
	}

	expected := "test notification error: test error"
	if err.Error() != expected {
		t.Errorf("Expected error '%s', got '%s'", expected, err.Error())
	}
}
