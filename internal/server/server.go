package server

import (
	"fmt"
	"net"
	"os"
	"sync"

	"github.com/mrinalxdev/go-ftp-server/internal/config"
	"github.com/mrinalxdev/go-ftp-server/internal/handler"
	"github.com/sirupsen/logrus"
)

type FTPServer struct {
	config   *config.Config
	listener net.Listener
	handler  *handler.FTPHandler
	shutdown chan struct{}
	wg       sync.WaitGroup
}

func NewFTPServer(cfg *config.Config) *FTPServer {
	return &FTPServer{
		config:   cfg,
		shutdown: make(chan struct{}),
	}
}

func (s *FTPServer) Start() error {
	if err := s.ensureRootDir(); err != nil {
		return fmt.Errorf("failed to create root directory: %w", err)
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port))
	if err != nil {
		return fmt.Errorf("failed to create listener: %w", err)
	}
	s.listener = listener

	s.handler = handler.NewFTPHandler(s.config)

	logrus.Infof("FTP server listening on %s:%d", s.config.Server.Host, s.config.Server.Port)

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

func (s *FTPServer) Shutdown() {
	logrus.Info("Shutting down FTP server...")
	close(s.shutdown)

	if s.listener != nil {
		s.listener.Close()
	}

	s.wg.Wait()
	logrus.Info("FTP server shutdown complete")
}

func (s *FTPServer) handleConnection(conn net.Conn) {
	defer func() {
		conn.Close()
		s.wg.Done()
	}()

	logrus.Debugf("New connection from %s", conn.RemoteAddr())
	s.handler.HandleConnection(conn)
}

func (s *FTPServer) ensureRootDir() error {
	if err := os.MkdirAll(s.config.FTP.RootDir, 0755); err != nil {
		return fmt.Errorf("failed to create root directory %s: %w", s.config.FTP.RootDir, err)
	}
	logrus.Infof("Root directory ready: %s", s.config.FTP.RootDir)
	return nil
}

func (s *FTPServer) isShutdown() bool {
	select {
	case <-s.shutdown:
		return true
	default:
		return false
	}
}
