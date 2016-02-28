package backend

import (
	"strings"

	"github.com/franela/goreq"

	"github.com/mindcastio/mindcastio/backend/environment"
)

func podcastAddToIndex(podcast *PodcastMetadata) error {

	uri := strings.Join([]string{environment.GetEnvironment().SearchServiceUrl(), "/search/podcast/", podcast.Uid}, "")
	payload := podcastMetadataToSearch(podcast)

	// post the payload to elasticsearch
	res, err := goreq.Request{
		Method:      "PUT",
		Uri:         uri,
		ContentType: "application/json",
		Body:        payload,
	}.Do()

	if res != nil {
		res.Body.Close()
	}
	return err
}

func episodeAddToIndex(episode *EpisodeMetadata) error {

	uri := strings.Join([]string{environment.GetEnvironment().SearchServiceUrl(), "/search/episode/", episode.Uid}, "")
	payload := episodeMetadataToSearch(episode)

	// post the payload to elasticsearch
	res, err := goreq.Request{
		Method:      "PUT",
		Uri:         uri,
		ContentType: "application/json",
		Body:        payload,
	}.Do()

	if res != nil {
		res.Body.Close()
	}
	return err
}
