# Contributing to go-derive

Thanks for considering a contribution. This project follows a few light
conventions to keep CI green and releases automated.

## Conventional Commits

Every commit on `master` must follow the [Conventional Commits](https://www.conventionalcommits.org/)
spec. release-please reads the commit log to compute the next semantic version
and to update [CHANGELOG.md](./CHANGELOG.md).

Common types this repo recognises (see [.github/release-please-config.json](.github/release-please-config.json)):

| Type     | Bump | Section in CHANGELOG |
|----------|------|----------------------|
| `feat`   | minor | Features |
| `fix`    | patch | Bug Fixes |
| `perf`   | patch | Performance Improvements |
| `deps`   | patch | Dependencies |
| `revert` | patch | Reverts |
| `docs`   | none  | Documentation |
| `refactor` | none | Code Refactoring |
| `test`   | none  | Tests |
| `build`  | none  | Build System |
| `ci`     | none  | Continuous Integration |
| `chore`  | none  | Miscellaneous |

Add `BREAKING CHANGE:` in the commit body or `!` after the type (e.g. `feat!:`)
to force a major bump.

Examples:

```text
feat(rest): add GetFundingRateHistory

fix(ws): drain pending RPCs on connection close

deps: bump go-ethereum to v1.14.12

feat(auth)!: rename Signer.Sign to Signer.SignAction

BREAKING CHANGE: rename to disambiguate from header signing.
```

## Git hooks (lefthook)

Hooks are managed by [lefthook](https://github.com/evilmartians/lefthook).
The configuration lives in [`lefthook.yml`](./lefthook.yml). One-time
setup after cloning:

```bash
make tools     # installs goimports, govulncheck, golangci-lint, lefthook
make hooks     # registers pre-commit, commit-msg, pre-push in .git/hooks
```

What each hook runs:

| Hook | Phase | What it does |
|---|---|---|
| `pre-commit` | before each commit, parallel | `gofmt -l/-w`, `goimports -l/-w` (auto-fix; commit fails so you re-stage), `go vet` on changed packages, `golangci-lint run --new-from-rev=HEAD`, hygiene (trailing whitespace, BOM, merge-conflict markers, files >500 KB) |
| `commit-msg` | once per commit | Conventional Commits regex (same regex as `.github/workflows/pr-title.yml` and the `master` branch ruleset) |
| `pre-push` | once per push, sequential | `go build ./...`, `go vet ./...`, `go test -race -count=1 ./...`, `govulncheck ./...` |

Bypass for one operation (use sparingly — CI will still enforce):

```bash
LEFTHOOK=0 git commit ...
LEFTHOOK=0 git push ...
```

Run a hook on demand without committing:

```bash
make hooks-run-pre-commit
make hooks-run-pre-push
```

## Local checks before pushing

If you skip `make hooks` for any reason, run the equivalent manually:

```bash
go fmt ./...
go vet ./...
go test -race ./...
golangci-lint run        # optional but matches CI
govulncheck ./...        # optional but matches CI
```

## Coverage

CI uploads coverage to [Codecov](https://codecov.io/gh/amiwrpremium/go-derive)
and Codacy. Patches that meaningfully drop coverage will get flagged on the PR.

## Releases

Releases are fully automated:

1. PRs merge to `master` with Conventional Commit messages.
2. release-please opens (or updates) a `chore: release X.Y.Z` PR with the next
   version and the changelog entries it would generate.
3. Merging that PR cuts the GitHub release and pushes the `vX.Y.Z` tag.

You should never tag manually.

## Code of conduct

This project follows the [Contributor Covenant 2.1](./CODE_OF_CONDUCT.md).
Reports go through GitHub's private advisory channel — see the
[Enforcement](./CODE_OF_CONDUCT.md#enforcement) section for the link.
