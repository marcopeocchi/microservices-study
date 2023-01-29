#!/usr/bin/env bash

export CGO_ENABLED=1

sudo docker buildx build \
 --push \
 --file Dockerfile \
 --platform linux/amd64 \
 --tag marcobaobao/fuu:latest .

echo "Done!"
exit 0
