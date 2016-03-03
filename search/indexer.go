package search

import (
	"strings"
	"gopkg.in/mgo.v2/bson"
	"strconv"

	"github.com/franela/goreq"

	"github.com/mindcastio/mindcastio/backend"
	"github.com/mindcastio/mindcastio/backend/services/datastore"
	"github.com/mindcastio/mindcastio/backend/services/logger"
	"github.com/mindcastio/mindcastio/backend/services/metrics"
	"github.com/mindcastio/mindcastio/backend/environment"
	"github.com/mindcastio/mindcastio/backend/util"
)

func PodcastAddToSearchIndex(podcast *backend.PodcastMetadata) error {

	uri := strings.Join([]string{environment.GetEnvironment().SearchServiceUrl(), "/search/podcast/", podcast.Uid}, "")
	payload := podcastMetadataToSearch(podcast)

	// post the payload to elasticsearch
	res, err := goreq.Request{
		Method:      "PUT",
		Uri:         uri,
		ContentType: "application/json",
		Body:        payload,
	}.Do()

	if res != nil {
		res.Body.Close()
	}
	return err
}

func EpisodeAddToSearchIndex(episode *backend.EpisodeMetadata) error {

	uri := strings.Join([]string{environment.GetEnvironment().SearchServiceUrl(), "/search/episode/", episode.Uid}, "")
	payload := episodeMetadataToSearch(episode)

	// post the payload to elasticsearch
	res, err := goreq.Request{
		Method:      "PUT",
		Uri:         uri,
		ContentType: "application/json",
		Body:        payload,
	}.Do()

	if res != nil {
		res.Body.Close()
	}
	return err
}

func SchedulePodcastIndexing() {

	logger.Log("schedule_podcast_indexing")

	// search for podcasts that are candidates for indexing
	notIndexed := podcastSearchNotIndexd(backend.DEFAULT_INDEX_UPDATE_BATCH, backend.SEARCH_REVISION)
	count := len(notIndexed)

	logger.Log("schedule_podcast_indexing.scheduling", strconv.FormatInt((int64)(count), 10))

	if count > 0 {
		ds := datastore.GetDataStore()
		defer ds.Close()

		podcast_metadata := ds.Collection(datastore.PODCASTS_COL)

		for i := 0; i < count; i++ {
			err := PodcastAddToSearchIndex(&notIndexed[i])
			if err != nil {
				logger.Error("schedule_podcast_indexing.error.1", err, notIndexed[i].Uid)
				metrics.Error("schedule_podcast_indexing.error.1", err.Error(), []string{notIndexed[i].Uid})
				// abort or disable at some point?
			}

			// update the metadata
			notIndexed[i].IndexVersion = backend.SEARCH_REVISION
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
	notIndexed := episodesSearchNotIndexd(backend.DEFAULT_INDEX_UPDATE_BATCH, backend.SEARCH_REVISION)
	count := len(notIndexed)

	logger.Log("schedule_episode_indexing.scheduling", strconv.FormatInt((int64)(count), 10))

	if count > 0 {
		ds := datastore.GetDataStore()
		defer ds.Close()

		episodes_metadata := ds.Collection(datastore.EPISODES_COL)

		for i := 0; i < count; i++ {
			err := EpisodeAddToSearchIndex(&notIndexed[i])
			if err != nil {
				logger.Error("schedule_episode_indexing.error.1", err, notIndexed[i].Uid)
				metrics.Error("schedule_episode_indexing.error.1", err.Error(), []string{notIndexed[i].Uid})
				// abort or disable at some point?
			}

			// update the metadata
			notIndexed[i].IndexVersion = backend.SEARCH_REVISION
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

func podcastSearchNotIndexd(limit int, version int) []backend.PodcastMetadata {

	ds := datastore.GetDataStore()
	defer ds.Close()

	podcast_metadata := ds.Collection(datastore.PODCASTS_COL)

	results := []backend.PodcastMetadata{}
	q := bson.M{"indexversion": bson.M{"$lt": version}}

	if limit <= 0 {
		// return all
		podcast_metadata.Find(q).All(&results)
	} else {
		// with a limit
		podcast_metadata.Find(q).Limit(limit).All(&results)
	}

	return results
}

func episodesSearchNotIndexd(limit int, version int) []backend.EpisodeMetadata {

	ds := datastore.GetDataStore()
	defer ds.Close()

	episodes_metadata := ds.Collection(datastore.EPISODES_COL)

	results := []backend.EpisodeMetadata{}
	q := bson.M{"indexversion": bson.M{"$lt": version}}

	if limit <= 0 {
		// return all
		episodes_metadata.Find(q).All(&results)
	} else {
		// with a limit
		episodes_metadata.Find(q).Limit(limit).All(&results)
	}

	return results
}

func podcastMetadataToSearch(p *backend.PodcastMetadata) PodcastMetadataSearch {
	return PodcastMetadataSearch{
		p.Uid,
		p.Title,
		p.Subtitle,
		p.Description,
		p.Published,
		p.Language,
		p.OwnerName,
		p.OwnerEmail,
		p.Tags,
	}
}

func episodeMetadataToSearch(e *backend.EpisodeMetadata) EpisodeMetadataSearch {
	return EpisodeMetadataSearch{
		e.Uid,
		e.Title,
		e.Url,
		e.Description,
		e.Published,
		e.Author,
		e.PodcastUid,
	}
}
