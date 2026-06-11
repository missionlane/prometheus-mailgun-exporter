FROM golang:alpine@sha256:a6a091eac01ceac4b97496fe2957a49b6cdd83365337d5f46f6f73710424e805 AS builder
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

FROM alpine:latest@sha256:a2d49ea686c2adfe3c992e47dc3b5e7fa6e6b5055609400dc2acaeb241c829f4
RUN apk --update --no-cache add ca-certificates
ENTRYPOINT ["/prometheus-mailgun-exporter"]
EXPOSE 9616/tcp
USER nobody
COPY --from=builder /go/bin/prometheus-mailgun-exporter .
