#!/bin/bash

# Script to send sample analytics events to the producer API

API_URL="${API_URL:-http://localhost:8080/event}"

echo "Sending sample analytics events to $API_URL"

# Send a page view event
echo "Sending page view event..."
curl -X POST "$API_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "page_view",
    "user_id": "user123",
    "session_id": "session456",
    "url": "https://example.com/home",
    "path": "/home",
    "referrer": "https://google.com",
    "user_agent": "Mozilla/5.0",
    "ip_address": "192.168.1.1",
    "metadata": {
      "page_title": "Home Page",
      "load_time": 1200
    }
  }'
echo ""

sleep 1

# Send a click event
echo "Sending click event..."
curl -X POST "$API_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "click",
    "user_id": "user123",
    "session_id": "session456",
    "url": "https://example.com/products",
    "path": "/products",
    "metadata": {
      "element_id": "buy-button",
      "element_type": "button",
      "element_text": "Buy Now"
    }
  }'
echo ""

sleep 1

# Send another page view
echo "Sending another page view event..."
curl -X POST "$API_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "page_view",
    "user_id": "user456",
    "session_id": "session789",
    "url": "https://example.com/about",
    "path": "/about",
    "user_agent": "Mozilla/5.0",
    "ip_address": "192.168.1.2"
  }'
echo ""

sleep 1

# Send a session event
echo "Sending session event..."
curl -X POST "$API_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "session",
    "user_id": "user123",
    "session_id": "session456",
    "url": "https://example.com",
    "path": "/",
    "metadata": {
      "duration": 300,
      "page_count": 5,
      "device": "desktop",
      "browser": "Chrome"
    }
  }'
echo ""

echo "All events sent successfully!"
