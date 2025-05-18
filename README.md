# notion dropbox syncer

## usage

docker:

```
docker run -d \
  --restart unless-stopped \
  --name syncer \
  -e NDSYNCER_DROPBOX__APP_KEY="" \
  -e NDSYNCER_DROPBOX__APP_SECRET="" \
  -e NDSYNCER_DROPBOX__REFRESH_TOKEN="" \
  -e NDSYNCER_DROPBOX__FOLDER_PATH="" \
  -e NDSYNCER_NOTION__API_KEY="" \
  -e NDSYNCER_NOTION__DATABASE_ID="" \
  bpsoos/notion-dropbox-syncer
```

docker compose:

use the docker-compose.yaml config file
