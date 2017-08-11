FROM golang:1.8.3 AS builder

RUN mkdir -p /go/src/app
WORKDIR /go/src/app

ENV LD_FLAGS="-w -X main.Version=${version}"
ENV CGO_ENABLED=0

COPY . /go/src/app
RUN go-wrapper download
RUN go-wrapper install

FROM scratch
MAINTAINER Robert Jacob <robert.jacob@holidaycheck.com>
EXPOSE 8080

COPY --from=builder /go/bin/goecho /goecho
CMD ["/goecho"]
