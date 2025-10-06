# Quick Start Guide

Get started with `notify` in 5 minutes! 🚀

## Installation

```bash
go get github.com/Doraverse-Workspace/notify
```

## 1. Telegram Notifications

### Setup
1. Create a bot with [@BotFather](https://t.me/botfather)
2. Get your bot token
3. Get your chat ID from `https://api.telegram.org/bot<TOKEN>/getUpdates`

### Code
```go
package main

import (
    "context"
    "log"
    "github.com/Doraverse-Workspace/notify"
)

func main() {
    telegram, err := notify.NewTelegramNotifier(notify.TelegramConfig{
        BotToken: "YOUR_BOT_TOKEN",
        ChatID:   "YOUR_CHAT_ID",
    })
    if err != nil {
        log.Fatal(err)
    }

    err = telegram.Send(context.Background(), "Hello from notify! 🎉")
    if err != nil {
        log.Fatal(err)
    }
}
```

## 2. Slack Notifications

### Setup
1. Create a Slack App at [api.slack.com/apps](https://api.slack.com/apps)
2. Add Bot Token Scope: `chat:write`
3. Install app to workspace
4. Copy Bot User OAuth Token

### Code
```go
package main

import (
    "context"
    "log"
    "github.com/Doraverse-Workspace/notify"
)

func main() {
    slack, err := notify.NewSlackNotifier(notify.SlackConfig{
        Token:          "xoxb-your-token",
        DefaultChannel: "#general",
    })
    if err != nil {
        log.Fatal(err)
    }

    err = slack.Send(context.Background(), "Hello from notify! 🎉")
    if err != nil {
        log.Fatal(err)
    }
}
```

## 3. Multiple Providers with Manager

```go
package main

import (
    "context"
    "log"
    "os"
    "github.com/Doraverse-Workspace/notify"
)

func main() {
    manager := notify.NewManager()

    // Add Telegram
    telegram, _ := notify.NewTelegramNotifier(notify.TelegramConfig{
        BotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
        ChatID:   os.Getenv("TELEGRAM_CHAT_ID"),
    })
    manager.Register(telegram)

    // Add Slack
    slack, _ := notify.NewSlackNotifier(notify.SlackConfig{
        Token:          os.Getenv("SLACK_BOT_TOKEN"),
        DefaultChannel: "#general",
    })
    manager.Register(slack)

    // Broadcast to all
    manager.Broadcast(context.Background(), "Broadcasting to all! 📢")

    // Or send to specific provider
    manager.Send(context.Background(), "telegram", "Only to Telegram")
}
```

## 4. Rich Messages

```go
msg := &notify.Message{
    Title:    "🚨 Alert",
    Text:     "CPU usage is high!",
    Priority: notify.PriorityHigh,
    Attachments: []notify.Attachment{
        {
            Title: "Details",
            Color: "danger",
            Fields: []notify.Field{
                {Title: "Server", Value: "prod-01", Short: true},
                {Title: "CPU", Value: "92%", Short: true},
            },
        },
    },
}

notifier.SendWithOptions(context.Background(), msg)
```

## Environment Variables

Create a `.env` file:

```bash
# Telegram
TELEGRAM_BOT_TOKEN=your_bot_token
TELEGRAM_CHAT_ID=your_chat_id

# Slack
SLACK_BOT_TOKEN=xoxb-your-token
SLACK_CHANNEL=#general
```

Load in your code:
```go
import _ "github.com/joho/godotenv/autoload"
```

## Common Use Cases

### 1. Server Monitoring
```go
if cpuUsage > 80 {
    notifier.Send(ctx, "⚠️ High CPU usage: " + fmt.Sprintf("%.1f%%", cpuUsage))
}
```

### 2. Deployment Notifications
```go
msg := &notify.Message{
    Title: "✅ Deployment Complete",
    Text:  "Version v1.2.3 deployed to production",
    Attachments: []notify.Attachment{
        {
            Fields: []notify.Field{
                {Title: "Version", Value: "v1.2.3"},
                {Title: "Duration", Value: "2m 15s"},
            },
        },
    },
}
manager.BroadcastWithOptions(ctx, msg)
```

### 3. Error Alerts
```go
if err != nil {
    telegram.Send(ctx, fmt.Sprintf("🔥 Error: %v", err))
}
```

### 4. Daily Reports
```go
func sendDailyReport() {
    msg := &notify.Message{
        Title: "📊 Daily Report",
        Text:  "Here's today's summary",
        Attachments: []notify.Attachment{
            {
                Fields: []notify.Field{
                    {Title: "Users", Value: "1,234"},
                    {Title: "Revenue", Value: "$5,678"},
                },
            },
        },
    }
    manager.BroadcastWithOptions(context.Background(), msg)
}
```

## Testing

Run the examples:

```bash
# Set environment variables first
export TELEGRAM_BOT_TOKEN="your_token"
export TELEGRAM_CHAT_ID="your_chat_id"

# Run example
cd examples/simple
go run main.go
```

## Next Steps

- Read the full [README.md](README.md) for detailed documentation
- Check out [examples/](examples/) for more complete examples
- Learn about [custom providers](examples/custom/main.go)
- Read [CONTRIBUTING.md](CONTRIBUTING.md) to add new providers

## Troubleshooting

### Telegram: "Unauthorized"
- Check your bot token is correct
- Make sure you've started a chat with your bot

### Slack: "not_in_channel"
- Invite your bot to the channel first
- Use `/invite @YourBot` in Slack

### Connection Timeout
- Add a timeout context:
  ```go
  ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
  defer cancel()
  ```

## Need Help?

- Open an issue on GitHub
- Check the [examples](examples/) directory
- Read the full documentation in [README.md](README.md)

Happy notifying! 🎉

