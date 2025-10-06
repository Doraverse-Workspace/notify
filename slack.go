package notify

import (
	"context"
	"fmt"

	"github.com/slack-go/slack"
)

// SlackNotifier sends notifications via Slack API
type SlackNotifier struct {
	client         *slack.Client
	defaultChannel string
	username       string
	iconEmoji      string
}

// SlackConfig holds configuration for Slack notifications
type SlackConfig struct {
	// Token is the Slack Bot or User OAuth token
	Token string

	// DefaultChannel is the default channel to send messages to (e.g., #general or @username)
	DefaultChannel string

	// Username is the bot username (optional)
	Username string

	// IconEmoji is the bot icon emoji (optional, e.g., :robot_face:)
	IconEmoji string

	// WebhookURL for incoming webhooks (alternative to Token)
	WebhookURL string
}

// NewSlackNotifier creates a new Slack notifier
func NewSlackNotifier(config SlackConfig) (*SlackNotifier, error) {
	if config.Token == "" && config.WebhookURL == "" {
		return nil, &NotificationError{
			Provider: "slack",
			Message:  "either token or webhook URL is required",
		}
	}

	var client *slack.Client
	if config.Token != "" {
		client = slack.New(config.Token)
	} else {
		// For webhook, we'll handle it differently in Send methods
		client = nil
	}

	return &SlackNotifier{
		client:         client,
		defaultChannel: config.DefaultChannel,
		username:       config.Username,
		iconEmoji:      config.IconEmoji,
	}, nil
}

// Name returns the name of the provider
func (s *SlackNotifier) Name() string {
	return "slack"
}

// Send sends a simple text message
func (s *SlackNotifier) Send(ctx context.Context, message string) error {
	return s.SendWithOptions(ctx, &Message{
		Text:    message,
		Channel: s.defaultChannel,
	})
}

// SendWithOptions sends a message with additional options
func (s *SlackNotifier) SendWithOptions(ctx context.Context, msg *Message) error {
	if s.client == nil {
		return &NotificationError{
			Provider: "slack",
			Message:  "slack client not initialized (webhook support not yet implemented for SendWithOptions)",
		}
	}

	if msg.Text == "" {
		return &NotificationError{
			Provider: "slack",
			Message:  "message text is required",
		}
	}

	channel := msg.Channel
	if channel == "" {
		channel = s.defaultChannel
	}

	if channel == "" {
		return &NotificationError{
			Provider: "slack",
			Message:  "channel is required",
		}
	}

	// Build message options
	options := []slack.MsgOption{
		slack.MsgOptionText(msg.Text, false),
	}

	if s.username != "" {
		options = append(options, slack.MsgOptionUsername(s.username))
	}

	if s.iconEmoji != "" {
		options = append(options, slack.MsgOptionIconEmoji(s.iconEmoji))
	}

	// Add attachments if present
	if len(msg.Attachments) > 0 {
		slackAttachments := s.convertAttachments(msg.Attachments)
		options = append(options, slack.MsgOptionAttachments(slackAttachments...))
	}

	// Add title as a block if present
	if msg.Title != "" {
		blocks := []slack.Block{
			slack.NewHeaderBlock(
				slack.NewTextBlockObject("plain_text", msg.Title, false, false),
			),
			slack.NewSectionBlock(
				slack.NewTextBlockObject("mrkdwn", msg.Text, false, false),
				nil, nil,
			),
		}
		options = append(options, slack.MsgOptionBlocks(blocks...))
		// Remove text option when using blocks
		options = options[1:]
	}

	_, _, err := s.client.PostMessageContext(ctx, channel, options...)
	if err != nil {
		return &NotificationError{
			Provider: "slack",
			Message:  "failed to send message",
			Err:      err,
		}
	}

	return nil
}

// SendRichMessage sends a message with blocks for rich formatting
func (s *SlackNotifier) SendRichMessage(ctx context.Context, channel string, blocks []slack.Block) error {
	if s.client == nil {
		return &NotificationError{
			Provider: "slack",
			Message:  "slack client not initialized",
		}
	}

	if channel == "" {
		channel = s.defaultChannel
	}

	_, _, err := s.client.PostMessageContext(
		ctx,
		channel,
		slack.MsgOptionBlocks(blocks...),
	)
	if err != nil {
		return &NotificationError{
			Provider: "slack",
			Message:  "failed to send rich message",
			Err:      err,
		}
	}

	return nil
}

// convertAttachments converts generic attachments to Slack attachments
func (s *SlackNotifier) convertAttachments(attachments []Attachment) []slack.Attachment {
	slackAttachments := make([]slack.Attachment, len(attachments))

	for i, att := range attachments {
		slackAtt := slack.Attachment{
			Title:      att.Title,
			Text:       att.Text,
			ImageURL:   att.ImageURL,
			Color:      att.Color,
			Footer:     att.Footer,
			FooterIcon: att.FooterIcon,
		}

		// Convert fields
		if len(att.Fields) > 0 {
			slackAtt.Fields = make([]slack.AttachmentField, len(att.Fields))
			for j, field := range att.Fields {
				slackAtt.Fields[j] = slack.AttachmentField{
					Title: field.Title,
					Value: field.Value,
					Short: field.Short,
				}
			}
		}

		slackAttachments[i] = slackAtt
	}

	return slackAttachments
}

// GetClient returns the underlying Slack client for advanced usage
func (s *SlackNotifier) GetClient() *slack.Client {
	return s.client
}

// SendFile uploads a file to Slack
func (s *SlackNotifier) SendFile(ctx context.Context, channel, filePath, title, comment string) error {
	if s.client == nil {
		return &NotificationError{
			Provider: "slack",
			Message:  "slack client not initialized",
		}
	}

	if channel == "" {
		channel = s.defaultChannel
	}

	params := slack.FileUploadParameters{
		File:           filePath,
		Channels:       []string{channel},
		Title:          title,
		InitialComment: comment,
	}

	_, err := s.client.UploadFileContext(ctx, params)
	if err != nil {
		return &NotificationError{
			Provider: "slack",
			Message:  fmt.Sprintf("failed to upload file: %v", err),
			Err:      err,
		}
	}

	return nil
}
