#!/bin/bash

VERSION=1.6

cd /tmp
wget https://storage.googleapis.com/golang/go$VERSION.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go$VERSION.linux-amd64.tar.gz

echo # "golang $VERSION" >> /home/vagrant/.bash_profile
echo 'export GOROOT=/usr/local/go' >> /home/vagrant/.bash_profile
echo 'export PATH=$PATH:$GOROOT/bin' >> /home/vagrant/.bash_profile

rm -rf go$VERSION*
