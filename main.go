package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mrinalxdev/go-ftp-server/internal/config"
	"github.com/mrinalxdev/go-ftp-server/internal/server"
	"github.com/sirupsen/logrus"
)

var (
	configFile = flag.String("config", "config.yaml", "Path to configuration file")
	port       = flag.Int("port", 2121, "Port to listen on")
	host       = flag.String("host", "localhost", "Host to bind to")
	rootDir    = flag.String("root", "./ftp_root", "Root directory for FTP server")
	verbose    = flag.Bool("verbose", false, "Enable verbose logging")
)

func main() {
	flag.Parse()

	logLevel := logrus.InfoLevel
	if *verbose {
		logLevel = logrus.DebugLevel
	}
	logrus.SetLevel(logLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	cfg, err := config.Load(*configFile)
	if err != nil {
		logrus.Warnf("Failed to load config file %s: %v, using defaults", *configFile, err)
		cfg = &config.Config{
			Server: config.ServerConfig{
				Host: *host,
				Port: *port,
			},
			FTP: config.FTPConfig{
				RootDir: *rootDir,
			},
		}
	}

	if *host != "localhost" {
		cfg.Server.Host = *host
	}
	if *port != 2121 {
		cfg.Server.Port = *port
	}
	if *rootDir != "./ftp_root" {
		cfg.FTP.RootDir = *rootDir
	}

	srv := server.NewFTPServer(cfg)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		logrus.Infof("Received signal %v, shutting down gracefully...", sig)
		srv.Shutdown()
	}()

	logrus.Infof("Starting FTP server on %s:%d", cfg.Server.Host, cfg.Server.Port)
	logrus.Infof("Root directory: %s", cfg.FTP.RootDir)

	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
