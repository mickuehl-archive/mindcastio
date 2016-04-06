package search

import (
	"strings"

	"github.com/mindcastio/mindcastio/backend/util"
)

const (
	ITUNES_SEARCH_URL string = "https://itunes.apple.com/search?media=podcast&entity=podcast&limit=200&term="
)

type (
	iTunesResponse struct {
		ResultCount int          `json:"resultCount"`
		Items       []iTunesItem `json:"results"`
	}

	iTunesItem struct {
		WrapperType            string   `json:"wrapperType"`
		Kinf                   string   `json:"kind"`
		CollectionId           int64    `json:"collectionId"`
		TrackId                int64    `json:"trackId"`
		ArtistName             string   `json:"artistName"`
		CollectionName         string   `json:"collectionName"`
		TrackName              string   `json:"trackName"`
		CollectionCensoredName string   `json:"collectionCensoredName"`
		TrackCensoredName      string   `json:"trackCensoredName"`
		CollectionViewUrl      string   `json:"collectionViewUrl"`
		FeedUrl                string   `json:"feedUrl"`
		TrackViewUrl           string   `json:"trackViewUrl"`
		ArtworkUrl30           string   `json:"artworkUrl30"`
		ArtworkUrl60           string   `json:"artworkUrl60"`
		ArtworkUrl100          string   `json:"artworkUrl100"`
		CollectionPrice        float32  `json:"collectionPrice"`
		TrackPrice             float32  `json:"trackPrice"`
		TrackRentalPrice       float32  `json:"trackRentalPrice"`
		CollectionHdPrice      float32  `json:"collectionHdPrice"`
		TrackHdPrice           float32  `json:"trackHdPrice"`
		TrackHdRentalPrice     float32  `json:"trackHdRentalPrice"`
		ReleaseDate            string   `json:"releaseDate"`
		CollectionExplicitness string   `json:"collectionExplicitness"`
		TrackExplicitness      string   `json:"trackExplicitness"`
		TrackCount             int      `json:"trackCount"`
		Country                string   `json:"country"`
		Currency               string   `json:"currency"`
		PrimaryGenreName       string   `json:"primaryGenreName"`
		RadioStationUrl        string   `json:"radioStationUrl"`
		ArtworkUrl600          string   `json:"artworkUrl600"`
		GenreIds               []string `json:"genreIds"`
		Genres                 []string `json:"genres"`
	}
)

func SearchITunes(q string) ([]*Result, error) {

	query := strings.Join([]string{ITUNES_SEARCH_URL, q}, "")

	response := iTunesResponse{}
	err := util.GetJson(query, &response)

	if err != nil {
		return nil, err
	} else {
		podcasts := make([]*Result, len(response.Items))

		for i, item := range response.Items {
			podcasts[i] = iTunesToResult(&item)
		}

		return podcasts, nil
	}
}

func iTunesToResult(item *iTunesItem) *Result {
	result := Result{
		util.UID(item.FeedUrl),
		"podcast",
		item.CollectionName,
		"",
		item.CollectionName,
		"",
		item.FeedUrl,
		item.ArtworkUrl100,
		0,
		0,
	}
	return &result
}
