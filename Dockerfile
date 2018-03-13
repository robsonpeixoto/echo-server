FROM python:3-alpine

RUN apk add --no-cache --virtual .build-dependencies gcc musl-dev && \
    pip install --no-cache-dir Flask==0.12.2 gunicorn==19.7.1 gevent==1.2.2 && \
    apk del .build-dependencies

WORKDIR /usr/src/app
ADD run.py .
EXPOSE 5000 5001

ENTRYPOINT ["gunicorn", "-k", "gevent", "-b", "0.0.0.0:5000", "-b", "[::]:5001", "run:app"]
