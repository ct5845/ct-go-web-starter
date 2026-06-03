package colorhandler

import (
	"context"
	"io"
	"log/slog"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorGray   = "\033[37m"
	colorCyan   = "\033[36m"
	colorGreen  = "\033[32m"
	colorDim    = "\033[2m"
)

type ColorHandler struct {
	opts slog.HandlerOptions
	out  io.Writer
}

func New(out io.Writer, opts *slog.HandlerOptions) *ColorHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	return &ColorHandler{
		opts: *opts,
		out:  out,
	}
}

func (h *ColorHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.opts.Level.Level()
}

func (h *ColorHandler) Handle(_ context.Context, r slog.Record) error {
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

	var buf strings.Builder

	buf.WriteString(color)
	buf.WriteString("[" + r.Level.String() + "]")
	buf.WriteString(colorReset)
	buf.WriteString(" ")

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

	buf.WriteString(r.Message)

	// Collect attrs, handling request_id and timing specially.
	var reqID, timing string
	var attrs []string

	r.Attrs(func(a slog.Attr) bool {
		switch a.Key {
		case "request_id":
			id := a.Value.String()
			if len(id) > 8 {
				id = id[len(id)-8:]
			}
			reqID = id
		case "timing":
			timing = a.Value.String()
		default:
			attrs = append(attrs, a.Key+"="+a.Value.String())
		}
		return true
	})

	if reqID != "" {
		buf.WriteString(" " + colorGreen + "[" + reqID + "]" + colorReset)
	}
	for _, attr := range attrs {
		buf.WriteString(" " + colorDim + attr + colorReset)
	}
	if timing != "" {
		buf.WriteString("\n    " + colorDim + timing + colorReset)
	}

	buf.WriteString("\n")
	_, err := h.out.Write([]byte(buf.String()))
	return err
}

func (h *ColorHandler) WithAttrs(_ []slog.Attr) slog.Handler { return h }
func (h *ColorHandler) WithGroup(_ string) slog.Handler      { return h }
