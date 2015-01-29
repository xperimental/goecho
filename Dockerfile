FROM scratch
MAINTAINER Robert Jacob <robert.jacob@holidaycheck.com>
EXPOSE 8080

ADD goecho /goecho
CMD ["/goecho"]
