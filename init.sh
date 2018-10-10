#!/bin/sh

PORT=${PORT:-5000}
PORT_IPV4=${PORT_IPV4:-$PORT}
PORT_IPV6=${PORT_IPV6:-5001}

exec waitress-serve \
  --listen=0.0.0.0:${PORT_IPV4} \
  --listen=[::]:${PORT_IPV6} \
  run:app
