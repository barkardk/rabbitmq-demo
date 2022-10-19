FROM --platform=${BUILDPLATFORM} golang:1.16-alpine as base
WORKDIR /src
COPY go.mod /src/
COPY it /src/it
RUN go mod download > /dev/null

FROM base AS build
ARG GOARCH
ARG GOOS
ARG CGO_ENABLED=0
ARG BUILD_TIME=$(date)
ARG VERSION
RUN --mount=type=cache,target=/root/.cache/go-build GO111MODULE=on CGO_ENABLED=$CGO_ENABLED GOARCH=$GOARCH GOOS=$GOOS \
    go build -o target/mq_test -tags $VERSION -ldflags "-s -w -X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME" ./it

FROM alpine
COPY --from=build /src/target/mq_test .
RUN uname -a
COPY passwd.minimal /etc/passwd
RUN chown rabbit mq_test
USER rabbit
CMD ["sh", "-c", "tail -f /dev/null"]