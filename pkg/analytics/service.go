package analytics

import (
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Hilina-t/go-kafka-analytics-pipeline/pkg/models"
)

// Service handles real-time analytics processing and aggregation
type Service struct {
	analytics *models.RealTimeAnalytics
	alerts    []models.AlertConfig
	mu        sync.RWMutex
}

// NewService creates a new analytics service
func NewService() *Service {
	return &Service{
		analytics: models.NewRealTimeAnalytics(),
		alerts:    make([]models.AlertConfig, 0),
	}
}

// ProcessEvent processes a single analytics event
func (s *Service) ProcessEvent(event *models.AnalyticsEvent) error {
	s.analytics.Mu.Lock()
	defer s.analytics.Mu.Unlock()

	// Add to recent events buffer (keep last 100)
	s.analytics.Events = append(s.analytics.Events, *event)
	if len(s.analytics.Events) > 100 {
		s.analytics.Events = s.analytics.Events[1:]
	}

	// Update total events counter
	s.analytics.TotalEvents++

	// Track event by type
	s.analytics.EventsByType[event.Type]++

	// Track unique users
	if event.UserID != "" {
		s.analytics.UniqueUsers[event.UserID] = true
	}

	// Update session activity
	if event.SessionID != "" {
		s.analytics.SessionsActive[event.SessionID] = event.Timestamp
	}

	// Track hourly data
	hour := event.Timestamp.Truncate(time.Hour).Unix()
	s.analytics.HourlyData[hour]++

	// Process specific event types
	switch event.Type {
	case models.PageView:
		s.processPageView(event)
	case models.Click:
		s.processClick(event)
	case models.Session:
		s.processSession(event)
	}

	// Extract traffic source from referrer
	if event.Referrer != "" {
		s.processReferrer(event.Referrer)
	}

	// Extract device and browser info from user agent
	if event.UserAgent != "" {
		s.processUserAgent(event.UserAgent)
	}

	// Periodic cleanup (every 5 minutes)
	if time.Since(s.analytics.LastCleanup) > 5*time.Minute {
		s.cleanup()
		s.analytics.LastCleanup = time.Now()
	}

	return nil
}

// processPageView handles page view specific processing
func (s *Service) processPageView(event *models.AnalyticsEvent) {
	s.analytics.PageViews[event.URL]++

	// Track unique visitors per page
	if s.analytics.PageVisitors[event.URL] == nil {
		s.analytics.PageVisitors[event.URL] = make(map[string]bool)
	}
	if event.UserID != "" {
		s.analytics.PageVisitors[event.URL][event.UserID] = true
	}

	// Extract load time from metadata
	if metadata, ok := event.Metadata["load_time"].(float64); ok {
		s.analytics.LoadTimes = append(s.analytics.LoadTimes, metadata)
		// Keep only last 1000 load times
		if len(s.analytics.LoadTimes) > 1000 {
			s.analytics.LoadTimes = s.analytics.LoadTimes[1:]
		}
	}
}

// processClick handles click event processing
func (s *Service) processClick(_ *models.AnalyticsEvent) {
	// Click events can be used for interaction tracking
	// Add specific click processing logic here if needed
}

// processSession handles session event processing
func (s *Service) processSession(event *models.AnalyticsEvent) {
	// Extract device info from metadata
	if device, ok := event.Metadata["device"].(string); ok && device != "" {
		s.analytics.DeviceTypes[device]++
	}
	if browser, ok := event.Metadata["browser"].(string); ok && browser != "" {
		s.analytics.BrowserTypes[browser]++
	}
}

// processReferrer extracts domain from referrer URL
func (s *Service) processReferrer(referrer string) {
	if u, err := url.Parse(referrer); err == nil && u.Host != "" {
		domain := u.Host
		if strings.HasPrefix(domain, "www.") {
			domain = domain[4:]
		}
		s.analytics.TrafficSources[domain]++
	}
}

// processUserAgent extracts browser and device info from user agent
func (s *Service) processUserAgent(userAgent string) {
	userAgent = strings.ToLower(userAgent)

	// Simple browser detection
	if strings.Contains(userAgent, "chrome") {
		s.analytics.BrowserTypes["Chrome"]++
	} else if strings.Contains(userAgent, "firefox") {
		s.analytics.BrowserTypes["Firefox"]++
	} else if strings.Contains(userAgent, "safari") {
		s.analytics.BrowserTypes["Safari"]++
	} else if strings.Contains(userAgent, "edge") {
		s.analytics.BrowserTypes["Edge"]++
	} else {
		s.analytics.BrowserTypes["Other"]++
	}

	// Simple device detection
	if strings.Contains(userAgent, "mobile") || strings.Contains(userAgent, "iphone") || strings.Contains(userAgent, "android") {
		s.analytics.DeviceTypes["Mobile"]++
	} else if strings.Contains(userAgent, "tablet") || strings.Contains(userAgent, "ipad") {
		s.analytics.DeviceTypes["Tablet"]++
	} else {
		s.analytics.DeviceTypes["Desktop"]++
	}
}

// cleanup removes old sessions and data
func (s *Service) cleanup() {
	now := time.Now()

	// Remove inactive sessions (older than 30 minutes)
	for sessionID, lastActivity := range s.analytics.SessionsActive {
		if now.Sub(lastActivity) > 30*time.Minute {
			delete(s.analytics.SessionsActive, sessionID)
		}
	}

	// Clean up old hourly data (keep last 48 hours)
	cutoff := now.Add(-48 * time.Hour).Truncate(time.Hour).Unix()
	for hour := range s.analytics.HourlyData {
		if hour < cutoff {
			delete(s.analytics.HourlyData, hour)
		}
	}
}

// GetSnapshot returns a complete analytics snapshot
func (s *Service) GetSnapshot() *models.MetricsSnapshot {
	s.analytics.Mu.RLock()
	defer s.analytics.Mu.RUnlock()

	snapshot := &models.MetricsSnapshot{
		Timestamp:          time.Now(),
		TotalEvents:        s.analytics.TotalEvents,
		UniqueUsers:        int64(len(s.analytics.UniqueUsers)),
		ActiveSessions:     int64(len(s.analytics.SessionsActive)),
		EventsByType:       make(map[models.EventType]int64),
		TopPages:           s.getTopPages(),
		TrafficSources:     s.getTrafficSources(),
		DeviceStats:        make(map[string]int64),
		BrowserStats:       make(map[string]int64),
		HourlyPageViews:    s.getHourlyPageViews(),
		RealTimeEvents:     s.getRecentEvents(),
		PerformanceMetrics: s.getPerformanceMetrics(),
	}

	// Copy event type stats
	for eventType, count := range s.analytics.EventsByType {
		snapshot.EventsByType[eventType] = count
	}

	// Copy device stats
	for device, count := range s.analytics.DeviceTypes {
		snapshot.DeviceStats[device] = count
	}

	// Copy browser stats
	for browser, count := range s.analytics.BrowserTypes {
		snapshot.BrowserStats[browser] = count
	}

	return snapshot
}

// getTopPages returns top pages sorted by views
func (s *Service) getTopPages() []models.PageMetric {
	type pageData struct {
		url      string
		views    int64
		visitors int64
	}

	pages := make([]pageData, 0, len(s.analytics.PageViews))
	for pageURL, views := range s.analytics.PageViews {
		visitors := int64(0)
		if s.analytics.PageVisitors[pageURL] != nil {
			visitors = int64(len(s.analytics.PageVisitors[pageURL]))
		}
		pages = append(pages, pageData{url: pageURL, views: views, visitors: visitors})
	}

	// Sort by views descending
	sort.Slice(pages, func(i, j int) bool {
		return pages[i].views > pages[j].views
	})

	// Convert to PageMetric (top 10)
	result := make([]models.PageMetric, 0, 10)
	for i, page := range pages {
		if i >= 10 {
			break
		}

		// Extract path from URL
		path := page.url
		if u, err := url.Parse(page.url); err == nil {
			path = u.Path
		}

		result = append(result, models.PageMetric{
			URL:            page.url,
			Path:           path,
			Views:          page.views,
			UniqueVisitors: page.visitors,
			BounceRate:     0, // TODO: Calculate bounce rate
		})
	}

	return result
}

// getTrafficSources returns top traffic sources
func (s *Service) getTrafficSources() []models.TrafficSource {
	type sourceData struct {
		source string
		count  int64
	}

	sources := make([]sourceData, 0, len(s.analytics.TrafficSources))
	totalTraffic := int64(0)

	for source, count := range s.analytics.TrafficSources {
		sources = append(sources, sourceData{source: source, count: count})
		totalTraffic += count
	}

	// Sort by count descending
	sort.Slice(sources, func(i, j int) bool {
		return sources[i].count > sources[j].count
	})

	// Convert to TrafficSource (top 10)
	result := make([]models.TrafficSource, 0, 10)
	for i, source := range sources {
		if i >= 10 {
			break
		}

		percent := float64(0)
		if totalTraffic > 0 {
			percent = float64(source.count) / float64(totalTraffic) * 100
		}

		result = append(result, models.TrafficSource{
			Source:  source.source,
			Count:   source.count,
			Percent: percent,
		})
	}

	return result
}

// getHourlyPageViews returns hourly page view data for the last 24 hours
func (s *Service) getHourlyPageViews() []models.HourlyMetric {
	now := time.Now()
	result := make([]models.HourlyMetric, 0, 24)

	for i := 23; i >= 0; i-- {
		hour := now.Add(-time.Duration(i) * time.Hour).Truncate(time.Hour)
		hourUnix := hour.Unix()

		count := int64(0)
		if hourlyCount, exists := s.analytics.HourlyData[hourUnix]; exists {
			count = hourlyCount
		}

		result = append(result, models.HourlyMetric{
			Hour:   hour,
			Events: count,
		})
	}

	return result
}

// getRecentEvents returns the most recent events for real-time display
func (s *Service) getRecentEvents() []models.RecentEvent {
	result := make([]models.RecentEvent, 0, len(s.analytics.Events))

	// Get last 20 events
	start := 0
	if len(s.analytics.Events) > 20 {
		start = len(s.analytics.Events) - 20
	}

	for i := start; i < len(s.analytics.Events); i++ {
		event := s.analytics.Events[i]
		result = append(result, models.RecentEvent{
			Timestamp: event.Timestamp,
			Type:      event.Type,
			URL:       event.URL,
			UserID:    event.UserID,
			Location:  s.extractLocation(event.IPAddress),
		})
	}

	return result
}

// getPerformanceMetrics calculates performance metrics from load times
func (s *Service) getPerformanceMetrics() models.PerformanceMetrics {
	if len(s.analytics.LoadTimes) == 0 {
		return models.PerformanceMetrics{}
	}

	// Calculate average load time
	sum := float64(0)
	for _, loadTime := range s.analytics.LoadTimes {
		sum += loadTime
	}
	avg := sum / float64(len(s.analytics.LoadTimes))

	// Calculate median (simple approach)
	sorted := make([]float64, len(s.analytics.LoadTimes))
	copy(sorted, s.analytics.LoadTimes)
	sort.Float64s(sorted)
	median := sorted[len(sorted)/2]

	// Count fast vs slow pages (threshold: 3 seconds = 3000ms)
	slowCount := int64(0)
	fastCount := int64(0)
	for _, loadTime := range s.analytics.LoadTimes {
		if loadTime > 3000 {
			slowCount++
		} else {
			fastCount++
		}
	}

	return models.PerformanceMetrics{
		AverageLoadTime: avg,
		MedianLoadTime:  median,
		SlowPagesCount:  slowCount,
		FastPagesCount:  fastCount,
	}
}

// extractLocation extracts location from IP address (simplified)
func (s *Service) extractLocation(ipAddress string) string {
	// This is a simplified implementation
	// In production, you'd use a GeoIP service
	if ipAddress == "" {
		return "Unknown"
	}

	// Check for common private IP ranges
	if strings.HasPrefix(ipAddress, "192.168.") ||
		strings.HasPrefix(ipAddress, "10.") ||
		strings.HasPrefix(ipAddress, "172.") {
		return "Local"
	}

	return "External"
}

// AddAlert adds a new alert configuration
func (s *Service) AddAlert(config models.AlertConfig) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.alerts = append(s.alerts, config)
}

// CheckAlerts evaluates all alert conditions and returns triggered alerts
func (s *Service) CheckAlerts() []models.Alert {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var triggeredAlerts []models.Alert
	snapshot := s.GetSnapshot()

	for _, alertConfig := range s.alerts {
		if !alertConfig.Enabled {
			continue
		}

		currentValue := s.getMetricValue(snapshot, alertConfig.Metric)
		triggered := s.evaluateAlertCondition(currentValue, alertConfig.Threshold, alertConfig.Operator)

		if triggered {
			alert := models.Alert{
				ID:           "alert_" + strconv.FormatInt(time.Now().Unix(), 10),
				Type:         alertConfig.Type,
				Message:      s.generateAlertMessage(alertConfig, currentValue),
				Severity:     s.getAlertSeverity(alertConfig.Type),
				Timestamp:    time.Now(),
				Resolved:     false,
				Threshold:    alertConfig.Threshold,
				CurrentValue: currentValue,
			}
			triggeredAlerts = append(triggeredAlerts, alert)
		}
	}

	return triggeredAlerts
}

// getMetricValue extracts a specific metric value from the snapshot
func (s *Service) getMetricValue(snapshot *models.MetricsSnapshot, metric string) float64 {
	switch metric {
	case "total_events":
		return float64(snapshot.TotalEvents)
	case "unique_users":
		return float64(snapshot.UniqueUsers)
	case "active_sessions":
		return float64(snapshot.ActiveSessions)
	case "average_load_time":
		return snapshot.PerformanceMetrics.AverageLoadTime
	default:
		return 0
	}
}

// evaluateAlertCondition checks if an alert condition is met
func (s *Service) evaluateAlertCondition(current, threshold float64, operator string) bool {
	switch operator {
	case "gt":
		return current > threshold
	case "lt":
		return current < threshold
	case "eq":
		return current == threshold
	default:
		return false
	}
}

// generateAlertMessage creates a human-readable alert message
func (s *Service) generateAlertMessage(config models.AlertConfig, currentValue float64) string {
	return fmt.Sprintf("Alert: %s - %s is %.2f (threshold: %.2f)",
		config.Name, config.Metric, currentValue, config.Threshold)
}

// getAlertSeverity determines alert severity based on type
func (s *Service) getAlertSeverity(alertType string) string {
	switch alertType {
	case "performance":
		return "medium"
	case "traffic":
		return "low"
	case "error":
		return "high"
	default:
		return "medium"
	}
}
