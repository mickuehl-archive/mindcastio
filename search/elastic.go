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

	ElasticOperatorQuery struct {
		Query MatchOperator `json:"query"`
	}

	ElasticPrecisionQuery struct {
		Query MatchPrecision `json:"query"`
	}

	MatchOperator struct {
		Match FieldOperator `json:"match"`
	}

	MatchPrecision struct {
		Match FieldPrecision `json:"match"`
	}

	FieldOperator struct {
		Field QueryOperator `json:"_all"`
	}

	FieldPrecision struct {
		Field QueryPrecision `json:"_all"`
	}

	QueryOperator struct {
		Q        string `json:"query"`
		Operator string `json:"operator"`
	}

	QueryPrecision struct {
		Q        string `json:"query"`
		Operator string `json:"minimum_should_match"`
	}
)

func SearchElastic(q string, page int, limit int) (*SearchResult, error) {

	// query url
	from := limit * (page - 1)
	url := strings.Join([]string{environment.GetEnvironment().SearchServiceUrl(), "podcasts/podcast/_search?size=", strconv.FormatInt((int64)(limit), 10), "&from=", strconv.FormatInt((int64)(from), 10)}, "")

	// query payload
	//query_body := ElasticOperatorQuery{MatchOperator{FieldOperator{QueryOperator{q, "and"}}}}
	query_body := ElasticPrecisionQuery{MatchPrecision{FieldPrecision{QueryPrecision{q, "2<75%"}}}}

	result := ElasticResponse{}
	err := util.PostJson(url, &query_body, &result)

	podcasts := make([]*Result, len(result.Hits.Hits))
	for i := range result.Hits.Hits {
		podcasts[i] = elasticToResult(&result.Hits.Hits[i])
	}

	return &SearchResult{"", result.Hits.Total, q, 0, podcasts}, err
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
