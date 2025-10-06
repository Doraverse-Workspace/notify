package notify

import (
	"context"
	"fmt"
	"sync"
)

// Manager manages multiple notification providers
type Manager struct {
	notifiers map[string]Notifier
	mu        sync.RWMutex
}

// NewManager creates a new notification manager
func NewManager() *Manager {
	return &Manager{
		notifiers: make(map[string]Notifier),
	}
}

// Register adds a notifier to the manager
func (m *Manager) Register(notifier Notifier) error {
	if notifier == nil {
		return fmt.Errorf("notifier cannot be nil")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	name := notifier.Name()
	if _, exists := m.notifiers[name]; exists {
		return fmt.Errorf("notifier with name %s already registered", name)
	}

	m.notifiers[name] = notifier
	return nil
}

// Unregister removes a notifier from the manager
func (m *Manager) Unregister(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.notifiers, name)
}

// Get retrieves a notifier by name
func (m *Manager) Get(name string) (Notifier, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	notifier, exists := m.notifiers[name]
	return notifier, exists
}

// List returns all registered notifier names
func (m *Manager) List() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	names := make([]string, 0, len(m.notifiers))
	for name := range m.notifiers {
		names = append(names, name)
	}
	return names
}

// Send sends a message to a specific notifier
func (m *Manager) Send(ctx context.Context, provider, message string) error {
	notifier, exists := m.Get(provider)
	if !exists {
		return fmt.Errorf("notifier %s not found", provider)
	}

	return notifier.Send(ctx, message)
}

// SendWithOptions sends a message with options to a specific notifier
func (m *Manager) SendWithOptions(ctx context.Context, provider string, msg *Message) error {
	notifier, exists := m.Get(provider)
	if !exists {
		return fmt.Errorf("notifier %s not found", provider)
	}

	return notifier.SendWithOptions(ctx, msg)
}

// Broadcast sends a message to all registered notifiers
func (m *Manager) Broadcast(ctx context.Context, message string) []error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var errors []error
	for name, notifier := range m.notifiers {
		if err := notifier.Send(ctx, message); err != nil {
			errors = append(errors, fmt.Errorf("%s: %w", name, err))
		}
	}

	return errors
}

// BroadcastWithOptions sends a message with options to all registered notifiers
func (m *Manager) BroadcastWithOptions(ctx context.Context, msg *Message) []error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var errors []error
	for name, notifier := range m.notifiers {
		if err := notifier.SendWithOptions(ctx, msg); err != nil {
			errors = append(errors, fmt.Errorf("%s: %w", name, err))
		}
	}

	return errors
}

// BroadcastAsync sends a message to all registered notifiers asynchronously
func (m *Manager) BroadcastAsync(ctx context.Context, message string) <-chan NotificationResult {
	m.mu.RLock()
	notifiers := make(map[string]Notifier, len(m.notifiers))
	for name, notifier := range m.notifiers {
		notifiers[name] = notifier
	}
	m.mu.RUnlock()

	resultChan := make(chan NotificationResult, len(notifiers))

	var wg sync.WaitGroup
	for name, notifier := range notifiers {
		wg.Add(1)
		go func(n string, nt Notifier) {
			defer wg.Done()
			err := nt.Send(ctx, message)
			resultChan <- NotificationResult{
				Provider: n,
				Success:  err == nil,
				Error:    err,
			}
		}(name, notifier)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	return resultChan
}

// BroadcastAsyncWithOptions sends a message with options to all registered notifiers asynchronously
func (m *Manager) BroadcastAsyncWithOptions(ctx context.Context, msg *Message) <-chan NotificationResult {
	m.mu.RLock()
	notifiers := make(map[string]Notifier, len(m.notifiers))
	for name, notifier := range m.notifiers {
		notifiers[name] = notifier
	}
	m.mu.RUnlock()

	resultChan := make(chan NotificationResult, len(notifiers))

	var wg sync.WaitGroup
	for name, notifier := range notifiers {
		wg.Add(1)
		go func(n string, nt Notifier) {
			defer wg.Done()
			err := nt.SendWithOptions(ctx, msg)
			resultChan <- NotificationResult{
				Provider: n,
				Success:  err == nil,
				Error:    err,
			}
		}(name, notifier)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	return resultChan
}

// NotificationResult represents the result of a notification attempt
type NotificationResult struct {
	Provider string
	Success  bool
	Error    error
}
