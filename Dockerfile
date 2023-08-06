FROM --platform=$BUILDPLATFORM golang:1.16.2-alpine AS builder

ARG TARGETOS
ARG TARGETARCH

ENV GOOS=$TARGETOS
ENV GOARCH=$TARGETARCH

RUN apk add --no-cache make git bash

WORKDIR /build

COPY go.mod go.sum /build/

RUN go mod download
RUN go mod verify

COPY . /build/
RUN make build-binary

FROM --platform=$TARGETPLATFORM busybox

LABEL maintainer="Robert Jacob <xperimental@solidproject.de>"
EXPOSE 8080
USER nobody

COPY --from=builder /build/goecho /bin/goecho
ENTRYPOINT ["/bin/goecho"]
