package util

import (
	"log/slog"
	"os"
)

type TeeWriter struct {
	stdout *os.File
	file   *os.File
}

func (t *TeeWriter) Write(p []byte) (n int, err error) {
	n, err = t.stdout.Write(p)
	if err != nil {
		return n, err
	}
	n, err = t.file.Write(p)
	return n, err
}

// CustomizeSlog initializes slog.
func CustomizeSlog(filename string) {
	opts := slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	file, _ := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)

	writer := &TeeWriter{
		stdout: os.Stdout,
		file:   file,
	}
	h := slog.NewTextHandler(writer, &opts)
	logger := slog.New(h)

	// set global logger with custom options
	slog.SetDefault(logger)
}
