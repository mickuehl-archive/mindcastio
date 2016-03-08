package search

import (
	"strconv"
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

func searchElastic(q string, page int, limit int) (*SearchResult, error) {

	from := limit * (page - 1)
	query1 := strings.Join([]string{environment.GetEnvironment().SearchServiceUrl(), "podcasts/podcast/_search?q=", q}, "")
	query := strings.Join([]string{query1, "&size=", strconv.FormatInt((int64)(limit), 10), "&from=", strconv.FormatInt((int64)(from), 10)}, "")

	// FIXME we currently only search the podcast index, episodes are ignored !

	result := ElasticResponse{}
	err := util.GetJson(query, &result)

	podcasts := make([]*Result, len(result.Hits.Hits))
	for i := range result.Hits.Hits {
		podcasts[i] = elasticToResult(&result.Hits.Hits[i])
	}

	return &SearchResult{"", result.Hits.Total, q, podcasts}, err
}

func elasticToResult(item *HitDetail) *Result {

	podcast := backend.PodcastLookup(item.Id)
	if podcast == nil {
		return &Result{
			item.Id,
			"podcast",
			"", "", "", "", "", "", 0, 0,
		}
	}

	return &Result{
		item.Id,
		"podcast",
		podcast.Title,
		podcast.Subtitle,
		podcast.Description,
		podcast.Url,
		podcast.Feed,
		podcast.ImageUrl,
		(int)(item.Score * 100),
		podcast.Published,
	}

}
