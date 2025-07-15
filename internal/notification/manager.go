package notification

import (
	"devopstoolkit/youtube-automation/internal/configuration"
	"devopstoolkit/youtube-automation/internal/storage"
	"devopstoolkit/youtube-automation/internal/video"
)

// NotificationManager handles notification setup and phase monitoring
type NotificationManager struct {
	eventBus *EventBus
	enabled  bool
}

// NewNotificationManager creates a new notification manager
func NewNotificationManager() *NotificationManager {
	nm := &NotificationManager{
		eventBus: NewEventBus(),
		enabled:  configuration.GlobalSettings.Notifications.Enabled,
	}
	
	// Set up channels from configuration
	nm.setupChannels()
	
	return nm
}

// setupChannels initializes notification channels from configuration
func (nm *NotificationManager) setupChannels() {
	// Email channel
	if emailChannel := NewEmailChannelFromConfig(); emailChannel != nil {
		nm.eventBus.AddChannel(emailChannel)
	}
	
	// Slack channel
	if slackChannel := NewSlackChannelFromConfig(); slackChannel != nil {
		nm.eventBus.AddChannel(slackChannel)
	}
}

// NotifyPhaseChange detects and notifies about phase changes
func (nm *NotificationManager) NotifyPhaseChange(oldVideo, newVideo storage.Video) {
	// Check if notifications are enabled globally and for phase transitions
	if !nm.enabled || !configuration.GlobalSettings.Notifications.PhaseTransitions {
		return
	}
	
	oldPhase := video.CalculateVideoPhase(oldVideo)
	newPhase := video.CalculateVideoPhase(newVideo)
	
	// Only notify if phase actually changed
	if oldPhase != newPhase {
		nm.eventBus.PublishPhaseTransition(newVideo, oldPhase, newPhase)
	}
}

// Global notification manager instance
var globalNotificationManager *NotificationManager

// GetNotificationManager returns the global notification manager instance
func GetNotificationManager() *NotificationManager {
	if globalNotificationManager == nil {
		globalNotificationManager = NewNotificationManager()
	}
	return globalNotificationManager
}

// InitializeNotifications sets up the global notification system
func InitializeNotifications() {
	globalNotificationManager = NewNotificationManager()
}

// NotifyPhaseChange is a convenience function to notify about phase changes
func NotifyPhaseChange(oldVideo, newVideo storage.Video) {
	GetNotificationManager().NotifyPhaseChange(oldVideo, newVideo)
} 