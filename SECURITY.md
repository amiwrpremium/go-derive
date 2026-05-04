# Security policy

## Supported versions

| Version | Supported |
|---------|-----------|
| `0.x`   | Latest minor only |
| `< 0.x` | No |

## Reporting a vulnerability

Please **do not** open a public GitHub issue for security problems.

Instead, use GitHub's [Private Vulnerability Reporting](https://github.com/amiwrpremium/go-derive/security/advisories/new)
to send a private advisory. Include:

- A description of the issue and its impact
- Steps to reproduce or a proof of concept
- Affected versions
- Any mitigations you're aware of

You should expect an acknowledgement within 72 hours and a triage update within
7 days. Critical issues will be patched and disclosed coordinated with the
reporter.

## Scope

In scope:

- Code in `pkg/` and `internal/` published as part of `go-derive`.
- The signing path (`pkg/auth`, `internal/codec`).
- The transports (`internal/transport`).

Out of scope:

- Vulnerabilities in upstream dependencies (please report to those projects).
- Issues in the Derive API itself — report those to the Derive team directly.
- Examples in `examples/` that are clearly demonstrations.
