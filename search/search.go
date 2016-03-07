package search

import (
	"github.com/mindcastio/mindcastio/backend"
	"github.com/mindcastio/mindcastio/backend/util"
)

const (
	MIN_RESULTS int = 15
)

func Search(q string) *SearchResult {
	uuid, _ := util.UUID()

	// log the search string first
	backend.LogSearchString(q)

	result1, _ := searchElastic(q)
	if len(result1.Results) < MIN_RESULTS {
		// search externally
		result2, _ := searchITunes(q)

		// send feeds to the crawler
		feeds := make([]string, len(result2))
		for i := range result2 {
			feeds[i] = result2[i].Feed
		}
		backend.BulkSubmitPodcastFeed(feeds)

		// combine both results
		if len(result1.Results) > 0 {
			result := make([]*Result, len(result1.Results) + len(result2))
			for i := range result1.Results {
				result[i] = result1.Results[i]
			}
			l := len(result1.Results)
			for i := range result2 {
				result[i + l] = result2[i]
			}

			return &SearchResult{uuid, result1.Count + len(result2), q, result}
		} else {
			return &SearchResult{uuid, len(result2), q, result2}
		}
	} else {
		return &SearchResult{uuid, result1.Count, q, result1.Results}
	}

}
