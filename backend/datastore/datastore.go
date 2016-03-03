package datastore

import (
	"gopkg.in/mgo.v2"

	"github.com/mindcastio/mindcastio/backend/environment"
	"github.com/mindcastio/mindcastio/backend/logger"
)

const (
	DATABASE string = "mindcast"

	// backend store constants
	META_COL     string = "meta"
	PODCASTS_COL string = "podcasts"
	EPISODES_COL string = "episodes"
)

var _session *mgo.Session

func Initialize(env *environment.Environment) {

	logger.Log("datastore.initialize", env.BackendServiceHosts()[0])

	_s, err := mgo.Dial(env.BackendServiceHosts()[0])
	if err != nil {
		panic(err)
	}
	_session = _s

	// make sure we have all the right structures in place
	createIndex()
}

func Shutdown() {
	logger.Log("datastore.shutdown")
	_session.Close()
}

//
// struct used to access the backend database
//

type DataStore struct {
	session *mgo.Session
}

func GetDataStore() *DataStore {
	return &DataStore{_session.Copy()}
}

func (ds *DataStore) Close() {
	ds.session.Close()
}

func (ds *DataStore) Collection(c string) *mgo.Collection {
	return ds.session.DB(DATABASE).C(c)
}
