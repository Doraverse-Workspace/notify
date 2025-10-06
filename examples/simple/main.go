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

	// Example 1: Simple Telegram notification
	fmt.Println("Example 1: Simple Telegram notification")
	telegramNotifier, err := notify.NewTelegramNotifier(notify.TelegramConfig{
		BotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		ChatID:   os.Getenv("TELEGRAM_CHAT_ID"),
	})
	if err != nil {
		log.Printf("Failed to create Telegram notifier: %v\n", err)
	} else {
		err = telegramNotifier.Send(ctx, "Hello from notify library! 🚀")
		if err != nil {
			log.Printf("Failed to send Telegram message: %v\n", err)
		} else {
			fmt.Println("✓ Telegram message sent successfully")
		}
	}

	// Example 2: Simple Slack notification
	fmt.Println("\nExample 2: Simple Slack notification")
	slackNotifier, err := notify.NewSlackNotifier(&notify.SlackConfig{
		Token:          os.Getenv("SLACK_BOT_TOKEN"),
		DefaultChannel: os.Getenv("SLACK_CHANNEL"),
	})
	if err != nil {
		log.Printf("Failed to create Slack notifier: %v\n", err)
	} else {
		err = slackNotifier.Send(ctx, "Hello from notify library! 🚀")
		if err != nil {
			log.Printf("Failed to send Slack message: %v\n", err)
		} else {
			fmt.Println("✓ Slack message sent successfully")
		}
	}

	// Example 3: Rich message with options
	fmt.Println("\nExample 3: Rich message with options")
	if telegramNotifier != nil {
		msg := &notify.Message{
			Title:    "System Alert",
			Text:     "Server CPU usage is at 85%",
			Priority: notify.PriorityHigh,
		}
		err = telegramNotifier.SendWithOptions(ctx, msg)
		if err != nil {
			log.Printf("Failed to send rich message: %v\n", err)
		} else {
			fmt.Println("✓ Rich message sent successfully")
		}
	}

	// Example 4: Message with attachments (Slack)
	fmt.Println("\nExample 4: Message with attachments")
	if slackNotifier != nil {
		msg := &notify.Message{
			Title: "Deployment Status",
			Text:  "Application deployed successfully",
			Attachments: []notify.Attachment{
				{
					Title: "Details",
					Color: "good",
					Fields: []notify.Field{
						{Title: "Version", Value: "v1.2.3", Short: true},
						{Title: "Environment", Value: "Production", Short: true},
						{Title: "Duration", Value: "2m 34s", Short: true},
					},
				},
			},
		}
		err = slackNotifier.SendWithOptions(ctx, msg)
		if err != nil {
			log.Printf("Failed to send message with attachments: %v\n", err)
		} else {
			fmt.Println("✓ Message with attachments sent successfully")
		}
	}
}
