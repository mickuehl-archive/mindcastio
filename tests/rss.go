package main

import (
	"fmt"

	"github.com/mindcastio/mindcastio/backend/feed"
)

func main() {

	url := "http://climate.wordpress.com/feed/"

	channel, _ := feed.RSS(url)
	fmt.Println(channel)
}
