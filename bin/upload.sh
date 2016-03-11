#!/bin/bash

function post {
	echo "Submitting '$1'"
	curl -i -H 'Content-Type: application/json' -d "{\"feed\":\"$1\"}" $2/api/1/submit
	echo ""
}

filename=$1
endpoint=$2

if [ -f $filename ]
then

  while read line
  do
	  post $line $endpoint
	  #sleep 1
  done < $filename

fi
