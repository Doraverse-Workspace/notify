# Quick Start Guide

Get started with `notify` in 5 minutes! 🚀

## Installation

```bash
go get github.com/Doraverse-Workspace/notify
```

## Recommended: Global Singleton Pattern ⭐

**Initialize once, use everywhere!** No need to pass configs around.

```go
package main

import (
    "context"
    "log"
    "os"
    "github.com/Doraverse-Workspace/notify"
)

func main() {
    // Initialize ONCE at startup
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
        log.Fatal(err)
    }
    
    // Use ANYWHERE in your app - no init needed!
    notify.Send(context.Background(), "telegram", "Hello! 🎉")
    notify.Broadcast(context.Background(), "Broadcasting to all! 📢")
    
    // Call from other functions
    sendAlert()
    processOrder()
}

func sendAlert() {
    // No need to pass config or manager!
    notify.Send(context.Background(), "telegram", "Alert message!")
}

func processOrder() {
    // Works everywhere!
    msg := &notify.Message{
        Title:    "New Order",
        Text:     "Order #123 processed",
        Priority: notify.PriorityHigh,
    }
    notify.BroadcastWithOptions(context.Background(), msg)
}
```

**Benefits:**
- ✅ Initialize once, use everywhere
- ✅ No need to pass configs around
- ✅ Clean and simple code
- ✅ Thread-safe

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

## Global Functions Reference

When using the global singleton pattern, these functions are available:

```go
// Setup & Management
notify.Setup(configs...)        // Initialize with configs
notify.Init()                   // Initialize empty manager
notify.Register(notifier)       // Register a notifier
notify.Unregister(name)         // Unregister a notifier
notify.Get(name)                // Get specific notifier
notify.List()                   // List all notifiers
notify.Global()                 // Get global manager instance

// Send Messages
notify.Send(ctx, provider, message)
notify.SendWithOptions(ctx, provider, msg)
notify.Broadcast(ctx, message)
notify.BroadcastWithOptions(ctx, msg)
notify.BroadcastAsync(ctx, message)
notify.BroadcastAsyncWithOptions(ctx, msg)
```

## Common Use Cases

### 1. Server Monitoring (using global functions)
```go
if cpuUsage > 80 {
    notify.Send(ctx, "telegram", 
        fmt.Sprintf("⚠️ High CPU usage: %.1f%%", cpuUsage))
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
notify.BroadcastWithOptions(ctx, msg)
```

### 3. Error Alerts
```go
if err != nil {
    notify.Send(ctx, "telegram", fmt.Sprintf("🔥 Error: %v", err))
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
    notify.BroadcastWithOptions(context.Background(), msg)
}
```

## Testing

Run the examples:

```bash
# Set environment variables first
export TELEGRAM_BOT_TOKEN="your_token"
export TELEGRAM_CHAT_ID="your_chat_id"
export SLACK_BOT_TOKEN="your_token"
export SLACK_CHANNEL="#general"

# Run global singleton example (recommended)
cd examples/global
go run main.go

# Or run simple example
cd examples/simple
go run main.go
```

## Next Steps

- Read the full [README.md](README.md) for detailed documentation
- Check out [examples/global/](examples/global/) for the global singleton pattern
- Explore [examples/](examples/) for more complete examples
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

