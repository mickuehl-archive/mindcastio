package main

import (
	"strings"
	"time"

	"github.com/mindcastio/go-json-rest/rest"
	"github.com/mindcastio/mindcastio/backend"

	"github.com/mindcastio/mindcastio/backend/search"
	"github.com/mindcastio/mindcastio/backend/services/logger"
	"github.com/mindcastio/mindcastio/backend/services/metrics"

	"github.com/mindcastio/mindcastio/backend/util"
)

func endpoint(w rest.ResponseWriter, r *rest.Request) {
	start := time.Now()

	if len(r.URL.Query()["q"]) == 0 {
		backend.JsonApiErrorResponse(w, "search.search.error", "missing parameter", nil)
		metrics.Error("search.search.error", "", nil)
		return
	}

	q := r.URL.Query()["q"][0]
	query := strings.Replace(q, " ", "+", -1)

	if len(query) == 0 {
		backend.JsonApiErrorResponse(w, "search.search.error", "missing query", nil)
		metrics.Error("search.search.error", "", nil)
		return
	}

	logger.Log("search.search.query", query)

	result := search.Search(query)
	backend.JsonApiResponse(w, result)

	// metrics
	metrics.Count("api.total", 1)
	metrics.Count("search.search", 1)
	metrics.Histogram("search.search.duration", (float64)(util.ElapsedTimeSince(start)))
}
