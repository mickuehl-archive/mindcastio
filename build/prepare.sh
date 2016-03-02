#!/bin/bash

export MINDCAST_HOME=/opt/data/build
export MINDCAST_SRC=/opt/data/build/src/github.com/mindcastio/mindcastio

export MINDCAST_REPO=https://github.com/mindcastio/mindcastio.git
export BRANCH=master

git clone $MINDCAST_REPO --branch $BRANCH --single-branch $MINDCAST_SRC
