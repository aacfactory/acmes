package server

import (
	"fmt"
	"github.com/aacfactory/logs"
	"github.com/go-acme/lego/v4/log"
	slog "log"
	"strings"
)

func createLog(level string, formatter string) (v logs.Logger, err error) {
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
	lf := logs.TextFormatter
	formatter = strings.TrimSpace(strings.ToLower(formatter))
	switch formatter {
	case "text_colorful":
		lf = logs.ColorTextFormatter
		break
	case "json":
		lf = logs.JsonFormatter
		break
	default:
		break
	}
	v, err = logs.New(
		logs.WithConsoleWriterFormatter(lf),
		logs.WithLevel(logLevel),
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
