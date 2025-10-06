package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Doraverse-Workspace/notify"
)

func main() {
	ctx := context.Background()

	// Create a notification manager
	manager := notify.NewManager()

	// Register Telegram notifier
	if botToken := os.Getenv("TELEGRAM_BOT_TOKEN"); botToken != "" {
		telegramNotifier, err := notify.NewTelegramNotifier(notify.TelegramConfig{
			BotToken: botToken,
			ChatID:   os.Getenv("TELEGRAM_CHAT_ID"),
		})
		if err != nil {
			log.Printf("Failed to create Telegram notifier: %v\n", err)
		} else {
			err = manager.Register(telegramNotifier)
			if err != nil {
				log.Printf("Failed to register Telegram notifier: %v\n", err)
			} else {
				fmt.Println("✓ Telegram notifier registered")
			}
		}
	}

	// Register Slack notifier
	if slackToken := os.Getenv("SLACK_BOT_TOKEN"); slackToken != "" {
		slackNotifier, err := notify.NewSlackNotifier(&notify.SlackConfig{
			Token:          slackToken,
			DefaultChannel: os.Getenv("SLACK_CHANNEL"),
		})
		if err != nil {
			log.Printf("Failed to create Slack notifier: %v\n", err)
		} else {
			err = manager.Register(slackNotifier)
			if err != nil {
				log.Printf("Failed to register Slack notifier: %v\n", err)
			} else {
				fmt.Println("✓ Slack notifier registered")
			}
		}
	}

	// List all registered notifiers
	fmt.Printf("\nRegistered notifiers: %v\n", manager.List())

	// Example 1: Send to specific provider
	fmt.Println("\nExample 1: Send to specific provider")
	err := manager.Send(ctx, "telegram", "Hello from manager!")
	if err != nil {
		log.Printf("Failed to send: %v\n", err)
	} else {
		fmt.Println("✓ Message sent to Telegram")
	}

	// Example 2: Broadcast to all providers (synchronous)
	fmt.Println("\nExample 2: Broadcast to all providers (synchronous)")
	errors := manager.Broadcast(ctx, "Broadcasting to all channels! 📢")
	if len(errors) > 0 {
		for _, err := range errors {
			log.Printf("Broadcast error: %v\n", err)
		}
	} else {
		fmt.Println("✓ Broadcast completed successfully")
	}

	// Example 3: Broadcast asynchronously
	fmt.Println("\nExample 3: Broadcast asynchronously")
	resultChan := manager.BroadcastAsync(ctx, "Async broadcast message!")

	successCount := 0
	failCount := 0
	for result := range resultChan {
		if result.Success {
			fmt.Printf("✓ %s: Success\n", result.Provider)
			successCount++
		} else {
			fmt.Printf("✗ %s: %v\n", result.Provider, result.Error)
			failCount++
		}
	}
	fmt.Printf("\nResults: %d succeeded, %d failed\n", successCount, failCount)

	// Example 4: Send with options to specific provider
	fmt.Println("\nExample 4: Send with options")
	msg := &notify.Message{
		Title:    "Important Alert",
		Text:     "This is a high priority message",
		Priority: notify.PriorityHigh,
	}
	err = manager.SendWithOptions(ctx, "slack", msg)
	if err != nil {
		log.Printf("Failed to send with options: %v\n", err)
	} else {
		fmt.Println("✓ Message with options sent successfully")
	}

	// Example 5: Broadcast with options asynchronously
	fmt.Println("\nExample 5: Broadcast with options asynchronously")
	richMsg := &notify.Message{
		Title:    "System Report",
		Text:     "Daily system status report",
		Priority: notify.PriorityNormal,
		Attachments: []notify.Attachment{
			{
				Title: "Metrics",
				Color: "good",
				Fields: []notify.Field{
					{Title: "CPU", Value: "45%", Short: true},
					{Title: "Memory", Value: "62%", Short: true},
					{Title: "Disk", Value: "78%", Short: true},
					{Title: "Uptime", Value: "99.9%", Short: true},
				},
			},
		},
	}

	resultChan = manager.BroadcastAsyncWithOptions(ctx, richMsg)
	for result := range resultChan {
		if result.Success {
			fmt.Printf("✓ %s: Rich message delivered\n", result.Provider)
		} else {
			fmt.Printf("✗ %s: %v\n", result.Provider, result.Error)
		}
	}
}
