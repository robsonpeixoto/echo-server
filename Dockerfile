FROM python:3-alpine

RUN apk add --no-cache --virtual .build-dependencies gcc musl-dev && \
    pip install --no-cache-dir Flask==1.0.2 gunicorn==19.9.0 gevent==1.3.4 && \
    apk del .build-dependencies

WORKDIR /usr/src/app
COPY run.py .
COPY init.sh .

ENTRYPOINT ["./init.sh"]
