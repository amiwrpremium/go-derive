# Support

Where to ask which kind of question.

## Bug or unexpected behaviour

Open a [bug report](https://github.com/amiwrpremium/go-derive/issues/new?template=bug_report.md).
Include the SDK version (`go list -m github.com/amiwrpremium/go-derive`),
Go version (`go version`), a minimal repro, and the error you hit.

## Feature request / API design question

Open a [feature request](https://github.com/amiwrpremium/go-derive/issues/new?template=feature_request.md).
A short example of the call site you'd like to write helps the discussion
move quickly.

## "How do I…?"

Use [GitHub Discussions](https://github.com/amiwrpremium/go-derive/discussions)
for usage questions. Things to try first:

1. The [godoc reference](https://pkg.go.dev/github.com/amiwrpremium/go-derive)
   on pkg.go.dev — every public identifier has documentation and most
   have an `Example_*`.
2. The runnable [`examples/`](./examples/) tree — 80 programs, one per
   API method.
3. The [`docs/`](./docs/) deep-dive guides.

## Security disclosure

**Do not** open a public issue for security problems. Use [GitHub's
private vulnerability reporting](https://github.com/amiwrpremium/go-derive/security/advisories/new).
Full policy: [SECURITY.md](./SECURITY.md).

## Code-of-conduct concern

Same private channel as security disclosure. Full policy:
[CODE_OF_CONDUCT.md](./CODE_OF_CONDUCT.md).

## Pull requests

See [CONTRIBUTING.md](./CONTRIBUTING.md). The repo enforces
[Conventional Commits](https://www.conventionalcommits.org/) on PR
titles via the `pr-title` workflow so release-please's auto-changelog
stays clean.

## Derive API itself

Issues with the upstream Derive API (wrong server response, rate-limit
behaviour) should go to the [Derive team](https://docs.derive.xyz/), not
here. The SDK only translates calls — the engine's behaviour is theirs.
