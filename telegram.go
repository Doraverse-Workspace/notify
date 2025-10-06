package notify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// TelegramNotifier sends notifications via Telegram Bot API
type TelegramNotifier struct {
	botToken  string
	chatID    string
	client    *http.Client
	parseMode string
}

// TelegramConfig holds configuration for Telegram notifications
type TelegramConfig struct {
	// BotToken is the Telegram Bot API token
	BotToken string

	// ChatID is the default chat ID to send messages to
	ChatID string

	// ParseMode defines how to parse the message (Markdown, HTML, or empty for plain text)
	ParseMode string

	// HTTPClient allows custom HTTP client (optional)
	HTTPClient *http.Client
}

// NewTelegramNotifier creates a new Telegram notifier
func NewTelegramNotifier(config TelegramConfig) (*TelegramNotifier, error) {
	if config.BotToken == "" {
		return nil, &NotificationError{
			Provider: "telegram",
			Message:  "bot token is required",
		}
	}

	if config.ChatID == "" {
		return nil, &NotificationError{
			Provider: "telegram",
			Message:  "chat ID is required",
		}
	}

	client := config.HTTPClient
	if client == nil {
		client = &http.Client{
			Timeout: 30 * time.Second,
		}
	}

	parseMode := config.ParseMode
	if parseMode == "" {
		parseMode = "Markdown"
	}

	return &TelegramNotifier{
		botToken:  config.BotToken,
		chatID:    config.ChatID,
		client:    client,
		parseMode: parseMode,
	}, nil
}

// Name returns the name of the provider
func (t *TelegramNotifier) Name() string {
	return "telegram"
}

// Send sends a simple text message
func (t *TelegramNotifier) Send(ctx context.Context, message string) error {
	return t.SendWithOptions(ctx, &Message{
		Text:    message,
		Channel: t.chatID,
	})
}

// SendWithOptions sends a message with additional options
func (t *TelegramNotifier) SendWithOptions(ctx context.Context, msg *Message) error {
	if msg.Text == "" {
		return &NotificationError{
			Provider: "telegram",
			Message:  "message text is required",
		}
	}

	chatID := msg.Channel
	if chatID == "" {
		chatID = t.chatID
	}

	// Build the message text
	messageText := msg.Text
	if msg.Title != "" {
		messageText = fmt.Sprintf("*%s*\n\n%s", msg.Title, msg.Text)
	}

	payload := map[string]interface{}{
		"chat_id":    chatID,
		"text":       messageText,
		"parse_mode": t.parseMode,
	}

	// Add priority-based notification settings
	if msg.Priority == PriorityLow {
		payload["disable_notification"] = true
	}

	return t.sendRequest(ctx, "sendMessage", payload)
}

// SendPhoto sends a photo with caption
func (t *TelegramNotifier) SendPhoto(ctx context.Context, chatID, photoURL, caption string) error {
	if chatID == "" {
		chatID = t.chatID
	}

	payload := map[string]interface{}{
		"chat_id": chatID,
		"photo":   photoURL,
		"caption": caption,
	}

	return t.sendRequest(ctx, "sendPhoto", payload)
}

// sendRequest sends a request to the Telegram Bot API
func (t *TelegramNotifier) sendRequest(ctx context.Context, method string, payload map[string]interface{}) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/%s", t.botToken, method)

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return &NotificationError{
			Provider: "telegram",
			Message:  "failed to marshal request",
			Err:      err,
		}
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return &NotificationError{
			Provider: "telegram",
			Message:  "failed to create request",
			Err:      err,
		}
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := t.client.Do(req)
	if err != nil {
		return &NotificationError{
			Provider: "telegram",
			Message:  "failed to send request",
			Err:      err,
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &NotificationError{
			Provider: "telegram",
			Message:  "failed to read response",
			Err:      err,
		}
	}

	if resp.StatusCode != http.StatusOK {
		return &NotificationError{
			Provider: "telegram",
			Message:  fmt.Sprintf("API request failed with status %d: %s", resp.StatusCode, string(body)),
		}
	}

	var result struct {
		Ok          bool   `json:"ok"`
		Description string `json:"description"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return &NotificationError{
			Provider: "telegram",
			Message:  "failed to parse response",
			Err:      err,
		}
	}

	if !result.Ok {
		return &NotificationError{
			Provider: "telegram",
			Message:  fmt.Sprintf("API returned error: %s", result.Description),
		}
	}

	return nil
}
