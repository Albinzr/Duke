# Duke Server

Duke is a minimal HTTP server written in Go.  It exposes a small set of endpoints
and demonstrates how to structure a web service using MongoDB, Redis and a set of
standalone modules.  The project bundles login and product modules which are
loaded during startup.

This repository only contains the glue code required to spin up these modules.
The actual modules are pulled in as Go dependencies.

## Requirements

- Go 1.17 or newer
- MongoDB instance
- Redis instance

## Configuration

All configuration is controlled through environment variables.  Two example files
are provided:

- `local.env` – defaults used during local development
- `production.env` – values used in production containers

The following variables are recognised:

| Variable       | Description                                         |
| -------------- | --------------------------------------------------- |
| `PORT`         | Port for the HTTP server (default `1000`)           |
| `RETRY`        | Retry count for external services                   |
| `MONGO_URL`    | MongoDB connection string                           |
| `DATABASE_NAME`| Name of the MongoDB database                        |
| `SECRET_KEY`   | Secret used for signing JWT tokens                  |
| `AUD`          | JWT audience claim                                  |
| `ISS`          | JWT issuer claim                                    |

The server chooses which file to load based on the `-env` flag.  When omitted it
assumes `development` and therefore reads `local.env`.

```bash
# Run using the values from local.env
$ go run main.go -env=development
```

```bash
# Run using the production configuration
$ go run main.go -env=production
```

## Docker

A `duke.dockerfile` is provided for containerising the application.  The
included `docker-compose.yaml` spins up the server along with a Redis service and
expects a MongoDB instance to be available.

```bash
# Build the image
$ docker build -f duke.dockerfile -t duke-image .

# Start the stack
$ docker-compose up
```

For production deployments `docker-compose-production.yaml` can be used.  It
pulls a prebuilt image and starts the same set of supporting services.

## Endpoints

Currently only a single `/profile` endpoint is registered.  It is protected by a
JWT based middleware.  Requests lacking a valid token will receive a
`401 Unauthorized` response.

```
GET /profile
```

## Development

- `go vet ./...` – run static analysis
- `go test ./...` – execute the test suite

The repository contains a small example test in `src/server_test.go` which uses
`httptest` to validate middleware behaviour.

## License

This project is licensed under the MIT License.
