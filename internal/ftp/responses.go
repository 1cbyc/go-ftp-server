package ftp

// Standard FTP response codes and messages
const (
	// Connection responses
	ResponseWelcome = "220 Welcome to Go FTP Server"
	ResponseGoodbye = "221 Goodbye"

	// Authentication responses
	ResponseUsernameOK  = "331 User name okay, need password"
	ResponseLoginOK     = "230 User logged in"
	ResponseLoginFailed = "530 Login failed"
	ResponseNotLoggedIn = "530 Please login with USER and PASS"

	// File system responses
	ResponseOK                = "200 OK"
	ResponseDirectoryChanged  = "250 Directory changed"
	ResponseDirectoryNotFound = "550 Directory not found"
	ResponseFileNotFound      = "550 File not found"
	ResponseFileExists        = "550 File already exists"
	ResponsePermissionDenied  = "550 Permission denied"

	// Data transfer responses
	ResponseDataConnection   = "150 Opening BINARY mode data connection"
	ResponseTransferComplete = "226 Transfer complete"
	ResponseTransferFailed   = "426 Transfer failed"

	// Command responses
	ResponseCommandNotImplemented = "502 Command not implemented"
	ResponseSyntaxError           = "501 Syntax error"
	ResponseParameterError        = "504 Parameter not implemented"

	// System responses
	ResponseSystemReady    = "220 Service ready"
	ResponseSystemShutdown = "421 Service shutting down"
)
