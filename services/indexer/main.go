package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mindcastio/mindcastio/search"
	
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
	background_channel := time.NewTicker(time.Second * time.Duration(backend.DEFAULT_INDEXER_SCHEDULE)).C

	// setup shutdown handling
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		shutdown()
		os.Exit(1)
	}()

	// start the scheduler
	logger.Log("mindcast.indexer.startup")
	metrics.Success("mindcast", "indexer.startup", nil)

	for {
		<-background_channel
		search.SchedulePodcastIndexing()
		search.ScheduleEpisodeIndexing()
	}
}

func shutdown() {
	logger.Log("mindcast.indexer.shutdown")
	metrics.Success("mindcast", "indexer.shutdown", nil)

	datastore.Shutdown()
	metrics.Shutdown()
}
