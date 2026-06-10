ARG GO_VERSION=1.26.4
ARG ALPINE_VERSION=3.23

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS builder
ARG TARGETOS
ARG TARGETARCH
ARG VERSION=dev
ARG COMMIT=none
WORKDIR /src
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=bind,source=go.mod,target=go.mod \
    --mount=type=bind,source=go.sum,target=go.sum \
    go mod download -x

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=bind,target=. \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build \
    -ldflags "-X main.version=${VERSION} -X main.commit=${COMMIT}" \
    -o /usr/local/bin/app .

FROM alpine:${ALPINE_VERSION}
COPY --from=builder /usr/local/bin/app /usr/local/bin/app
CMD ["/usr/local/bin/app"]
