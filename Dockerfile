FROM golang:alpine@sha256:ac09a5f469f307e5da71e766b0bd59c9c49ea460a528cc3e6686513d64a6f1fb AS builder
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

FROM alpine:latest@sha256:865b95f46d98cf867a156fe4a135ad3fe50d2056aa3f25ed31662dff6da4eb62
RUN apk --update --no-cache add ca-certificates
ENTRYPOINT ["/prometheus-mailgun-exporter"]
EXPOSE 9616/tcp
USER nobody
COPY --from=builder /go/bin/prometheus-mailgun-exporter .
