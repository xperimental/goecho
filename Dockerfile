FROM golang:1.12.6 AS builder

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y upx

WORKDIR /build

COPY go.mod go.sum /build/

RUN go mod download
RUN go mod verify

COPY . /build/

ARG VERSION

ENV LD_FLAGS="-w -X main.Version=${VERSION}"
ENV CGO_ENABLED=0

RUN go install -v -tags netgo -ldflags "${LD_FLAGS}" .
RUN upx -9 /go/bin/goecho

FROM busybox

LABEL maintainer="Robert Jacob <xperimental@solidproject.de>"
EXPOSE 8080
USER nobody

COPY --from=builder /go/bin/goecho /bin/goecho
ENTRYPOINT ["/bin/goecho"]
