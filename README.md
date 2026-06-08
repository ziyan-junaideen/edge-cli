# Edge CLI

Go-based command line client for the Edge Payment Technologies JSON:API API.

## Prerequisite

Install project tools with mise:

```sh
mise install
```

## Installation

```sh
brew tap ziyan-junaideen/tap
brew install edge-cli
edge --help
```

## Development

```sh
mise exec -- go mod tidy
mise exec -- go test ./...
mise exec -- go run ./cmd/edge --help
```

## Profiles

Production is the default profile and targets `https://api.tryedge.io/v2`.

Local development can use the Phoenix dev API host:

```sh
edge profiles set dev --api-url https://api.tryedge.test:4001/v2 --ca-cert /Volumes/Dev/Work/Edge/edge/ept/priv/cert/_wildcard.tryedge.test+3.pem
edge profiles use dev
edge auth login
edge merchants list
edge customers show <customer-id> --include addresses --json
```

For a custom local host:

```sh
edge profiles set local-dashboard --api-url https://dashboard.tryedge.test:4001/v2 --ca-cert /Volumes/Dev/Work/Edge/edge/ept/priv/cert/_wildcard.tryedge.test+3.pem
```

Use `--insecure-skip-verify` only for local development endpoints.

## Resources

```sh
edge merchants list
edge merchants show <merchant-id>

edge customers list
edge customers show <customer-id> --include addresses
edge customers show <customer-id> --preload addresses --json

edge consumer-addresses list
edge consumer-addresses show <address-id> --include customer

edge payment-demands list
edge payment-demands show <payment-demand-id> --include payer,billing_address,payment_method

edge payment-subscriptions list
edge payment-subscriptions show <payment-subscription-id> --include payer,payment_method

edge payment-methods list
edge payment-methods show <payment-method-id> --include customer,address

edge refund-demands list
edge refund-demands show <refund-demand-id> --include payment_demand

edge account-alerts list
edge account-alerts show <account-alert-id> --include red_flag

edge accounts list
edge accounts show <account-id> --include memberships

edge memberships list
edge memberships show <membership-id> --include account,permissions

edge merchant-punitive-actions list
edge merchant-punitive-actions show <action-id> --include merchant,red_flag

edge permissions list
edge permissions show <permission-id> --include merchant_tokens

edge red-flags list
edge red-flags show <red-flag-id> --include merchant
```

`--preload` is an alias for JSON:API `--include`. JSON output returns the full JSON:API document, including `included`, `links`, and `meta`.
