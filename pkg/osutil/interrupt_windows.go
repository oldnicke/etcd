//go:build windows
// +build windows

package osutil

import (
	"os"

	"go.uber.org/zap"
)

type InterruptHandler func()

// RegisterInterruptHandler is a no-op on windows
func RegisterInterruptHandler(h InterruptHandler) {}

// HandleInterrupts is a no-op on windows
func HandleInterrupts(*zap.Logger) {}

// Exit calls os.Exit
func Exit(code int) {
	os.Exit(code)
}
