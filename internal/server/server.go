package server

import (
	"fmt"
	"net"
	"sync"

	"github.com/mrinalxdev/go-ftp-server/internal/config"
	"github.com/mrinalxdev/go-ftp-server/internal/handler"
	"github.com/sirupsen/logrus"
)

// FTPServer represents the main FTP server
type FTPServer struct {
	config   *config.Config
	listener net.Listener
	handler  *handler.FTPHandler
	shutdown chan struct{}
	wg       sync.WaitGroup
}

// NewFTPServer creates a new FTP server instance
func NewFTPServer(cfg *config.Config) *FTPServer {
	return &FTPServer{
		config:   cfg,
		shutdown: make(chan struct{}),
	}
}

// Start starts the FTP server
func (s *FTPServer) Start() error {
	// Create root directory if it doesn't exist
	if err := s.ensureRootDir(); err != nil {
		return fmt.Errorf("failed to create root directory: %w", err)
	}

	// Create listener
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port))
	if err != nil {
		return fmt.Errorf("failed to create listener: %w", err)
	}
	s.listener = listener

	// Create handler
	s.handler = handler.NewFTPHandler(s.config)

	logrus.Infof("FTP server listening on %s:%d", s.config.Server.Host, s.config.Server.Port)

	// Accept connections
	for {
		select {
		case <-s.shutdown:
			return nil
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				if s.isShutdown() {
					return nil
				}
				logrus.Errorf("Failed to accept connection: %v", err)
				continue
			}

			s.wg.Add(1)
			go s.handleConnection(conn)
		}
	}
}

// Shutdown gracefully shuts down the server
func (s *FTPServer) Shutdown() {
	logrus.Info("Shutting down FTP server...")
	close(s.shutdown)

	if s.listener != nil {
		s.listener.Close()
	}

	s.wg.Wait()
	logrus.Info("FTP server shutdown complete")
}

// handleConnection handles a single client connection
func (s *FTPServer) handleConnection(conn net.Conn) {
	defer func() {
		conn.Close()
		s.wg.Done()
	}()

	logrus.Debugf("New connection from %s", conn.RemoteAddr())
	s.handler.HandleConnection(conn)
}

// ensureRootDir creates the root directory if it doesn't exist
func (s *FTPServer) ensureRootDir() error {
	// This will be implemented when we add the file system utilities
	return nil
}

// isShutdown checks if the server is shutting down
func (s *FTPServer) isShutdown() bool {
	select {
	case <-s.shutdown:
		return true
	default:
		return false
	}
}
