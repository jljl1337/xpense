package log

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/jljl1337/xpense/internal/env"
)

func SetCustomLogger() {
	slog.SetDefault(newCustomLogger())
}

func newCustomLogger() *slog.Logger {
	defaultLogLevel := env.MustGetInt("LOG_LEVEL", 0)

	return slog.New(&customHandler{
		level: slog.Level(defaultLogLevel),
	})
}

type customHandler struct {
	level slog.Leveler
}

func (h *customHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.level.Level()
}

func (h *customHandler) Handle(ctx context.Context, r slog.Record) error {
	// Format timestamp like default log
	ts := time.Now().Format("2006/01/02 15:04:05.000")

	lvl := r.Level.String()
	switch lvl {
	case "DEBUG":
		lvl = "DBG"
	case "INFO":
		lvl = "INF"
	case "WARN":
		lvl = "WRN"
	case "ERROR":
		lvl = "ERR"
	}

	msg := r.Message
	formatted := fmt.Sprintf("%s %s %s\n", ts, lvl, msg)

	// Write using standard log or directly to stderr
	_, err := fmt.Fprint(os.Stdout, formatted)
	return err
}

func (h *customHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h // no grouping in this simple example
}

func (h *customHandler) WithGroup(name string) slog.Handler {
	return h
}

func (h *customHandler) Close() error { return nil }

func (h *customHandler) Flush() error { return nil }
