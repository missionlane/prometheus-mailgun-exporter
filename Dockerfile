FROM golang:alpine@sha256:26111811bc967321e7b6f852e914d14bede324cd1accb7f81811929a6a57fea9 AS builder
WORKDIR /go/src/app

ARG BUILD_DATE
ARG BUILD_USER
ARG GIT_BRANCH
ARG GIT_REVISION
ARG GO111MODULE
ARG VERSION

COPY . .
RUN apk --update --no-cache add git && \
        go mod tidy && \
        go install \
            -ldflags "-X github.com/prometheus/common/version.BuildDate=${BUILD_DATE} \
                        -X github.com/prometheus/common/version.BuildUser=${BUILD_USER} \
                        -X github.com/prometheus/common/version.Branch=${GIT_BRANCH} \
                        -X github.com/prometheus/common/version.Revision=${GIT_REVISION} \
                        -X github.com/prometheus/common/version.Version=${VERSION}"

FROM alpine:latest@sha256:51183f2cfa6320055da30872f211093f9ff1d3cf06f39a0bdb212314c5dc7375
RUN apk --update --no-cache add ca-certificates
ENTRYPOINT ["/prometheus-mailgun-exporter"]
EXPOSE 9616/tcp
USER nobody
COPY --from=builder /go/bin/prometheus-mailgun-exporter .
