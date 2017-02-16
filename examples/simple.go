package main

import (
	"github.com/Komly/logga"
	"os"
)

func main() {
	logga := logga.NewLogger(
		logga.WithLevel(logga.Error),
		logga.WithMessageTemplate("{{.Level}} - {{.Time}} -  {{.Message}}\n"),
		logga.WithOutput(os.Stdout),
	)

	logga.Debugf("Debug message: %d", 1)
	logga.Debugf("Info message: %d", 2)
	logga.Warningf("Warning message: %d", 3)
	logga.Errorf("Error message: %d", 4)
	logga.Fatalf("Fatal message: %d", 5)
}
