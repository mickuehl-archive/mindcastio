package search

import (
	"github.com/mindcastio/mindcastio/backend"
	"github.com/mindcastio/mindcastio/backend/util"
)

type ResultSorter []*Result

func (r ResultSorter) Len() int           { return len(r) }
func (r ResultSorter) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r ResultSorter) Less(i, j int) bool { return r[i].Score > r[j].Score }

func Search(q string, page int, size int) *SearchResult {
	var result []*Result
	var ll int

	uuid, _ := util.UUID()

	// log the search string first
	backend.LogSearchString(q)

	result1, _ := searchElastic(q, page, size)
	if result1.Count < size {
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
			result = make([]*Result, len(result1.Results)+len(result2))
			for i := range result1.Results {
				result[i] = result1.Results[i]
			}
			l := len(result1.Results)
			for i := range result2 {
				result[i+l] = result2[i]
			}
			ll = len(result)
		} else {
			result = result2
			ll = len(result2)
		}
	} else {
		result = result1.Results
		ll = result1.Count
	}

	// sort the results first

	return &SearchResult{uuid, ll, q, result}

}
