#!/bin/bash

# updated dependencies
# go get -u "github.com/mindcastio/podcast-feed"
# go get -u "github.com/mindcastio/go-json-rest"

# build binaries
cd $MINDCASTIO_HOME/services/crawler
go build

cd $MINDCASTIO_HOME/services/indexer
go build

cd $MINDCASTIO_HOME/services/search
go build
