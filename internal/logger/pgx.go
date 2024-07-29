package logger

import (
	"context"
	"github.com/jackc/pgx/v5/tracelog"
	"log/slog"
)

type PgxLogger struct {
	log *slog.Logger
}

func NewPgxLogger(log *slog.Logger) *PgxLogger {
	return &PgxLogger{log: log.With(slog.String("lib", "pgx"))}
}

func (l *PgxLogger) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]interface{}) {
	attrs := make([]slog.Attr, 0, len(data))
	for k, v := range data {
		attrs = append(attrs, slog.Any(k, v))
	}

	slogLevel := PgxLevelToSlog(level)
	if level == tracelog.LogLevelTrace {
		attrs = append(attrs, slog.Any("PGX_LOG_LEVEL", level))
	}

	l.log.LogAttrs(ctx, slogLevel, msg, attrs...)
}

func (l *PgxLogger) TraceLog(ctx context.Context) *tracelog.TraceLog {
	return &tracelog.TraceLog{
		Logger:   l,
		LogLevel: l.GetPgxLevel(ctx),
	}
}

func (l *PgxLogger) GetPgxLevel(ctx context.Context) tracelog.LogLevel {
	return SlogLevelToPgx(GetLevel(ctx, l.log))
}

var (
	pgxLevelToSlog = map[tracelog.LogLevel]slog.Level{
		tracelog.LogLevelTrace: slog.LevelDebug,
		tracelog.LogLevelDebug: slog.LevelDebug,
		tracelog.LogLevelInfo:  slog.LevelInfo,
		tracelog.LogLevelWarn:  slog.LevelWarn,
		tracelog.LogLevelError: slog.LevelError,
	}
	slogLevelToPgx = map[slog.Level]tracelog.LogLevel{
		slog.LevelDebug: tracelog.LogLevelTrace,
		slog.LevelInfo:  tracelog.LogLevelInfo,
		slog.LevelWarn:  tracelog.LogLevelWarn,
		slog.LevelError: tracelog.LogLevelError,
	}
)

func PgxLevelToSlog(lvl tracelog.LogLevel) slog.Level {
	slogLevel, ok := pgxLevelToSlog[lvl]
	if !ok {
		return slog.LevelError
	}
	return slogLevel
}

func SlogLevelToPgx(lvl slog.Level) tracelog.LogLevel {
	pgxLevel, ok := slogLevelToPgx[lvl]
	if !ok {
		return tracelog.LogLevelNone
	}
	return pgxLevel
}
