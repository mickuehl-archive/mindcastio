package main

import (
	"strconv"
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
	defer metrics.Shutdown()
	datastore.Initialize(env)
	defer datastore.Shutdown()

	ds := datastore.GetDataStore()
	defer ds.Close()

	search_keywords := ds.Collection(datastore.KEYWORDS_COL)

	results := []backend.SearchKeyword{}
	search_keywords.Find(nil).Sort("-frequency").All(&results)

	var total int = 0

	for i := range results {

		time.Sleep(3000 * time.Millisecond)

		// search iTunes and submit the results to the crawler
		itunes_result, _ := search.SearchITunes(results[i].Word)

		// send feeds to the crawler
		feeds := make([]string, len(itunes_result))
		for i := range itunes_result {
			feeds[i] = itunes_result[i].Feed
		}

		count, _ := backend.BulkSubmitPodcastFeed(feeds)
		total = total + count

		logger.Log("re_search.search", results[i].Word, strconv.FormatInt((int64)(count), 10))
	}

	logger.Log("re_search.done", strconv.FormatInt((int64)(len(results)), 10), strconv.FormatInt((int64)(total), 10))

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
