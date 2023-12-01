//go:build windows

package utils

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

func IsHidden(path string) bool {
	// Get the base name of the file/folder
	base := filepath.Base(path)

	// Check if the base name starts with a dot (hidden on Unix-like systems)
	if strings.HasPrefix(base, ".") {
		return true
	}

	// Check if the file/folder has the hidden attribute on Windows
	info, err := os.Stat(path)
	if err == nil {
		attrs := info.Sys().(*syscall.Win32FileAttributeData).FileAttributes
		return attrs&syscall.FILE_ATTRIBUTE_HIDDEN != 0
	}

	return false
}
