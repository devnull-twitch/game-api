#!/bin/sh

docker volume create 2donlinerpg

docker run --rm \
    --name game-db-docker \
    -p 5434:5432 \
    -e POSTGRES_USER=game \
    -e POSTGRES_PASSWORD=game \
    -v 2donlinerpg:/var/lib/postgresql/data \
    postgres:14-alpine