package models

import "time"

// EventType represents the type of analytics event
type EventType string

const (
	PageView  EventType = "page_view"
	Click     EventType = "click"
	Session   EventType = "session"
	UserEvent EventType = "user_event"
)

// AnalyticsEvent represents a website analytics event
type AnalyticsEvent struct {
	ID        string                 `json:"id"`
	Type      EventType              `json:"type"`
	Timestamp time.Time              `json:"timestamp"`
	UserID    string                 `json:"user_id"`
	SessionID string                 `json:"session_id"`
	URL       string                 `json:"url"`
	Path      string                 `json:"path"`
	Referrer  string                 `json:"referrer,omitempty"`
	UserAgent string                 `json:"user_agent,omitempty"`
	IPAddress string                 `json:"ip_address,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// PageViewEvent represents a page view event
type PageViewEvent struct {
	AnalyticsEvent
	PageTitle    string `json:"page_title,omitempty"`
	LoadTime     int64  `json:"load_time,omitempty"` // in milliseconds
	ScreenWidth  int    `json:"screen_width,omitempty"`
	ScreenHeight int    `json:"screen_height,omitempty"`
}

// ClickEvent represents a click event
type ClickEvent struct {
	AnalyticsEvent
	ElementID   string `json:"element_id,omitempty"`
	ElementType string `json:"element_type,omitempty"`
	ElementText string `json:"element_text,omitempty"`
	XPosition   int    `json:"x_position,omitempty"`
	YPosition   int    `json:"y_position,omitempty"`
}

// SessionEvent represents a user session event
type SessionEvent struct {
	AnalyticsEvent
	Duration int64  `json:"duration,omitempty"` // in seconds
	PageCount int   `json:"page_count,omitempty"`
	Device    string `json:"device,omitempty"`
	Browser   string `json:"browser,omitempty"`
}
