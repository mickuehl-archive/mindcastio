package main

import (

	"gopkg.in/mgo.v2/bson"

	"github.com/mindcastio/mindcastio/backend"
	"github.com/mindcastio/mindcastio/backend/datastore"
	"github.com/mindcastio/mindcastio/backend/environment"
	"github.com/mindcastio/mindcastio/backend/logger"
	"github.com/mindcastio/mindcastio/backend/metrics"
	"github.com/mindcastio/mindcastio/backend/util"
)

/*
	migration_01

	Change the update frequency and rebalance the update schedule.
*/

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
	main_index.Find(nil).All(&results)

	for i := range results {
		results[i].UpdateRate = backend.DEFAULT_UPDATE_RATE
		results[i].Next = util.IncT(util.Timestamp(), util.Random(backend.DEFAULT_UPDATE_RATE))
		results[i].Updated = util.Timestamp()

		main_index.Update(bson.M{"uid": results[i].Uid}, &results[i])
	}

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
