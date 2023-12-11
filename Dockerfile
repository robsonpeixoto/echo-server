FROM golang:1.21-alpine3.18 AS builder
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN --mount=type=cache,target=/root/.cache/go-build go build -o  /usr/local/bin/app ./...

FROM alpine:3.19
COPY --from=builder /usr/local/bin/app /usr/local/bin/app
CMD ["app"]
