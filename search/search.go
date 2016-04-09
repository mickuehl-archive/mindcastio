package search

import (
	"time"

	"github.com/mindcastio/mindcastio/backend"
	"github.com/mindcastio/mindcastio/backend/metrics"
	"github.com/mindcastio/mindcastio/backend/util"
)

// TODO keep this for now, we may need it for sorting of results at some point
type ResultSorter []*Result

func (r ResultSorter) Len() int           { return len(r) }
func (r ResultSorter) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r ResultSorter) Less(i, j int) bool { return r[i].Score > r[j].Score }

func Search(q string, page int, limit int) *SearchResult {

	start := time.Now()
	uuid, _ := util.UUID()

	// log the search string first
	backend.LogSearchString(q)

	// search our own index
	start_1 := time.Now()
	result, _ := SearchElastic(q, page, limit)
	metrics.Histogram("search.internal.duration", (float64)(util.ElapsedTimeSince(start_1)))

	// trigger external search in iTunes if there is not enough in our own index ...
	if result.Count < MIN_RESULTS {

		go func() {
			start_2 := time.Now()
			result2, _ := SearchITunes(q)
			metrics.Histogram("search.external.duration", (float64)(util.ElapsedTimeSince(start_2)))

			// send feeds to the crawler
			feeds := make([]string, len(result2))
			for i := range result2 {
				feeds[i] = result2[i].Feed
			}
			backend.BulkSubmitPodcastFeed(feeds)

			metrics.Count("search.external.count", len(result2))
		}()
	}

	metrics.Count("search.internal.count", result.Count)
	metrics.Histogram("search.duration", (float64)(util.ElapsedTimeSince(start)))

	return &SearchResult{uuid, result.Count, q, util.ElapsedTimeSince(start), result.Results}

}
