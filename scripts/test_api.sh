#!/bin/bash

# Test script for blog API
BASE_URL="http://localhost:8080"

echo "=== Testing Blog API ==="
echo ""

# Test 1: Get public posts (should work without auth)
echo "1. Testing GET /api/posts"
curl -s "$BASE_URL/api/posts" | jq '.' || echo "No posts yet or server not running"
echo ""

# Test 2: Login as admin
echo "2. Testing POST /api/admin/login"
COOKIE_FILE=$(mktemp)
curl -s -c "$COOKIE_FILE" -X POST "$BASE_URL/api/admin/login" \
  -H "Content-Type: application/json" \
  -d '{"password":"123456"}' | jq '.'
echo ""

# Test 3: Create a post with JSON
echo "3. Testing POST /api/admin/posts (JSON)"
curl -s -b "$COOKIE_FILE" -X POST "$BASE_URL/api/admin/posts" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Test Post",
    "content": "# Hello World\n\nThis is a test post with **markdown**.",
    "excerpt": "A test post",
    "published": true
  }' | jq '.'
echo ""

# Test 4: Get all admin posts
echo "4. Testing GET /api/admin/posts"
curl -s -b "$COOKIE_FILE" "$BASE_URL/api/admin/posts" | jq '.'
echo ""

# Test 5: Get public posts again
echo "5. Testing GET /api/posts (should show new post)"
curl -s "$BASE_URL/api/posts" | jq '.'
echo ""

# Test 6: Get post by slug
echo "6. Testing GET /api/posts/test-post"
curl -s "$BASE_URL/api/posts/test-post" | jq '.'
echo ""

# Cleanup
rm -f "$COOKIE_FILE"

echo "=== Tests Complete ==="
