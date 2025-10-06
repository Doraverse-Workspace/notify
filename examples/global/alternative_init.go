package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Doraverse-Workspace/notify"
)

// Example showing alternative initialization methods

func exampleMethod2ManualInit() {
	fmt.Println("\n=== Method 2: Manual Registration ===")

	// Initialize the global manager
	notify.Init()

	// Create notifiers manually and register them
	telegramNotifier, err := notify.NewTelegramNotifier(notify.TelegramConfig{
		BotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		ChatID:   os.Getenv("TELEGRAM_CHAT_ID"),
	})
	if err == nil {
		err = notify.Register(telegramNotifier)
		if err != nil {
			log.Printf("Failed to register: %v\n", err)
		} else {
			fmt.Println("✓ Telegram notifier registered")
		}
	}

	slackNotifier, err := notify.NewSlackNotifier(&notify.SlackConfig{
		Token:          os.Getenv("SLACK_BOT_TOKEN"),
		DefaultChannel: os.Getenv("SLACK_CHANNEL"),
	})
	if err == nil {
		err = notify.Register(slackNotifier)
		if err != nil {
			log.Printf("Failed to register: %v\n", err)
		} else {
			fmt.Println("✓ Slack notifier registered")
		}
	}

	// Use the global functions
	ctx := context.Background()
	err = notify.Send(ctx, "telegram", "Hello from method 2!")
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Println("✓ Message sent successfully")
	}
}

func exampleMethod3DirectGlobalAccess() {
	fmt.Println("\n=== Method 3: Direct Global Manager Access ===")

	// Get direct access to the global manager
	manager := notify.Global()

	// Use it like a regular manager
	ctx := context.Background()

	notifier, exists := manager.Get("telegram")
	if exists {
		err := notifier.Send(ctx, "Direct access message!")
		if err != nil {
			log.Printf("Error: %v\n", err)
		} else {
			fmt.Println("✓ Message sent via direct access")
		}
	}

	// List all providers
	providers := manager.List()
	fmt.Printf("Available providers: %v\n", providers)
}

func exampleLazyInitialization() {
	fmt.Println("\n=== Method 4: Lazy Initialization ===")

	// No explicit Init() or Setup() call
	// The first call to any global function will auto-initialize
	ctx := context.Background()

	// This will auto-initialize if not already done
	providers := notify.List()
	fmt.Printf("Providers (auto-init): %v\n", providers)

	// Register a notifier even without explicit Init()
	telegramNotifier, err := notify.NewTelegramNotifier(notify.TelegramConfig{
		BotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		ChatID:   os.Getenv("TELEGRAM_CHAT_ID"),
	})
	if err == nil {
		notify.Register(telegramNotifier)
		fmt.Println("✓ Notifier registered via lazy init")

		notify.Send(ctx, "telegram", "Lazy init message!")
	}
}

// Example of using in a package
type NotificationService struct {
	// No need to store manager or config!
}

func (s *NotificationService) SendAlert(ctx context.Context, message string) error {
	// Just use the global notify functions
	return notify.Send(ctx, "telegram", message)
}

func (s *NotificationService) BroadcastAnnouncement(ctx context.Context, title, text string) error {
	msg := &notify.Message{
		Title:    title,
		Text:     text,
		Priority: notify.PriorityNormal,
	}

	errors := notify.BroadcastWithOptions(ctx, msg)
	if len(errors) > 0 {
		return errors[0]
	}
	return nil
}

func exampleInService() {
	fmt.Println("\n=== Example: Using in a Service ===")

	service := &NotificationService{}
	ctx := context.Background()

	err := service.SendAlert(ctx, "Alert from service!")
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Println("✓ Alert sent from service")
	}

	err = service.BroadcastAnnouncement(ctx, "Announcement", "System maintenance in 1 hour")
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Println("✓ Announcement broadcast from service")
	}
}
