package main

import (
	"fmt"

	"github.com/mindcastio/mindcastio/backend/feed"
)

func main() {

	url := "http://jerrywho.podOmatic.com/rss2.xml"

	channel, _ := feed.RSS(url)
	fmt.Println(channel)
}
