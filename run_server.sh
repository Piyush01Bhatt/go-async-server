#!/bin/sh

set -e

(
  cd "$(dirname "$0")"
  go build -o /tmp/go-async-server cmd/server.go
)

exec /tmp/go-async-server "$@"