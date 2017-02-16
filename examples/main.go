package main

import (
	"github.com/Komly/logga"
)

func main() {
	logga := logga.NewLogger(
		logga.WithLevel(logga.Error),
		logga.WithMessageTemplate("{{.Level}} - {{.Time}} -  {{.Message}}\n"),
	)

	logga.Debugf("Debug message: %d", 2)
	logga.Warningf("Warning message: %d", 2)
	logga.Errorf("Error message: %d", 2)
	logga.Fatalf("Fatal message: %d", 2)
}
