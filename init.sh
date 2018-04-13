#!/bin/sh

PORT_IPV4=${PORT_IPV4:-5000}
PORT_IPV6=${PORT_IPV6:-5001}

exec gunicorn -k gevent \
  -b "0.0.0.0:${PORT_IPV4}" \
  -b "[::]:${PORT_IPV6}" \
  run:app

