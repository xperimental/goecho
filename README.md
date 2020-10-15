# goecho

This is a simple HTTP server built in Go that suited us good in testing our orchestration so far. It will answer to every path and method and will echo information about the request back to the client.

Currently it has the following endpoints:

- `/metrics` – prometheus metrics about the service
- `/version` – Version information
- `/_ready` - Readiness check
- `/*` – Echo handler (will echo information about the request back)

The echo handler understands one query parameter: "env". When set to any non-zero-length value the response will also contain the environment variables set for the server.

## Installation

For our purposes we typically use the Docker image, which can be found [on Docker Hub](https://hub.docker.com/r/xperimental/goecho/):

```bash
docker run --rm --interactive --tty --publish 8080:8080 xperimental/goecho:$TAG
```

If you have Go installed it is also very simple to build the binary yourself:

```bash
git clone https://github.com/xperimental/goecho.git
cd goecho
go build
```

If you want to build the Docker image instead, download the sources and run the provided script:

```bash
git clone https://github.com/xperimental/goecho.git
cd goecho
./build.sh
```

The image will be called "xperimental/goecho" and have a tag based on the git repository ("latest" if you have made any uncommitted changes).

## Usage

Some configuration options can be either set by using command-line options or environment variables. If both are set, environment variables will override the command-line options:

|      Option       | Environment Variable | Default |                           Description                            |
| :---------------- | :------------------- | :------ | :--------------------------------------------------------------- |
| `-addr`           | `LISTEN_ADDR`        | `:8080` | Address and port to listen on.                                   |
| `-allow-env`      | `ALLOW_ENV`          |         | Allow retrieval of environment variables.                        |
| `-graceful-delay` | `GRACEFUL_DELAY`     | `2s`    | Delay between receiving a shutdown signal and starting shutdown. |
| `-tls-cert`       | `TLS_CERT`           |         | Path to TLS certificate file.                                    |
| `-tls-key`        | `TLS_KEY`            |         | Path to TLS key file.                                            |

### TLS Support

goecho will by default only run a HTTP server, but if you provide a TLS certificate and key either using command-line flags or environment variables it will instead provide a HTTPS server.

When TLS is enabled and TLS requests reach the server the echo handler will also return information about the TLS version and the hostname used for the connection.
