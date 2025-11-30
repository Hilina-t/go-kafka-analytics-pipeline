#!/bin/bash

# Enhanced test script for the Real-Time Analytics Dashboard
# This script sends various types of analytics events to test the new features

BASE_URL="http://localhost:8080"

echo "🚀 Testing Enhanced Real-Time Analytics Dashboard"
echo "================================================"

# Function to send an analytics event
send_event() {
    local event_data="$1"
    local description="$2"

    echo "📊 Sending: $description"
    curl -s -X POST "$BASE_URL/event" \
        -H "Content-Type: application/json" \
        -d "$event_data" | jq '.'
    echo
}

# Test connectivity
echo "🔗 Testing API connectivity..."
curl -s "$BASE_URL/health" | jq '.'
echo

# Generate realistic test data
USER_IDS=("user123" "user456" "user789" "user101" "user202")
SESSION_IDS=("session_$(date +%s)_1" "session_$(date +%s)_2" "session_$(date +%s)_3")
PAGES=("/home" "/products" "/about" "/contact" "/blog" "/pricing")
REFERRERS=("https://google.com" "https://facebook.com" "https://twitter.com" "direct" "")
USER_AGENTS=(
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/91.0.4472.124"
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 Safari/537.36"
    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 Chrome/91.0.4472.124"
    "Mozilla/5.0 (iPhone; CPU iPhone OS 14_6 like Mac OS X) AppleWebKit/605.1.15"
)

echo "🎯 Sending realistic page view events..."

for i in {1..15}; do
    user_id=${USER_IDS[$((RANDOM % ${#USER_IDS[@]}))]}
    session_id=${SESSION_IDS[$((RANDOM % ${#SESSION_IDS[@]}))]}
    page=${PAGES[$((RANDOM % ${#PAGES[@]}))]}
    referrer=${REFERRERS[$((RANDOM % ${#REFERRERS[@]}))]}
    user_agent=${USER_AGENTS[$((RANDOM % ${#USER_AGENTS[@]}))]}
    load_time=$((RANDOM % 5000 + 500)) # 500-5500ms

    event_data=$(cat <<EOF
{
  "type": "page_view",
  "user_id": "$user_id",
  "session_id": "$session_id",
  "url": "https://example.com$page",
  "path": "$page",
  "referrer": "$referrer",
  "user_agent": "$user_agent",
  "ip_address": "192.168.1.$((RANDOM % 255 + 1))",
  "metadata": {
    "page_title": "Example Page - $(echo $page | tr '[:lower:]' '[:upper:]')",
    "load_time": $load_time,
    "screen_width": $((RANDOM % 2000 + 800)),
    "screen_height": $((RANDOM % 1500 + 600))
  }
}
EOF
    )

    send_event "$event_data" "Page view: $page (Load time: ${load_time}ms)"
    sleep 0.5
done

echo "👆 Sending click events..."

for i in {1..8}; do
    user_id=${USER_IDS[$((RANDOM % ${#USER_IDS[@]}))]}
    session_id=${SESSION_IDS[$((RANDOM % ${#SESSION_IDS[@]}))]}
    page=${PAGES[$((RANDOM % ${#PAGES[@]}))]}

    elements=("buy-button" "signup-btn" "menu-link" "product-card" "cta-button")
    element_id=${elements[$((RANDOM % ${#elements[@]}))]}

    event_data=$(cat <<EOF
{
  "type": "click",
  "user_id": "$user_id",
  "session_id": "$session_id",
  "url": "https://example.com$page",
  "path": "$page",
  "user_agent": "${USER_AGENTS[$((RANDOM % ${#USER_AGENTS[@]}))]}",
  "ip_address": "192.168.1.$((RANDOM % 255 + 1))",
  "metadata": {
    "element_id": "$element_id",
    "element_type": "button",
    "element_text": "$(echo $element_id | tr '-' ' ' | tr '[:lower:]' '[:upper:]')",
    "x_position": $((RANDOM % 1000)),
    "y_position": $((RANDOM % 800))
  }
}
EOF
    )

    send_event "$event_data" "Click event: $element_id on $page"
    sleep 0.3
done

echo "👥 Sending session events..."

for i in {1..5}; do
    user_id=${USER_IDS[$((RANDOM % ${#USER_IDS[@]}))]}
    session_id=${SESSION_IDS[$((RANDOM % ${#SESSION_IDS[@]}))]}

    devices=("Desktop" "Mobile" "Tablet")
    browsers=("Chrome" "Safari" "Firefox" "Edge")

    device=${devices[$((RANDOM % ${#devices[@]}))]}
    browser=${browsers[$((RANDOM % ${#browsers[@]}))]}
    duration=$((RANDOM % 1800 + 60)) # 1-30 minutes
    page_count=$((RANDOM % 20 + 1))

    event_data=$(cat <<EOF
{
  "type": "session",
  "user_id": "$user_id",
  "session_id": "$session_id",
  "url": "https://example.com",
  "path": "/",
  "user_agent": "${USER_AGENTS[$((RANDOM % ${#USER_AGENTS[@]}))]}",
  "ip_address": "192.168.1.$((RANDOM % 255 + 1))",
  "metadata": {
    "duration": $duration,
    "page_count": $page_count,
    "device": "$device",
    "browser": "$browser"
  }
}
EOF
    )

    send_event "$event_data" "Session: $device/$browser (${duration}s, $page_count pages)"
    sleep 0.2
done

echo "🏃‍♂️ Sending rapid events to test real-time updates..."

for i in {1..10}; do
    user_id=${USER_IDS[$((RANDOM % ${#USER_IDS[@]}))]}
    page=${PAGES[$((RANDOM % ${#PAGES[@]}))]}

    event_data=$(cat <<EOF
{
  "type": "page_view",
  "user_id": "$user_id",
  "session_id": "rapid_session_$i",
  "url": "https://example.com$page",
  "path": "$page",
  "user_agent": "${USER_AGENTS[0]}",
  "metadata": {
    "page_title": "Rapid Test Page",
    "load_time": $((RANDOM % 1000 + 200))
  }
}
EOF
    )

    curl -s -X POST "$BASE_URL/event" \
        -H "Content-Type: application/json" \
        -d "$event_data" > /dev/null

    echo -n "."
    sleep 0.1
done

echo ""
echo ""

# Test analytics API endpoint
echo "📈 Testing Analytics API endpoint..."
curl -s "$BASE_URL/analytics" | jq '.' | head -20
echo "... (response truncated)"
echo ""

echo "✅ Test complete!"
echo ""
echo "🎯 Open your browser and go to:"
echo "   📊 Dashboard: $BASE_URL"
echo "   📈 Analytics API: $BASE_URL/analytics"
echo "   🔗 WebSocket: ws://localhost:8080/ws"
echo ""
echo "📱 The dashboard should show:"
echo "   • Real-time event updates"
echo "   • Interactive charts and metrics"
echo "   • Live event stream"
echo "   • Performance analytics"
echo "   • Device and browser stats"
echo ""
echo "🔔 Check consumer logs for alerts and detailed analytics"
