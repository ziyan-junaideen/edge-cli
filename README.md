# Edge CLI

Go-based command line client for the Edge Payment Technologies JSON:API API.

## Prerequisite

Install Go 1.22 or newer. The current local machine did not have `go` on `PATH` when this project was scaffolded.

## Development

```sh
go mod tidy
go test ./...
go run ./cmd/edge --help
```

## Profiles

Production is the default profile and targets `https://api.tryedge.io/v2`.

Local development can use the Phoenix dev API host:

```sh
edge profiles set dev --api-url https://api.tryedge.test:4001/v2 --ca-cert /Volumes/Dev/Work/Edge/edge/ept/priv/cert/_wildcard.tryedge.test+3.pem
edge profiles use dev
edge auth login
edge merchants list
```

For a custom local host:

```sh
edge profiles set local-dashboard --api-url https://dashboard.tryedge.test:4001/v2 --ca-cert /Volumes/Dev/Work/Edge/edge/ept/priv/cert/_wildcard.tryedge.test+3.pem
```

Use `--insecure-skip-verify` only for local development endpoints.
