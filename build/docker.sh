#!/bin/bash

docker stop mongod
docker stop elasticd

docker rm mongod elasticd

docker pull mongo
docker pull elasticsearch

docker create --name mongod -p $PRIVATE_IP4:27017:27017 mongo --storageEngine=wiredTiger
docker create --name elasticd -p $PRIVATE_IP4:9200:9200 -p $PRIVATE_IP4:9300:9300 elasticsearch
#docker create --name natsd -p $PRIVATE_IPV4:4222:4222 -p 6222:6222 nats

docker start mongod elasticd
