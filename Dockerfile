FROM python:3-alpine

WORKDIR /usr/src/app
COPY requirements.txt .

RUN apk add --no-cache --virtual .build-dependencies gcc musl-dev && \
    pip install --no-cache-dir -r requirements.txt && \
    apk del .build-dependencies

COPY run.py .
COPY init.sh .

ENTRYPOINT ["./init.sh"]
