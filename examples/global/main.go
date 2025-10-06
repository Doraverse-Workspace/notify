package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Doraverse-Workspace/notify"
)

// This example demonstrates using the global singleton pattern
// You initialize once and use anywhere in your application

func main() {
	ctx := context.Background()

	// Method 1: Initialize with Setup() - one-time initialization
	fmt.Println("=== Method 1: Using Setup() ===")
	err := notify.Setup(
		notify.TelegramConfig{
			BotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
			ChatID:   os.Getenv("TELEGRAM_CHAT_ID"),
		},
		&notify.SlackConfig{
			Token:          os.Getenv("SLACK_BOT_TOKEN"),
			DefaultChannel: os.Getenv("SLACK_CHANNEL"),
		},
	)
	if err != nil {
		log.Printf("Warning: Failed to setup notifiers: %v\n", err)
	} else {
		fmt.Println("✓ Global notifiers initialized successfully")
	}

	// Now you can use notify functions anywhere without passing manager around
	fmt.Println("\nRegistered notifiers:", notify.List())

	// Example: Send to specific provider
	fmt.Println("\n--- Sending to Telegram ---")
	err = notify.Send(ctx, "telegram", "Hello from global notify! 🌍")
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Println("✓ Sent successfully")
	}

	// Example: Broadcast to all providers
	fmt.Println("\n--- Broadcasting to all providers ---")
	errors := notify.Broadcast(ctx, "Global broadcast message! 📢")
	if len(errors) > 0 {
		for _, e := range errors {
			log.Printf("Broadcast error: %v\n", e)
		}
	} else {
		fmt.Println("✓ Broadcast successful")
	}

	// Example: Send rich message
	fmt.Println("\n--- Sending rich message ---")
	msg := &notify.Message{
		Title:    "System Alert",
		Text:     "This is a rich message from global notify",
		Priority: notify.PriorityHigh,
	}
	err = notify.SendWithOptions(ctx, "slack", msg)
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Println("✓ Rich message sent successfully")
	}

	// Example: Async broadcast
	fmt.Println("\n--- Async broadcast ---")
	resultChan := notify.BroadcastAsync(ctx, "Async message from global notify!")

	successCount := 0
	for result := range resultChan {
		if result.Success {
			fmt.Printf("✓ %s: Success\n", result.Provider)
			successCount++
		} else {
			fmt.Printf("✗ %s: %v\n", result.Provider, result.Error)
		}
	}
	fmt.Printf("Total successful: %d\n", successCount)

	// Demonstrate using in different functions
	fmt.Println("\n--- Using in different functions ---")
	sendFromAnotherFunction(ctx)
	sendFromYetAnotherFunction(ctx)
}

// You can use notify anywhere in your codebase without initialization
func sendFromAnotherFunction(ctx context.Context) {
	// No need to pass manager or config!
	err := notify.Send(ctx, "telegram", "Message from another function! 🎉")
	if err != nil {
		log.Printf("Error in another function: %v\n", err)
	} else {
		fmt.Println("✓ Sent from another function")
	}
}

func sendFromYetAnotherFunction(ctx context.Context) {
	// Still no initialization needed!
	msg := &notify.Message{
		Text:     "Message from yet another function",
		Priority: notify.PriorityNormal,
	}
	errors := notify.BroadcastWithOptions(ctx, msg)
	if len(errors) == 0 {
		fmt.Println("✓ Broadcast from yet another function")
	} else {
		log.Printf("Errors: %v\n", errors)
	}
}
