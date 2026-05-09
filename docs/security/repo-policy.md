# Repository security settings

The code-side controls (SAST, dependency scanning, signed releases, pinned
actions, fuzzing) all live in this repo. Several **OpenSSF Scorecard**
checks also depend on GitHub repository settings — these are now
declared as files in the repo and applied automatically:

| Surface | File | Applied by |
|---|---|---|
| Repo metadata, features, merge policy, full label palette, fallback branch protection | [`.github/settings.yml`](../../.github/settings.yml) | `repository-settings/app` (Probot — install once) |
| `master` branch ruleset (PR + reviews + status checks + signed commits + linear history + Conventional-Commit regex) | [`.github/rulesets/branch-master.json`](../../.github/rulesets/branch-master.json) | `gh api repos/.../rulesets` import (one-time) |
| `v*` tag ruleset (only release-please[bot] and admins create/move/delete) | [`.github/rulesets/tags-v.json`](../../.github/rulesets/tags-v.json) | same as above |
| Path-based PR labels | [`.github/labeler.yml`](../../.github/labeler.yml) + [`labeler.yml`](../../.github/workflows/labeler.yml) workflow | `actions/labeler` on every PR |
| Dependabot auto-merge (patch/minor) | [`.github/workflows/auto-merge.yml`](../../.github/workflows/auto-merge.yml) | `dependabot/fetch-metadata` + `gh pr merge --auto` |
| Stale issue/PR triage | [`.github/workflows/stale.yml`](../../.github/workflows/stale.yml) | weekly cron |
| Auto-assign reviewers on PRs | [`.github/auto_assign.yml`](../../.github/auto_assign.yml) + [`auto-assign.yml`](../../.github/workflows/auto-assign.yml) workflow | every PR |
| Renovate (Go modules + GitHub Actions) | [`renovate.json`](../../renovate.json) | the Renovate GitHub App |
| Dependabot (Go modules + GitHub Actions, fallback) | [`.github/dependabot.yml`](../../.github/dependabot.yml) | the built-in Dependabot |

The bullet lists below describe **what those files encode** — they are
authoritative; this prose is documentation of the resulting policy.

## Settings encoded by the files above

### Branch protection on `master`

Encoded in [`.github/rulesets/branch-master.json`](../../.github/rulesets/branch-master.json)
(authoritative) and [`.github/settings.yml`](../../.github/settings.yml)
`branches:` block (fallback for first-fork bootstrap):

- [x] Require a pull request before merging
- [x] Require approvals (`1` minimum)
- [x] Dismiss stale pull request approvals when new commits are pushed
- [x] Require review from Code Owners
- [x] Require status checks to pass before merging
  - [x] Require branches to be up to date before merging
  - Required checks:
    - `required checks` (rollup job in `ci.yml` — depends on every other job in that workflow)
    - `golangci-lint` / `staticcheck` (from `lint.yml`)
    - `analyze` (from `codeql.yml`)
    - `gosec security scanner` (from `gosec.yml`)
    - `conventional commit title` (from `pr-title.yml`)
    - `scorecard analysis` (from `scorecard.yml`) *(if enabled — runs on push, not PR; only required when set up to run on PR via the workflow's `branch_protection_rule` trigger)*
- [x] Require conversation resolution before merging
- [x] Require signed commits
- [x] Require linear history
- [ ] Allow force pushes — keep **disabled**
- [ ] Allow deletions — keep **disabled**
- [x] Restrict who can push to matching branches → CODEOWNERS only

### Repository general

`Settings → General`:

- [x] Default branch: `master`
- [ ] Allow merge commits — disable in favour of squash
- [x] Allow squash merging — preferred (keeps `master` history linear)
- [ ] Allow rebase merging — disable (encourages linear history above)
- [x] Automatically delete head branches after merge

### Pull requests

`Settings → General → Pull requests`:

- [x] Always suggest updating pull request branches
- [x] Allow auto-merge

### Actions permissions

`Settings → Actions → General`:

- [x] Allow `amiwrpremium` actions and reusable workflows — and select
  third-parties listed below.
- [x] Approval required for outside-collaborator workflows (a.k.a.
  *Require approval for first-time contributors*)
- [x] Workflow permissions: **Read repository contents and packages
  permissions** (the workflows in this repo grant write at the job level
  where needed).

Third-party actions allowed (all pinned by SHA in
[`.github/workflows/`](../../.github/workflows/)):

- `actions/checkout`, `actions/setup-go`, `actions/upload-artifact`
- `codecov/codecov-action`
- `github/codeql-action/*`
- `golangci/golangci-lint-action`
- `googleapis/release-please-action`
- `securego/gosec`
- `codacy/codacy-analysis-cli-action`
- `ossf/scorecard-action`
- `google/osv-scanner-action`
- `step-security/harden-runner`
- `anchore/sbom-action`
- `sigstore/cosign-installer`
- `slsa-framework/slsa-github-generator`

### Security & analysis

`Settings → Code security and analysis`:

- [x] **Private vulnerability reporting** — enabled. Used by `SECURITY.md`
  and `CODE_OF_CONDUCT.md`.
- [x] **Dependency graph** — enabled.
- [x] **Dependabot alerts** — enabled.
- [x] **Dependabot security updates** — enabled.
- [x] **Code scanning** — CodeQL + gosec + Codacy + Scorecard +
  OSV-Scanner all upload SARIF here.
- [x] **Secret scanning** + **push protection** — enabled.

### Two-factor authentication

`https://github.com/orgs/<org>/settings/security`:

- [x] Require 2FA for everyone in the org. Required for OpenSSF Scorecard
  to mark `Maintained` and `Code-Review` as fully passing on org-owned repos.

## Why these settings

| Scorecard check | Settings that satisfy it |
|---|---|
| `Branch-Protection` | Branch protection rules above; required checks rolled up via `ci.yml`'s `required` job |
| `Code-Review` | "Require approvals" + "Dismiss stale" |
| `Token-Permissions` | Workflow permissions read-only by default + per-job write where needed |
| `Maintained` | Recent commits + active issue triage (out of repo's control) |
| `Pinned-Dependencies` | All actions pinned by SHA (already in workflows); Go modules pinned by `go.sum` |
| `Signed-Releases` | `release.yml` signs SBOM with cosign and emits SLSA L3 provenance |
| `Vulnerabilities` | Dependabot alerts + govulncheck CI job + OSV-Scanner |
| `SAST` | CodeQL + gosec + staticcheck |
| `Dangerous-Workflow` | The workflows that use `pull_request_target` (`labeler`, `auto-merge`, `auto-assign`, `pr-title`) only ever check out and execute the base branch — they never run untrusted PR-head code. Scorecard's heuristic is satisfied. |
| `Dependency-Update-Tool` | Renovate (primary) + Dependabot |
| `Security-Policy` | `SECURITY.md` + `SECURITY-INSIGHTS.yml` |
| `Fuzzing` | Native Go `Fuzz*` tests in `types.go`, `auth.go`, `errors.go`, `internal/jsonrpc` |
| `License` | `LICENSE` (MIT) |

## Verifying

After enabling the settings above, the public Scorecard at

```text
https://api.scorecard.dev/projects/github.com/amiwrpremium/go-derive/badge
```

should report 9.0+. Re-runs happen automatically on schedule
(`scorecard.yml`); to force a refresh, dispatch the workflow manually.
