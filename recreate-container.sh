#!/usr/bin/env bash
set -e

docker-compose down --rmi all | true
docker-compose up -d --force-recreate