package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Doraverse-Workspace/notify"
)

// CustomNotifier is an example of a custom notification provider
type CustomNotifier struct {
	name string
}

// NewCustomNotifier creates a new custom notifier
func NewCustomNotifier(name string) *CustomNotifier {
	return &CustomNotifier{name: name}
}

// Name returns the name of the provider
func (c *CustomNotifier) Name() string {
	return c.name
}

// Send sends a simple text message
func (c *CustomNotifier) Send(ctx context.Context, message string) error {
	// Implement your custom notification logic here
	fmt.Printf("[%s] Sending message: %s\n", c.name, message)
	return nil
}

// SendWithOptions sends a message with additional options
func (c *CustomNotifier) SendWithOptions(ctx context.Context, msg *notify.Message) error {
	// Implement your custom notification logic here
	fmt.Printf("[%s] Sending message with options:\n", c.name)
	if msg.Title != "" {
		fmt.Printf("  Title: %s\n", msg.Title)
	}
	fmt.Printf("  Text: %s\n", msg.Text)
	if msg.Priority != "" {
		fmt.Printf("  Priority: %s\n", msg.Priority)
	}
	return nil
}

func main() {
	ctx := context.Background()

	// Create a custom notifier
	customNotifier := NewCustomNotifier("console")

	// Create a manager and register the custom notifier
	manager := notify.NewManager()
	err := manager.Register(customNotifier)
	if err != nil {
		log.Fatalf("Failed to register custom notifier: %v", err)
	}

	fmt.Println("Custom Notifier Example")
	fmt.Println("======================")

	// Example 1: Simple message
	fmt.Println("\nExample 1: Simple message")
	err = manager.Send(ctx, "console", "Hello from custom notifier!")
	if err != nil {
		log.Printf("Error: %v\n", err)
	}

	// Example 2: Message with options
	fmt.Println("\nExample 2: Message with options")
	msg := &notify.Message{
		Title:    "System Alert",
		Text:     "This is a custom notification",
		Priority: notify.PriorityHigh,
	}
	err = manager.SendWithOptions(ctx, "console", msg)
	if err != nil {
		log.Printf("Error: %v\n", err)
	}

	// Example 3: Register multiple custom notifiers
	fmt.Println("\nExample 3: Multiple custom notifiers")
	emailNotifier := NewCustomNotifier("email")
	smsNotifier := NewCustomNotifier("sms")

	manager.Register(emailNotifier)
	manager.Register(smsNotifier)

	fmt.Printf("Registered notifiers: %v\n", manager.List())

	// Broadcast to all notifiers
	fmt.Println("\nBroadcasting to all notifiers:")
	errors := manager.Broadcast(ctx, "Broadcast message to all custom notifiers!")
	if len(errors) > 0 {
		for _, err := range errors {
			log.Printf("Error: %v\n", err)
		}
	}
}
