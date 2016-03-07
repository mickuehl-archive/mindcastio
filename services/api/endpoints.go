package main

import (
	"strings"
	"time"

	"github.com/mindcastio/go-json-rest/rest"

	"github.com/mindcastio/mindcastio/search"

	"github.com/mindcastio/mindcastio/backend"
	"github.com/mindcastio/mindcastio/backend/logger"
	"github.com/mindcastio/mindcastio/backend/metrics"

	"github.com/mindcastio/mindcastio/backend/util"
)

func endpoint(w rest.ResponseWriter, r *rest.Request) {
	start := time.Now()

	if len(r.URL.Query()["q"]) == 0 {
		backend.JsonApiErrorResponse(w, "api.search.error", "missing parameter", nil)
		metrics.Error("api.search.error", "", nil)
		return
	}

	q := r.URL.Query()["q"][0]
	query := strings.Replace(q, " ", "+", -1)

	if len(query) == 0 {
		backend.JsonApiErrorResponse(w, "api.search.error", "missing query", nil)
		metrics.Error("api.search.error", "", nil)
		return
	}

	logger.Log("api.search.query", query)

	result := search.Search(query)
	backend.JsonApiResponse(w, result)

	// metrics
	metrics.Count("api.total", 1)
	metrics.Count("api.search", 1)
	metrics.Histogram("api.search.duration", (float64)(util.ElapsedTimeSince(start)))
}
