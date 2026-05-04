<!--
  Title: keep it short, use Conventional Commits format:
    feat(scope): add X
    fix(scope): handle Y when Z
    deps: bump <pkg> to <ver>
-->

## Summary

<!-- 1-3 bullets on what changed and why. -->
-

## Type of change

<!-- Tick the relevant box. -->
- [ ] feat — new user-facing feature
- [ ] fix — bug fix
- [ ] perf — performance improvement
- [ ] refactor — code change that neither fixes a bug nor adds a feature
- [ ] docs — documentation only
- [ ] test — adds/updates tests
- [ ] ci — CI/build changes
- [ ] deps — dependency bump
- [ ] chore — other

## Test plan

<!-- How did you verify this works? -->
- [ ] `go test -race ./...` passes locally
- [ ] `go vet ./...` clean
- [ ] `golangci-lint run` clean
- [ ] Manually verified against testnet (describe below)

## Breaking changes

<!--
  If yes, describe what breaks and how callers should migrate. Include
  "BREAKING CHANGE:" in your commit body so release-please bumps the major.
-->
- [ ] No breaking changes
- [ ] Yes — see description below

## Checklist

- [ ] Conventional-commit title and body
- [ ] Updated docs / godoc where relevant
- [ ] Added/updated tests
- [ ] No new lint warnings
