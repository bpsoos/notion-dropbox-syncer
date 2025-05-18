dev_image := "notion-dropbox-syncer-dev"
image := "bpsoos/notion-dropbox-syncer"

build:
    docker build -f Dockerfile.dev -t {{dev_image}} .

shell: build
    docker run --rm -it {{dev_image}} /bin/bash

upd: build
    SYNCER_DEV_IMAGE={{dev_image}} docker compose -f docker-compose-dev.yaml up -d

up: build
    SYNCER_DEV_IMAGE={{dev_image}} docker compose -f docker-compose-dev.yaml up

down:
    SYNCER_DEV_IMAGE={{dev_image}} docker compose -f docker-compose-dev.yaml down --remove-orphans

build-prod:
    docker buildx build \
        --platform linux/amd64,linux/arm64 \
        -t {{image}}:$(date +%s) \
        --push .
