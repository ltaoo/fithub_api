#!/bin/sh
set -eu
cd "$(dirname "$0")"
pnpm --dir h5 run build
pnpm --dir admin run build
go build -o bin/fithub-server ./cmd/server
