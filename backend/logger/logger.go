package logger

import (
	kitlog "github.com/go-kit/kit/log"
	stdlog "log"
	"os"
	"time"
)

var logger kitlog.Logger

func Initialize() {
	// initialize logging
	logger = kitlog.NewJSONLogger(os.Stdout)
	stdlog.SetOutput(kitlog.NewStdlibAdapter(logger))
}

func Log(msg string, m ...string) {
	logger.Log("msg", msg, "ts", time.Now(), "level", "info", "params", m)
}

func Warn(msg string, m ...string) {
	logger.Log("msg", msg, "ts", time.Now(), "level", "warn", "params", m)
}

func Error(msg string, err error, m ...string) {
	logger.Log("msg", msg, "ts", time.Now(), "level", "error", "error", err.Error(), "params", m)
}
