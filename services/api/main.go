package main

import (
	stdlog "log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/mindcastio/go-json-rest/rest"
	"github.com/mindcastio/mindcastio/backend"
	"github.com/mindcastio/mindcastio/backend/environment"
	"github.com/mindcastio/mindcastio/backend/search"
	"github.com/mindcastio/mindcastio/backend/services/logger"
	"github.com/mindcastio/mindcastio/backend/services/metrics"
	"github.com/mindcastio/mindcastio/backend/services/datastore"
	"github.com/mindcastio/mindcastio/backend/util"
)

const (
	SEARCH_ENDPOINT string = "/api/1/search"
)

func main() {

	// environment setup
	env := environment.GetEnvironment()

	logger.Initialize()
	metrics.Initialize(env)
	datastore.Initialize(env)

	// initilize the REST API router
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)

	router, err := rest.MakeRouter(
		rest.Get(SEARCH_ENDPOINT, endpoint),
	)

	if err != nil {
		stdlog.Fatal(err)
	}
	api.SetApp(router)

	// setup shutdown handling
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		shutdown()
		os.Exit(1)
	}()

	// start the REST api
	logger.Log("mindcast.search.startup")
	metrics.Success("mindcast", "search.startup", nil)

	http.ListenAndServe(env.ListenPort(), api.MakeHandler())

}

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

func shutdown() {
	logger.Log("mindcast.search.shutdown")
	metrics.Success("mindcast", "search.shutdown", nil)

	// shutdown of services
	datastore.Shutdown()
	metrics.Shutdown()
}
