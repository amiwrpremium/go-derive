# Continuous integration

Twenty-four workflows under [`.github/workflows/`](../.github/workflows/).
All actions are SHA-pinned with the version as a comment; Renovate keeps
that pattern current. `pin-check.yml` enforces this on every push.

Repo policy (settings, labels, branch + tag rulesets) is also expressed
as files — see [security/repo-policy.md](./security/repo-policy.md) for
the table of declarative surfaces (`.github/settings.yml`,
`.github/rulesets/*.json`, `.github/labeler.yml`, `.github/auto_assign.yml`,
`renovate.json`, `.github/dependabot.yml`).

The same checks run locally before code reaches CI — see
[`lefthook.yml`](../lefthook.yml). `pre-commit` mirrors `fmt` + `vet` +
`lint`; `commit-msg` mirrors `pr-title`; `pre-push` mirrors the `required`
rollup (`build` + `vet` + `test -race` + `govulncheck`). Install with
`make hooks` after cloning. See [CONTRIBUTING.md](../CONTRIBUTING.md#git-hooks-lefthook).

## Per-trigger summary

| Trigger | Workflows |
|---|---|
| push to `master` | ci, lint, extra-lint, codeql, gosec, semgrep, codacy, scorecard, osv-scanner, trivy, gitleaks, trufflehog, license-check, pin-check |
| pull request | ci, lint, extra-lint, codeql, gosec, semgrep, codacy, osv-scanner, trivy, gitleaks, trufflehog, dependency-review, license-check, pin-check, pr-title, labeler, auto-assign |
| pull request from dependabot | auto-merge |
| weekly cron | codeql, codacy, scorecard, osv-scanner, trivy, gitleaks, trufflehog, semgrep, verify-release, stale |
| release published | release, verify-release |
| manual dispatch | integration, scorecard, release, stale |

## Workflow details

### `ci.yml` — the main check

Runs on every push and PR.

| Job | What it does |
|---|---|
| `fmt` | `gofmt -l` (fail if any file is not formatted) |
| `vet` | `go vet ./...` |
| `build` | `go build` matrix on Linux/macOS/Windows × Go 1.25 + 1.26 |
| `test` | `go test -race -coverprofile` matrix on Go 1.25 + 1.26; uploads to Codecov + Codacy |
| `govulncheck` | scans the dependency graph for known CVEs |
| `tidy` | `go mod tidy` and fail if it produced a diff |
| `examples` | `go build ./examples/...` (compiles all 91 programs) |
| `hooks-config` | `lefthook dump` — validates `lefthook.yml` parses |
| `required` | rollup that `needs:` every job above. Branch protection requires only this one check |

### `lint.yml`

| Job | Tool |
|---|---|
| `golangci-lint` | golangci-lint v2-schema config; pinned via Renovate |
| `staticcheck` | honnef.co/go/tools/cmd/staticcheck |

### `extra-lint.yml`

Markdown / YAML / shell / typo linting that doesn't fit golangci-lint.

| Job | Tool |
|---|---|
| markdownlint | DavidAnson/markdownlint-cli2 against `**/*.md` |
| yamllint | adrienverge/yamllint against `.github/**/*.yml`, root `.yml` |
| actionlint | rhysd/actionlint against `.github/workflows/*.yml` |
| editorconfig-checker | editorconfig-checker against the repo per `.editorconfig` |
| typos | crate-ci/typos against the whole tree per `.typos.toml` |

### `codeql.yml`

GitHub CodeQL with the `security-and-quality` query suite. Runs on push,
PR, and weekly Monday at 03:17 UTC.

### `gosec.yml`

[securego/gosec](https://github.com/securego/gosec) producing SARIF for
the GitHub code-scanning UI.

### `codacy.yml`

[Codacy Analysis CLI](https://github.com/codacy/codacy-analysis-cli)
SARIF upload. Runs gosec, staticcheck, revive, gofmt under the hood.

### `scorecard.yml`

[OpenSSF Scorecard](https://scorecard.dev) v2.4.3. Runs weekly + on push.
Publishes results to scorecard.dev so the public README badge works.
Uses [step-security/harden-runner](https://github.com/step-security/harden-runner)
in egress-audit mode.

### `osv-scanner.yml`

[google/osv-scanner](https://github.com/google/osv-scanner) cross-checks
the dependency graph against the OSV database (osv.dev). Complements
govulncheck by covering the full transitive graph including non-Go
components. Uploads SARIF on scheduled runs.

### `semgrep.yml`

[Semgrep](https://semgrep.dev/) SAST against `p/security-audit`,
`p/golang`, and `p/secrets` rule packs. SARIF upload to GitHub
code-scanning.

### `trivy.yml`

[Aqua Trivy](https://github.com/aquasecurity/trivy) in three modes:
filesystem-vulnerability, secret, and IaC/config scan. Three jobs, three
SARIF uploads.

### `gitleaks.yml`

[gitleaks](https://github.com/gitleaks/gitleaks) scans the full git
history for committed secrets. PRs and weekly cron.

### `trufflehog.yml`

[TruffleHog](https://github.com/trufflesecurity/trufflehog) entropy-based
secret scan. Configured `--only-verified --fail` so it only fires on
secrets it could authenticate against (low false-positive).

### `dependency-review.yml`

`actions/dependency-review-action` blocks PRs that introduce dependencies
with high-severity vulnerabilities or license-deny violations. PR-only.

### `license-check.yml`

[google/go-licenses](https://github.com/google/go-licenses) enforces a
manual allow-list (Apache-2.0, BSD-2/3-Clause, ISC, MIT, MPL-2.0,
Unlicense). Fails on any other SPDX identifier.

### `pin-check.yml`

Greps every `uses:` line in `.github/workflows/*.yml` and fails if any
action ref isn't a 40-char SHA pin. Fast — no external tools needed.

### `verify-release.yml`

After a release publishes (and weekly), re-verifies the cosign
signatures and the SLSA L3 provenance on the release artifacts.
Catches retroactive tampering of release assets.

### `release.yml`

Triggered when release-please publishes a tag. For each release:

1. **CycloneDX SBOM** via `anchore/sbom-action` — what dependencies the
   release was built with.
2. **SLSA Level 3 provenance** via the
   `slsa-framework/slsa-github-generator/.github/workflows/generator_generic_slsa3.yml`
   reusable workflow — proof of which builder produced the artifact.
3. **Cosign keyless OIDC** signature on the SBOM.

All three are uploaded as release assets so consumers can verify what
they're pulling.

### `release-please.yml`

Reads the conventional-commit history on `master`, opens (or updates) a
`chore: release X.Y.Z` PR with the changelog. Merging the PR cuts the
GitHub release and tags `vX.Y.Z`. The `release.yml` workflow then
attaches the SBOM and SLSA provenance.

### `pr-title.yml`

Validates that PR titles follow Conventional Commits using
[`amannn/action-semantic-pull-request`](https://github.com/amannn/action-semantic-pull-request).
With squash-merge enabled (the default), the PR title becomes the commit
title that release-please reads — a typo here would silently break the
changelog. The allowed type list mirrors `.github/release-please-config.json`.
The release-please bump PRs themselves are skipped via the
`autorelease: pending` / `autorelease: tagged` labels.

### `labeler.yml`

[`actions/labeler`](https://github.com/actions/labeler) v6 applies
`area:*` labels by changed-file path on every PR. Mapping lives in
[`.github/labeler.yml`](../.github/labeler.yml); the label palette
itself is declared in [`.github/settings.yml`](../.github/settings.yml).
Runs on `pull_request_target` so PRs from forks are labelled too — only
the base-branch workflow file runs, not PR-head code.

### `auto-merge.yml`

Auto-merges Dependabot PRs once the required checks are green:

| Update type | Action |
|---|---|
| `version-update:semver-patch` | enable `--auto --squash` |
| `version-update:semver-minor` | enable `--auto --squash` |
| `version-update:semver-major` | leave alone, comment a reminder |
| any PR labeled `review-required` | leave alone (manual review) |

Renovate handles its own auto-merge via `renovate.json`; this workflow
only runs for `dependabot[bot]` actor.

### `auto-assign.yml`

Adds the maintainer as reviewer + the author as assignee on every
human PR. Bots (`dependabot[bot]`, `renovate[bot]`,
`release-please[bot]`, `github-actions[bot]`) are skipped via the
job's `if:` condition. Config in
[`.github/auto_assign.yml`](../.github/auto_assign.yml).

### `stale.yml`

`actions/stale` v10 weekly cron. Issues stale after 60 days, closed
14 days later; PRs stale after 30 days, closed 14 days later. Exempt
labels: `security`, `release-blocker`, `pinned`, `good first issue`,
`help wanted`. 30 operations per run (well under the default 1000).

### `integration.yml` — manual dispatch only

Live testnet integration tests. Three jobs:

| Job | Runs when |
|---|---|
| public subset | always (no secrets needed) |
| private subset | when `DERIVE_SESSION_KEY` secret is present |
| live-order subset | only when dispatched with `run_live_orders: true` |

See [testing.md](./testing.md).

## Required secrets

| Secret | Used by |
|---|---|
| `CODECOV_TOKEN` | `ci.yml` |
| `CODACY_PROJECT_TOKEN` | `ci.yml`, `codacy.yml` (optional) |
| `DERIVE_SESSION_KEY` | `integration.yml` (optional) |
| `DERIVE_OWNER` | `integration.yml` (optional) |

## Required vars (non-secret)

| Var | Used by |
|---|---|
| `DERIVE_SUBACCOUNT` | `integration.yml` |
| `DERIVE_INSTRUMENT` | `integration.yml` |
| `DERIVE_BASE_ASSET` | `integration.yml` (live orders only) |

`GITHUB_TOKEN` is provided automatically by Actions. release-please uses it.

## Pinning policy

Every action is pinned to a 40-char commit SHA with the semver tag in a
trailing comment:

```yaml
- uses: actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd # v6
```

Renovate's `helpers:pinGitHubActionDigestsToSemver` preset (in
`renovate.json`) keeps both the SHA and the comment current when bumping.

`grep -E "uses: [^@]+@v[0-9]+( |$)" .github/workflows/*.yml` should
return nothing. The `pin-check.yml` workflow enforces this on every
push and PR.
