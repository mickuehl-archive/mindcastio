#!/bin/bash

# build binaries
cd $MINDCAST_SRC/services/crawler
go get && go build

cd $MINDCAST_SRC/services/indexer
go get && go build

cd $MINDCAST_SRC/services/search
go get && go build

if [ ! -L "/usr/local/bin/mindcast-crawler" ]; then
	sudo ln -s "$MINDCAST_SRC/services/crawler/crawler" /usr/local/bin/mindcast-crawler
fi

if [ ! -L "/usr/local/bin/mindcast-indexer" ]; then
	sudo ln -s "$MINDCAST_SRC/services/indexer/indexer" /usr/local/bin/mindcast-indexer
fi

if [ ! -L "/usr/local/bin/mindcast-search" ]; then
	sudo ln -s "$MINDCAST_SRC/services/search/search" /usr/local/bin/mindcast-search
fi
