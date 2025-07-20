#!/bin/bash

# Go FTP Server Test Script
# This script tests the FTP server using curl commands

HOST="localhost"
PORT="2121"
USERNAME="anonymous"
PASSWORD="anonymous"

echo "Go FTP Server Test Script"
echo "========================="
echo "Host: $HOST"
echo "Port: $PORT"
echo "Username: $USERNAME"
echo ""

# Test 1: Basic connection and directory listing
echo "Test 1: Directory listing"
echo "-------------------------"
curl -s -u "$USERNAME:$PASSWORD" "ftp://$HOST:$PORT/" || {
    echo "❌ Failed to connect to FTP server"
    echo "Make sure the server is running on $HOST:$PORT"
    exit 1
}
echo "✅ Directory listing successful"
echo ""

# Test 2: Create a test file
echo "Test 2: File upload"
echo "-------------------"
TEST_FILE="test_upload_$(date +%s).txt"
TEST_CONTENT="This is a test file created at $(date)"

echo "$TEST_CONTENT" > "$TEST_FILE"
curl -s -u "$USERNAME:$PASSWORD" -T "$TEST_FILE" "ftp://$HOST:$PORT/$TEST_FILE" && {
    echo "✅ File upload successful"
} || {
    echo "❌ File upload failed"
    rm -f "$TEST_FILE"
    exit 1
}
echo ""

# Test 3: Download the test file
echo "Test 3: File download"
echo "---------------------"
DOWNLOAD_FILE="test_download_$(date +%s).txt"
curl -s -u "$USERNAME:$PASSWORD" "ftp://$HOST:$PORT/$TEST_FILE" -o "$DOWNLOAD_FILE" && {
    echo "✅ File download successful"
} || {
    echo "❌ File download failed"
    rm -f "$TEST_FILE" "$DOWNLOAD_FILE"
    exit 1
}
echo ""

# Test 4: Verify file content
echo "Test 4: Content verification"
echo "----------------------------"
if [ -f "$DOWNLOAD_FILE" ]; then
    DOWNLOADED_CONTENT=$(cat "$DOWNLOAD_FILE")
    if [ "$DOWNLOADED_CONTENT" = "$TEST_CONTENT" ]; then
        echo "✅ Content verification successful"
    else
        echo "❌ Content verification failed"
        echo "Expected: $TEST_CONTENT"
        echo "Got: $DOWNLOADED_CONTENT"
        rm -f "$TEST_FILE" "$DOWNLOAD_FILE"
        exit 1
    fi
else
    echo "❌ Downloaded file not found"
    rm -f "$TEST_FILE"
    exit 1
fi
echo ""

# Test 5: List files again to see our uploaded file
echo "Test 5: Verify file in directory listing"
echo "----------------------------------------"
curl -s -u "$USERNAME:$PASSWORD" "ftp://$HOST:$PORT/" | grep -q "$TEST_FILE" && {
    echo "✅ File found in directory listing"
} || {
    echo "❌ File not found in directory listing"
}
echo ""

# Cleanup
echo "Cleaning up test files..."
rm -f "$TEST_FILE" "$DOWNLOAD_FILE"
echo "✅ Cleanup complete"
echo ""

echo "🎉 All tests passed successfully!"
echo "The Go FTP Server is working correctly."
echo ""
echo "You can now use any FTP client to connect to:"
echo "  Host: $HOST"
echo "  Port: $PORT"
echo "  Username: $USERNAME"
echo "  Password: $PASSWORD" 