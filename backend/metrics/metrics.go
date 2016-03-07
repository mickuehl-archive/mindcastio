package metrics

import (
	"strings"

	"github.com/ooyala/go-dogstatsd"

	"github.com/mindcastio/mindcastio/backend/environment"
	"github.com/mindcastio/mindcastio/backend/logger"
)

var _statsd *dogstatsd.Client

func Initialize(env *environment.Environment) {

	if env.StatsServiceEnabled() == false {
		logger.Warn("metrics.initialize", env.StatsServiceHost(), env.StatsServicePort(), "DISABLED")

		_statsd = nil
		return
	}

	logger.Log("metrics.initialize", env.StatsServiceHost(), env.StatsServicePort())

	host := strings.Join([]string{env.StatsServiceHost(), env.StatsServicePort()}, ":")
	c, err := dogstatsd.New(host)

	if err != nil {
		_statsd = nil
		logger.Error("metrics.initialize.error.1", err, host)
	} else {
		// test the connection to be sure ...
		err := c.Success("mindcastio", "metrics.initialize", nil)

		if err == nil {
			_statsd = c
			_statsd.Namespace = "mindcastio."
		} else {
			_statsd = nil
			logger.Error("metrics.initialize.error.2", err, host)
		}
	}
}

func Shutdown() {

	logger.Log("metrics.shutdown")

	if _statsd != nil {
		_statsd.Close()
	}
}

// events
func Info(title string, text string, tags []string) {
	if _statsd != nil {
		err := _statsd.Info(title, text, tags)
		if err != nil {
			logger.Error("metrics.info", err)
		}
	}
}

func Success(title string, text string, tags []string) {
	if _statsd != nil {
		err := _statsd.Success(title, text, tags)
		if err != nil {
			logger.Error("metrics.success", err)
		}
	}
}

func Warning(title string, text string, tags []string) {
	if _statsd != nil {
		err := _statsd.Warning(title, text, tags)
		if err != nil {
			logger.Error("metrics.warning", err)
		}
	}
}

func Error(title string, text string, tags []string) {
	if _statsd != nil {
		err := _statsd.Error(title, text, tags)
		if err != nil {
			logger.Error("metrics.error", err)
		}
	}
}

// metrics

func Count(name string, value int) {
	if _statsd != nil {
		err := _statsd.Count(name, (int64)(value), nil, 1)
		if err != nil {
			logger.Error("metrics.count", err)
		}
	}
}

func Gauge(name string, value float64) {
	if _statsd != nil {
		err := _statsd.Gauge(name, value, nil, 1)
		if err != nil {
			logger.Error("metrics.gauge", err)
		}
	}
}

func Histogram(name string, value float64) {
	if _statsd != nil {
		err := _statsd.Histogram(name, value, nil, 1)
		if err != nil {
			logger.Error("metrics.histogram", err)
		}
	}
}
