package logga

import (
	"fmt"
	"io"
	"os"
	"text/template"
	"time"
)

type Level int
type Logger struct {
	level      Level
	formatter  Formatter
	timeFormat string
}

type Formatter interface {
	Format(message interface{}, out io.Writer)
}

type textFormatter struct {
	template *template.Template
}

func newTextFormatter(tmplText string) *textFormatter {
	tmpl, _ := template.New("").Parse(tmplText)
	return &textFormatter{
		template: tmpl,
	}
}

func (f textFormatter) Format(message interface{}, out io.Writer) {
	f.template.Execute(out, message)
}

type Option func(*Logger) error

const (
	All Level = iota
	Debug
	Warning
	Error
	Fatal
	Off
)

var levelDescription = map[Level]string{
	Debug:   "DEBUG",
	Warning: "WARN ",
	Error:   "ERROR",
	Fatal:   "FATAL",
}

func WithLevel(level Level) Option {
	return func(l *Logger) error {
		l.level = level
		return nil
	}
}

func WithMessageTemplate(tmplText string) Option {
	return func(l *Logger) error {
		l.formatter = newTextFormatter(tmplText)
		return nil
	}
}

func NewLogger(opts ...Option) *Logger {
	l := &Logger{}
	l.timeFormat = time.RFC3339
	l.formatter = newTextFormatter("{{.Level}} - {{.Time}} - {{.Message}}\n")
	for _, opt := range opts {
		opt(l)
	}
	return l
}

func (l Logger) Debugf(format string, args ...interface{}) {
	l.printf(Debug, format, args...)
}

func (l Logger) Warningf(format string, args ...interface{}) {
	l.printf(Warning, format, args...)
}

func (l Logger) Errorf(format string, args ...interface{}) {
	l.printf(Error, format, args...)
}

func (l Logger) Fatalf(format string, args ...interface{}) {
	l.printf(Fatal, format, args...)
}

func (l Logger) printf(level Level, format string, args ...interface{}) {
	if level >= l.level {
		message := fmt.Sprintf(format, args...)
		l.formatter.Format(struct {
			Message string
			Time    string
			Level   string
		}{
			Message: message,
			Time:    time.Now().Format(l.timeFormat),
			Level:   levelDescription[level],
		}, os.Stderr)
	}
}
