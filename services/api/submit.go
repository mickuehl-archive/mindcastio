package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/ant0ine/go-json-rest/rest"

	"github.com/mindcastio/mindcastio/backend"
	"github.com/mindcastio/mindcastio/backend/metrics"

	"github.com/mindcastio/mindcastio/backend/util"
)

type feedType struct {
	Feed string
}

func submit_endpoint(w rest.ResponseWriter, r *rest.Request) {
	start := time.Now()

	ft := feedType{}

	err := r.DecodeJsonPayload(&ft)
	if err != nil {
		backend.JsonApiErrorResponse(w, "api.submit.error", "missing parameter", err)

		metrics.Error("api.submit.error", err.Error(), nil)
		metrics.Count("api.total.count", 1)
		metrics.Count("api.submit.count", 1)

		return
	}

	// remove leading & trailing whitespace
	feed := strings.TrimSpace(ft.Feed)

	// test if the feed can be crawled, otherwise discard it already
	err = util.ValidateUrl(feed)
	if err != nil {
		backend.JsonApiErrorResponse(w, "api.submit.error", "invalid url", err)
		metrics.Error("api.submit.error", err.Error(), nil)
	} else {
		// response
		err = backend.SubmitPodcastFeed(feed)
		backend.StatusResponse(w, http.StatusOK)
	}

	// metrics
	metrics.Count("api.total.count", 1)
	metrics.Count("api.submit.count", 1)
	metrics.Histogram("api.submit.duration", (float64)(util.ElapsedTimeSince(start)))
}
