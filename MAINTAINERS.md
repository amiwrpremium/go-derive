# Maintainers

The list below is normative — `.github/CODEOWNERS` references usernames
in the same order. To add a maintainer, see [GOVERNANCE.md](./GOVERNANCE.md#adding-a-maintainer).

## Current

| Maintainer | Areas | Contact |
|---|---|---|
| [@amiwrpremium](https://github.com/amiwrpremium) | all | private channel via [SECURITY.md](./SECURITY.md) |

## Emeritus

(none yet)

## Maintainer responsibilities

- Triage issues within ~7 days of opening.
- Review PRs within ~7 days of submission.
- Cut releases when conventional-commit history warrants one (release-please opens the PR; a maintainer reviews and merges).
- Respond to security advisories within 72 hours per
  [SECURITY.md](./SECURITY.md).
- Keep dependencies current — Renovate auto-merges patch/minor; majors
  for security-critical packages (`go-ethereum`, `gorilla/websocket`)
  require manual review per `renovate.json`.
