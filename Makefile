# Go FTP Server Makefile

# Variables
BINARY_NAME=ftp-server
MAIN_FILE=main.go
CONFIG_FILE=config.yaml

# Build targets
.PHONY: build clean test run help

# Default target
all: build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) $(MAIN_FILE)
	@echo "Build complete!"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME).exe
	@echo "Clean complete!"

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Vet code
vet:
	@echo "Vetting code..."
	go vet ./...

# Lint code (requires golangci-lint)
lint:
	@echo "Linting code..."
	golangci-lint run

# Run the server
run: build
	@echo "Starting FTP server..."
	./$(BINARY_NAME)

# Run the server with verbose logging
run-verbose: build
	@echo "Starting FTP server with verbose logging..."
	./$(BINARY_NAME) -verbose

# Run the server on a different port
run-port: build
	@echo "Starting FTP server on port 2122..."
	./$(BINARY_NAME) -port 2122

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

# Build for different platforms
build-linux:
	@echo "Building for Linux..."
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)-linux $(MAIN_FILE)

build-windows:
	@echo "Building for Windows..."
	GOOS=windows GOARCH=amd64 go build -o $(BINARY_NAME).exe $(MAIN_FILE)

build-mac:
	@echo "Building for macOS..."
	GOOS=darwin GOARCH=amd64 go build -o $(BINARY_NAME)-mac $(MAIN_FILE)

# Build all platforms
build-all: build-linux build-windows build-mac
	@echo "Build for all platforms complete!"

# Create sample files in ftp_root
setup-sample:
	@echo "Setting up sample files..."
	mkdir -p ftp_root
	echo "Hello, this is a sample file!" > ftp_root/sample.txt
	echo "This is another sample file." > ftp_root/readme.txt
	mkdir -p ftp_root/subdir
	echo "File in subdirectory" > ftp_root/subdir/test.txt
	@echo "Sample files created in ftp_root/"

# Show help
help:
	@echo "Available targets:"
	@echo "  build          - Build the application"
	@echo "  clean          - Clean build artifacts"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  fmt            - Format code"
	@echo "  vet            - Vet code"
	@echo "  lint           - Lint code (requires golangci-lint)"
	@echo "  run            - Build and run the server"
	@echo "  run-verbose    - Build and run with verbose logging"
	@echo "  run-port       - Build and run on port 2122"
	@echo "  deps           - Install dependencies"
	@echo "  build-linux    - Build for Linux"
	@echo "  build-windows  - Build for Windows"
	@echo "  build-mac      - Build for macOS"
	@echo "  build-all      - Build for all platforms"
	@echo "  setup-sample   - Create sample files in ftp_root"
	@echo "  help           - Show this help message" 