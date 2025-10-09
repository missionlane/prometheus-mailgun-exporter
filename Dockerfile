FROM golang:alpine@sha256:b6ed3fd0452c0e9bcdef5597f29cc1418f61672e9d3a2f55bf02e7222c014abd AS builder
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

FROM alpine:latest@sha256:4b7ce07002c69e8f3d704a9c5d6fd3053be500b7f1c69fc0d80990c2ad8dd412
RUN apk --update --no-cache add ca-certificates
ENTRYPOINT ["/prometheus-mailgun-exporter"]
EXPOSE 9616/tcp
USER nobody
COPY --from=builder /go/bin/prometheus-mailgun-exporter .
