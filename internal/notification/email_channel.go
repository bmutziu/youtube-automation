package notification

import (
	"context"
	"fmt"

	"devopstoolkit/youtube-automation/internal/configuration"
)

// EmailChannel implements NotificationChannel for email notifications
type EmailChannel struct {
	emailService *Email
	from         string
	to           string
}

// NewEmailChannel creates a new email notification channel
func NewEmailChannel(password, from, to string) *EmailChannel {
	return &EmailChannel{
		emailService: NewEmail(password),
		from:         from,
		to:           to,
	}
}

// Send implements NotificationChannel interface
func (ec *EmailChannel) Send(ctx context.Context, event NotificationEvent) error {
	if event.Type != EventPhaseTransition {
		return nil // Only handle phase transitions in this PoC
	}

	subject := fmt.Sprintf("Video Phase Changed: %s", event.VideoName)
	body := fmt.Sprintf(`
<h2>Video Phase Transition</h2>
<p><strong>Video:</strong> %s</p>
<p><strong>Category:</strong> %s</p>
<p><strong>Phase Change:</strong> %s → %s</p>
<p><strong>Time:</strong> %s</p>
<hr>
<p><em>This is an automated notification from YouTube Automation Tool</em></p>
`, 
		event.VideoName,
		event.Category,
		GetPhaseName(event.OldPhase),
		GetPhaseName(event.NewPhase),
		event.Timestamp.Format("2006-01-02 15:04:05"),
	)

	return ec.emailService.Send(ec.from, []string{ec.to}, subject, body, "")
}

// GetChannelType implements NotificationChannel interface
func (ec *EmailChannel) GetChannelType() string {
	return "email"
}

// NewEmailChannelFromConfig creates email channel from global configuration
func NewEmailChannelFromConfig() *EmailChannel {
	config := configuration.GlobalSettings.Email
	if config.Password == "" || config.From == "" || config.ThumbnailTo == "" {
		return nil // Not configured
	}
	
	return NewEmailChannel(config.Password, config.From, config.ThumbnailTo)
} 