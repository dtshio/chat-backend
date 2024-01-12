#!/bin/sh

docker-compose -f docker-compose.yml -f db/docker-compose.yml --env-file=.env.local up
