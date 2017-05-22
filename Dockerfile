FROM python:3-alpine
RUN apk add --no-cache --virtual .build-dependencies gcc musl-dev && \
    pip install Flask==0.12.2	gevent==1.2.1 && \
    apk del .build-dependencies

ADD . /usr/src/app
WORKDIR /usr/src/app
EXPOSE 5000

ENTRYPOINT ["python", "run.py"]
