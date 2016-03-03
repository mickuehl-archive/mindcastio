package search

type (
	SearchResult struct {
		Uid        string    `jsonapi:"primary,search"`
		Count      int       `jsonapi:"attr,count"`
		SearchTerm string    `jsonapi:"attr,search_term"`
		Results    []*Result `jsonapi:"relation,results"`
	}

	Result struct {
		Uid         string `jsonapi:"primary,result"`
		Kind        string `jsonapi:"attr,kind"` // podcast | episode
		Title       string `jsonapi:"attr,title"`
		Subtitle    string `jsonapi:"attr,subtitle"`
		Description string `jsonapi:"attr,description"`
		Url         string `jsonapi:"attr,url"`
		Feed        string `jsonapi:"attr,feed"`
		ImageUrl    string `jsonapi:"attr,image_url"`

		// metadata
		Score     int   `jsonapi:"attr,score"` // scaled to [0..100]
		Published int64 `jsonapi:"attr,published"`
	}

	PodcastMetadataSearch struct {
		Uid         string `json:"uid"`
		Title       string `json:"title"`
		Subtitle    string `json:"subtitle"`
		Description string `json:"description"`
		Published   int64  `json:"published"`
		Language    string `json:"language"`
		OwnerName   string `json:"owner_name"`
		OwnerEmail  string `json:"owner_email"`
		Tags        string `json:"tags"`
	}

	EpisodeMetadataSearch struct {
		Uid         string `json:"uid"`
		Title       string `json:"title"`
		Link        string `json:"link"`
		Description string `json:"description"`
		Published   int64  `json:"published"`
		Author      string `json:"author"`
		PodcastUid  string `json:"puid"`
	}
	
)
