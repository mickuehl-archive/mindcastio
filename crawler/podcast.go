package crawler

import (
	"strconv"
	"strings"

	"github.com/kennygrant/sanitize"

	"github.com/mindcastio/mindcastio/backend/feed"
	"github.com/mindcastio/mindcastio/backend/util"
)

type (
	Podcast struct {
		Title       string       `json:"title"`
		Subtitle    string       `json:"subtitle"`
		Url         string       `json:"url"`
		Feed        string       `json:"feed"`
		Uid         string       `json:"uid"`
		Description string       `json:"description"`
		Published   int64        `json:"published"`
		Language    string       `json:"language"`
		Image       string       `json:"image"`
		Owner       PodcastOwner `json:"owner"`
		Episodes    []Episode    `json:"episodes"`
	}

	Episode struct {
		Title       string     `json:"title"`
		Url         string     `json:"url"`
		Uid         string     `json:"uid"`
		Description string     `json:"description"`
		Text        string     `json:"text"`
		Published   int64      `json:"published"`
		Duration    int64      `json:"duration"`
		Author      string     `json:"author"`
		Content     MediaAsset `json:"content"`
	}

	PodcastOwner struct {
		Name  string `json:"title"`
		Email string `json:"title"`
	}

	MediaAsset struct {
		Url  string `json:"url"`
		Type string `json:"type"`
		Size int    `json:"size"`
	}

	Chapter struct {
		Start string `json:"start"`
		Title string `json:"title"`
	}
)

func ParsePodcastFeed(url string) (*Podcast, error) {
	// parse the podcast feed
	channel, err := feed.RSS(url)
	if err != nil {
		return nil, err
	}

	return channelToPodcast(channel, url), nil
}

func channelToPodcast(channel *feed.Channel, url string) *Podcast {

	uid := util.UID(url)

	// construct the return struct
	owner := PodcastOwner{
		channel.Owner.Name,
		channel.Owner.Email,
	}

	// items, i.e. episodes
	episodes := make([]Episode, len(channel.Item))
	for i, item := range channel.Item {

		// media assets
		content := MediaAsset{"", "", 0}
		if item.Enclosure != nil {
			content = MediaAsset{
				item.Enclosure[0].URL,
				item.Enclosure[0].Type,
				item.Enclosure[0].Length,
			}
		}

		e := Episode{
			item.Title,
			item.Link,
			util.Fingerprint(item.Title, uid),
			sanitize.HTML(item.Description),
			sanitize.HTML(item.Text),
			convertDateToUnix(item.PubDate),
			duration(item.Duration),
			item.Author,
			content,
		}

		episodes[i] = e
	}

	// podcast
	var lastBuildDate int64 = 0
	if len(episodes) != 0 {
		lastBuildDate = episodes[0].Published
	}

	p := Podcast{
		channel.Title,
		channel.Subtitle,
		channel.Link,
		url,
		uid,
		sanitize.HTML(channel.Description),
		lastBuildDate,
		language(channel.Language),
		channel.Image.URL,
		owner,
		episodes,
	}

	return &p
}

func duration(d string) int64 {
	var ss = strings.Split(d, ":")

	switch len(ss) {
	case 3:
		h, _ := strconv.Atoi(ss[0])
		m, _ := strconv.Atoi(ss[1])
		s, _ := strconv.Atoi(ss[2])
		return (int64)(s + m*60 + h*3600)
	case 2:
		m, _ := strconv.Atoi(ss[0])
		s, _ := strconv.Atoi(ss[1])
		return (int64)(s + m*60)
	default:
		return 0
	}
}

func convertDateToUnix(d feed.RSSDate) int64 {
	t, _ := d.Parse()
	return t.Unix()
}

func language(l string) string {
	return strings.ToUpper(strings.Split(l, "-")[0])
}
