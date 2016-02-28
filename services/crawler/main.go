package main

import (
	"github.com/nats-io/nats"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/mindcastio/mindcastio/backend"
	"github.com/mindcastio/mindcastio/backend/datastore"
	"github.com/mindcastio/mindcastio/backend/environment"
	"github.com/mindcastio/mindcastio/backend/logger"
	"github.com/mindcastio/mindcastio/backend/messaging"
	"github.com/mindcastio/mindcastio/backend/metrics"
)

func main() {

	// environment setup
	env := environment.GetEnvironment()

	logger.Initialize()
	metrics.Initialize(env)
	datastore.Initialize(env)
	messaging.Initialize(env)

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

	// subscribe to crawler queue
	messaging.QueueSubscribe(backend.Q_CRAWLER_REQUEST, backend.CRAWLER_QUEUE, func(msg *nats.Msg) {
		podcastCrawl(string(msg.Data))
	})

	for {
		<-background_channel
		schedulePodcastCrawling()
	}
}

func schedulePodcastCrawling() {
	logger.Log("mindcast.crawler.schedule_podcast_crawling")

	// search for podcasts that are candidates for crawling
	expired := backend.SearchExpiredPodcasts(backend.DEFAULT_UPDATE_BATCH)
	count := len(expired)

	logger.Log("mindcast.crawler.schedule_podcast_crawling.scheduling", strconv.FormatInt((int64)(count), 10))

	if count > 0 {
		for i := 0; i < count; i++ {
			messaging.Send(backend.Q_CRAWLER_REQUEST, expired[i].Uid)
		}
		metrics.Count("crawler.scheduled", count)
	}

	logger.Log("mindcast.crawler.schedule_podcast_crawling.done")
}

func podcastCrawl(uid string) {
	logger.Log("mindcast.crawler.podcast_crawl", uid)

	// crawl the the feed
	go backend.CrawlPodcastFeed(uid)

	metrics.Count("crawler.crawled", 1)
	logger.Log("mindcast.crawler.podcast_crawl.done", uid)
}

func shutdown() {
	logger.Log("mindcast.crawler.shutdown")
	metrics.Success("mindcast", "crawler.shutdown", nil)

	messaging.Shutdown()
	datastore.Shutdown()
	metrics.Shutdown()
}
