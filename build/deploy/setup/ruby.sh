#!/usr/bin/env bash

# add dev dependencies
sudo apt-get -y install ruby libsqlite3-dev

VERSION='2.3'
PATCH='0'

# build ruby from source
cd /tmp
wget http://ftp.ruby-lang.org/pub/ruby/$VERSION/ruby-$VERSION.$PATCH.tar.gz
tar -xzvf ruby-$VERSION.$PATCH.tar.gz
cd ruby-$VERSION.$PATCH
sudo ./configure --disable-install-rdoc
sudo make
sudo apt-get -y remove ruby
sudo make install
sudo rm -rf /tmp/ruby*

# path to new ruby, for all users (system-wide)
sudo echo "# ruby $VERSION.$PATCH" >> ~/.bash_profile
sudo echo 'export PATH=$PATH:/usr/local/bin' >> ~/.bash_profile

# the path for .gems is 2.2.0 even if we have ruby 2.2.x
sudo echo 'export PATH=$PATH:$HOME/.gem/ruby/$VERSION.0/bin:/usr/local/lib/ruby/gems/$VERSION.0' >> ~/.bash_profile

# Disable document generation, update RubyGems and install Bundler 
sudo echo "gem: --no-document" >> ~/.gemrc
#sudo gem update --system
#sudo gem install bundler