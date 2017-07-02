FROM golang:1.8.3-onbuild AS builder

FROM scratch
MAINTAINER Robert Jacob <robert.jacob@holidaycheck.com>
EXPOSE 8080

COPY --from=builder /go/bin/goecho /goecho
CMD ["/goecho"]
