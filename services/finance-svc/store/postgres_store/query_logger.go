package postgres_store

import (
	"context"
	"fmt"
	"strings"

	"finance-svc/store/postgres_store/sqlc"
	"finance-svc/util/logging"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// LoggingDBTX wraps a DBTX interface to log all SQL queries with parameters substituted
type LoggingDBTX struct {
	db     sqlc.DBTX
	logger *zerolog.Logger
	tracer trace.Tracer
}

// NewLoggingDBTX creates a new LoggingDBTX wrapper
func NewLoggingDBTX(db sqlc.DBTX, logger *zerolog.Logger, tracer trace.Tracer) *LoggingDBTX {
	return &LoggingDBTX{
		db:     db,
		logger: logger,
		tracer: tracer,
	}
}

// Exec executes a query and logs it with parameters substituted
func (l *LoggingDBTX) Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	const op = "postgres_store.LoggingDBTX.Exec"

	// Start span
	ctx, span := l.tracer.Start(ctx, op)
	defer span.End()

	// Format query with parameters
	loggedQuery := l.formatQuery(query, args...)

	// Get query name from sqlc comment if available
	queryName := l.extractQueryName(query)

	// Set span attributes
	span.SetAttributes(
		attribute.String("db.operation", "exec"),
		attribute.String("db.query.name", queryName),
		attribute.String("db.query", loggedQuery),
	)

	// Get logger with trace ID
	log := logging.LogWithTrace(ctx, l.logger)
	log.Debug().
		Str("op", op).
		Str("query", loggedQuery).
		Str("type", "exec").
		Msg("Executing SQL query")

	// Execute query
	result, err := l.db.Exec(ctx, query, args...)

	// Set span status
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetAttributes(
			attribute.String("db.rows_affected", result.String()),
		)
		span.SetStatus(codes.Ok, "success")
	}

	return result, err
}

// Query executes a query and logs it with parameters substituted
func (l *LoggingDBTX) Query(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	const op = "postgres_store.LoggingDBTX.Query"

	// Start span
	ctx, span := l.tracer.Start(ctx, op)
	defer span.End()

	// Format query with parameters
	loggedQuery := l.formatQuery(query, args...)

	// Get query name from sqlc comment if available
	queryName := l.extractQueryName(query)

	// Set span attributes
	span.SetAttributes(
		attribute.String("db.operation", "query"),
		attribute.String("db.query.name", queryName),
		attribute.String("db.query", loggedQuery),
	)

	// Get logger with trace ID
	log := logging.LogWithTrace(ctx, l.logger)
	log.Debug().
		Str("op", op).
		Str("query", loggedQuery).
		Str("type", "query").
		Msg("Executing SQL query")

	// Execute query
	rows, err := l.db.Query(ctx, query, args...)

	// Set span status
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetStatus(codes.Ok, "success")
	}

	return rows, err
}

// QueryRow executes a query and logs it with parameters substituted
func (l *LoggingDBTX) QueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	const op = "postgres_store.LoggingDBTX.QueryRow"

	// Start span
	ctx, span := l.tracer.Start(ctx, op)
	defer span.End()

	// Format query with parameters
	loggedQuery := l.formatQuery(query, args...)

	// Get query name from sqlc comment if available
	queryName := l.extractQueryName(query)

	// Set span attributes
	span.SetAttributes(
		attribute.String("db.operation", "query_row"),
		attribute.String("db.query.name", queryName),
		attribute.String("db.query", loggedQuery),
	)

	// Get logger with trace ID
	log := logging.LogWithTrace(ctx, l.logger)
	log.Debug().
		Str("op", op).
		Str("query", loggedQuery).
		Str("type", "query_row").
		Msg("Executing SQL query")

	// Execute query
	row := l.db.QueryRow(ctx, query, args...)

	// Note: QueryRow doesn't return error immediately, so we can't set span status here
	// The error will be returned when Scan() is called

	return row
}

// formatQuery formats a SQL query by replacing $1, $2, etc. with actual parameter values
func (l *LoggingDBTX) formatQuery(query string, args ...any) string {
	// Clean up the query first (remove sqlc comments, normalize whitespace)
	cleanedQuery := l.cleanQuery(query)

	// Replace parameters if any
	if len(args) == 0 {
		return cleanedQuery
	}

	result := cleanedQuery
	for i, arg := range args {
		placeholder := fmt.Sprintf("$%d", i+1)
		formattedValue := l.formatValue(arg)
		result = strings.ReplaceAll(result, placeholder, formattedValue)
	}

	return result
}

// cleanQuery removes sqlc comment headers and converts query to single line for cleaner logs
func (l *LoggingDBTX) cleanQuery(query string) string {
	lines := strings.Split(query, "\n")
	var cleanedLines []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Skip sqlc comment headers (-- name: ... :exec/:many/:one)
		if strings.HasPrefix(trimmed, "-- name:") {
			continue
		}

		// Skip empty lines
		if trimmed == "" {
			continue
		}

		cleanedLines = append(cleanedLines, trimmed)
	}

	// Join lines with single space to create a single-line query
	result := strings.Join(cleanedLines, " ")

	// Normalize multiple consecutive spaces to single space
	for strings.Contains(result, "  ") {
		result = strings.ReplaceAll(result, "  ", " ")
	}

	// Trim any leading/trailing whitespace
	result = strings.TrimSpace(result)

	return result
}

// extractQueryName extracts the query name from sqlc comment header
// Example: "-- name: InsertBRJCorrection :exec" -> "InsertBRJCorrection"
func (l *LoggingDBTX) extractQueryName(query string) string {
	lines := strings.Split(query, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "-- name:") {
			// Extract name between "-- name: " and " :"
			parts := strings.Split(trimmed, " ")
			if len(parts) >= 3 {
				return parts[2] // The query name
			}
		}
	}
	return "unknown"
}

// formatValue formats a value for SQL logging
func (l *LoggingDBTX) formatValue(v any) string {
	if v == nil {
		return "NULL"
	}

	switch val := v.(type) {
	case string:
		return fmt.Sprintf("'%s'", strings.ReplaceAll(val, "'", "''"))
	case *string:
		if val == nil {
			return "NULL"
		}
		return fmt.Sprintf("'%s'", strings.ReplaceAll(*val, "'", "''"))
	case pgtype.Text:
		if !val.Valid {
			return "NULL"
		}
		return fmt.Sprintf("'%s'", strings.ReplaceAll(val.String, "'", "''"))
	case pgtype.Int4:
		if !val.Valid {
			return "NULL"
		}
		return fmt.Sprintf("%d", val.Int32)
	case pgtype.Int8:
		if !val.Valid {
			return "NULL"
		}
		return fmt.Sprintf("%d", val.Int64)
	case pgtype.Float8:
		if !val.Valid {
			return "NULL"
		}
		return fmt.Sprintf("%f", val.Float64)
	case pgtype.Bool:
		if !val.Valid {
			return "NULL"
		}
		if val.Bool {
			return "TRUE"
		}
		return "FALSE"
	case pgtype.Date:
		if !val.Valid {
			return "NULL"
		}
		return fmt.Sprintf("'%s'", val.Time.Format("2006-01-02"))
	case pgtype.Timestamp:
		if !val.Valid {
			return "NULL"
		}
		return fmt.Sprintf("'%s'", val.Time.Format("2006-01-02 15:04:05"))
	case pgtype.Timestamptz:
		if !val.Valid {
			return "NULL"
		}
		return fmt.Sprintf("'%s'", val.Time.Format("2006-01-02 15:04:05"))
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", val)
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", val)
	case float32, float64:
		return fmt.Sprintf("%f", val)
	case bool:
		if val {
			return "TRUE"
		}
		return "FALSE"
	default:
		// For unknown types, use fmt.Sprintf and wrap in quotes if it looks like a string
		str := fmt.Sprintf("%v", v)
		// If it contains spaces or special chars, quote it
		if strings.ContainsAny(str, " \t\n\r") || strings.Contains(str, "'") {
			return fmt.Sprintf("'%s'", strings.ReplaceAll(str, "'", "''"))
		}
		return str
	}
}
