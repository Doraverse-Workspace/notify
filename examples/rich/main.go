package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Doraverse-Workspace/notify"
	"github.com/slack-go/slack"
)

func main() {
	ctx := context.Background()

	// Example 1: Using global interface with SendRichMessage for Slack
	fmt.Println("Example 1: Global SendRichMessage for Slack")
	
	// Setup global notifiers
	err := notify.Setup(&notify.SlackConfig{
		Token:          os.Getenv("SLACK_BOT_TOKEN"),
		DefaultChannel: os.Getenv("SLACK_CHANNEL"),
	})
	if err != nil {
		log.Printf("Failed to setup Slack notifier: %v\n", err)
		return
	}

	// Create rich Slack blocks
	blocks := []slack.Block{
		slack.NewHeaderBlock(
			slack.NewTextBlockObject("plain_text", "🚀 Deployment Notification", false, false),
		),
		slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn", "Your application has been successfully deployed!", false, false),
			nil, nil,
		),
		slack.NewDividerBlock(),
		slack.NewSectionBlock(
			nil,
			[]*slack.TextBlockObject{
				slack.NewTextBlockObject("mrkdwn", "*Version:*\nv1.2.3", false, false),
				slack.NewTextBlockObject("mrkdwn", "*Environment:*\nProduction", false, false),
				slack.NewTextBlockObject("mrkdwn", "*Duration:*\n2m 34s", false, false),
			},
			nil,
		),
		slack.NewContextBlock("", slack.NewTextBlockObject("mrkdwn", "Deployed by: @john.doe", false, false)),
	}

	// Send rich message using global interface
	err = notify.SendRichMessage(ctx, "slack", "", blocks)
	if err != nil {
		log.Printf("Failed to send rich message: %v\n", err)
	} else {
		fmt.Println("✓ Rich Slack message sent successfully")
	}

	// Example 2: Using global interface with SendRichMessage for Telegram
	fmt.Println("\nExample 2: Global SendRichMessage for Telegram")
	
	// Setup Telegram notifier
	err = notify.Setup(notify.TelegramConfig{
		BotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		ChatID:   os.Getenv("TELEGRAM_CHAT_ID"),
	})
	if err != nil {
		log.Printf("Failed to setup Telegram notifier: %v\n", err)
		return
	}

	// For Telegram, we can pass simple text or structured data
	telegramBlocks := []string{
		"🚀 Deployment Notification",
		"Your application has been successfully deployed!",
		"",
		"Version: v1.2.3",
		"Environment: Production", 
		"Duration: 2m 34s",
		"",
		"Deployed by: @john.doe",
	}

	// Send rich message using global interface
	err = notify.SendRichMessage(ctx, "telegram", "", telegramBlocks)
	if err != nil {
		log.Printf("Failed to send rich message: %v\n", err)
	} else {
		fmt.Println("✓ Rich Telegram message sent successfully")
	}

	// Example 3: Direct notifier usage
	fmt.Println("\nExample 3: Direct notifier usage")
	
	slackNotifier, err := notify.NewSlackNotifier(&notify.SlackConfig{
		Token:          os.Getenv("SLACK_BOT_TOKEN"),
		DefaultChannel: os.Getenv("SLACK_CHANNEL"),
	})
	if err != nil {
		log.Printf("Failed to create Slack notifier: %v\n", err)
		return
	}

	// Create a simple rich message
	simpleBlocks := []slack.Block{
		slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn", "*Hello from direct notifier!*", false, false),
			nil, nil,
		),
	}

	err = slackNotifier.SendRichMessage(ctx, "", simpleBlocks)
	if err != nil {
		log.Printf("Failed to send rich message: %v\n", err)
	} else {
		fmt.Println("✓ Direct rich message sent successfully")
	}
}
