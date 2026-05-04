# Threat model

This page documents what the SDK protects against, what it doesn't, and
which attacker scenarios were considered when designing the security
controls.

## In scope

The SDK is a client library that:

1. Produces signed JSON-RPC requests for Derive's matching engine.
2. Maintains a WebSocket connection for streaming data.
3. Holds the user's session-key private key in memory while running.

The threat model focuses on those three responsibilities.

## Out of scope

- **Mainnet smart-contract bugs.** The SDK signs structured data; it
  doesn't audit the contracts that consume those signatures. Any use
  of `pkg/contracts` should be preceded by independent contract review.
- **The user's wallet OS / hardware.** If the host is compromised, the
  session key is compromised — the SDK can't help.
- **Derive's API itself.** Server-side bugs are reported to the Derive
  team, not us.
- **Network-layer eavesdropping.** TLS is provided by Go's stdlib
  `net/http` and `gorilla/websocket`; we don't roll our own crypto.

## Assets

| Asset | Where it lives | Compromise impact |
|---|---|---|
| Session-key private key | `auth.LocalSigner.key` | attacker can sign on the user's behalf for the session-key's lifetime |
| Owner address | `auth.SessionKeySigner.owner` | not secret — public on-chain |
| Subaccount id | `methods.API.Subaccount` | not secret |
| Open WS connection | inside `transport.WSTransport` | hijack lets an attacker observe stream contents |

## Attackers and their capabilities

### A1 — Network attacker (passive + active MitM)

- Can read or modify network traffic between the SDK and Derive.
- Mitigation: TLS for both REST and WSS. The SDK never opens unencrypted
  channels. EIP-712 signed actions are tamper-evident at the engine.

### A2 — Compromised dependency

- Attacker controls one of the SDK's transitive Go modules.
- Mitigations:
  - All Go modules pinned by `go.sum` content hash.
  - Renovate auto-merges patch/minor; major bumps for security-critical
    deps (`go-ethereum`, `gorilla/websocket`) require human review per
    `renovate.json`.
  - `govulncheck` (Go-affecting CVEs) and `osv-scanner` (full graph)
    run on every push and weekly.
  - Dependabot security alerts at the GitHub level.
  - SBOM published with every release so downstream users can audit
    transitively.

### A3 — Compromised CI / supply chain

- Attacker compromises a GitHub Action used in CI.
- Mitigations:
  - Every action pinned to a 40-char commit SHA, not a tag.
  - `step-security/harden-runner` runs in egress-audit mode in critical
    workflows (Scorecard).
  - Workflows declare `permissions: read-all` at the top and grant
    write only at the per-job level.
  - SLSA Level 3 provenance on every release proves which workflow run
    produced the artifact.
  - cosign keyless OIDC signs the SBOM so consumers can verify origin.

### A4 — Malicious PR

- Outside contributor opens a PR with a backdoor.
- Mitigations:
  - Branch protection on `master` requires a reviewing approver
    (documented in [repo-policy.md](./repo-policy.md)).
  - "Require approval for first-time contributors" prevents PRs from
    auto-running the full workflow set on submission.
  - CodeQL `security-and-quality` query suite, gosec, and Codacy scan
    every PR diff.
  - The `Dangerous-Workflow` Scorecard check verifies no workflow uses
    `pull_request_target` with `${{ github.event.pull_request.head }}`.
  - Required signed commits on `master`.

### A5 — Compromised maintainer credential

- Attacker steals the maintainer's GitHub PAT or SSH key.
- Mitigations:
  - 2FA mandatory on the org (documented in repo-policy.md).
  - Branch protection prevents force-push and direct pushes to `master`.
  - Any release artifact must come through `release.yml`, which is
    triggered only by a release event — a stolen PAT can't unilaterally
    publish a backdoored binary because the cosign + SLSA flow is tied
    to the GitHub Actions OIDC identity.

### A6 — Application-level: replay / signature reuse

- Attacker captures a signed action and replays it.
- Mitigations:
  - Every signed action carries a `nonce` (strictly-increasing, sourced
    from `auth.NonceGen`) and an `expiry` Unix timestamp. The matching
    engine rejects duplicates and expired actions:
    - code 11017 `NonUniqueNonce`
    - code 11018 `InvalidNonceDate`
    - code 11005 `AlreadyExpired`

### A7 — Application-level: signing the wrong domain

- Attacker tricks a user into signing a request against the wrong
  network's domain.
- Mitigations:
  - The SDK pins `chainId` and `verifyingContract` per network in
    `internal/netconf`, not from user input.
  - The matching engine rejects with code 14024 `ChainIDMismatch`.

## What the user must do

The SDK can't carry the whole burden. Operators are responsible for:

1. **Storing the session-key private key safely.** Outside-process is
   ideal (KMS, HSM); inside the trading process is acceptable only when
   the host is hardened and credentials rotate frequently.
2. **Rotating session keys.** Register new ones, revoke old ones. Don't
   reuse a leaked key by hoping nobody used it yet.
3. **Running on the right network.** The SDK defaults to testnet for
   examples and integration tests; production deployments must
   explicitly opt into mainnet via `WithMainnet()`.
4. **Reporting vulnerabilities** through the GitHub private advisory
   flow, not as public issues.

## Reviewing this document

This threat model is a living artifact. PRs that add new components
(transports, signers, on-chain helpers) should also update this page.
