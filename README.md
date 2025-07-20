# Go FTP Server

A High performance File Transfer Protocol (FTP) server implementation built in Go. This server provides enterprise-grade file transfer capabilities with support for concurrent connections, secure authentication, comprehensive logging, and high-performance file operations.

## Features

- **High Performance**: Built with Go's concurrency primitives for efficient handling of multiple connections
- **Enterprise Security**: Path traversal protection, input validation, and configurable authentication
- **Production Configuration**: YAML-based configuration with command-line overrides and environment variable support
- **Comprehensive Logging**: Structured logging with multiple levels and configurable output formats
- **Authentication System**: Support for both anonymous and user-based authentication with role-based access
- **RFC Compliant**: Implements core FTP commands according to RFC 959 specifications
- **Graceful Shutdown**: Proper cleanup and connection handling for production deployments
- **Cross-Platform**: Native support for Linux, Windows, and macOS

## Supported FTP Commands

| Command | Description | Status |
|---------|-------------|--------|
| `USER` | Set username for authentication | Implemented |
| `PASS` | Set password for authentication | Implemented |
| `CWD` | Change working directory | Implemented |
| `PWD` | Print working directory | Implemented |
| `LIST` | List directory contents | Implemented |
| `RETR` | Retrieve (download) file | Implemented |
| `STOR` | Store (upload) file | Implemented |
| `QUIT` | Quit connection | Implemented |
| `NOOP` | No operation (keep-alive) | Implemented |

## Installation

### Prerequisites

- Go 1.21 or higher
- Git

### Building from Source

1. Clone the repository:
   ```bash
   git clone https://github.com/1cbyc/go-ftp-server.git
   cd go-ftp-server
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Build the project:
   ```bash
   go build -o ftp-server main.go
   ```

## Usage

### Basic Usage

Start the server with default settings:
```bash
./ftp-server
```

The server will start on `localhost:2121` with the root directory set to `./ftp_root`.

### Command Line Options

```bash
./ftp-server [options]

Options:
  -config string
        Path to configuration file (default "config.yaml")
  -host string
        Host to bind to (default "localhost")
  -port int
        Port to listen on (default 2121)
  -root string
        Root directory for FTP server (default "./ftp_root")
  -verbose
        Enable verbose logging
```

### Configuration File

The server uses a YAML configuration file (`config.yaml`):

```yaml
server:
  host: "localhost"
  port: 2121

ftp:
  root_dir: "./ftp_root"
  max_connections: 100
  timeout: 300

auth:
  anonymous: true
  users:
    anonymous: "anonymous"
    admin: "admin123"

log:
  level: "info"
  format: "text"
```

### Connecting with FTP Clients

#### Using Command Line FTP Client

```bash
ftp localhost 2121
```

#### Using FileZilla

1. Open FileZilla
2. Enter host: `localhost`
3. Enter port: `2121`
4. Enter username: `anonymous` (or configured user)
5. Enter password: `anonymous` (or configured password)
6. Click "Quickconnect"

#### Using curl

```bash
# List directory
curl -u anonymous:anonymous ftp://localhost:2121/

# Download file
curl -u anonymous:anonymous ftp://localhost:2121/filename.txt -o localfile.txt

# Upload file
curl -u anonymous:anonymous -T localfile.txt ftp://localhost:2121/remotefile.txt
```

<!-- ## 📁 Project Structure

```
go-ftp-server/
├── main.go                 # Application entry point
├── config.yaml            # Default configuration
├── go.mod                 # Go module file
├── go.sum                 # Go module checksums
├── .gitignore             # Git ignore rules
├── README.md              # This file
├── QUICKSTART.md          # Quick start guide
├── Makefile               # Development tasks
├── docs/
│   ├── explanation.md     # Technical explanation
│   └── whats-next.md      # Development roadmap
├── examples/
│   ├── test_ftp.py        # Python test script
│   └── test_ftp.sh        # Shell test script
└── internal/
    ├── config/            # Configuration management
    │   └── config.go
    ├── server/            # Main server logic
    │   └── server.go
    ├── handler/           # FTP command handler
    │   └── handler.go
    └── ftp/               # FTP protocol constants
        └── responses.go
``` -->

## Development

### Running Tests

```bash
go test ./...
```

### Building for Different Platforms

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o ftp-server-linux main.go

# Windows
GOOS=windows GOARCH=amd64 go build -o ftp-server.exe main.go

# macOS
GOOS=darwin GOARCH=amd64 go build -o ftp-server-mac main.go
```

### Code Style

This project follows Go's standard formatting and linting rules. Use:

```bash
go fmt ./...
go vet ./...
```

## Security Considerations

- **Path Traversal Protection**: The server validates all file paths to prevent directory traversal attacks
- **Input Validation**: All user inputs are validated and sanitized
- **Authentication**: Support for user-based authentication with configurable credentials
- **Anonymous Access**: Configurable anonymous access for public file sharing
- **Session Management**: Proper session handling and timeout management

## Performance

- **Concurrent Connections**: Each client connection runs in its own goroutine
- **Memory Efficient**: Minimal memory footprint per connection
- **Non-blocking I/O**: Efficient handling of file transfers
- **Connection Pooling**: Configurable maximum connection limits
- **Buffer Management**: Optimized buffer handling for large file transfers

## Troubleshooting

### Common Issues

1. **Connection Refused**
   - Check if the server is running
   - Verify the port is not in use by another application
   - Ensure firewall allows connections on the specified port

2. **Permission Denied**
   - Check file and directory permissions
   - Ensure the server has read/write access to the root directory

3. **File Transfer Issues**
   - Verify the file exists (for downloads)
   - Check available disk space (for uploads)
   - Ensure proper file permissions

### Debug Mode

Enable verbose logging for debugging:

```bash
./ftp-server -verbose
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with Go's excellent standard library
- Inspired by RFC 959 FTP specification
- Uses [logrus](https://github.com/sirupsen/logrus) for structured logging
- Uses [yaml.v3](https://gopkg.in/yaml.v3) for configuration management

## Support

For support and questions:

- Create an issue on GitHub
- Check the [documentation](docs/)
- Review the [roadmap](docs/whats-next.md) for upcoming features

---

**Built with ❤️ by Isaac using Go** 