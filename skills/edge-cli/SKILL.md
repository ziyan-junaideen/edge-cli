---
name: edge-cli
description: Operate Edge Payments Technologies APIs with the `edge` command-line client. Use when an agent needs to authenticate, configure API profiles, inspect Edge merchants or customer/payment/account resources, request JSON:API relationships, or obtain machine-readable JSON from Edge environments.
---

# Edge CLI

Use the installed `edge` binary instead of constructing HTTP requests directly.

## Work safely

- Run `edge --help` and `edge <resource> --help` when command details are unclear.
- Confirm the active target with `edge profiles list` before reading sensitive or production data.
- Never print, persist, or request access tokens. Use `edge auth login`; credentials are stored in the operating-system keyring.
- Treat `--insecure-skip-verify` as local-development-only. Prefer a configured `--ca-cert`.
- Start with human-readable output. Add `--json` when structured output is needed for inspection or processing.
- Ask before changing profiles or authentication state unless the user explicitly requested that change.

## Configure and authenticate

```sh
edge profiles list
edge profiles set dev --api-url https://api.tryedge.test:4001/v2 --ca-cert /path/to/ca.pem
edge profiles use dev
edge auth login
```

Production is the default profile and uses `https://api.tryedge.io/v2`. Override a single invocation with `--profile <name>` or `--api-url <url>` when appropriate.

## Read resources

Resources use consistent `list` and `show <id>` operations:

```sh
edge merchants list
edge customers show <customer-id>
edge payment-demands show <payment-demand-id> --include payer,billing_address,payment_method
edge accounts list --json
```

Available resources include `merchants`, `customers`, `consumer-addresses`, `payment-demands`, `payment-subscriptions`, `payment-methods`, `refund-demands`, `account-alerts`, `accounts`, `memberships`, `merchant-punitive-actions`, `permissions`, `red-flags`, and `webhook-deliveries`.

Use `--include` for JSON:API relationships; repeat it or comma-separate values. `--preload` is an alias. Allowed relationships differ by resource, so consult command help when uncertain.

With `--json`, expect the full JSON:API document, including `data`, `included`, `links`, and `meta`; do not assume the output is a bare array or object.

## Replay webhooks locally

Replay an existing delivery to an endpoint reachable from the current machine:

```sh
edge webhook-deliveries replay <delivery-id> --to http://localhost:4000/webhooks/edge
```

This fetches the delivery's event and subscription secret, creates a fresh v3 signature, and posts the reconstructed webhook body. Use `--dry-run` to inspect the signed request without sending it. Use `--delivery-version v1` or `v2` only when the integration explicitly requires a legacy payload and signature.
