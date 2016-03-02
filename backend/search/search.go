package search

import (

	"github.com/mindcastio/mindcastio/backend"
	"github.com/mindcastio/mindcastio/backend/util"
)

func Search(q string) *SearchResult {
	uuid, _ := util.UUID()

	result, _ := searchElastic(q)
	if len(result) == 0 {
		// search externally
		result, _ = searchITunes(q)

		// send feeds to the crawler
		feeds := make([]string, len(result))
		for i := range result {
			feeds[i] = result[i].Feed
		}
		backend.BulkSubmitPodcastFeed(feeds)
	}

	return &SearchResult{uuid, len(result), q, result}

}
