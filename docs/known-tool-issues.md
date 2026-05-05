# Known tool issues

This document catalogues code-scanning alerts on this repository that
remain **open** because the underlying issue is in the analysing tool
itself, not in our code. We do not dismiss or suppress these alerts —
they stay visible so that:

- Auditors can see the reasoning rather than encountering silent
  suppressions.
- Future contributors are not surprised by an unfixed alert.
- Each entry is a reminder to check whether the upstream tool has
  fixed its bug, in which case the alert will close on the next
  scan with no action needed from us.

When a tool ships a fix, the corresponding row here can be removed
in a follow-up PR.

## Alerts that stay open

### Codacy Deadcode — `pkg/ws/subscribe.go:34`

**Symptom.** Codacy reports a `deadcode_deadcode` finding with the
message `expected '(', found '['` against `Subscribe[T any](...)`.

**Root cause.** Codacy's bundled deadcode analyser fails to parse
Go 1.18+ generic-function syntax. The `[T any]` type-parameter
list is interpreted as a syntax error and, because deadcode bails
out early, every function in the file becomes "unreachable".

**Why we accept it.** The function is reachable, used by tests and
examples, and is part of the public WS API. Removing generics would
degrade the API for a tool bug. Reporting upstream is the path
forward; there is no source change that closes this alert.

### Codacy remark-lint — CODE_OF_CONDUCT.md (~20 alerts)

**Symptom.** `remark-lint-no-undefined-references` flags 20 reference
links in `CODE_OF_CONDUCT.md`.

**Root cause.** Every flagged reference (`[homepage]`, `[v2.1]`,
`[Mozilla CoC]`, `[FAQ]`, `[translations]`) IS defined at the foot
of the file (lines 80–84). The remark-lint resolver mis-handles
the cross-paragraph layout used by the upstream Contributor Covenant
template.

**Why we accept it.** The content is the canonical Contributor
Covenant text. Restructuring it to avoid the linter false positive
would diverge from the published template. Codacy's remark-lint
version is older than upstream; future Codacy updates may include
the fix.

### Checkov CKV_GHA_7 — release.yml, verify-release.yml, integration.yml

**Symptom.** Checkov flags every workflow that defines
`workflow_dispatch.inputs`, regardless of how those inputs are used.

**Root cause.** The rule is presence-based: any input block fails
the check, with no way to express that the input is bounded or
validated.

**Why we accept it.** All three of our `workflow_dispatch.inputs`
are load-bearing operational tooling:

- `release.yml` — `tag` input enables manual re-run when a release
  flow fails (documented in `docs/release-process.md`). Now validated
  against `^v[0-9]+\.[0-9]+\.[0-9]+(-[A-Za-z0-9.-]+)?$` in a
  `validate-input` job before any code runs.
- `verify-release.yml` — `tag` input lets a maintainer re-verify a
  specific historical release. Validated inline in the "Resolve
  tag" step before any download or verification.
- `integration.yml` — `run_live_orders` boolean is bounded by
  GitHub's UI (`true`/`false` only) and is consumed via an `if:`
  gate at the test step, never interpolated into shell.

In each case the rule's underlying intent (no arbitrary user input
into build output) is satisfied; the rule's presence-only check
is not. We chose to keep the inputs and add validation rather than
remove documented operational capabilities.

### Scorecard `BranchProtectionID`

**Symptom.** OpenSSF Scorecard reports that branch protection is
not configured.

**Root cause.** Scorecard predates the GitHub Rulesets API
(`/repos/{owner}/{repo}/rulesets`) and only inspects the legacy
`/repos/{owner}/{repo}/branches/{branch}/protection` endpoint. Our
master-branch protection lives in ruleset id `15981309`
("master-branch-protection") which Scorecard cannot see.

**Why we accept it.** Master is fully protected — required PRs,
linear history, status checks, signing — but via a ruleset that
Scorecard does not yet support. The path forward is upstream:
contribute ruleset support to <https://github.com/ossf/scorecard>.

### Scorecard `MaintainedID`, `SASTID`, `CITestsID`

**Symptom.** Three soft-positive Scorecard checks each fire one
alert claiming the project is unmaintained, lacks SAST coverage, or
lacks CI tests.

**Root cause.** Scorecard's heuristics are conservative and produce
false positives on actively-maintained projects with comprehensive
SAST and CI coverage. Specifically:

- **Maintained**: requires N commits in the last 90 days; this
  repo will accumulate that history naturally.
- **SAST**: detects only a fixed set of tools / actions; our SAST
  posture is CodeQL + gosec + Semgrep + Trivy + Codacy + osv-scanner
  + govulncheck + gitleaks + trufflehog, which the heuristic does
  not fully match.
- **CITests**: looks for test invocation patterns in CI; ours run
  via `go test ./... -race` across a `os × go-version` matrix.

**Why we accept it.** All three have full source-level fixes
already in place. The Scorecard signal lags reality.
