package search

import (
	"strconv"
	"time"

	"github.com/mindcastio/mindcastio/backend"
	"github.com/mindcastio/mindcastio/backend/logger"
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
	externalCount := 0

	// log the search string first
	backend.LogSearchString(q)

	// search our own index
	result, _ := searchElastic(q, page, limit)

	// trigger external search in iTunes if there is not enough in our own index ...
	if result.Count < MIN_RESULTS {

		result2, _ := searchITunes(q)

		// send feeds to the crawler
		feeds := make([]string, len(result2))
		for i := range result2 {
			feeds[i] = result2[i].Feed
		}
		backend.BulkSubmitPodcastFeed(feeds)
		externalCount = len(result2)
	}

	logger.Log("backend.search.hits", q, strconv.FormatInt((int64)(result.Count), 10), strconv.FormatInt((int64)(externalCount), 10))

	return &SearchResult{uuid, result.Count, q, util.ElapsedTimeSince(start), result.Results}

}
