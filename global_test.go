package notify

import (
	"context"
	"testing"
)

// MockNotifier for testing
type mockGlobalNotifier struct {
	name        string
	sendCalled  bool
	lastMessage string
	shouldFail  bool
}

func (m *mockGlobalNotifier) Name() string {
	return m.name
}

func (m *mockGlobalNotifier) Send(ctx context.Context, message string) error {
	m.sendCalled = true
	m.lastMessage = message
	if m.shouldFail {
		return &NotificationError{Provider: m.name, Message: "mock error"}
	}
	return nil
}

func (m *mockGlobalNotifier) SendWithOptions(ctx context.Context, msg *Message) error {
	m.sendCalled = true
	m.lastMessage = msg.Text
	if m.shouldFail {
		return &NotificationError{Provider: m.name, Message: "mock error"}
	}
	return nil
}

func TestGlobalInit(t *testing.T) {
	// Reset before test
	Reset()

	// Test Init
	Init()
	if globalManager == nil {
		t.Error("Expected globalManager to be initialized")
	}

	// Test multiple Init calls (should only init once)
	oldManager := globalManager
	Init()
	if globalManager != oldManager {
		t.Error("Expected globalManager to remain the same after multiple Init calls")
	}
}

func TestGlobalSetup(t *testing.T) {
	Reset()

	// Setup with mock configs
	telegramConfig := TelegramConfig{
		BotToken: "test-token",
		ChatID:   "test-chat",
	}

	err := Setup(telegramConfig)
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	// Check that notifier was registered
	notifiers := List()
	if len(notifiers) != 1 {
		t.Errorf("Expected 1 notifier, got %d", len(notifiers))
	}

	if notifiers[0] != "telegram" {
		t.Errorf("Expected telegram notifier, got %s", notifiers[0])
	}
}

func TestGlobalRegister(t *testing.T) {
	Reset()

	mock := &mockGlobalNotifier{name: "mock"}
	err := Register(mock)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	// Check registration
	notifier, exists := Get("mock")
	if !exists {
		t.Error("Expected notifier to be registered")
	}

	if notifier != mock {
		t.Error("Expected to get the same notifier instance")
	}
}

func TestGlobalSend(t *testing.T) {
	Reset()

	mock := &mockGlobalNotifier{name: "mock"}
	Register(mock)

	ctx := context.Background()
	err := Send(ctx, "mock", "test message")
	if err != nil {
		t.Fatalf("Send failed: %v", err)
	}

	if !mock.sendCalled {
		t.Error("Expected Send to be called on notifier")
	}

	if mock.lastMessage != "test message" {
		t.Errorf("Expected message 'test message', got '%s'", mock.lastMessage)
	}
}

func TestGlobalSendWithOptions(t *testing.T) {
	Reset()

	mock := &mockGlobalNotifier{name: "mock"}
	Register(mock)

	ctx := context.Background()
	msg := &Message{
		Text:     "test message",
		Title:    "Test",
		Priority: PriorityHigh,
	}

	err := SendWithOptions(ctx, "mock", msg)
	if err != nil {
		t.Fatalf("SendWithOptions failed: %v", err)
	}

	if !mock.sendCalled {
		t.Error("Expected SendWithOptions to be called on notifier")
	}
}

func TestGlobalBroadcast(t *testing.T) {
	Reset()

	mock1 := &mockGlobalNotifier{name: "mock1"}
	mock2 := &mockGlobalNotifier{name: "mock2"}

	Register(mock1)
	Register(mock2)

	ctx := context.Background()
	errors := Broadcast(ctx, "broadcast message")

	if len(errors) != 0 {
		t.Errorf("Expected no errors, got %d", len(errors))
	}

	if !mock1.sendCalled || !mock2.sendCalled {
		t.Error("Expected both notifiers to receive the message")
	}

	if mock1.lastMessage != "broadcast message" {
		t.Errorf("Mock1 expected 'broadcast message', got '%s'", mock1.lastMessage)
	}
}

func TestGlobalBroadcastWithErrors(t *testing.T) {
	Reset()

	mock1 := &mockGlobalNotifier{name: "mock1", shouldFail: true}
	mock2 := &mockGlobalNotifier{name: "mock2"}

	Register(mock1)
	Register(mock2)

	ctx := context.Background()
	errors := Broadcast(ctx, "test")

	if len(errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(errors))
	}

	// mock2 should still succeed
	if !mock2.sendCalled {
		t.Error("Expected mock2 to be called despite mock1 failure")
	}
}

func TestGlobalBroadcastAsync(t *testing.T) {
	Reset()

	mock1 := &mockGlobalNotifier{name: "mock1"}
	mock2 := &mockGlobalNotifier{name: "mock2"}

	Register(mock1)
	Register(mock2)

	ctx := context.Background()
	resultChan := BroadcastAsync(ctx, "async message")

	successCount := 0
	for result := range resultChan {
		if result.Success {
			successCount++
		}
	}

	if successCount != 2 {
		t.Errorf("Expected 2 successful sends, got %d", successCount)
	}
}

func TestGlobalBroadcastAsyncWithOptions(t *testing.T) {
	Reset()

	mock1 := &mockGlobalNotifier{name: "mock1"}
	mock2 := &mockGlobalNotifier{name: "mock2"}

	Register(mock1)
	Register(mock2)

	ctx := context.Background()
	msg := &Message{
		Text:     "async message",
		Priority: PriorityNormal,
	}
	resultChan := BroadcastAsyncWithOptions(ctx, msg)

	results := make([]NotificationResult, 0)
	for result := range resultChan {
		results = append(results, result)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}
}

func TestGlobalUnregister(t *testing.T) {
	Reset()

	mock := &mockGlobalNotifier{name: "mock"}
	Register(mock)

	// Verify it's registered
	_, exists := Get("mock")
	if !exists {
		t.Error("Expected notifier to be registered")
	}

	// Unregister
	Unregister("mock")

	// Verify it's unregistered
	_, exists = Get("mock")
	if exists {
		t.Error("Expected notifier to be unregistered")
	}
}

func TestGlobalList(t *testing.T) {
	Reset()

	mock1 := &mockGlobalNotifier{name: "mock1"}
	mock2 := &mockGlobalNotifier{name: "mock2"}

	Register(mock1)
	Register(mock2)

	list := List()
	if len(list) != 2 {
		t.Errorf("Expected 2 notifiers, got %d", len(list))
	}
}

func TestGlobalReset(t *testing.T) {
	Reset()

	mock := &mockGlobalNotifier{name: "mock"}
	Register(mock)

	// Verify registration
	if len(List()) != 1 {
		t.Error("Expected 1 notifier before reset")
	}

	// Reset
	Reset()

	// After reset, should be able to init fresh
	Init()
	if len(List()) != 0 {
		t.Error("Expected 0 notifiers after reset")
	}
}

func TestGlobalAutoInit(t *testing.T) {
	Reset()

	// Don't call Init() or Setup()
	// Just use Global() directly - should auto-initialize
	manager := Global()
	if manager == nil {
		t.Error("Expected Global() to auto-initialize manager")
	}

	// Verify it works
	mock := &mockGlobalNotifier{name: "mock"}
	err := Register(mock)
	if err != nil {
		t.Errorf("Register failed after auto-init: %v", err)
	}
}

func TestGlobalSetupWithInvalidConfig(t *testing.T) {
	Reset()

	// Try to setup with unsupported config type
	err := Setup("invalid config")
	if err == nil {
		t.Error("Expected error for invalid config type")
	}
}

func TestGlobalSetupWithMultipleConfigs(t *testing.T) {
	Reset()

	telegramConfig := TelegramConfig{
		BotToken: "test-token",
		ChatID:   "test-chat",
	}

	slackConfig := &SlackConfig{
		Token:          "test-token",
		DefaultChannel: "#test",
	}

	err := Setup(telegramConfig, slackConfig)
	if err != nil {
		t.Fatalf("Setup with multiple configs failed: %v", err)
	}

	list := List()
	if len(list) != 2 {
		t.Errorf("Expected 2 notifiers, got %d", len(list))
	}
}

func TestGlobalSetupWithCustomNotifier(t *testing.T) {
	Reset()

	mock := &mockGlobalNotifier{name: "custom"}
	err := Setup(mock)
	if err != nil {
		t.Fatalf("Setup with custom notifier failed: %v", err)
	}

	notifier, exists := Get("custom")
	if !exists {
		t.Error("Expected custom notifier to be registered")
	}

	if notifier != mock {
		t.Error("Expected to get the same custom notifier instance")
	}
}
