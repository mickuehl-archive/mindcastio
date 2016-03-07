package search

import (
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"strings"
	"time"

	"github.com/mindcastio/mindcastio/backend"
	"github.com/mindcastio/mindcastio/backend/datastore"
	"github.com/mindcastio/mindcastio/backend/environment"
	"github.com/mindcastio/mindcastio/backend/logger"
	"github.com/mindcastio/mindcastio/backend/metrics"
	"github.com/mindcastio/mindcastio/backend/util"
)

type (

	PodcastSearchMetadata struct {
		Uid         string `json:"uid"`
		Title       string `json:"title"`
		Subtitle    string `json:"subtitle"`
		Description string `json:"description"`
		Language    string `json:"language"`
		OwnerName   string `json:"owner_name"`
		OwnerEmail  string `json:"owner_email"`
	}

)

func SchedulePodcastIndexing() {

	start := time.Now()
	logger.Log("schedule_podcast_indexing")

	// search for podcasts that are candidates for indexing
	notIndexed := podcastSearchNotIndexed(backend.DEFAULT_INDEX_UPDATE_BATCH, backend.SEARCH_REVISION)
	count := len(notIndexed)

	logger.Log("schedule_podcast_indexing.scheduling", strconv.FormatInt((int64)(count), 10))

	if count > 0 {
		ds := datastore.GetDataStore()
		defer ds.Close()

		podcast_metadata := ds.Collection(datastore.PODCASTS_COL)

		for i := 0; i < count; i++ {
			err := podcastAddToSearchIndex(&notIndexed[i])
			if err != nil {
				logger.Error("schedule_podcast_indexing.error.1", err, notIndexed[i].Uid)
				metrics.Error("schedule_podcast_indexing.error.1", err.Error(), []string{notIndexed[i].Uid})
				// abort or disable at some point?
			}

			// update the metadata
			notIndexed[i].Version = backend.SEARCH_REVISION
			notIndexed[i].Updated = util.Timestamp()
			err = podcast_metadata.Update(bson.M{"uid": notIndexed[i].Uid}, &notIndexed[i])
			if err != nil {
				logger.Error("schedule_podcast_indexing.error.2", err, notIndexed[i].Uid)
				metrics.Error("schedule_podcast_indexing.error.2", err.Error(), []string{notIndexed[i].Uid})
				// abort or disable at some point?
			}
		}
		metrics.Count("indexer.podcasts.new", count)
	}

	logger.Log("schedule_podcast_indexing.done")
	metrics.Histogram("indexer.schedule_podcast_indexing.duration", (float64)(util.ElapsedTimeSince(start)))
}

func ScheduleEpisodeIndexing() {

	start := time.Now()
	logger.Log("schedule_episode_indexing")

	// search for podcasts that are candidates for indexing
	notIndexed := episodesSearchNotIndexed(backend.DEFAULT_INDEX_UPDATE_BATCH, backend.SEARCH_REVISION)
	count := len(notIndexed)

	logger.Log("schedule_episode_indexing.scheduling", strconv.FormatInt((int64)(count), 10))

	if count > 0 {
		ds := datastore.GetDataStore()
		defer ds.Close()

		episodes_metadata := ds.Collection(datastore.EPISODES_COL)

		for i := 0; i < count; i++ {
			err := episodeAddToSearchIndex(&notIndexed[i])
			if err != nil {
				logger.Error("schedule_episode_indexing.error.1", err, notIndexed[i].Uid)
				metrics.Error("schedule_episode_indexing.error.1", err.Error(), []string{notIndexed[i].Uid})
				// abort or disable at some point?
			}

			// update the metadata
			notIndexed[i].Version = backend.SEARCH_REVISION
			notIndexed[i].Updated = util.Timestamp()
			err = episodes_metadata.Update(bson.M{"uid": notIndexed[i].Uid}, &notIndexed[i])
			if err != nil {
				logger.Error("schedule_episode_indexing.error.2", err, notIndexed[i].Uid)
				metrics.Error("schedule_episode_indexing.error.2", err.Error(), []string{notIndexed[i].Uid})
				// abort or disable at some point?
			}
		}
		metrics.Count("indexer.episodes.new", count)
	}

	logger.Log("schedule_episode_indexing.done")
	metrics.Histogram("indexer.schedule_episode_indexing.duration", (float64)(util.ElapsedTimeSince(start)))
}

func podcastAddToSearchIndex(podcast *backend.PodcastMetadata) error {

	uri := strings.Join([]string{environment.GetEnvironment().SearchServiceUrl(), "/podcasts/podcast/", podcast.Uid}, "")

	payload := PodcastSearchMetadata{
		podcast.Uid,
		podcast.Title,
		podcast.Subtitle,
		podcast.Description,
		podcast.Language,
		podcast.OwnerName,
		podcast.OwnerEmail,
	}

	return util.PutJson(uri, payload)

}

func episodeAddToSearchIndex(episode *backend.EpisodeMetadata) error {

	podcast := backend.PodcastLookup(episode.PodcastUid)
	// FIXME we simply assume no errors here !!

	id := strings.Join([]string{episode.PodcastUid, episode.Uid}, "-")
	uri := strings.Join([]string{environment.GetEnvironment().SearchServiceUrl(), "/podcasts/episode/", id}, "")

	payload := PodcastSearchMetadata{
		episode.Uid,
		episode.Title,
		"",
		episode.Description,
		podcast.Language,
		episode.Author,
		"",
	}

	return util.PutJson(uri, payload)

}

func podcastSearchNotIndexed(limit int, version int) []backend.PodcastMetadata {

	ds := datastore.GetDataStore()
	defer ds.Close()

	podcast_metadata := ds.Collection(datastore.PODCASTS_COL)

	results := []backend.PodcastMetadata{}
	q := bson.M{"version": bson.M{"$lt": version}}

	if limit <= 0 {
		// return all
		podcast_metadata.Find(q).All(&results)
	} else {
		// with a limit
		podcast_metadata.Find(q).Limit(limit).All(&results)
	}

	return results
}

func episodesSearchNotIndexed(limit int, version int) []backend.EpisodeMetadata {

	ds := datastore.GetDataStore()
	defer ds.Close()

	episodes_metadata := ds.Collection(datastore.EPISODES_COL)

	results := []backend.EpisodeMetadata{}
	q := bson.M{"version": bson.M{"$lt": version}}

	if limit <= 0 {
		// return all
		episodes_metadata.Find(q).All(&results)
	} else {
		// with a limit
		episodes_metadata.Find(q).Limit(limit).All(&results)
	}

	return results
}
