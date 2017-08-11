FROM golang:1.8.3 AS builder

ARG PACKAGE

RUN mkdir -p /go/src/${PACKAGE}
WORKDIR /go/src/${PACKAGE}

ARG VERSION

ENV LD_FLAGS="-w -X main.Version=${VERSION}"
ENV CGO_ENABLED=0

COPY . /go/src/${PACKAGE}
RUN go get -d -v .
RUN go install -a -v -tags netgo -ldflags "${LD_FLAGS}" .

FROM scratch
MAINTAINER Robert Jacob <robert.jacob@holidaycheck.com>
EXPOSE 8080

COPY --from=builder /go/bin/goecho /goecho
CMD ["/goecho"]
