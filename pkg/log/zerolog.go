package log

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	log zerolog.Logger
}

func NewLogger(cfg Config) ILogger {
	level, _ := zerolog.ParseLevel(cfg.Level)
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	writer := buildWriter(cfg)

	z := zerolog.New(writer).
		Level(level).
		With().
		Timestamp().
		Str("service", cfg.Service).
		Logger()

	return &Logger{log: z}
}

func (l *Logger) Debug(msg string, fields ...Field) {
	l.log.Debug().Fields(toMap(fields)).Msg(msg)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.log.Debug().Msgf(format, args...)
}

func (l *Logger) Info(msg string, fields ...Field) {
	l.log.Info().Fields(toMap(fields)).Msg(msg)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.log.Info().Msgf(format, args...)
}

func (l *Logger) Warn(msg string, fields ...Field) {
	l.log.Warn().Fields(toMap(fields)).Msg(msg)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.log.Warn().Msgf(format, args...)
}

func (l *Logger) Error(msg string, err error, fields ...Field) {
	l.log.Error().
		Stack().
		Err(err).
		Fields(toMap(fields)).
		Msg(msg)
}

func (l *Logger) Errorf(err error, format string, args ...interface{}) {
	l.log.Error().Stack().Err(err).Msgf(format, args...)
}

func (l *Logger) With(fields ...Field) ILogger {
	return &Logger{
		log: l.log.With().Fields(toMap(fields)).Logger(),
	}
}

func (l *Logger) WithContext(ctx context.Context) ILogger {
	return l
}

func toMap(fields []Field) map[string]any {
	m := make(map[string]any, len(fields))
	for _, f := range fields {
		m[f.Key] = f.Value
	}
	return m
}

func buildWriter(cfg Config) io.Writer {
	var w io.Writer

	switch cfg.Output {
	case OutputFile:
		w = &lumberjack.Logger{
			Filename: cfg.FilePath,
			MaxSize:  50,
		}
	default:
		w = os.Stdout
	}

	if cfg.Format == FormatText {
		return zerolog.ConsoleWriter{
			Out:        w,
			TimeFormat: time.RFC3339,
		}
	}

	return w // JSON
}
