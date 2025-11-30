package models

import (
	"sync"
	"time"
)

// MetricsSnapshot represents a point-in-time analytics snapshot
type MetricsSnapshot struct {
	Timestamp          time.Time           `json:"timestamp"`
	TotalEvents        int64               `json:"total_events"`
	UniqueUsers        int64               `json:"unique_users"`
	ActiveSessions     int64               `json:"active_sessions"`
	EventsByType       map[EventType]int64 `json:"events_by_type"`
	TopPages           []PageMetric        `json:"top_pages"`
	TrafficSources     []TrafficSource     `json:"traffic_sources"`
	DeviceStats        map[string]int64    `json:"device_stats"`
	BrowserStats       map[string]int64    `json:"browser_stats"`
	HourlyPageViews    []HourlyMetric      `json:"hourly_page_views"`
	RealTimeEvents     []RecentEvent       `json:"real_time_events"`
	PerformanceMetrics PerformanceMetrics  `json:"performance_metrics"`
}

// PageMetric represents page visit statistics
type PageMetric struct {
	URL            string  `json:"url"`
	Path           string  `json:"path"`
	Views          int64   `json:"views"`
	UniqueVisitors int64   `json:"unique_visitors"`
	AverageTime    float64 `json:"average_time_seconds"`
	BounceRate     float64 `json:"bounce_rate"`
}

// TrafficSource represents referrer statistics
type TrafficSource struct {
	Source  string  `json:"source"`
	Count   int64   `json:"count"`
	Percent float64 `json:"percent"`
}

// HourlyMetric represents hourly aggregated data
type HourlyMetric struct {
	Hour   time.Time `json:"hour"`
	Events int64     `json:"events"`
}

// RecentEvent represents recent real-time events for display
type RecentEvent struct {
	Timestamp time.Time `json:"timestamp"`
	Type      EventType `json:"type"`
	URL       string    `json:"url"`
	UserID    string    `json:"user_id"`
	Location  string    `json:"location"`
}

// PerformanceMetrics represents performance analytics
type PerformanceMetrics struct {
	AverageLoadTime float64 `json:"average_load_time_ms"`
	MedianLoadTime  float64 `json:"median_load_time_ms"`
	SlowPagesCount  int64   `json:"slow_pages_count"`
	FastPagesCount  int64   `json:"fast_pages_count"`
}

// Alert represents a system alert
type Alert struct {
	ID           string    `json:"id"`
	Type         string    `json:"type"`
	Message      string    `json:"message"`
	Severity     string    `json:"severity"`
	Timestamp    time.Time `json:"timestamp"`
	Resolved     bool      `json:"resolved"`
	Threshold    float64   `json:"threshold"`
	CurrentValue float64   `json:"current_value"`
}

// AlertConfig represents alert configuration
type AlertConfig struct {
	Name          string  `json:"name"`
	Type          string  `json:"type"`
	Metric        string  `json:"metric"`
	Threshold     float64 `json:"threshold"`
	Operator      string  `json:"operator"` // "gt", "lt", "eq"
	Enabled       bool    `json:"enabled"`
	WindowMinutes int     `json:"window_minutes"`
}

// WebSocketMessage represents a message sent to WebSocket clients
type WebSocketMessage struct {
	Type      string      `json:"type"`
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data"`
}

// RealTimeAnalytics handles real-time analytics aggregation with time windows
type RealTimeAnalytics struct {
	Mu             sync.RWMutex
	Events         []AnalyticsEvent     // Recent events buffer
	PageViews      map[string]int64     // URL -> count
	UniqueUsers    map[string]bool      // UserID -> exists
	SessionsActive map[string]time.Time // SessionID -> last activity
	EventsByType   map[EventType]int64
	HourlyData     map[int64]int64            // Unix hour -> event count
	LoadTimes      []float64                  // Page load times
	TrafficSources map[string]int64           // Referrer domain -> count
	DeviceTypes    map[string]int64           // Device type -> count
	BrowserTypes   map[string]int64           // Browser -> count
	PageVisitors   map[string]map[string]bool // URL -> set of user IDs
	LastCleanup    time.Time
	StartTime      time.Time
	TotalEvents    int64
}

// NewRealTimeAnalytics creates a new real-time analytics instance
func NewRealTimeAnalytics() *RealTimeAnalytics {
	return &RealTimeAnalytics{
		Events:         make([]AnalyticsEvent, 0, 1000),
		PageViews:      make(map[string]int64),
		UniqueUsers:    make(map[string]bool),
		SessionsActive: make(map[string]time.Time),
		EventsByType:   make(map[EventType]int64),
		HourlyData:     make(map[int64]int64),
		LoadTimes:      make([]float64, 0, 1000),
		TrafficSources: make(map[string]int64),
		DeviceTypes:    make(map[string]int64),
		BrowserTypes:   make(map[string]int64),
		PageVisitors:   make(map[string]map[string]bool),
		LastCleanup:    time.Now(),
		StartTime:      time.Now(),
	}
}
