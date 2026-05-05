# go-derive

[![CI](https://github.com/amiwrpremium/go-derive/actions/workflows/ci.yml/badge.svg)](https://github.com/amiwrpremium/go-derive/actions/workflows/ci.yml)
[![Lint](https://github.com/amiwrpremium/go-derive/actions/workflows/lint.yml/badge.svg)](https://github.com/amiwrpremium/go-derive/actions/workflows/lint.yml)
[![CodeQL](https://github.com/amiwrpremium/go-derive/actions/workflows/codeql.yml/badge.svg)](https://github.com/amiwrpremium/go-derive/actions/workflows/codeql.yml)
[![codecov](https://codecov.io/gh/amiwrpremium/go-derive/branch/master/graph/badge.svg)](https://codecov.io/gh/amiwrpremium/go-derive)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/1cfcd38cd2b043a1bbba2bdc7b188026)](https://app.codacy.com/gh/amiwrpremium/go-derive/dashboard)
[![Codacy coverage](https://app.codacy.com/project/badge/Coverage/1cfcd38cd2b043a1bbba2bdc7b188026)](https://app.codacy.com/gh/amiwrpremium/go-derive/dashboard)
[![Go Reference](https://pkg.go.dev/badge/github.com/amiwrpremium/go-derive.svg)](https://pkg.go.dev/github.com/amiwrpremium/go-derive)
[![Go Report Card](https://goreportcard.com/badge/github.com/amiwrpremium/go-derive)](https://goreportcard.com/report/github.com/amiwrpremium/go-derive)
[![Go Version](https://img.shields.io/github/go-mod/go-version/amiwrpremium/go-derive)](https://github.com/amiwrpremium/go-derive/blob/master/go.mod)
[![Release](https://img.shields.io/github/v/release/amiwrpremium/go-derive?include_prereleases&sort=semver)](https://github.com/amiwrpremium/go-derive/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Renovate enabled](https://img.shields.io/badge/renovate-enabled-brightgreen.svg)](https://renovatebot.com/)
[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg)](https://conventionalcommits.org)
[![release-please](https://img.shields.io/badge/release-please-blue)](https://github.com/googleapis/release-please)
[![govulncheck](https://img.shields.io/badge/security-govulncheck-success)](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck)
[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/amiwrpremium/go-derive/badge)](https://scorecard.dev/viewer/?uri=github.com/amiwrpremium/go-derive)
[![OSV-Scanner](https://github.com/amiwrpremium/go-derive/actions/workflows/osv-scanner.yml/badge.svg)](https://github.com/amiwrpremium/go-derive/actions/workflows/osv-scanner.yml)
[![SLSA 3](https://slsa.dev/images/gh-badge-level3.svg)](https://slsa.dev)
[![Maintained](https://img.shields.io/badge/maintained-yes-brightgreen.svg)](https://github.com/amiwrpremium/go-derive/commits/master)
[![Last commit](https://img.shields.io/github/last-commit/amiwrpremium/go-derive/master)](https://github.com/amiwrpremium/go-derive/commits/master)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](./CONTRIBUTING.md)
<!-- After enrolling at https://www.bestpractices.dev/, replace BADGE_ID below.
[![OpenSSF Best Practices](https://www.bestpractices.dev/projects/BADGE_ID/badge)](https://www.bestpractices.dev/projects/BADGE_ID)
-->

A Go SDK for the [Derive](https://docs.derive.xyz/) exchange (formerly Lyra) — a layer-2 derivatives venue with perps, options, and spot.

Covers REST (public + private), WebSocket (public + private + subscriptions), and EIP-712 order signing with session keys.

## Status

`v0.1.0-dev` — under active development. API may change before `v1.0.0`.

## Install

```bash
go get github.com/amiwrpremium/go-derive
```

Requires Go 1.25+.

## Quick start

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/amiwrpremium/go-derive/pkg/auth"
    "github.com/amiwrpremium/go-derive/pkg/derive"
    "github.com/amiwrpremium/go-derive/pkg/enums"
)

func main() {
    signer, err := auth.NewLocalSigner(os.Getenv("DERIVE_SESSION_KEY"))
    if err != nil {
        log.Fatal(err)
    }

    c, err := derive.NewClient(
        derive.WithMainnet(),
        derive.WithSigner(signer),
        derive.WithSubaccount(123),
    )
    if err != nil {
        log.Fatal(err)
    }
    defer c.Close()

    ctx := context.Background()

    instruments, err := c.REST.GetInstruments(ctx, "BTC", enums.InstrumentTypePerp)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(len(instruments), "BTC perps")
}
```

See [`examples/`](./examples/) for more.

## Architecture

```text
pkg/derive               top-level facade (REST + WS)
pkg/rest, pkg/ws         dedicated clients
pkg/channels             typed WebSocket subscriptions
pkg/auth                 EIP-712 signing, session keys
pkg/types, pkg/enums     domain types, named-string enums
pkg/errors               sentinel errors + APIError

internal/jsonrpc         JSON-RPC 2.0 framing
internal/transport       HTTP + WS transports (shared interface)
internal/methods         RPC method definitions (shared by REST + WS)
internal/netconf         endpoints + EIP-712 domains per network
internal/codec           decimal/u256/address encoding
internal/retry           exponential backoff
```

The Derive API is JSON-RPC 2.0 over both HTTP and WebSocket — same method names, same params. The SDK reflects that: a single `Transport` interface backs both `pkg/rest` and `pkg/ws`, so each method is defined once.

See [`docs/architecture.md`](./docs/architecture.md) for the full design.

## Documentation

The full doc set lives under [`docs/`](./docs/):

| Topic | |
|---|---|
| Concepts | [getting-started](./docs/getting-started.md) · [architecture](./docs/architecture.md) · [transports](./docs/transports.md) · [auth](./docs/auth.md) · [subscriptions](./docs/subscriptions.md) · [numerics](./docs/numerics.md) · [error handling](./docs/error-handling.md) · [rate limiting](./docs/rate-limiting.md) · [reconnection](./docs/reconnection.md) |
| Process | [examples](./docs/examples.md) · [testing](./docs/testing.md) · [ci](./docs/ci.md) · [release process](./docs/release-process.md) |
| Security | [security index](./docs/security/README.md) · [repo policy](./docs/security/repo-policy.md) · [threat model](./docs/security/threat-model.md) |

## Continuous integration

Every push and pull request runs:

| Check | Tool | Workflow |
|---|---|---|
| Format | `gofmt -l` | [ci.yml](.github/workflows/ci.yml) |
| Vet | `go vet` | [ci.yml](.github/workflows/ci.yml) |
| Build | `go build` on Linux/macOS/Windows × Go 1.25/1.26 | [ci.yml](.github/workflows/ci.yml) |
| Tests | `go test -race -coverprofile` | [ci.yml](.github/workflows/ci.yml) |
| Mod tidy | `go mod tidy` diff check | [ci.yml](.github/workflows/ci.yml) |
| Vulnerabilities | `govulncheck` | [ci.yml](.github/workflows/ci.yml) |
| Linters | `golangci-lint`, `staticcheck` | [lint.yml](.github/workflows/lint.yml) |
| Extra linters | markdownlint, yamllint, actionlint, editorconfig-checker, typos | [extra-lint.yml](.github/workflows/extra-lint.yml) |
| Security (SAST) | CodeQL, gosec, Semgrep (security-audit + golang + secrets), Codacy | [codeql.yml](.github/workflows/codeql.yml), [gosec.yml](.github/workflows/gosec.yml), [semgrep.yml](.github/workflows/semgrep.yml), [codacy.yml](.github/workflows/codacy.yml) |
| Filesystem / IaC scan | Trivy (filesystem + secret + config) | [trivy.yml](.github/workflows/trivy.yml) |
| Secret scanning | Gitleaks (git history) + TruffleHog (entropy, verified-only) | [gitleaks.yml](.github/workflows/gitleaks.yml), [trufflehog.yml](.github/workflows/trufflehog.yml) |
| Dependency review | PR-time license + vulnerability gate | [dependency-review.yml](.github/workflows/dependency-review.yml) |
| License compliance | `go-licenses` allow-list (Apache-2.0, BSD, ISC, MIT, MPL-2.0, Unlicense) | [license-check.yml](.github/workflows/license-check.yml) |
| Action SHA pinning | enforces every `uses:` is a 40-char SHA | [pin-check.yml](.github/workflows/pin-check.yml) |
| Coverage | Codecov + Codacy upload | [ci.yml](.github/workflows/ci.yml) |
| Releases | release-please (Conventional Commits → CHANGELOG + tag) | [release-please.yml](.github/workflows/release-please.yml) |
| Dependencies | Renovate (primary), Dependabot (fallback) | [renovate.json](renovate.json), [dependabot.yml](.github/dependabot.yml) |
| Integration | live testnet smoke tests, manual dispatch only | [integration.yml](.github/workflows/integration.yml) |
| OpenSSF Scorecard | weekly + on push, publishes public score | [scorecard.yml](.github/workflows/scorecard.yml) |
| OSV-Scanner | weekly + on push/PR, transitive dep CVE scan | [osv-scanner.yml](.github/workflows/osv-scanner.yml) |
| SLSA + SBOM | runs on every published release | [release.yml](.github/workflows/release.yml) |
| Post-release re-verify | cosign + slsa-verifier, weekly + on release | [verify-release.yml](.github/workflows/verify-release.yml) |

All workflows additionally run `step-security/harden-runner` in `audit`
mode for egress monitoring.

### Required repository secrets

| Secret | Used by | Required? |
|---|---|---|
| `CODECOV_TOKEN` | Codecov upload in `ci.yml` | Yes for private repos; public repos can omit |
| `CODACY_PROJECT_TOKEN` | Codacy coverage upload in `ci.yml` | Optional — coverage upload silently skipped if missing |
| `RELEASE_PLEASE_TOKEN` | release-please uses this PAT to publish releases that auto-trigger `release.yml` | Recommended — if missing, falls back to `GITHUB_TOKEN`, but releases won't auto-fire `release.yml` (artefacts need manual `gh workflow run release.yml -f tag=vX.Y.Z`) |

`GITHUB_TOKEN` is provided by Actions automatically. `RELEASE_PLEASE_TOKEN`
should be a fine-grained PAT scoped to this repo with **Contents: write**
and **Pull requests: write** permissions — needed because GitHub's
anti-loop protection doesn't fire `release.yml` on releases published by
`GITHUB_TOKEN`.

## Security

This project follows the [OpenSSF best practices](https://openssf.org/) and
publishes a public Scorecard at
[scorecard.dev](https://scorecard.dev/viewer/?uri=github.com/amiwrpremium/go-derive).

| What | Where |
|---|---|
| Vulnerability disclosure | [SECURITY.md](./SECURITY.md) — uses GitHub private advisories |
| Code of conduct | [CODE_OF_CONDUCT.md](./CODE_OF_CONDUCT.md) (Contributor Covenant 2.1) |
| Security metadata | [SECURITY-INSIGHTS.yml](./SECURITY-INSIGHTS.yml) (OpenSSF spec 1.0.0) |
| Required repo settings | [docs/security/repo-policy.md](./docs/security/repo-policy.md) |
| Static analysis (SAST) | CodeQL, gosec, Semgrep, staticcheck, Codacy |
| Filesystem / IaC scanning | Trivy (filesystem + secret + config modes) |
| Secret scanning | Gitleaks (history) + TruffleHog (verified-only) |
| Dependency scanning | govulncheck, OSV-Scanner, Trivy filesystem, dependency-review |
| Dependency updates | Renovate (primary) + Dependabot (fallback) |
| License compliance | `go-licenses` allow-list enforced in CI |
| Egress audit | `step-security/harden-runner` on every workflow (audit mode) |
| Action pinning enforcement | `pin-check` workflow rejects unpinned `uses:` lines |
| Release integrity | SLSA Level 3 provenance + CycloneDX & SPDX SBOMs + license inventory, all cosign-signed — [release.yml](.github/workflows/release.yml) |
| Post-release verification | cosign signatures + SLSA provenance re-checked weekly + on every release — [verify-release.yml](.github/workflows/verify-release.yml) |
| Fuzzing | Go-native `Fuzz*` tests in `pkg/types`, `pkg/auth`, `pkg/errors`, `internal/jsonrpc` |
| Pinned actions | every action pinned by SHA with the version as a comment |

To report a vulnerability or code-of-conduct violation, use [GitHub's
private vulnerability reporting](https://github.com/amiwrpremium/go-derive/security/advisories/new).
The same channel handles both so reports go through one triage pipeline.

## Running integration tests

Live-network tests live under [`test/`](./test/) and are gated by the
`integration` build tag, so the default `go test ./...` is unaffected.

```bash
# Public-only subset (no creds needed) against testnet.
make test-integration

# All tests except live order placement.
DERIVE_SESSION_KEY=0x... DERIVE_SUBACCOUNT=123 \
  go test -tags=integration -count=1 ./test/...

# Add live order placement (testnet only — never against mainnet).
DERIVE_RUN_LIVE_ORDERS=1 DERIVE_BASE_ASSET=0x... \
  DERIVE_SESSION_KEY=0x... DERIVE_SUBACCOUNT=123 \
  go test -tags=integration -count=1 -run='^TestPrivate_PlaceAndCancel' ./test/...
```

See [test/README.md](./test/README.md) for the full env-var list and what
each subset covers.

## Project files

| File | Purpose |
|---|---|
| [CONTRIBUTING.md](./CONTRIBUTING.md) | how to submit changes; Conventional Commits |
| [CODE_OF_CONDUCT.md](./CODE_OF_CONDUCT.md) | Contributor Covenant 2.1 |
| [SECURITY.md](./SECURITY.md) | vulnerability disclosure |
| [SUPPORT.md](./SUPPORT.md) | where to ask which kind of question |
| [GOVERNANCE.md](./GOVERNANCE.md) | how decisions get made |
| [MAINTAINERS.md](./MAINTAINERS.md) | who reviews and merges |
| [CHANGELOG.md](./CHANGELOG.md) | every release, generated by release-please |
| [AUTHORS](./AUTHORS) | contributors in chronological order |
| [SECURITY-INSIGHTS.yml](./SECURITY-INSIGHTS.yml) | OpenSSF security metadata |
| [.github/settings.yml](./.github/settings.yml) | declarative repo settings + label palette (Probot Settings) |
| [.github/rulesets/](./.github/rulesets/) | branch + tag rulesets, importable via `gh api` |

## Contributing

Commits must follow [Conventional Commits](https://www.conventionalcommits.org/) so release-please can derive the next version and update [CHANGELOG.md](./CHANGELOG.md). See [CONTRIBUTING.md](./CONTRIBUTING.md).

## License

[MIT](./LICENSE).

## Acknowledgements

This SDK exists thanks to:

- **[Derive](https://docs.derive.xyz/)** (formerly Lyra) — the public REST + WebSocket reference underpins every method and sentinel in this module.
- **Upstream Go libraries** the SDK builds on: [`go-ethereum`](https://github.com/ethereum/go-ethereum) for crypto and EIP-712 hashing, [`gorilla/websocket`](https://github.com/gorilla/websocket) for the WS transport, [`shopspring/decimal`](https://github.com/shopspring/decimal) for fixed-point arithmetic, [`stretchr/testify`](https://github.com/stretchr/testify) for the test helpers.
- **[OpenSSF](https://openssf.org/)** — the [Scorecard](https://scorecard.dev), [SLSA](https://slsa.dev), [Security Insights](https://github.com/ossf/security-insights-spec), and [best-practices](https://www.bestpractices.dev/) projects shaped most of the security plumbing in this repo.

## Disclaimer

This software is provided **"as is"**, without warranty of any kind, express or implied, including but not limited to the warranties of merchantability, fitness for a particular purpose, and non-infringement (see [LICENSE](./LICENSE) for the full terms).

This is an **independent, unofficial** project. It is **not** affiliated with, endorsed by, or sponsored by Derive, Lyra Finance, the Lyra DAO, the Ethereum Foundation, or any other organisation or person. All product names, logos, and brands referenced are the property of their respective owners; their use here is for identification purposes only.

Trading derivatives carries financial risk. Nothing in this repository is financial advice. You are solely responsible for any orders submitted, keys generated or stored, and integrations built on top of this code. Use at your own risk.
