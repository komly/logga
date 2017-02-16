# Logga - simple, "convention over configuration" logger for golang

## Installation

```sh
go get "github.com/Komly/logga"
```

## Usage
```go
logga := logga.NewLogger(
    logga.WithLevel(logga.Error),
    logga.WithMessageTemplate("{{.Level}} - {{.Time}} -  {{.Message}}\n"),
)
logga.Debugf("Debug message: %d", 1)
logga.Warningf("Warning message: %d", 2)
logga.Errorf("Error message: %d", 3)
logga.Fatalf("Fatal message: %d", 4)
```
