# Governance

Lightweight project. The maintainer (currently
[@amiwrpremium](https://github.com/amiwrpremium)) has final say on
direction, releases and merges. Contributors are welcome to open issues
and PRs.

## Decisions

| Type | Who decides | How |
|---|---|---|
| Bug fix | Maintainer | PR review + merge |
| Small feature (additive, no API break) | Maintainer | PR review + merge |
| Breaking API change | Maintainer | Discussion issue first, then PR with `feat!:` or `BREAKING CHANGE:` |
| Security policy | Maintainer | Edits to `SECURITY.md` / `docs/security/threat-model.md` |
| New maintainer | Existing maintainers, by consensus | A PR adding the username to `.github/CODEOWNERS` |

## Conflict resolution

If you disagree with a decision, open a Discussion. We'd rather have a
conversation in the open than a private one.

## Releases

Releases are automated via release-please — see [release-process.md](./docs/release-process.md).
The maintainer's only manual step is reviewing and merging the
release-please PR.

## Out of scope

This project doesn't have a formal RFC process or steering committee.
That can change if the project grows; for now, the bar is "small enough
that one maintainer can keep it coherent."

## Adding a maintainer

Open a PR that:

1. Adds the new username as a CODEOWNER for the relevant directories.
2. Updates [MAINTAINERS.md](./MAINTAINERS.md).

Approval requires sign-off from every existing maintainer.
