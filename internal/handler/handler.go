package handler

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/mrinalxdev/go-ftp-server/internal/config"
	"github.com/mrinalxdev/go-ftp-server/internal/ftp"
	"github.com/sirupsen/logrus"
)

// FTPHandler handles FTP client connections and commands
type FTPHandler struct {
	config *config.Config
	mu     sync.RWMutex
}

// NewFTPHandler creates a new FTP handler
func NewFTPHandler(cfg *config.Config) *FTPHandler {
	return &FTPHandler{
		config: cfg,
	}
}

// HandleConnection handles a single client connection
func (h *FTPHandler) HandleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	// Send welcome message
	h.sendResponse(writer, ftp.ResponseWelcome)

	// Create session for this connection
	session := &FTPSession{
		conn:          conn,
		reader:        reader,
		writer:        writer,
		config:        h.config,
		rootDir:       h.config.FTP.RootDir,
		currentDir:    ".",
		authenticated: false,
	}

	// Process commands
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			logrus.Debugf("Client disconnected: %v", err)
			return
		}

		message = strings.TrimSpace(message)
		if message == "" {
			continue
		}

		logrus.Debugf("Received command: %s", message)
		h.processCommand(session, message)
	}
}

// processCommand processes a single FTP command
func (h *FTPHandler) processCommand(session *FTPSession, message string) {
	parts := strings.SplitN(message, " ", 2)
	command := strings.ToUpper(parts[0])
	var args string
	if len(parts) > 1 {
		args = parts[1]
	}

	switch command {
	case "USER":
		h.handleUser(session, args)
	case "PASS":
		h.handlePass(session, args)
	case "CWD":
		h.handleCwd(session, args)
	case "PWD":
		h.handlePwd(session)
	case "LIST":
		h.handleList(session, args)
	case "RETR":
		h.handleRetr(session, args)
	case "STOR":
		h.handleStor(session, args)
	case "QUIT":
		h.handleQuit(session)
	case "NOOP":
		h.handleNoop(session)
	default:
		h.sendResponse(session.writer, ftp.ResponseCommandNotImplemented)
	}
}

// handleUser handles the USER command
func (h *FTPHandler) handleUser(session *FTPSession, username string) {
	session.username = username
	h.sendResponse(session.writer, ftp.ResponseUsernameOK)
}

// handlePass handles the PASS command
func (h *FTPHandler) handlePass(session *FTPSession, password string) {
	if h.authenticate(session.username, password) {
		session.authenticated = true
		h.sendResponse(session.writer, ftp.ResponseLoginOK)
	} else {
		h.sendResponse(session.writer, ftp.ResponseLoginFailed)
	}
}

// handleCwd handles the CWD command
func (h *FTPHandler) handleCwd(session *FTPSession, dir string) {
	if !session.authenticated {
		h.sendResponse(session.writer, ftp.ResponseNotLoggedIn)
		return
	}

	newPath := filepath.Join(session.rootDir, session.currentDir, dir)
	if !h.isValidPath(newPath, session.rootDir) {
		h.sendResponse(session.writer, ftp.ResponseDirectoryNotFound)
		return
	}

	if info, err := os.Stat(newPath); err != nil || !info.IsDir() {
		h.sendResponse(session.writer, ftp.ResponseDirectoryNotFound)
		return
	}

	session.currentDir = filepath.Join(session.currentDir, dir)
	h.sendResponse(session.writer, ftp.ResponseDirectoryChanged)
}

// handlePwd handles the PWD command
func (h *FTPHandler) handlePwd(session *FTPSession) {
	if !session.authenticated {
		h.sendResponse(session.writer, ftp.ResponseNotLoggedIn)
		return
	}

	response := fmt.Sprintf("257 \"%s\" is current directory", session.currentDir)
	h.sendResponse(session.writer, response)
}

// handleList handles the LIST command
func (h *FTPHandler) handleList(session *FTPSession, path string) {
	if !session.authenticated {
		h.sendResponse(session.writer, ftp.ResponseNotLoggedIn)
		return
	}

	listPath := filepath.Join(session.rootDir, session.currentDir, path)
	if !h.isValidPath(listPath, session.rootDir) {
		h.sendResponse(session.writer, ftp.ResponseDirectoryNotFound)
		return
	}

	files, err := os.ReadDir(listPath)
	if err != nil {
		h.sendResponse(session.writer, ftp.ResponseDirectoryNotFound)
		return
	}

	h.sendResponse(session.writer, ftp.ResponseDataConnection)

	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			continue
		}

		line := fmt.Sprintf("%s\t%d\t%s",
			file.Name(),
			info.Size(),
			info.ModTime().Format("Jan 02 15:04"))
		h.sendResponse(session.writer, line)
	}

	h.sendResponse(session.writer, ftp.ResponseTransferComplete)
}

// handleRetr handles the RETR command
func (h *FTPHandler) handleRetr(session *FTPSession, filename string) {
	if !session.authenticated {
		h.sendResponse(session.writer, ftp.ResponseNotLoggedIn)
		return
	}

	filePath := filepath.Join(session.rootDir, session.currentDir, filename)
	if !h.isValidPath(filePath, session.rootDir) {
		h.sendResponse(session.writer, ftp.ResponseFileNotFound)
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		h.sendResponse(session.writer, ftp.ResponseFileNotFound)
		return
	}
	defer file.Close()

	h.sendResponse(session.writer, ftp.ResponseDataConnection)

	// Simple file transfer - in a real implementation, you'd use a data connection
	buffer := make([]byte, 1024)
	for {
		n, err := file.Read(buffer)
		if err != nil || n == 0 {
			break
		}
		session.writer.Write(buffer[:n])
		session.writer.Flush()
	}

	h.sendResponse(session.writer, ftp.ResponseTransferComplete)
}

// handleStor handles the STOR command
func (h *FTPHandler) handleStor(session *FTPSession, filename string) {
	if !session.authenticated {
		h.sendResponse(session.writer, ftp.ResponseNotLoggedIn)
		return
	}

	filePath := filepath.Join(session.rootDir, session.currentDir, filename)
	if !h.isValidPath(filePath, session.rootDir) {
		h.sendResponse(session.writer, ftp.ResponseFileNotFound)
		return
	}

	file, err := os.Create(filePath)
	if err != nil {
		h.sendResponse(session.writer, ftp.ResponseFileNotFound)
		return
	}
	defer file.Close()

	h.sendResponse(session.writer, ftp.ResponseDataConnection)

	// Simple file receive - in a real implementation, you'd use a data connection
	buffer := make([]byte, 1024)
	for {
		n, err := session.reader.Read(buffer)
		if err != nil || n == 0 {
			break
		}
		file.Write(buffer[:n])
	}

	h.sendResponse(session.writer, ftp.ResponseTransferComplete)
}

// handleQuit handles the QUIT command
func (h *FTPHandler) handleQuit(session *FTPSession) {
	h.sendResponse(session.writer, ftp.ResponseGoodbye)
	session.conn.Close()
}

// handleNoop handles the NOOP command
func (h *FTPHandler) handleNoop(session *FTPSession) {
	h.sendResponse(session.writer, ftp.ResponseOK)
}

// authenticate validates user credentials
func (h *FTPHandler) authenticate(username, password string) bool {
	if h.config.Auth.Anonymous && username == "anonymous" {
		return true
	}

	if storedPassword, exists := h.config.Auth.Users[username]; exists {
		return storedPassword == password
	}

	return false
}

// isValidPath checks if a path is valid and within the root directory
func (h *FTPHandler) isValidPath(path, rootDir string) bool {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	absRoot, err := filepath.Abs(rootDir)
	if err != nil {
		return false
	}

	return strings.HasPrefix(absPath, absRoot)
}

// sendResponse sends a response to the client
func (h *FTPHandler) sendResponse(writer *bufio.Writer, response string) {
	writer.WriteString(response + "\r\n")
	writer.Flush()
	logrus.Debugf("Sent response: %s", response)
}

// FTPSession represents a client session
type FTPSession struct {
	conn          net.Conn
	reader        *bufio.Reader
	writer        *bufio.Writer
	config        *config.Config
	rootDir       string
	currentDir    string
	username      string
	authenticated bool
}
