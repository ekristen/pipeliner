#!/bin/bash

set -euo pipefail

cd ui || exit 1

if [ -z "$(which npm)" ]; then
    docker run -i --init --rm -e CHOKIDAR_USEPOLLING=true --network=host -v "$PWD":/src --workdir /src node:12 npm ci
    docker run -i --init --rm -e CHOKIDAR_USEPOLLING=true --network=host -v "$PWD":/src --workdir /src node:12 npm run build
else
    npm ci
    npm run build
fi
