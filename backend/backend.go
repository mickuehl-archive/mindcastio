package backend

import (
	"math"
	"strconv"
	"strings"

	"gopkg.in/mgo.v2/bson"

	"github.com/mindcastio/mindcastio/backend/datastore"
	"github.com/mindcastio/mindcastio/backend/logger"
	"github.com/mindcastio/mindcastio/backend/metrics"
	"github.com/mindcastio/mindcastio/backend/util"
)

func SubmitPodcastFeed(feed string) error {

	logger.Log("submit_podcast_feed", feed)

	// check if the podcast is already in the index
	uid := util.UID(feed)
	idx := IndexLookup(uid)

	if idx == nil {
		err := IndexAdd(uid, feed)
		if err != nil {

			logger.Error("submit_podcast_feed.error", err, feed)
			metrics.Error("submit_podcast_feed.error", err.Error(), []string{feed})

			return err
		}
	} else {
		logger.Warn("submit_podcast_feed.duplicate", uid, feed)
		metrics.Warning("submit_podcast_feed.duplicate", "", []string{feed})
	}

	logger.Log("submit_podcast_feed.done", uid, feed)
	return nil
}

func BulkSubmitPodcastFeed(urls []string) error {

	logger.Log("bulk_submit_podcast_feed")

	count := 0
	feed := ""

	for i := 0; i < len(urls); i++ {
		feed = urls[i]

		// check if the podcast is already in the index
		uid := util.UID(feed)
		idx := IndexLookup(uid)

		if idx == nil {
			err := IndexAdd(uid, feed)
			if err != nil {

				logger.Error("bulk_submit_podcast_feed.error", err, feed)
				metrics.Error("bulk_submit_podcast_feed.error", err.Error(), []string{feed})

				return err
			} else {
				count++
			}
		}
	}

	logger.Log("bulk_submit_podcast_feed.done", strconv.FormatInt((int64)(count), 10))
	return nil
}

func IndexLookup(uid string) *PodcastIndex {

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

func IndexAdd(uid string, url string) error {
	ds := datastore.GetDataStore()
	defer ds.Close()

	main_index := ds.Collection(datastore.META_COL)

	// add some random element to the first update point in time
	next := util.IncT(util.Timestamp(), 2+util.Random(FIRST_UPDATE_RATE))

	i := PodcastIndex{uid, url, DEFAULT_UPDATE_RATE, next, 0, 0, util.Timestamp(), 0}
	return main_index.Insert(&i)
}

func IndexUpdate(uid string) error {

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

func IndexBackoff(uid string) (bool, error) {

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

func PodcastLookup(uid string) *PodcastMetadata {

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

func EpisodeLookup(uid string) *EpisodeMetadata {

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

func LogSearchString(s string) {
	ds := datastore.GetDataStore()
	defer ds.Close()

	search_term := ds.Collection(datastore.SEARCH_TERM_COM)
	search_term.Insert(&SearchTerm{strings.Replace(s, "+", " ", -1), util.Timestamp()})

	// split into keywords and update the dictionary
	search_keywords := ds.Collection(datastore.KEYWORDS_COL)

	tt := strings.Split(s, "+")
	//if len(tt) == 0 {
	//	tt := make([]string, 1)
	//	tt[0] = s
	//}

	for i := range tt {
		t := SearchKeyword{}
		search_keywords.Find(bson.M{"word": tt[i]}).One(&t)
		if t.Word == "" {
			t.Word = tt[i]
			t.Frequency = 1
			err := search_keywords.Insert(&t)
			if err != nil {
				logger.Error("log_search_string.error", err, s)
			}

		} else {
			t.Frequency = t.Frequency + 1
			err := search_keywords.Update(bson.M{"word": tt[i]}, &t)
			if err != nil {
				logger.Error("log_search_string.error", err, s)
			}
		}

	}

}



/*
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

func simpleStats() (*ApiInfo, error) {

	ds := GetDataStore()
	defer ds.Close()

	podcast_metadata := ds.Collection(PODCASTS_COL)
	podcasts, _ := podcast_metadata.Count()

	episodes_metadata := ds.Collection(EPISODES_COL)
	episodes, _ := episodes_metadata.Count()

	info := ApiInfo{
		BACKEND_VERSION,
		podcasts,
		episodes,
	}

	return &info, nil
}

*/
