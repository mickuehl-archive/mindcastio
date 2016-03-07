package backend

const (
	FIRST_UPDATE_RATE          int   = 30   // min.
	DEFAULT_UPDATE_RATE        int   = 720  // min.
	DEFAULT_CRAWLER_SCHEDULE   int64 = 60   // sec
	DEFAULT_INDEXER_SCHEDULE   int64 = 60   // sec
	DEFAULT_UPDATE_BATCH       int   = 50   // how many podcasts to update per crawler run
	DEFAULT_INDEX_UPDATE_BATCH int   = 1000 // how many podcasts or episodes to send to elasicsearch each batch
	MAX_ERRORS                 int   = 4
	SEARCH_REVISION            int   = 1
)

type (
	PodcastIndex struct {
		Uid        string `json:"uid"`
		Feed       string `json:"feed"`
		UpdateRate int    `json:"update_rate"` // update rate in minutes
		Next       int64  `json:"next"`        // next scheduled update (unix time)
		N          int64  `json:"n"`
		Errors     int    `json:errors`
		Created    int64  `json:"created"`
		Updated    int64  `json:"updated"`
	}

	/*
		PodcastCollection struct {
			Count    int              `json:"count"`
			Podcasts []PodcastSummary `json:"podcasts"`
		}

		PodcastSummary struct {
			Uid         string `json:"uid"`
			Title       string `json:"title"`
			Author      string `json:"author"`
			Description string `json:"description"`
			Url         string `json:"url"`
			Feed        string `json:"feed"`
			ImageUrl    string `json:"image_url"`

			// internal admin stuff

			Published int64 `json:"published"`
		}
	*/

	PodcastMetadata struct {
		Uid         string `json:"uid"`
		Title       string `json:"title"`
		Subtitle    string `json:"subtitle"`
		Url         string `json:"url"`
		Feed        string `json:"feed"`
		Description string `json:"description"`
		Published   int64  `json:"published"`
		Language    string `json:"language"`
		ImageUrl    string `json:"image_url"`
		OwnerName   string `json:"owner_name"`
		OwnerEmail  string `json:"owner_email"`
		Tags        string `json:"tags"`

		// internal admin stuff

		Score1  int64 `json:"score1"` // scores, not defined yet
		Score2  int64 `json:"score2"`
		Score3  int64 `json:"score3"`
		Version int   `json:"version"`

		Created int64 `json:"created"`
		Updated int64 `json:"updated"`
	}

	EpisodeMetadata struct {
		Uid         string `json:"uid"`
		Title       string `json:"title"`
		Url         string `json:"url"`
		Description string `json:"description"`
		Published   int64  `json:"published"`
		Duration    int64  `json:"duration"`
		Author      string `json:"author"`
		AssetUrl    string `json:"asset_url"`
		AssetType   string `json:"asset_type"`
		AssetSize   int    `json:"asset_size"`

		// internal admin stuff

		PodcastUid string `json:"puid"`
		Version    int    `json:"version"`

		Created int64 `json:"created"`
		Updated int64 `json:"updated"`
	}

	SearchTerm struct {
		Term    string `json:"term"`
		Created int64  `json:"created"`
	}

	SearchKeyword struct {
		Word      string `json:"word"`
		Frequency int64  `json:"frequency"`
	}
)
