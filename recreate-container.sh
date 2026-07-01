#!/usr/bin/env bash
set -e

# Use Compose v2 ("docker compose"). Compose v1 ("docker-compose") has a
# KeyError: 'ContainerConfig' bug when recreating containers whose image was
# built with BuildKit, which breaks deploys.
docker compose down --rmi all || true
docker compose up -d --force-recreate
