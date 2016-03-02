#!/bin/bash

docker stop natsd
docker stop mongod
docker stop elasticd

docker rm natsd mongod elasticd

docker pull mongo
docker pull elasticsearch
docker pull nats

docker create --name mongod -p $PRIVATE_IPV4:27017:27017 mongo --storageEngine=wiredTiger
docker create --name elasticd -p $PRIVATE_IPV4:9200:9200 -p $PRIVATE_IPV4:9300:9300 elasticsearch
docker create --name natsd -p $PRIVATE_IPV4:4222:4222 -p 6222:6222 nats

docker start natsd mongod elasticd
