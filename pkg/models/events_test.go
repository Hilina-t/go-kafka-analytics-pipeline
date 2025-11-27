package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestAnalyticsEventJSON(t *testing.T) {
	event := AnalyticsEvent{
		ID:        "test-123",
		Type:      PageView,
		Timestamp: time.Now(),
		UserID:    "user-456",
		SessionID: "session-789",
		URL:       "https://example.com",
		Path:      "/home",
		Metadata: map[string]interface{}{
			"test": "value",
		},
	}

	// Test marshaling
	data, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("Failed to marshal event: %v", err)
	}

	// Test unmarshaling
	var decoded AnalyticsEvent
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal event: %v", err)
	}

	// Verify fields
	if decoded.ID != event.ID {
		t.Errorf("ID mismatch: got %s, want %s", decoded.ID, event.ID)
	}
	if decoded.Type != event.Type {
		t.Errorf("Type mismatch: got %s, want %s", decoded.Type, event.Type)
	}
	if decoded.UserID != event.UserID {
		t.Errorf("UserID mismatch: got %s, want %s", decoded.UserID, event.UserID)
	}
}

func TestEventTypes(t *testing.T) {
	tests := []struct {
		name      string
		eventType EventType
		expected  string
	}{
		{"PageView", PageView, "page_view"},
		{"Click", Click, "click"},
		{"Session", Session, "session"},
		{"UserEvent", UserEvent, "user_event"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.eventType) != tt.expected {
				t.Errorf("EventType mismatch: got %s, want %s", tt.eventType, tt.expected)
			}
		})
	}
}

func TestPageViewEvent(t *testing.T) {
	event := PageViewEvent{
		AnalyticsEvent: AnalyticsEvent{
			ID:        "pv-123",
			Type:      PageView,
			Timestamp: time.Now(),
			UserID:    "user-123",
			SessionID: "session-123",
			URL:       "https://example.com/home",
			Path:      "/home",
		},
		PageTitle:    "Home Page",
		LoadTime:     1200,
		ScreenWidth:  1920,
		ScreenHeight: 1080,
	}

	// Test marshaling
	data, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("Failed to marshal page view event: %v", err)
	}

	// Verify the data can be unmarshaled
	var decoded PageViewEvent
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal page view event: %v", err)
	}

	if decoded.PageTitle != event.PageTitle {
		t.Errorf("PageTitle mismatch: got %s, want %s", decoded.PageTitle, event.PageTitle)
	}
	if decoded.LoadTime != event.LoadTime {
		t.Errorf("LoadTime mismatch: got %d, want %d", decoded.LoadTime, event.LoadTime)
	}
}

func TestClickEvent(t *testing.T) {
	event := ClickEvent{
		AnalyticsEvent: AnalyticsEvent{
			ID:        "click-123",
			Type:      Click,
			Timestamp: time.Now(),
			UserID:    "user-123",
			SessionID: "session-123",
			URL:       "https://example.com/products",
			Path:      "/products",
		},
		ElementID:   "buy-button",
		ElementType: "button",
		ElementText: "Buy Now",
		XPosition:   100,
		YPosition:   200,
	}

	// Test marshaling
	data, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("Failed to marshal click event: %v", err)
	}

	// Verify the data can be unmarshaled
	var decoded ClickEvent
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal click event: %v", err)
	}

	if decoded.ElementID != event.ElementID {
		t.Errorf("ElementID mismatch: got %s, want %s", decoded.ElementID, event.ElementID)
	}
	if decoded.ElementText != event.ElementText {
		t.Errorf("ElementText mismatch: got %s, want %s", decoded.ElementText, event.ElementText)
	}
}

func TestSessionEvent(t *testing.T) {
	event := SessionEvent{
		AnalyticsEvent: AnalyticsEvent{
			ID:        "session-123",
			Type:      Session,
			Timestamp: time.Now(),
			UserID:    "user-123",
			SessionID: "session-123",
			URL:       "https://example.com",
			Path:      "/",
		},
		Duration:  300,
		PageCount: 5,
		Device:    "desktop",
		Browser:   "Chrome",
	}

	// Test marshaling
	data, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("Failed to marshal session event: %v", err)
	}

	// Verify the data can be unmarshaled
	var decoded SessionEvent
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal session event: %v", err)
	}

	if decoded.Duration != event.Duration {
		t.Errorf("Duration mismatch: got %d, want %d", decoded.Duration, event.Duration)
	}
	if decoded.Device != event.Device {
		t.Errorf("Device mismatch: got %s, want %s", decoded.Device, event.Device)
	}
}
