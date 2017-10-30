# goecho

This is a simple HTTP server built in Go that suited us good in testing our orchestration so far. It will answer to every path and method and will echo information about the request back to the client.

Currently it has the following endpoints:

- `/metrics` – prometheus metrics about the service
- `/version` – Version information
- `/*` – Echo handler (will echo information about the request back)

The echo handler understands one query parameter: "env". When set to any non-zero-length value the response will also contain the environment variables set for the server.

## Installation

For our purposes we typically use the Docker image, which can be found [on Docker Hub](https://hub.docker.com/r/xperimental/goecho/):

```bash
docker run --rm --interactive --tty --publish 8080:8080 xperimental/goecho:$TAG
```

If you have Go installed it is also very simple to build the binary yourself:

```bash
go get github.com/xperimental/goecho
```

If you want to build the Docker image instead, download the sources (`go get -d`) and run the provided script:

```bash
./build.sh
```

The image will be called "xperimental/goecho" and have a tag based on the git repository ("latest" if you have made any uncommitted changes).

## Usage

The program understands one parameter `-addr` which sets the address it listens on for connections. By default it listens on port `8080`:

```bash
goecho -addr :8080
```
