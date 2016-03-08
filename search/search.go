package search

import (
	"github.com/mindcastio/mindcastio/backend"
	"github.com/mindcastio/mindcastio/backend/util"
)

type ResultSorter []*Result

func (r ResultSorter) Len() int           { return len(r) }
func (r ResultSorter) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r ResultSorter) Less(i, j int) bool { return r[i].Score > r[j].Score }

func Search(q string, page int, limit int) *SearchResult {
	var result []*Result
	var ll int

	uuid, _ := util.UUID()

	// log the search string first
	backend.LogSearchString(q)

	result1, _ := searchElastic(q, page, limit)
	if result1.Count < limit {
		// search externally
		result2, _ := searchITunes(q)

		// send feeds to the crawler
		feeds := make([]string, len(result2))
		for i := range result2 {
			feeds[i] = result2[i].Feed
		}
		backend.BulkSubmitPodcastFeed(feeds)

		// return either internal result or a subset from the external search
		if len(result1.Results) > 0 {
			// just return what we alreday got ...
			result = result1.Results
			ll = result1.Count
		} else {
			// limit the result set ...
			if len(result2) > limit {
				result = make([]*Result, limit)
				for i := 0; i < limit; i++ {
						result[i] = result2[i]
				}
				ll = len(result)
			} else {
				result = result2
				ll = len(result2)
			}
		}
	} else {
		result = result1.Results
		ll = result1.Count
	}

	return &SearchResult{uuid, ll, q, result}

}
