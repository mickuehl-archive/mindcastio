package backend

import (
	"gopkg.in/mgo.v2/bson"
	"math"

	"github.com/mindcastio/podcast-feed"

	"github.com/mindcastio/mindcastio/backend/util"
	"github.com/mindcastio/mindcastio/backend/services/datastore"
)

func indexLookup(uid string) *PodcastIndex {

	ds := datastore.GetDataStore()
	defer ds.Close()

	main_index := ds.Collection(datastore.META_COL)

	i := PodcastIndex{}
	main_index.Find(bson.M{"uid": uid}).One(&i)

	if i.Feed == "" {
		return nil
	} else {
		return &i
	}
}

func indexAdd(uid string, url string) error {
	ds := datastore.GetDataStore()
	defer ds.Close()

	main_index := ds.Collection(datastore.META_COL)

	// add some random element to the first update point in time
	next := util.IncT(util.Timestamp(), 10+util.Random(DEFAULT_UPDATE_RATE))

	i := PodcastIndex{uid, url, DEFAULT_UPDATE_RATE, next, 0, 0, util.Timestamp(), 0}
	return main_index.Insert(&i)
}

func indexUpdate(uid string) error {

	ds := datastore.GetDataStore()
	defer ds.Close()

	main_index := ds.Collection(datastore.META_COL)

	i := PodcastIndex{}
	err := main_index.Find(bson.M{"uid": uid}).One(&i)

	if i.Feed == "" || err != nil {
		return err
	} else {
		i.Updated = util.Timestamp()
		i.Next = util.IncT(i.Next, i.UpdateRate+util.RandomPlusMinus(15))
		i.Errors = 0 // reset in case there was an erro

		// update the DB
		err = main_index.Update(bson.M{"uid": uid}, &i)
	}

	return err
}

func indexBackoff(uid string) (bool, error) {

	ds := datastore.GetDataStore()
	defer ds.Close()

	suspended := false
	main_index := ds.Collection(datastore.META_COL)

	i := PodcastIndex{}
	err := main_index.Find(bson.M{"uid": uid}).One(&i)

	if i.Feed == "" || err != nil {
		return suspended, err
	} else {
		i.Updated = util.Timestamp()
		i.Errors++

		if i.Errors > MAX_ERRORS {
			// just disable the UID by using a LAAAARGE next time
			i.Next = math.MaxInt64
			suspended = true
		} else {
			// + 10, 100, 1000, 10000 min ...
			i.Next = util.IncT(i.Updated, (int)(math.Pow(10, (float64)(i.Errors))))
		}

		// update the DB
		err = main_index.Update(bson.M{"uid": uid}, &i)
	}

	return suspended, err
}

func podcastLookup(uid string) *PodcastMetadata {

	ds := datastore.GetDataStore()
	defer ds.Close()

	podcast_metadata := ds.Collection(datastore.PODCASTS_COL)

	p := PodcastMetadata{}
	podcast_metadata.Find(bson.M{"uid": uid}).One(&p)

	if p.Uid == "" {
		return nil
	} else {
		return &p
	}
}

func podcastAdd(podcast *podcast.PodcastDetails) (bool, error) {
	p := podcastLookup(podcast.Uid)
	if p != nil {
		return false, nil
	}

	ds := datastore.GetDataStore()
	defer ds.Close()

	podcast_metadata := ds.Collection(datastore.PODCASTS_COL)

	meta := podcastDetailsToMetadata(podcast)

	// fix the published timestamp
	now := util.Timestamp()
	if podcast.Published > now {
		meta.Published = now // prevents dates in the future
	}

	err := podcast_metadata.Insert(&meta)

	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func podcastUpdateTimestamp(podcast *podcast.PodcastDetails) (bool, error) {

	ds := datastore.GetDataStore()
	defer ds.Close()

	podcast_metadata := ds.Collection(datastore.PODCASTS_COL)

	p := PodcastMetadata{}
	podcast_metadata.Find(bson.M{"uid": podcast.Uid}).One(&p)

	if p.Uid == "" {
		return false, nil
	} else {
		now := util.Timestamp()
		p.Updated = now
		if podcast.Published > now {
			p.Published = now // prevents dates in the future
		} else {
			p.Published = podcast.Published
		}

		// update the DB
		err := podcast_metadata.Update(bson.M{"uid": podcast.Uid}, &p)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func episodeLookup(uid string) *EpisodeMetadata {

	ds := datastore.GetDataStore()
	defer ds.Close()

	episodes_metadata := ds.Collection(datastore.EPISODES_COL)

	e := EpisodeMetadata{}
	episodes_metadata.Find(bson.M{"uid": uid}).One(&e)

	if e.Uid == "" {
		return nil
	} else {
		return &e
	}
}

func episodeAdd(episode *podcast.EpisodeDetails, puid string) (bool, error) {
	e := episodeLookup(episode.Uid)
	if e != nil {
		return false, nil
	}

	ds := datastore.GetDataStore()
	defer ds.Close()

	episodes_metadata := ds.Collection(datastore.EPISODES_COL)

	meta := episodeDetailsToMetadata(episode, puid)
	err := episodes_metadata.Insert(&meta)

	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func episodesAddAll(podcast *podcast.PodcastDetails) (int, error) {
	count := 0

	for i := 0; i < len(podcast.Episodes); i++ {
		added, err := episodeAdd(&podcast.Episodes[i], podcast.Uid)
		if err != nil {
			return 0, err
		}
		if added {
			count++
		}
	}
	return count, nil
}

func latestUpdatedPodcasts(limit int, page int) (*PodcastCollection, error) {

	ds := datastore.GetDataStore()
	defer ds.Close()

	podcast_metadata := ds.Collection(datastore.PODCASTS_COL)

	results := []PodcastMetadata{}
	err := podcast_metadata.Find(nil).Limit(limit).Sort("-published").All(&results)
	if err != nil {
		return nil, err
	}

	podcasts := make([]PodcastSummary, len(results))
	for i := 0; i < len(results); i++ {
		podcasts[i] = podcastMetadataToSummary(&results[i])
	}

	podcastCollection := PodcastCollection{
		len(results),
		podcasts,
	}

	return &podcastCollection, nil
}
