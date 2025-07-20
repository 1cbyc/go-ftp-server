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

type FTPHandler struct {
	config *config.Config
	mu     sync.RWMutex
}

func NewFTPHandler(cfg *config.Config) *FTPHandler {
	return &FTPHandler{
		config: cfg,
	}
}

func (h *FTPHandler) HandleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	h.sendResponse(writer, ftp.ResponseWelcome)

	session := &FTPSession{
		conn:          conn,
		reader:        reader,
		writer:        writer,
		config:        h.config,
		rootDir:       h.config.FTP.RootDir,
		currentDir:    ".",
		authenticated: false,
	}

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

func (h *FTPHandler) handleUser(session *FTPSession, username string) {
	session.username = username
	h.sendResponse(session.writer, ftp.ResponseUsernameOK)
}

func (h *FTPHandler) handlePass(session *FTPSession, password string) {
	if h.authenticate(session.username, password) {
		session.authenticated = true
		h.sendResponse(session.writer, ftp.ResponseLoginOK)
	} else {
		h.sendResponse(session.writer, ftp.ResponseLoginFailed)
	}
}

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

func (h *FTPHandler) handlePwd(session *FTPSession) {
	if !session.authenticated {
		h.sendResponse(session.writer, ftp.ResponseNotLoggedIn)
		return
	}

	response := fmt.Sprintf("257 \"%s\" is current directory", session.currentDir)
	h.sendResponse(session.writer, response)
}

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

func (h *FTPHandler) handleQuit(session *FTPSession) {
	h.sendResponse(session.writer, ftp.ResponseGoodbye)
	session.conn.Close()
}

func (h *FTPHandler) handleNoop(session *FTPSession) {
	h.sendResponse(session.writer, ftp.ResponseOK)
}

func (h *FTPHandler) authenticate(username, password string) bool {
	if h.config.Auth.Anonymous && username == "anonymous" {
		return true
	}

	if storedPassword, exists := h.config.Auth.Users[username]; exists {
		return storedPassword == password
	}

	return false
}

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

func (h *FTPHandler) sendResponse(writer *bufio.Writer, response string) {
	writer.WriteString(response + "\r\n")
	writer.Flush()
	logrus.Debugf("Sent response: %s", response)
}

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
