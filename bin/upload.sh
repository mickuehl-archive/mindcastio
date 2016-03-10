#!/bin/bash

function post {
	echo "Submitting '$1'"
	curl -i -H 'Content-Type: application/json' -d "{\"feed\":\"$1\"}" $2/api/1/submit
	echo ""
}

if [ -f $1 ]
then

  while read line
  do
	  post $line $2
	  #sleep 1
  done < $1

fi
