#!/bin/bash

docker-compose up -d 
sleep 3
docker run -v $(pwd)/migrations:/migrations --network vk-feed_default migrate/migrate -path=/migrations/ -database "postgres://postgres:example@db/postgres?sslmode=disable" up 1
