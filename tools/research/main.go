package main

import (

	"fmt"
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

	for i := range results {
		fmt.Println("Keyword: ", results[i].Word, results[i].Frequency)
		time.Sleep(1000 * time.Millisecond)

		// search iTunes and submit the results to the crawler
		itunes_result, _ := search.SearchITunes(results[i].Word)

		// send feeds to the crawler
		feeds := make([]string, len(itunes_result))
		for i := range itunes_result {
			feeds[i] = itunes_result[i].Feed
		}
		backend.BulkSubmitPodcastFeed(feeds)

	}

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
