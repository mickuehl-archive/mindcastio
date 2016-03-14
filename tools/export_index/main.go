package main

import (
	"os"

	"github.com/mindcastio/mindcastio/backend"
	"github.com/mindcastio/mindcastio/backend/datastore"
	"github.com/mindcastio/mindcastio/backend/environment"
	"github.com/mindcastio/mindcastio/backend/logger"
	"github.com/mindcastio/mindcastio/backend/metrics"
)

func main() {

	// environment setup
	env := environment.GetEnvironment()
	logger.Initialize()
	metrics.Initialize(env)
	defer metrics.Shutdown()
	datastore.Initialize(env)
	defer datastore.Shutdown()

	ds := datastore.GetDataStore()
	defer ds.Close()

	main_index := ds.Collection(datastore.META_COL)

	results := []backend.PodcastIndex{}
	main_index.Find(nil).Sort("-created").All(&results)

	// open the file
	f, err := os.Create("export_index.txt")
	defer f.Close()

	check(err)

	for i := range results {
		f.WriteString(results[i].Feed + "\n")
	}

	// flush writes
	f.Sync()

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
