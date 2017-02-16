package logga

import (
	"fmt"
	"io"
	"os"
	"text/template"
	"time"
)

type Level int
type Logger interface {
	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Warningf(string, ...interface{})
	Errorf(string, ...interface{})
	Fatalf(string, ...interface{})
}
type Option func(Logger) error

type logger struct {
	level      Level
	formatter  Formatter
	timeFormat string
	out        io.Writer
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

const (
	All Level = iota
	Debug
	Info
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
	return func(l Logger) error {
		l.(*logger).level = level
		return nil
	}
}

func WithMessageTemplate(tmplText string) Option {
	return func(l Logger) error {
		l.(*logger).formatter = newTextFormatter(tmplText)
		return nil
	}
}

func WithOutput(out io.Writer) Option {
	return func(l Logger) error {
		l.(*logger).out = out
		return nil
	}
}

func NewLogger(opts ...Option) Logger {
	l := &logger{}
	l.out = os.Stderr
	l.timeFormat = time.RFC3339
	l.formatter = newTextFormatter("{{.Level}} - {{.Time}} - {{.Message}}\n")
	for _, opt := range opts {
		opt(l)
	}
	return l
}

func (l logger) Debugf(format string, args ...interface{}) {
	l.printf(Debug, format, args...)
}

func (l logger) Infof(format string, args ...interface{}) {
	l.printf(Info, format, args...)
}

func (l logger) Warningf(format string, args ...interface{}) {
	l.printf(Warning, format, args...)
}

func (l logger) Errorf(format string, args ...interface{}) {
	l.printf(Error, format, args...)
}

func (l logger) Fatalf(format string, args ...interface{}) {
	l.printf(Fatal, format, args...)
}

func (l logger) printf(level Level, format string, args ...interface{}) {
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
		}, l.out)
	}
}
