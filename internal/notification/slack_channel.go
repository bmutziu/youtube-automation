package notification

import (
	"context"
	"fmt"

	"devopstoolkit/youtube-automation/internal/slack"

	slackapi "github.com/slack-go/slack"
)

// SlackChannel implements NotificationChannel for Slack notifications
type SlackChannel struct {
	slackService *slack.SlackService
	channelIDs   []string
}

// NewSlackChannel creates a new Slack notification channel
func NewSlackChannel(slackService *slack.SlackService, channelIDs []string) *SlackChannel {
	return &SlackChannel{
		slackService: slackService,
		channelIDs:   channelIDs,
	}
}

// Send implements NotificationChannel interface
func (sc *SlackChannel) Send(ctx context.Context, event NotificationEvent) error {
	if event.Type != EventPhaseTransition {
		return nil // Only handle phase transitions in this PoC
	}

	// Create simple message text
	messageText := fmt.Sprintf("📹 Video Phase Changed: %s\n%s → %s\nCategory: %s",
		event.VideoName,
		GetPhaseName(event.OldPhase),
		GetPhaseName(event.NewPhase),
		event.Category,
	)

	// Send to all configured channels
	var lastError error
	successCount := 0
	
	for _, channelID := range sc.channelIDs {
		// Create Slack client directly for simplicity (PoC)
		if err := sc.sendToChannel(channelID, messageText); err != nil {
			lastError = err
			slack.LogSlackError(slack.CategorizeError(err), fmt.Sprintf("Failed to send phase notification to channel %s", channelID))
			continue
		}
		successCount++
	}

	if successCount == 0 && lastError != nil {
		return fmt.Errorf("failed to send to any Slack channel: %w", lastError)
	}

	return nil
}

// sendToChannel sends message to a specific Slack channel
func (sc *SlackChannel) sendToChannel(channelID, messageText string) error {
	// Load Slack config and create client
	if err := slack.LoadAndValidateSlackConfig(""); err != nil {
		return fmt.Errorf("failed to load Slack config: %w", err)
	}

	auth, err := slack.NewSlackAuth(slack.GlobalSlackConfig.Token)
	if err != nil {
		return fmt.Errorf("failed to create Slack auth: %w", err)
	}

	client, err := slack.NewSlackClient(auth)
	if err != nil {
		return fmt.Errorf("failed to create Slack client: %w", err)
	}

	_, _, err = client.PostMessage(
		channelID,
		slackapi.MsgOptionText(messageText, false),
	)

	return err
}

// GetChannelType implements NotificationChannel interface
func (sc *SlackChannel) GetChannelType() string {
	return "slack"
}

// NewSlackChannelFromConfig creates Slack channel from global configuration
func NewSlackChannelFromConfig() *SlackChannel {
	// Load Slack config
	if err := slack.LoadAndValidateSlackConfig(""); err != nil {
		return nil // Not configured
	}

	// Get target channels from configuration
	channelIDs := slack.GetTargetChannels()
	if len(channelIDs) == 0 {
		return nil // No channels configured
	}

	// Create service (though we'll use client directly for simplicity)
	service, err := slack.NewSlackService(slack.GlobalSlackConfig)
	if err != nil {
		return nil // Failed to create service
	}

	return NewSlackChannel(service, channelIDs)
} 