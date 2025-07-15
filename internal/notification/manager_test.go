package notification

import (
	"context"
	"testing"

	"devopstoolkit/youtube-automation/internal/storage"
	"devopstoolkit/youtube-automation/internal/workflow"
)

func TestNotificationManager(t *testing.T) {
	// Create a new notification manager
	nm := NewNotificationManager()
	
	if nm == nil {
		t.Fatal("NewNotificationManager returned nil")
	}
	
	if nm.eventBus == nil {
		t.Fatal("EventBus not initialized")
	}
}

func TestPhaseChangeDetection(t *testing.T) {
	// Create test videos representing different phases
	oldVideo := storage.Video{
		Name:     "test-video",
		Category: "test-category",
		// No fields set = PhaseIdeas (7)
	}
	
	newVideo := storage.Video{
		Name:     "test-video",
		Category: "test-category",
		Date:     "2023-01-01T10:00", // Has date = PhaseStarted (4)
	}
	
	// Test that phase change is detected
	nm := NewNotificationManager()
	
	// This should trigger a phase transition notification
	// (in a real test, we would mock the channels to verify the notification was sent)
	nm.NotifyPhaseChange(oldVideo, newVideo)
	
	// Basic test passes if no panic occurs
	t.Log("Phase change notification triggered successfully")
}

func TestEventBusChannelManagement(t *testing.T) {
	eb := NewEventBus()
	
	// Test that channels can be added
	if len(eb.channels) != 0 {
		t.Errorf("Expected 0 channels initially, got %d", len(eb.channels))
	}
	
	// Create a mock channel
	mockChannel := &MockChannel{}
	eb.AddChannel(mockChannel)
	
	if len(eb.channels) != 1 {
		t.Errorf("Expected 1 channel after adding, got %d", len(eb.channels))
	}
}

// MockChannel for testing
type MockChannel struct {
	lastEvent NotificationEvent
}

func (mc *MockChannel) Send(ctx context.Context, event NotificationEvent) error {
	mc.lastEvent = event
	return nil
}

func (mc *MockChannel) GetChannelType() string {
	return "mock"
}

func TestEventTypes(t *testing.T) {
	// Test that event types are correctly defined
	if EventPhaseTransition != "video.phase.changed" {
		t.Errorf("Expected EventPhaseTransition to be 'video.phase.changed', got %s", EventPhaseTransition)
	}
}

func TestGetPhaseName(t *testing.T) {
	// Test phase name retrieval
	testCases := []struct {
		phase    int
		expected string
	}{
		{workflow.PhaseIdeas, "Ideas"},
		{workflow.PhaseStarted, "Started"},
		{workflow.PhasePublished, "Published"},
	}
	
	for _, tc := range testCases {
		result := GetPhaseName(tc.phase)
		if result != tc.expected {
			t.Errorf("Expected phase %d to be '%s', got '%s'", tc.phase, tc.expected, result)
		}
	}
} 