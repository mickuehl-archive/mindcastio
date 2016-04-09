package main

import (
	"strings"
	"time"

	"github.com/mindcastio/go-json-rest/rest"

	"github.com/mindcastio/mindcastio/backend"
	"github.com/mindcastio/mindcastio/backend/metrics"

	"github.com/mindcastio/mindcastio/backend/util"
)

// req.PathParam("host")
// "/lookup/#host"
// 2a38720c9b2d51bde2a1dcfa49eb1690

func podcast_endpoint(w rest.ResponseWriter, r *rest.Request) {
	start := time.Now()

	// get the id first
	uid := strings.Trim(r.PathParam("id"), " ")
	if uid == "" {
		backend.JsonApiErrorResponse(w, "api.podcast.error", "missing parameter", nil)
		metrics.Error("api.podcast.error", "", nil)
		return
	}

	result := backend.PodcastLookup(uid)
	if result == nil {
		backend.JsonApiErrorResponse(w, "api.podcast.error", "podcast not found", nil)

		metrics.Error("api.podcast.error", "podcast not found", []string{uid})

		metrics.Count("api.total.count", 1)
		metrics.Count("api.podcast.count", 1)

		return
	}

	// create an 'outside' view
	podcast := backend.Podcast{
		result.Uid,
		result.Title,
		result.Subtitle,
		result.Url,
		result.Feed,
		result.Description,
		result.Published,
		result.Language,
		result.ImageUrl,
		result.OwnerName,
		result.OwnerEmail,
	}
	backend.JsonApiResponse(w, &podcast)

	// metrics
	metrics.Count("api.total.count", 1)
	metrics.Count("api.podcast.count", 1)
	metrics.Histogram("api.podcast.duration", (float64)(util.ElapsedTimeSince(start)))
}

func episode_endpoint(w rest.ResponseWriter, r *rest.Request) {
	start := time.Now()

	// get the id first
	uid := strings.Trim(r.PathParam("id"), " ")
	if uid == "" {
		backend.JsonApiErrorResponse(w, "api.episode.error", "missing parameter", nil)
		metrics.Error("api.episode.error", "", nil)
		return
	}

	result := backend.EpisodeLookup(uid)
	if result == nil {
		backend.JsonApiErrorResponse(w, "api.episode.error", "episode not found", nil)

		metrics.Error("api.episode.error", "episode not found", []string{uid})

		metrics.Count("api.total.count", 1)
		metrics.Count("api.episode.count", 1)

		return
	}

	// create an 'outside' view
	episode := backend.Episode{
		result.Uid,
		result.PodcastUid,
		result.Title,
		result.Url,
		result.Description,
		result.Published,
		result.Duration,
		result.Author,
		result.AssetUrl,
		result.AssetType,
		result.AssetSize,
	}
	backend.JsonApiResponse(w, &episode)

	// metrics
	metrics.Count("api.total.count", 1)
	metrics.Count("api.episode.count", 1)
	metrics.Histogram("api.episode.duration", (float64)(util.ElapsedTimeSince(start)))
}
