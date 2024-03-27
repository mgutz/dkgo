package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// HOME is user's home dir
var HOME string

func init() {
	dir, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Errorf("could not get user home dir: %w", err))
	}
	HOME = dir
}

// Untildify replaces leading "~" with user's home directory.
func Untildify(path string) string {
	return Unhomify(path, HOME)
}

// Tildify replace $HOME with "~".
func Tildify(p string) string {
	return Homify(p, HOME)
}

// Homify replace home with "~".
func Homify(p string, home string) string {
	if strings.HasPrefix(p, home) {
		return "~" + p[len(home):]
	}

	return p
}

// Unhomify replaces leading "~" with user's home directory.
func Unhomify(path string, home string) string {
	if strings.HasPrefix(path, "~/") || strings.HasPrefix(path, "~\\") {
		return strings.Replace(path, "~", home, 1)
	}
	return path
}

// BaseNoExt returns the base without its extension from the last path segment.
func BaseNoExt(path string) string {
	if path == "" {
		return ""
	}

	base := filepath.Base(path)
	ext := filepath.Ext(path)
	return base[0 : len(base)-len(ext)]
}
