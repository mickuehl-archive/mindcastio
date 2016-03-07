package crawler

import (
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"time"

	"github.com/mindcastio/mindcastio/backend"

	"github.com/mindcastio/mindcastio/backend/datastore"
	"github.com/mindcastio/mindcastio/backend/logger"
	"github.com/mindcastio/mindcastio/backend/metrics"
	"github.com/mindcastio/mindcastio/backend/util"
)

func SchedulePodcastCrawling() {
	logger.Log("mindcast.crawler.schedule_podcast_crawling")

	// search for podcasts that are candidates for crawling
	expired := searchExpiredPodcasts(backend.DEFAULT_UPDATE_BATCH)
	count := len(expired)

	logger.Log("crawler.schedule_podcast_crawling.scheduling", strconv.FormatInt((int64)(count), 10))

	if count > 0 {
		for i := 0; i < count; i++ {
			//go CrawlPodcastFeed(expired[i].Uid)
			CrawlPodcastFeed(expired[i].Uid)
			// FIXME use go routine or not ?
		}
		metrics.Count("crawler.scheduled", count)
	}

	logger.Log("crawler.schedule_podcast_crawling.done")
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
	podcast, err := ParsePodcastFeed(idx.Feed)
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
		backend.IndexUpdate(uid)

		if count > 0 {
			// update stats and metrics
			if is_new {
				metrics.Count("crawler.podcast.new", 1)
				metrics.Count("crawler.episodes.new", count)
			} else {
				// new episodes added -> update the podcast.published timestamp
				podcastUpdateTimestamp(podcast)
				metrics.Count("crawler.episodes.update", count)
			}
		}

		logger.Log("crawl_podcast_feed.done", uid, idx.Feed, strconv.FormatInt((int64)(count), 10))
		metrics.Histogram("crawler.podcast_feed.duration", (float64)(util.ElapsedTimeSince(start_1)))
	}
}

func podcastAdd(podcast *Podcast) (bool, error) {
	p := backend.PodcastLookup(podcast.Uid)
	if p != nil {
		return false, nil
	}

	ds := datastore.GetDataStore()
	defer ds.Close()

	podcast_metadata := ds.Collection(datastore.PODCASTS_COL)

	meta := podcastDetailsToMetadata(podcast)

	// fix the published timestamp
	now := util.Timestamp()
	if podcast.Published > now {
		meta.Published = now // prevents dates in the future
	}

	err := podcast_metadata.Insert(&meta)

	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func podcastUpdateTimestamp(podcast *Podcast) (bool, error) {

	ds := datastore.GetDataStore()
	defer ds.Close()

	podcast_metadata := ds.Collection(datastore.PODCASTS_COL)

	p := backend.PodcastMetadata{}
	podcast_metadata.Find(bson.M{"uid": podcast.Uid}).One(&p)

	if p.Uid == "" {
		return false, nil
	} else {
		now := util.Timestamp()
		p.Updated = now
		if podcast.Published > now {
			p.Published = now // prevents dates in the future
		} else {
			p.Published = podcast.Published
		}

		// update the DB
		err := podcast_metadata.Update(bson.M{"uid": podcast.Uid}, &p)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func episodeAdd(episode *Episode, puid string) (bool, error) {
	e := backend.EpisodeLookup(episode.Uid)
	if e != nil {
		return false, nil
	}

	ds := datastore.GetDataStore()
	defer ds.Close()

	episodes_metadata := ds.Collection(datastore.EPISODES_COL)

	meta := episodeDetailsToMetadata(episode, puid)
	err := episodes_metadata.Insert(&meta)

	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func episodesAddAll(podcast *Podcast) (int, error) {
	count := 0

	for i := 0; i < len(podcast.Episodes); i++ {
		added, err := episodeAdd(&podcast.Episodes[i], podcast.Uid)
		if err != nil {
			return 0, err
		}
		if added {
			count++
		}
	}
	return count, nil
}

func searchExpiredPodcasts(limit int) []backend.PodcastIndex {

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

func podcastDetailsToMetadata(podcast *Podcast) *backend.PodcastMetadata {
	meta := backend.PodcastMetadata{
		podcast.Uid,
		podcast.Title,
		podcast.Subtitle,
		podcast.Url,
		podcast.Feed,
		podcast.Description,
		podcast.Published,
		podcast.Language,
		podcast.Image,
		podcast.Owner.Name,
		podcast.Owner.Email,
		"",
		0,
		0,
		0,
		0,
		util.Timestamp(),
		0,
	}
	return &meta
}

func episodeDetailsToMetadata(episode *Episode, puid string) *backend.EpisodeMetadata {
	meta := backend.EpisodeMetadata{
		episode.Uid,
		episode.Title,
		episode.Url,
		episode.Description,
		episode.Published,
		episode.Duration,
		episode.Author,
		episode.Content.Url,
		episode.Content.Type,
		episode.Content.Size,
		puid,
		0,
		util.Timestamp(),
		0,
	}
	return &meta
}
