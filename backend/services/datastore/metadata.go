package datastore

import (
	"gopkg.in/mgo.v2"

	"github.com/mindcastio/mindcastio/backend/services/logger"
)

func createIndex() {

	ds := GetDataStore()
	defer ds.Close()

	// main_index
	main_index := ds.Collection(META_COL)
	// main_index.uid
	err := main_index.EnsureIndex(mgo.Index{Key: []string{"uid"}, Unique: true, DropDups: true, Background: true, Sparse: true})
	if err != nil {
		logger.Error("backend.datastore.create_index", err, "")
	}
	// main_index.next
	err = main_index.EnsureIndex(mgo.Index{Key: []string{"next"}, Unique: false, DropDups: true, Background: true, Sparse: true})
	if err != nil {
		logger.Error("backend.datastore.create_index", err, "")
	}
	// main_index.errors
	err = main_index.EnsureIndex(mgo.Index{Key: []string{"errors"}, Unique: false, DropDups: true, Background: true, Sparse: true})
	if err != nil {
		logger.Error("backend.datastore.create_index", err, "")
	}

	// podcast metadata
	podcast_metadata := ds.Collection(PODCASTS_COL)
	// podcast_metadata.uid
	err = podcast_metadata.EnsureIndex(mgo.Index{Key: []string{"uid"}, Unique: true, DropDups: true, Background: true, Sparse: true})
	if err != nil {
		logger.Error("backend.datastore.create_index", err, "")
	}
	// podcast_metadata.published
	err = podcast_metadata.EnsureIndex(mgo.Index{Key: []string{"published"}, Unique: false, DropDups: true, Background: true, Sparse: true})
	if err != nil {
		logger.Error("backend.datastore.create_index", err, "")
	}
	// podcast_metadata.created
	err = podcast_metadata.EnsureIndex(mgo.Index{Key: []string{"created"}, Unique: false, DropDups: true, Background: true, Sparse: true})
	if err != nil {
		logger.Error("backend.datastore.create_index", err, "")
	}
	// podcast_metadata.score1
	err = podcast_metadata.EnsureIndex(mgo.Index{Key: []string{"score1"}, Unique: false, DropDups: true, Background: true, Sparse: true})
	if err != nil {
		logger.Error("backend.datastore.create_index", err, "")
	}
	// podcast_metadata.score2
	err = podcast_metadata.EnsureIndex(mgo.Index{Key: []string{"score2"}, Unique: false, DropDups: true, Background: true, Sparse: true})
	if err != nil {
		logger.Error("backend.datastore.create_index", err, "")
	}
	// podcast_metadata.score3
	err = podcast_metadata.EnsureIndex(mgo.Index{Key: []string{"score3"}, Unique: false, DropDups: true, Background: true, Sparse: true})
	if err != nil {
		logger.Error("backend.datastore.create_index", err, "")
	}
	// podcast_metadata.index_version
	err = podcast_metadata.EnsureIndex(mgo.Index{Key: []string{"index_version"}, Unique: false, DropDups: true, Background: true, Sparse: true})
	if err != nil {
		logger.Error("backend.datastore.create_index", err, "")
	}

	// episode metadata
	episodes_metadata := ds.Collection(EPISODES_COL)
	// episodes_metadata.uid
	err = episodes_metadata.EnsureIndex(mgo.Index{Key: []string{"uid"}, Unique: true, DropDups: true, Background: true, Sparse: true})
	if err != nil {
		logger.Error("backend.datastore.create_index", err, "")
	}
	// episodes_metadata.published
	err = episodes_metadata.EnsureIndex(mgo.Index{Key: []string{"published"}, Unique: false, DropDups: true, Background: true, Sparse: true})
	if err != nil {
		logger.Error("backend.datastore.create_index", err, "")
	}
	// episodes_metadata.asset_type
	err = episodes_metadata.EnsureIndex(mgo.Index{Key: []string{"asset_type"}, Unique: false, DropDups: true, Background: true, Sparse: true})
	if err != nil {
		logger.Error("backend.datastore.create_index", err, "")
	}
	// episodes_metadata.puid
	err = episodes_metadata.EnsureIndex(mgo.Index{Key: []string{"puid"}, Unique: false, DropDups: true, Background: true, Sparse: true})
	if err != nil {
		logger.Error("backend.datastore.create_index", err, "")
	}
	// episodes_metadata.index_version
	err = episodes_metadata.EnsureIndex(mgo.Index{Key: []string{"index_version"}, Unique: false, DropDups: true, Background: true, Sparse: true})
	if err != nil {
		logger.Error("backend.datastore.create_index", err, "")
	}
}
