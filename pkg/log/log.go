package log

import (
	"context"
)

type ILogger interface {
	Debug(msg string, fields ...Field)
	Debugf(format string, args ...interface{})
	Info(msg string, fields ...Field)
	Infof(format string, args ...interface{})
	Warn(msg string, fields ...Field)
	Warnf(format string, args ...interface{})
	Error(msg string, err error, fields ...Field)
	Errorf(err error, format string, args ...interface{})

	With(fields ...Field) ILogger
	WithContext(ctx context.Context) ILogger
}

type Field struct {
	Key   string
	Value any
}

func F(key string, value any) Field {
	return Field{Key: key, Value: value}
}

type Format string
type Output string

const (
	FormatJSON Format = "json"
	FormatText Format = "text"

	OutputStdout Output = "stdout"
	OutputFile   Output = "file"
)

type Config struct {
	Level   string `env:"LOG_LEVEL" envDefault:"info"`
	Service string `env:"LOG_SERVICE" envDefault:"layer8-server"`

	Format Format `env:"LOG_FORMAT" envDefault:"text"`
	Output Output `env:"LOG_OUTPUT" envDefault:"stdout"`

	FilePath string `env:"LOG_FILE_PATH"` // used only if OutputFile
}
