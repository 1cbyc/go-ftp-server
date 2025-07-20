package handler

import (
	"bufio"
	"net"
	"strings"
	"testing"

	"github.com/mrinalxdev/go-ftp-server/internal/config"
)

func TestFTPHandler_Authentication(t *testing.T) {
	cfg := &config.Config{
		Auth: config.AuthConfig{
			Anonymous: true,
			Users: map[string]string{
				"testuser": "testpass",
			},
		},
	}

	handler := NewFTPHandler(cfg)

	// Test anonymous authentication
	if !handler.authenticate("anonymous", "anonymous") {
		t.Error("Anonymous authentication should succeed")
	}

	// Test valid user authentication
	if !handler.authenticate("testuser", "testpass") {
		t.Error("Valid user authentication should succeed")
	}

	// Test invalid user authentication
	if handler.authenticate("invaliduser", "invalidpass") {
		t.Error("Invalid user authentication should fail")
	}
}

func TestFTPHandler_PathValidation(t *testing.T) {
	cfg := &config.Config{}
	handler := NewFTPHandler(cfg)

	// Test valid paths
	validPaths := []string{
		"./ftp_root/file.txt",
		"./ftp_root/subdir/file.txt",
		"./ftp_root/../ftp_root/file.txt",
	}

	for _, path := range validPaths {
		if !handler.isValidPath(path, "./ftp_root") {
			t.Errorf("Path %s should be valid", path)
		}
	}

	// Test invalid paths (path traversal attempts)
	invalidPaths := []string{
		"./ftp_root/../../../etc/passwd",
		"./ftp_root/..\\..\\..\\windows\\system32\\config\\sam",
		"./ftp_root/%2e%2e/%2e%2e/%2e%2e/etc/passwd",
	}

	for _, path := range invalidPaths {
		if handler.isValidPath(path, "./ftp_root") {
			t.Errorf("Path %s should be invalid", path)
		}
	}
}

func TestFTPHandler_CommandParsing(t *testing.T) {
	cfg := &config.Config{}
	handler := NewFTPHandler(cfg)

	// Create a mock connection
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	// Start handler in goroutine
	go func() {
		handler.HandleConnection(server)
	}()

	// Test basic command parsing
	commands := []string{
		"USER anonymous\r\n",
		"PASS anonymous\r\n",
		"QUIT\r\n",
	}

	reader := bufio.NewReader(client)
	writer := bufio.NewWriter(client)

	for _, cmd := range commands {
		writer.WriteString(cmd)
		writer.Flush()

		// Read response
		response, err := reader.ReadString('\n')
		if err != nil {
			t.Errorf("Failed to read response for command %s: %v", cmd, err)
			continue
		}

		if !strings.Contains(response, "220") && !strings.Contains(response, "331") &&
			!strings.Contains(response, "230") && !strings.Contains(response, "221") {
			t.Errorf("Unexpected response for command %s: %s", cmd, strings.TrimSpace(response))
		}
	}
}
