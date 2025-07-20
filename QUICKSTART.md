# Quick Start Guide

If you want to get your Go FTP Server up and running in minutes, just read my steps here:

## 5-Minute Setup

### 1. Build the Server
```bash
go build -o ftp-server main.go
```

### 2. Start the Server
```bash
./ftp-server
```

### 3. Connect with an FTP Client
- **Host**: `localhost`
- **Port**: `2121`
- **Username**: `anonymous`
- **Password**: `anonymous`

## Test the Server

### Using the Python Test Script
```bash
python3 examples/test_ftp.py
```

### Using the Shell Test Script
```bash
chmod +x examples/test_ftp.sh
./examples/test_ftp.sh
```

### Using curl
```bash
# List files
curl -u anonymous:anonymous ftp://localhost:2121/

# Upload a file
echo "Hello World" > test.txt
curl -u anonymous:anonymous -T test.txt ftp://localhost:2121/hello.txt

# Download a file
curl -u anonymous:anonymous ftp://localhost:2121/hello.txt -o downloaded.txt
```

## Common Commands

### Start with verbose logging
```bash
./ftp-server -verbose
```

### Start on a different port
```bash
./ftp-server -port 2122
```

### Start with custom root directory
```bash
./ftp-server -root /path/to/your/files
```

### Start with custom configuration
```bash
./ftp-server -config my-config.yaml
```

## File Structure

After starting the server, a `ftp_root` directory will be created automatically. This is where your FTP files will be stored.

```
ftp_root/
├── (your uploaded files will appear here)
└── (create subdirectories as needed)
```

## Troubleshooting

### Server won't start?
- Check if port 2121 is already in use
- Try a different port: `./ftp-server -port 2122`

### Can't connect?
- Make sure the server is running
- Check firewall settings
- Try connecting to `127.0.0.1` instead of `localhost`

### Permission errors?
- Ensure the server has write access to the current directory
- Check file permissions in the `ftp_root` directory

## Next Steps

1. **Read the full documentation**: [README.md](README.md)
2. **Explore the configuration**: [config.yaml](config.yaml)
3. **Check the roadmap**: [docs/whats-next.md](docs/whats-next.md)
4. **Understand the architecture**: [docs/explanation.md](docs/explanation.md)

## Pro Tips

- Use `-verbose` flag for debugging
- The server supports concurrent connections
- Files are stored in the `ftp_root` directory
- Anonymous access is enabled by default
- The server gracefully handles disconnections

---

**That's it! Your FTP server is ready to use.** 🎉 