package notification

import (
	"context"
	"fmt"
	"time"

	"devopstoolkit/youtube-automation/internal/storage"
	"devopstoolkit/youtube-automation/internal/workflow"
)

// EventType represents different types of events
type EventType string

const (
	EventPhaseTransition EventType = "video.phase.changed"
)

// NotificationEvent represents a notification event
type NotificationEvent struct {
	Type      EventType     `json:"type"`
	VideoID   string        `json:"video_id"`
	VideoName string        `json:"video_name"`
	Category  string        `json:"category"`
	OldPhase  int           `json:"old_phase,omitempty"`
	NewPhase  int           `json:"new_phase,omitempty"`
	Timestamp time.Time     `json:"timestamp"`
	Video     storage.Video `json:"video"`
}

// NotificationChannel defines the interface for notification channels
type NotificationChannel interface {
	Send(ctx context.Context, event NotificationEvent) error
	GetChannelType() string
}

// EventBus manages event distribution
type EventBus struct {
	channels []NotificationChannel
	enabled  bool
}

// NewEventBus creates a new event bus
func NewEventBus() *EventBus {
	return &EventBus{
		channels: make([]NotificationChannel, 0),
		enabled:  true, // Simple enable/disable for PoC
	}
}

// AddChannel adds a notification channel to the bus
func (eb *EventBus) AddChannel(channel NotificationChannel) {
	eb.channels = append(eb.channels, channel)
}

// PublishPhaseTransition publishes a phase transition event
func (eb *EventBus) PublishPhaseTransition(video storage.Video, oldPhase, newPhase int) {
	if !eb.enabled || len(eb.channels) == 0 {
		return
	}

	event := NotificationEvent{
		Type:      EventPhaseTransition,
		VideoID:   fmt.Sprintf("%s/%s", video.Category, video.Name),
		VideoName: video.Name,
		Category:  video.Category,
		OldPhase:  oldPhase,
		NewPhase:  newPhase,
		Timestamp: time.Now(),
		Video:     video,
	}

	ctx := context.Background()
	for _, channel := range eb.channels {
		go func(ch NotificationChannel) {
			if err := ch.Send(ctx, event); err != nil {
				// Basic error logging - could be enhanced
				fmt.Printf("Failed to send notification via %s: %v\n", ch.GetChannelType(), err)
			}
		}(channel)
	}
}

// GetPhaseName returns human-readable phase name
func GetPhaseName(phase int) string {
	return workflow.PhaseNames[phase]
} 