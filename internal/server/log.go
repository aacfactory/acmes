package server

import (
	"fmt"
	"github.com/aacfactory/logs"
	"github.com/go-acme/lego/v4/log"
	slog "log"
	"os"
	"strings"
)

func createLog(level string) (v logs.Logger, err error) {
	formatter := logs.ConsoleFormatter
	logLevel := logs.ErrorLevel
	levelValue := strings.ToLower(level)
	switch levelValue {
	case "debug":
		logLevel = logs.DebugLevel
	case "info":
		logLevel = logs.InfoLevel
	case "warn":
		logLevel = logs.WarnLevel
	default:
		logLevel = logs.ErrorLevel
	}
	v, err = logs.New(
		logs.WithFormatter(formatter),
		logs.Name("ACMES"),
		logs.WithLevel(logLevel),
		logs.Writer(os.Stdout),
		logs.Color(true),
	)
	if err != nil {
		return
	}
	log.Logger = slog.New(&writer{
		core: v,
	}, "", slog.LstdFlags)
	return
}

type writer struct {
	core logs.Logger
}

func (w *writer) Write(p []byte) (n int, err error) {
	if w.core.DebugEnabled() {
		w.core.Debug().Message(fmt.Sprintf("acmes: %s", string(p)))
		n = len(p)
	}
	return
}
