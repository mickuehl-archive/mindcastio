package main

import (
	"strings"

	"github.com/mindcastio/mindcastio/backend"
	"github.com/mindcastio/mindcastio/backend/environment"
	"github.com/mindcastio/mindcastio/backend/util"
)

type (
	ElasticResponse struct {
		Took    int   `json:"took"`
		TimeOut bool  `json:"time_out"`
		Shards  Shard `json:"_shards"`
		Hits    HitsInfo
	}

	Shard struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Failed     int `json:"failed"`
	}

	HitsInfo struct {
		Total    int         `json:"total"`
		MaxScore float32     `json:"max_score"`
		Hits     []HitDetail `json:"hits"`
	}

	HitDetail struct {
		Index string  `json:"_index"`
		Kind  string  `json:"_type"`
		Id    string  `json:"_id"`
		Score float32 `json:"_score"`
	}
)

func searchElastic(q string) ([]*Result, error) {

	query := strings.Join([]string{environment.GetEnvironment().SearchServiceUrl(), "search/_search?q=", q}, "")

	result := ElasticResponse{}
	err := util.GetJson(query, &result)

	podcasts := make([]*Result, len(result.Hits.Hits))
	for i := range result.Hits.Hits {
		podcasts[i] = elasticToResult(&result.Hits.Hits[i])
	}
	return podcasts, err

}

func elasticToResult(item *HitDetail) *Result {

	podcast := backend.PodcastLookup(item.Id)
	if podcast == nil {
		return &Result{
			item.Id,
			"podcast",
			"", "", "", "", "", 0, 0,
		}
	}

	return &Result{
		item.Id,
		"podcast",
		podcast.Title,
		podcast.Description,
		podcast.Url,
		podcast.Feed,
		podcast.ImageUrl,
		(int)(item.Score * 100),
		0,
	}

}
