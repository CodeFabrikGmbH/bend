#!/usr/bin/env bash

set -ex

CURRENTDATE=`date +"%Y-%m-%d %T"`

curl -v \
  -H "Authorization: Bearer bearerDataStuff" \
  -H "Content-Type: application/json" \
  -d "{ \"happenedAt\": \"${CURRENTDATE}\", \"description\": \"Something happened.\" }" \
  "http://localhost:8080/other-endpoint?testParameter=Y"