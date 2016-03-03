package main

import (
	stdlog "log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/mindcastio/go-json-rest/rest"

	"github.com/mindcastio/mindcastio/backend/datastore"
	"github.com/mindcastio/mindcastio/backend/environment"
	"github.com/mindcastio/mindcastio/backend/logger"
	"github.com/mindcastio/mindcastio/backend/metrics"
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

func shutdown() {
	logger.Log("mindcast.search.shutdown")
	metrics.Success("mindcast", "search.shutdown", nil)

	// shutdown of services
	datastore.Shutdown()
	metrics.Shutdown()
}
