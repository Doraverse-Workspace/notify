package notify

import (
	"context"
	"testing"
)

func TestNewManager(t *testing.T) {
	manager := NewManager()
	if manager == nil {
		t.Fatal("Expected manager to be created")
	}

	if len(manager.List()) != 0 {
		t.Errorf("Expected empty manager, got %d notifiers", len(manager.List()))
	}
}

func TestManagerRegister(t *testing.T) {
	manager := NewManager()
	notifier := NewMockNotifier("test")

	err := manager.Register(notifier)
	if err != nil {
		t.Fatalf("Failed to register notifier: %v", err)
	}

	if len(manager.List()) != 1 {
		t.Errorf("Expected 1 notifier, got %d", len(manager.List()))
	}

	// Test duplicate registration
	err = manager.Register(notifier)
	if err == nil {
		t.Error("Expected error when registering duplicate notifier")
	}
}

func TestManagerUnregister(t *testing.T) {
	manager := NewManager()
	notifier := NewMockNotifier("test")

	manager.Register(notifier)
	manager.Unregister("test")

	if len(manager.List()) != 0 {
		t.Errorf("Expected 0 notifiers after unregister, got %d", len(manager.List()))
	}
}

func TestManagerGet(t *testing.T) {
	manager := NewManager()
	notifier := NewMockNotifier("test")

	manager.Register(notifier)

	retrieved, exists := manager.Get("test")
	if !exists {
		t.Error("Expected notifier to exist")
	}

	if retrieved.Name() != "test" {
		t.Errorf("Expected notifier name 'test', got '%s'", retrieved.Name())
	}

	_, exists = manager.Get("nonexistent")
	if exists {
		t.Error("Expected notifier to not exist")
	}
}

func TestManagerSend(t *testing.T) {
	manager := NewManager()
	notifier := NewMockNotifier("test")
	manager.Register(notifier)

	ctx := context.Background()
	err := manager.Send(ctx, "test", "Hello")
	if err != nil {
		t.Fatalf("Failed to send: %v", err)
	}

	if !notifier.sendCalled {
		t.Error("Expected Send to be called")
	}

	if notifier.lastMessage != "Hello" {
		t.Errorf("Expected message 'Hello', got '%s'", notifier.lastMessage)
	}
}

func TestManagerBroadcast(t *testing.T) {
	manager := NewManager()
	notifier1 := NewMockNotifier("test1")
	notifier2 := NewMockNotifier("test2")

	manager.Register(notifier1)
	manager.Register(notifier2)

	ctx := context.Background()
	errors := manager.Broadcast(ctx, "Broadcast message")

	if len(errors) != 0 {
		t.Errorf("Expected no errors, got %d", len(errors))
	}

	if !notifier1.sendCalled || !notifier2.sendCalled {
		t.Error("Expected both notifiers to be called")
	}
}

func TestManagerBroadcastWithErrors(t *testing.T) {
	manager := NewManager()
	notifier1 := NewMockNotifier("test1")
	notifier2 := NewMockNotifier("test2")
	notifier2.shouldFail = true

	manager.Register(notifier1)
	manager.Register(notifier2)

	ctx := context.Background()
	errors := manager.Broadcast(ctx, "Broadcast message")

	if len(errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(errors))
	}
}

func TestManagerBroadcastAsync(t *testing.T) {
	manager := NewManager()
	notifier1 := NewMockNotifier("test1")
	notifier2 := NewMockNotifier("test2")

	manager.Register(notifier1)
	manager.Register(notifier2)

	ctx := context.Background()
	resultChan := manager.BroadcastAsync(ctx, "Async message")

	successCount := 0
	failCount := 0

	for result := range resultChan {
		if result.Success {
			successCount++
		} else {
			failCount++
		}
	}

	if successCount != 2 {
		t.Errorf("Expected 2 successes, got %d", successCount)
	}

	if failCount != 0 {
		t.Errorf("Expected 0 failures, got %d", failCount)
	}
}

func TestManagerSendWithOptions(t *testing.T) {
	manager := NewManager()
	notifier := NewMockNotifier("test")
	manager.Register(notifier)

	ctx := context.Background()
	msg := &Message{
		Text:     "Test message",
		Title:    "Test Title",
		Priority: PriorityHigh,
	}

	err := manager.SendWithOptions(ctx, "test", msg)
	if err != nil {
		t.Fatalf("Failed to send with options: %v", err)
	}

	if !notifier.sendCalled {
		t.Error("Expected SendWithOptions to be called")
	}
}
