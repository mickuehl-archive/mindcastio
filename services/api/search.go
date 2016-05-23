package main

import (
	"strconv"
	"time"

	"github.com/ant0ine/go-json-rest/rest"

	"github.com/mindcastio/mindcastio/search"

	"github.com/mindcastio/mindcastio/backend"
	"github.com/mindcastio/mindcastio/backend/logger"
	"github.com/mindcastio/mindcastio/backend/metrics"

	"github.com/mindcastio/mindcastio/backend/util"
)

func search_endpoint(w rest.ResponseWriter, r *rest.Request) {
	start := time.Now()

	var size int = search.PAGE_SIZE
	var page int = 1

	// &size=25
	if len(r.URL.Query()["size"]) != 0 {
		ss, _ := strconv.ParseInt(r.URL.Query()["size"][0], 10, 64)
		size = (int)(ss)
		if size < 1 {
			size = search.PAGE_SIZE
		}
	}

	// &page=1
	if len(r.URL.Query()["page"]) != 0 {
		pp, _ := strconv.ParseInt(r.URL.Query()["page"][0], 10, 64)
		page = (int)(pp)
		if page < 1 {
			page = 1
		}
	}

	// &q=Harry+Potter
	if len(r.URL.Query()["q"]) == 0 {
		backend.JsonApiErrorResponse(w, "api.search.error", "missing parameter", nil)
		metrics.Error("api.search.error", "", nil)
		return
	}

	q := r.URL.Query()["q"][0]
	if len(q) == 0 {
		backend.JsonApiErrorResponse(w, "api.search.error", "missing query", nil)
		metrics.Error("api.search.error", "", nil)
		return
	}

	query := util.NormalizeSearchString(q)
	logger.Log("api.search.query", query)

	result := search.Search(query, page, size)
	backend.JsonApiResponse(w, result)

	// metrics
	metrics.Count("api.total.count", 1)
	metrics.Count("api.search.count", 1)
	metrics.Histogram("api.search.duration", (float64)(util.ElapsedTimeSince(start)))
}
