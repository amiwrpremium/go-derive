# go-derive documentation

Pick a topic. Most users start with [getting-started](./getting-started.md);
if you've used a derivatives SDK before, jump straight to
[architecture](./architecture.md).

## Concepts

| Doc | Covers |
|---|---|
| [architecture.md](./architecture.md) | how the packages compose; layering; why pkg/ + internal/ |
| [getting-started.md](./getting-started.md) | first runnable example; env setup; testnet vs mainnet |
| [transports.md](./transports.md) | REST vs WebSocket; when to pick each; cross-transport guarantees |
| [auth.md](./auth.md) | EIP-191 + EIP-712 deep dive; LocalSigner vs SessionKeySigner; `Action` struct |
| [subscriptions.md](./subscriptions.md) | typed channel descriptors; `Subscribe[T]` and `SubscribeFunc[T]` |
| [numerics.md](./numerics.md) | `Decimal`, fixed-point scaling, why no `float64` |
| [error-handling.md](./error-handling.md) | the 136 server codes; sentinel mapping; canonical messages |
| [rate-limiting.md](./rate-limiting.md) | token bucket, defaults, custom configs |
| [reconnection.md](./reconnection.md) | WS lifecycle; pumps; auto-reconnect; resubscribe on recover |

## Process

| Doc | Covers |
|---|---|
| [examples.md](./examples.md) | guided tour of `examples/` (91 programs) |
| [testing.md](./testing.md) | unit, fuzz, integration, coverage, `-race` |
| [ci.md](./ci.md) | every workflow under `.github/workflows/` |
| [release-process.md](./release-process.md) | release-please + SLSA + SBOM + cosign |

## Security

| Doc | Covers |
|---|---|
| [security/README.md](./security/README.md) | security index |
| [security/repo-policy.md](./security/repo-policy.md) | required GitHub repo settings |
| [security/threat-model.md](./security/threat-model.md) | what the SDK protects against, and what it doesn't |

For vulnerability disclosure see the root [SECURITY.md](../SECURITY.md).

## Project meta (root)

| File | Purpose |
|---|---|
| [../README.md](../README.md) | quickstart, badges, project files index |
| [../CONTRIBUTING.md](../CONTRIBUTING.md) | how to submit changes |
| [../CODE_OF_CONDUCT.md](../CODE_OF_CONDUCT.md) | Contributor Covenant 2.1 |
| [../SECURITY.md](../SECURITY.md) | vulnerability disclosure |
| [../SUPPORT.md](../SUPPORT.md) | where to ask which kind of question |
| [../GOVERNANCE.md](../GOVERNANCE.md) | how decisions get made |
| [../MAINTAINERS.md](../MAINTAINERS.md) | maintainer list and responsibilities |
| [../CHANGELOG.md](../CHANGELOG.md) | release history |
| [../SECURITY-INSIGHTS.yml](../SECURITY-INSIGHTS.yml) | OpenSSF security metadata |
