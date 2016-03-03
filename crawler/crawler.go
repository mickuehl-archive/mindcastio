package crawler

import (
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"time"

	"github.com/mindcastio/podcast-feed"

	"github.com/mindcastio/mindcastio/backend"

	"github.com/mindcastio/mindcastio/backend/datastore"
	"github.com/mindcastio/mindcastio/backend/logger"
	"github.com/mindcastio/mindcastio/backend/metrics"
	"github.com/mindcastio/mindcastio/backend/util"
)

func SearchExpiredPodcasts(limit int) []backend.PodcastIndex {

	ds := datastore.GetDataStore()
	defer ds.Close()

	main_index := ds.Collection(datastore.META_COL)

	results := []backend.PodcastIndex{}
	q := bson.M{"next": bson.M{"$lte": util.Timestamp()}, "errors": bson.M{"$lte": backend.MAX_ERRORS}}

	if limit <= 0 {
		// return all
		main_index.Find(q).All(&results)
	} else {
		// with a limit
		main_index.Find(q).Limit(limit).All(&results)
	}

	return results
}

func CrawlPodcastFeed(uid string) {

	start_1 := time.Now()
	logger.Log("crawl_podcast_feed", uid)

	idx := backend.IndexLookup(uid)
	if idx == nil {
		logger.Error("crawl_podcast_feed.error.1", nil, uid)
		metrics.Error("crawl_podcast_feed.error", "", []string{uid})
		return
	}

	// HINT: ignore the fact that the item might be disables idx.erros > ...

	// fetch the podcast feed
	start_2 := time.Now()
	podcast, err := podcast.ParsePodcastFeed(idx.Feed)
	metrics.Histogram("crawler.parse_feed", (float64)(util.ElapsedTimeSince(start_2)))

	if err != nil {
		suspended, _ := backend.IndexBackoff(uid)

		if suspended {
			logger.Error("crawl_podcast_feed.suspended", err, uid, idx.Feed)
			metrics.Error("crawl_podcast_feed.suspended", err.Error(), []string{uid, idx.Feed})
		}

		return
	}

	// add to podcast metadata index
	is_new, err := backend.PodcastAdd(podcast)
	if err != nil {
		logger.Error("crawl_podcast_feed.error.3", err, uid, idx.Feed)
		metrics.Error("crawl_podcast_feed.error", err.Error(), []string{uid, idx.Feed})

		return
	}

	// add to the episodes metadata index
	count, err := backend.EpisodesAddAll(podcast)
	if err != nil {
		logger.Error("crawl_podcast_feed.error.4", err, uid, idx.Feed)
		metrics.Error("crawl_podcast_feed.error", err.Error(), []string{uid, idx.Feed})

		return
	} else {
		// update main metadata index
		backend.IndexUpdate(uid)

		if count > 0 {
			// update stats and metrics
			if is_new {
				metrics.Count("index.podcast.new", 1)
				metrics.Count("index.episodes.new", count)
			} else {
				// new episodes added -> update the podcast.published timestamp
				backend.PodcastUpdateTimestamp(podcast)

				metrics.Count("index.episodes.update", count)
			}
		}

		logger.Log("crawl_podcast_feed.done", uid, idx.Feed, strconv.FormatInt((int64)(count), 10))
		metrics.Histogram("crawler.crawl", (float64)(util.ElapsedTimeSince(start_1)))
	}
}
