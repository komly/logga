package logga

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
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
	SetOption(Option)
}
type Option func(Logger) error

type LogRecord struct {
	Message string `json:"message"`
	Time    string `json:"time"`
	Level   string `json:"level"`
}

type logger struct {
	level      Level
	formatter  Formatter
	timeFormat string
	out        io.Writer
	mutex      sync.Mutex
}

type Formatter interface {
	Format(message *LogRecord, out io.Writer)
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

func (f textFormatter) Format(message *LogRecord, out io.Writer) {
	f.template.Execute(out, message)
}

type JSONFormatter struct {
}

func (j JSONFormatter) Format(message *LogRecord, out io.Writer) {
	encoder := json.NewEncoder(out)
	encoder.Encode(message)
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
	Warning: "WARNING",
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

func WithFormatter(fmt Formatter) Option {
	return func(l Logger) error {
		l.(*logger).formatter = fmt
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

func (l *logger) SetOption(opt Option) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	opt(l)
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
	os.Exit(1)
}

func (l logger) printf(level Level, format string, args ...interface{}) {
	if level >= l.level {
		message := fmt.Sprintf(format, args...)
		l.formatter.Format(&LogRecord{
			Message: message,
			Time:    time.Now().Format(l.timeFormat),
			Level:   levelDescription[level],
		}, l.out)
	}
}
