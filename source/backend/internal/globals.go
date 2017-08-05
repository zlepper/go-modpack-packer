package internal

import (
	"github.com/labstack/echo"
	"strings"
	"sync"
)

var (
	// The running echo instance
	EchoInstance       *echo.Echo
	OutstandingProcess sync.WaitGroup
)

// Gets the port the http process is running on
func GetPort() string {
	parts := strings.Split(EchoInstance.Server.Addr, ":")

	return parts[1]
}
