#!/bin/sh

PORT=${PORT:-5000}

echo "STARTING"
exec waitress-serve \
  --host="0.0.0.0" \
  --port="${PORT}" \
  run:app
