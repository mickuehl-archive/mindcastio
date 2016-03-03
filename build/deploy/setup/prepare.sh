#!/usr/bin/env bash

# make sure we are on a defined local
export LANGUAGE=en_US.UTF-8
export LC_ALL=en_US.UTF-8
export LANG=en_US.UTF-8
export LC_TYPE=en_US.UTF-8

# update the base image first
sudo apt-get -y update && sudo apt-get -y upgrade

# install some basics
sudo apt-get -y install git unzip curl sysstat wget locales

sudo locale-gen UTF-8

# install development basics
sudo apt-get -y install build-essential autoconf bison \
  zlib1g-dev libssl-dev libreadline6-dev libyaml-dev \
  libncurses5-dev zlib1g-dev libffi-dev libxslt1-dev libxml2-dev \
  software-properties-common

# mysql client libs
sudo apt-get -y install mariadb-client-5.5 libmariadbd-dev

# 'disable' the firewall
sudo iptables -P INPUT ACCEPT
sudo iptables -P OUTPUT ACCEPT
sudo iptables -P FORWARD ACCEPT
sudo iptables -F
