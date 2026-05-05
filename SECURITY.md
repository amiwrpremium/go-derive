# Security policy

## Supported versions

| Version | Supported |
|---------|-----------|
| `0.x`   | Latest minor only |
| `< 0.x` | No |

## Reporting a vulnerability

Please **do not** open a public GitHub issue for security problems.

Instead, use GitHub's [Private Vulnerability Reporting](https://github.com/amiwrpremium/go-derive/security/advisories/new)
to send a private advisory.

### Required (we cannot triage without these)

- A description of the issue and its impact
- Steps to reproduce or a proof of concept
- Affected versions

### Optional but appreciated

- Any mitigations you're aware of
- Suggested fix

### Response SLAs

| Severity | First acknowledgement | Triage update |
|---|---|---|
| **Critical** (remote code execution, key disclosure, unauthorised funds movement) | within 72 hours | within 7 days |
| **High** (auth bypass, information disclosure of session material) | within 72 hours | within 7 days |
| **Medium** (DoS, inconsistent state without funds risk) | within 5 business days | best effort |
| **Low** (informational, hardening suggestion) | within 5 business days | best effort |

Critical and High issues will be patched and disclosed in coordination
with the reporter.

## Scope

In scope:

- Code in `pkg/` and `internal/` published as part of `go-derive`.
- The signing path (`pkg/auth`, `internal/codec`).
- The transports (`internal/transport`).

Out of scope:

- Vulnerabilities in upstream dependencies (please report to those projects).
- Issues in the Derive API itself — report those to the Derive team directly.
- Examples in `examples/` that are clearly demonstrations.
