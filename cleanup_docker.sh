#!/usr/bin/env bash

docker compose rm -f db migrate
docker volume rm "$(docker volume ls -q)"
