FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:alpine AS builder
ARG GOPROXY
ARG GIT_REV
ARG VERSION
ARG TARGETARCH
ARG TARGETOS
COPY / /go/src/baetyl-cloud
WORKDIR /go/src/baetyl-cloud
ENV GOPROXY=${GOPROXY:-direct}
RUN --mount=type=cache,id=gomod,target=/go/pkg/mod \
    go mod download -x
RUN --mount=type=cache,id=gomod,target=/go/pkg/mod \
    GOARCH=${TARGETARCH} CGO_ENABLED=0 \
    go build -ldflags "-s -w -X github.com/baetyl/baetyl-go/v2/utils.REVISION=${GIT_REV} -X github.com/baetyl/baetyl-go/v2/utils.VERSION=${VERSION}" .

FROM scratch
COPY /scripts/native/templates /etc/templates
COPY --from=builder /go/src/baetyl-cloud/baetyl-cloud /bin/baetyl-cloud
ENTRYPOINT ["/bin/baetyl-cloud"]
