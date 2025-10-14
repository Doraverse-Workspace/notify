package notify

import (
	"context"
	"fmt"
	"sync"
)

var (
	// globalManager is the singleton instance
	globalManager *Manager
	once          sync.Once
	mu            sync.RWMutex
)

// Init initializes the global notification manager
// This should be called once at application startup
func Init() {
	once.Do(func() {
		globalManager = NewManager()
	})
}

// Setup initializes and configures the global notification manager with providers
// This is a convenience function that calls Init and registers providers
func Setup(configs ...interface{}) error {
	Init()

	for _, config := range configs {
		var notifier Notifier
		var err error

		switch cfg := config.(type) {
		case *SlackConfig:
			notifier, err = NewSlackNotifier(cfg)
		case SlackConfig:
			notifier, err = NewSlackNotifier(&cfg)
		case *TelegramConfig:
			notifier, err = NewTelegramNotifier(*cfg)
		case TelegramConfig:
			notifier, err = NewTelegramNotifier(cfg)
		case Notifier:
			// Allow custom notifiers to be passed directly
			notifier = cfg
		default:
			return fmt.Errorf("unsupported config type: %T", config)
		}

		if err != nil {
			return fmt.Errorf("failed to create notifier: %w", err)
		}

		if err := globalManager.Register(notifier); err != nil {
			return fmt.Errorf("failed to register notifier: %w", err)
		}
	}

	return nil
}

// Global returns the global manager instance
// If Init() or Setup() has not been called, it will initialize a new manager
func Global() *Manager {
	mu.RLock()
	if globalManager != nil {
		mu.RUnlock()
		return globalManager
	}
	mu.RUnlock()

	Init()
	return globalManager
}

// Register adds a notifier to the global manager
func Register(notifier Notifier) error {
	return Global().Register(notifier)
}

// Unregister removes a notifier from the global manager
func Unregister(name string) {
	Global().Unregister(name)
}

// Get retrieves a notifier by name from the global manager
func Get(name string) (Notifier, bool) {
	return Global().Get(name)
}

// List returns all registered notifier names from the global manager
func List() []string {
	return Global().List()
}

// Send sends a message to a specific provider using the global manager
func Send(ctx context.Context, provider, message string) error {
	return Global().Send(ctx, provider, message)
}

// SendWithOptions sends a message with options to a specific provider using the global manager
func SendWithOptions(ctx context.Context, provider string, msg *Message) error {
	return Global().SendWithOptions(ctx, provider, msg)
}

// SendRichMessage sends a rich message to a specific provider using the global manager
func SendRichMessage(ctx context.Context, provider, channel string, blocks interface{}) error {
	return Global().SendRichMessage(ctx, provider, channel, blocks)
}

// Broadcast sends a message to all registered notifiers using the global manager
func Broadcast(ctx context.Context, message string) []error {
	return Global().Broadcast(ctx, message)
}

// BroadcastWithOptions sends a message with options to all registered notifiers using the global manager
func BroadcastWithOptions(ctx context.Context, msg *Message) []error {
	return Global().BroadcastWithOptions(ctx, msg)
}

// BroadcastAsync sends a message to all registered notifiers asynchronously using the global manager
func BroadcastAsync(ctx context.Context, message string) <-chan NotificationResult {
	return Global().BroadcastAsync(ctx, message)
}

// BroadcastAsyncWithOptions sends a message with options to all registered notifiers asynchronously using the global manager
func BroadcastAsyncWithOptions(ctx context.Context, msg *Message) <-chan NotificationResult {
	return Global().BroadcastAsyncWithOptions(ctx, msg)
}

// Reset clears the global manager (useful for testing)
func Reset() {
	mu.Lock()
	defer mu.Unlock()
	globalManager = nil
	once = sync.Once{}
}
