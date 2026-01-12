package slogpretty

import (
	"context"
	"fmt"
	"io"
	"time"

	"log/slog"

	"github.com/bytedance/sonic"
	"github.com/fatih/color"
)

type PrettyHandlerOptions struct {
	SlogOpts *slog.HandlerOptions
}

type PrettyHandler struct {
	slog.Handler

	opts  PrettyHandlerOptions
	out   io.Writer
	attrs []slog.Attr
}

func (opts PrettyHandlerOptions) NewPrettyHandler(out io.Writer) *PrettyHandler {
	return &PrettyHandler{
		Handler: slog.NewJSONHandler(out, opts.SlogOpts),
		out:     out,
	}
}

func (h *PrettyHandler) Handle(_ context.Context, r slog.Record) error {
	level := formatLevel(r.Level)
	fields := make(map[string]any, r.NumAttrs())

	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()
		return true
	})

	for _, a := range h.attrs {
		fields[a.Key] = a.Value.Any()
	}

	var fieldsJSON string
	if len(fields) > 0 {
		b, err := sonic.ConfigFastest.MarshalIndent(fields, "", "  ")
		if err == nil {
			fieldsJSON = color.WhiteString(string(b))
		}
	}

	timeStr := color.HiBlackString("[%s]", r.Time.Format(time.StampMilli))
	msg := color.CyanString(r.Message)

	_, _ = fmt.Fprintf(h.out, "%s %s %s %s\n", timeStr, level, msg, fieldsJSON)
	return nil
}

func formatLevel(level slog.Level) string {
	switch level {
	case slog.LevelDebug:
		return color.MagentaString("DEBUG:")
	case slog.LevelInfo:
		return color.BlueString("INFO:")
	case slog.LevelWarn:
		return color.YellowString("WARN:")
	case slog.LevelError:
		return color.RedString("ERROR:")
	default:
		return fmt.Sprintf("%s:", level.String())
	}
}

func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrs := make([]slog.Attr, len(h.attrs)+len(attrs))
	copy(newAttrs, h.attrs)
	copy(newAttrs[len(h.attrs):], attrs)

	return &PrettyHandler{
		Handler: h.Handler,
		out:     h.out,
		attrs:   newAttrs,
	}
}

func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	return &PrettyHandler{
		Handler: h.Handler.WithGroup(name),
		out:     h.out,
		attrs:   h.attrs,
	}
}
