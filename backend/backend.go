package backend

import (
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"time"

	"github.com/mindcastio/podcast-feed"

	"github.com/mindcastio/mindcastio/backend/datastore"
	"github.com/mindcastio/mindcastio/backend/logger"
	"github.com/mindcastio/mindcastio/backend/messaging"
	"github.com/mindcastio/mindcastio/backend/metrics"
	"github.com/mindcastio/mindcastio/backend/util"
)

func PodcastLookup(uid string) *PodcastMetadata {
	// TODO implement caching
	return podcastLookup(uid)
}

func SubmitPodcastFeed(feed string) error {

	logger.Log("submit_podcast_feed", feed)

	// check if the podcast is already in the index
	uid := util.UID(feed)
	idx := indexLookup(uid)

	if idx == nil {
		err := indexAdd(uid, feed)
		if err != nil {

			logger.Error("submit_podcast_feed.error", err, feed)
			metrics.Error("submit_podcast_feed.error", err.Error(), []string{feed})

			return err
		} else {
			// initiate immediate crawling
			messaging.Send(Q_CRAWLER_REQUEST, uid)
		}
	} else {
		logger.Warn("submit_podcast_feed.duplicate", uid, feed)
		metrics.Warning("submit_podcast_feed.duplicate", "", []string{feed})
	}

	logger.Log("submit_podcast_feed.done", uid, feed)
	return nil
}

func BulkSubmitPodcastFeed(urls []string) error {

	logger.Log("bulk_submit_podcast_feed")

	count := 0
	feed := ""

	for i := 0; i < len(urls); i++ {
		feed = urls[i]

		// check if the podcast is already in the index
		uid := util.UID(feed)
		idx := indexLookup(uid)

		if idx == nil {
			err := indexAdd(uid, feed)
			if err != nil {

				logger.Error("bulk_submit_podcast_feed.error", err, feed)
				metrics.Error("bulk_submit_podcast_feed.error", err.Error(), []string{feed})

				return err
			} else {
				// initiate immediate crawling
				messaging.Send(Q_CRAWLER_REQUEST, uid)

				count++
			}
		}
	}

	logger.Log("bulk_submit_podcast_feed.done", strconv.FormatInt((int64)(count), 10))
	return nil
}

func SearchExpiredPodcasts(limit int) []PodcastIndex {

	ds := datastore.GetDataStore()
	defer ds.Close()

	main_index := ds.Collection(datastore.META_COL)

	results := []PodcastIndex{}
	q := bson.M{"next": bson.M{"$lte": util.Timestamp()}, "errors": bson.M{"$lte": MAX_ERRORS}}

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

	idx := indexLookup(uid)
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
		suspended, _ := indexBackoff(uid)

		if suspended {
			logger.Error("crawl_podcast_feed.suspended", err, uid, idx.Feed)
			metrics.Error("crawl_podcast_feed.suspended", err.Error(), []string{uid, idx.Feed})
		}

		return
	}

	// add to podcast metadata index
	is_new, err := podcastAdd(podcast)
	if err != nil {
		logger.Error("crawl_podcast_feed.error.3", err, uid, idx.Feed)
		metrics.Error("crawl_podcast_feed.error", err.Error(), []string{uid, idx.Feed})

		return
	}

	// add to the episodes metadata index
	count, err := episodesAddAll(podcast)
	if err != nil {
		logger.Error("crawl_podcast_feed.error.4", err, uid, idx.Feed)
		metrics.Error("crawl_podcast_feed.error", err.Error(), []string{uid, idx.Feed})

		return
	} else {
		// update main metadata index
		indexUpdate(uid)

		if count > 0 {
			// update stats and metrics
			if is_new {
				metrics.Count("index.podcast.new", 1)
				metrics.Count("index.episodes.new", count)
			} else {
				// new episodes added -> update the podcast.published timestamp
				podcastUpdateTimestamp(podcast)

				metrics.Count("index.episodes.update", count)
			}
		}

		logger.Log("crawl_podcast_feed.done", uid, idx.Feed, strconv.FormatInt((int64)(count), 10))
		metrics.Histogram("crawler.crawl", (float64)(util.ElapsedTimeSince(start_1)))
	}
}

func SchedulePodcastIndexing() {

	logger.Log("schedule_podcast_indexing")

	// search for podcasts that are candidates for indexing
	notIndexed := podcastSearchNotIndexd(DEFAULT_INDEX_UPDATE_BATCH, SEARCH_REVISION)
	count := len(notIndexed)

	logger.Log("schedule_podcast_indexing.scheduling", strconv.FormatInt((int64)(count), 10))

	if count > 0 {
		ds := datastore.GetDataStore()
		defer ds.Close()

		podcast_metadata := ds.Collection(datastore.PODCASTS_COL)

		for i := 0; i < count; i++ {
			err := podcastAddToIndex(&notIndexed[i])
			if err != nil {
				logger.Error("schedule_podcast_indexing.error.1", err, notIndexed[i].Uid)
				metrics.Error("schedule_podcast_indexing.error.1", err.Error(), []string{notIndexed[i].Uid})
				// abort or disable at some point?
			}

			// update the metadata
			notIndexed[i].IndexVersion = SEARCH_REVISION
			notIndexed[i].Updated = util.Timestamp()
			err = podcast_metadata.Update(bson.M{"uid": notIndexed[i].Uid}, &notIndexed[i])
			if err != nil {
				logger.Error("schedule_podcast_indexing.error.2", err, notIndexed[i].Uid)
				metrics.Error("schedule_podcast_indexing.error.2", err.Error(), []string{notIndexed[i].Uid})
				// abort or disable at some point?
			}
		}
		metrics.Count("indexer.podcasts.scheduled", count)
	}

	logger.Log("schedule_podcast_indexing.done")
}

func ScheduleEpisodeIndexing() {

	logger.Log("schedule_episode_indexing")

	// search for podcasts that are candidates for indexing
	notIndexed := episodesSearchNotIndexd(DEFAULT_INDEX_UPDATE_BATCH, SEARCH_REVISION)
	count := len(notIndexed)

	logger.Log("schedule_episode_indexing.scheduling", strconv.FormatInt((int64)(count), 10))

	if count > 0 {
		ds := datastore.GetDataStore()
		defer ds.Close()

		episodes_metadata := ds.Collection(datastore.EPISODES_COL)

		for i := 0; i < count; i++ {
			err := episodeAddToIndex(&notIndexed[i])
			if err != nil {
				logger.Error("schedule_episode_indexing.error.1", err, notIndexed[i].Uid)
				metrics.Error("schedule_episode_indexing.error.1", err.Error(), []string{notIndexed[i].Uid})
				// abort or disable at some point?
			}

			// update the metadata
			notIndexed[i].IndexVersion = SEARCH_REVISION
			notIndexed[i].Updated = util.Timestamp()
			err = episodes_metadata.Update(bson.M{"uid": notIndexed[i].Uid}, &notIndexed[i])
			if err != nil {
				logger.Error("schedule_episode_indexing.error.2", err, notIndexed[i].Uid)
				metrics.Error("schedule_episode_indexing.error.2", err.Error(), []string{notIndexed[i].Uid})
				// abort or disable at some point?
			}
		}
		metrics.Count("indexer.episodes.scheduled", count)
	}

	logger.Log("schedule_episode_indexing.done")
}
