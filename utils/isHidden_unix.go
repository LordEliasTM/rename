//go:build unix

package utils

import (
	"path/filepath"
	"strings"
)

func IsHidden(path string) bool {
	// Get the base name of the file/folder
	base := filepath.Base(path)

	// Check if the base name starts with a dot (hidden on Unix-like systems)
	return strings.HasPrefix(base, ".")
}
