package main

import (
	"github.com/mindcastio/mindcastio/crawler"
	"github.com/mindcastio/mindcastio/backend/util"
)

func main() {

	url := "http://climate.wordpress.com/feed/"

	podcast, _ := crawler.ParsePodcastFeed(url)
	util.PrettyPrintJson(podcast)
}
