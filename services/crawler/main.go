package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mindcastio/mindcastio/crawler"

	"github.com/mindcastio/mindcastio/backend"
	"github.com/mindcastio/mindcastio/backend/datastore"
	"github.com/mindcastio/mindcastio/backend/environment"
	"github.com/mindcastio/mindcastio/backend/logger"
	"github.com/mindcastio/mindcastio/backend/metrics"
)

func main() {

	// environment setup
	env := environment.GetEnvironment()

	logger.Initialize()
	metrics.Initialize(env)
	datastore.Initialize(env)

	// periodic background processes
	background_channel := time.NewTicker(time.Second * time.Duration(backend.DEFAULT_CRAWLER_SCHEDULE)).C

	// setup shutdown handling
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		shutdown()
		os.Exit(1)
	}()

	// start the scheduler
	logger.Log("mindcast.crawler.startup")
	metrics.Success("mindcast", "crawler.startup", nil)

	for {
		<-background_channel
		crawler.SchedulePodcastCrawling()
	}
}

func shutdown() {
	logger.Log("mindcast.crawler.shutdown")
	metrics.Success("mindcast", "crawler.shutdown", nil)

	datastore.Shutdown()
	metrics.Shutdown()
}
