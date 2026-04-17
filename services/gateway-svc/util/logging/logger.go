package logging

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

// ─── Level (pick one and pass in Options) ───────────────────────────────────

// Supported log levels. Use these or zerolog.Level in Options.
var (
	LevelDebug = zerolog.DebugLevel
	LevelInfo  = zerolog.InfoLevel
	LevelWarn  = zerolog.WarnLevel
	LevelError = zerolog.ErrorLevel
)

// DefaultLevel is used when Options.Level is zero.
const DefaultLevel = zerolog.DebugLevel

// ─── Time zone (pick one and pass in Options) ───────────────────────────────

// Supported time zones for log timestamps. Use these or *time.Location in Options.
var (
	TZUTC = time.UTC
	TZWIB = time.FixedZone("WIB", 7*3600) // UTC+7
)

// DefaultTimeZone is used when Options.TimeZone is nil.
var DefaultTimeZone = TZWIB

// ─── Options and NewLogger ──────────────────────────────────────────────────

// Options configures the logger. Pick level and time zone from the constants above.
// Zero values use DefaultLevel and DefaultTimeZone. No config or env; pass options directly.
type Options struct {
	Level    zerolog.Level
	TimeZone *time.Location
}

// NewLogger returns a zerolog.Logger that outputs JSON to stdout.
// Level and time zone come from opts; if zero, defaults are used.
// Set zerolog.TimestampFunc so timestamps use opts.TimeZone (global for the process).
func NewLogger(opts Options) zerolog.Logger {
	level := opts.Level
	if level == 0 {
		level = DefaultLevel
	}
	tz := opts.TimeZone
	if tz == nil {
		tz = DefaultTimeZone
	}
	zerolog.TimestampFunc = func() time.Time { return time.Now().In(tz) }
	return zerolog.New(os.Stdout).
		With().Timestamp().
		Logger().
		Level(level)
}

// LogWithTrace returns a zerolog.Logger enriched with trace_id from the Otel span context.
// If the context has no valid span, the returned logger gets trace_id "unknown".
// Use the returned logger for the rest of the request/operation so logs are correlated.
func LogWithTrace(ctx context.Context, logger *zerolog.Logger) zerolog.Logger {
	sc := trace.SpanContextFromContext(ctx)
	if !sc.IsValid() {
		return logger.With().Str("trace_id", "unknown").Logger()
	}
	return logger.With().Str("trace_id", sc.TraceID().String()).Logger()
}
