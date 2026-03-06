package config

import (
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorGray   = "\033[37m"
	colorCyan   = "\033[36m"
)

type ColorHandler struct {
	opts slog.HandlerOptions
	out  io.Writer
}

func NewColorHandler(out io.Writer, opts *slog.HandlerOptions) *ColorHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	return &ColorHandler{
		opts: *opts,
		out:  out,
	}
}

func (h *ColorHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.opts.Level.Level()
}

func (h *ColorHandler) Handle(ctx context.Context, r slog.Record) error {
	// Choose color based on level
	var color string
	switch r.Level {
	case slog.LevelDebug:
		color = colorGray
	case slog.LevelInfo:
		color = colorBlue
	case slog.LevelWarn:
		color = colorYellow
	case slog.LevelError:
		color = colorRed
	default:
		color = colorReset
	}

	// Format: [LEVEL] source:line message key=value
	var buf strings.Builder

	// Level with color
	buf.WriteString(color)
	buf.WriteString("[" + r.Level.String() + "]")
	buf.WriteString(colorReset)
	buf.WriteString(" ")

	// Source location if enabled
	if h.opts.AddSource && r.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		if f.File != "" {
			buf.WriteString(colorCyan)
			buf.WriteString(filepath.Base(f.File) + ":" + strconv.Itoa(f.Line))
			buf.WriteString(colorReset)
			buf.WriteString(" ")
		}
	}

	// Message
	buf.WriteString(r.Message)

	// Attributes
	r.Attrs(func(a slog.Attr) bool {
		buf.WriteString(" ")
		buf.WriteString(a.Key)
		buf.WriteString("=")
		buf.WriteString(a.Value.String())
		return true
	})

	buf.WriteString("\n")
	_, err := h.out.Write([]byte(buf.String()))
	return err
}

func (h *ColorHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// For simplicity, return the same handler
	return h
}

func (h *ColorHandler) WithGroup(name string) slog.Handler {
	// For simplicity, return the same handler
	return h
}

func init() {
	// Configure slog with colored output and source location information
	opts := &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	}
	handler := NewColorHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
