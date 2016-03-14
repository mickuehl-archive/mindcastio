#!/bin/bash

MINDCAST_SRC=`pwd`

# build binaries
cd $MINDCAST_SRC/services/crawler
echo "Building the crawler ..."
go get && go build

cd $MINDCAST_SRC/services/indexer
echo "Building the indexer ..."
go get && go build

cd $MINDCAST_SRC/services/api
echo "Building the api service ..."
go get && go build

echo "Building tools"

cd $MINDCAST_SRC/tools/export_index
go get && go build
