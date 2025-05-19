
package main

import (
	"path/filepath"
	"strings"
)

func IsSupportedFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".pdf", ".docx", ".txt":
		return true
	default:
		return false
	}
}
