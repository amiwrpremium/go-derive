# Security

This directory holds the security-specific guidance. Vulnerability
reporting itself goes through GitHub's private advisory channel — see
the root [SECURITY.md](../../SECURITY.md).

## Contents

| Doc | Covers |
|---|---|
| [repo-policy.md](./repo-policy.md) | required GitHub repo settings (branch protection, signed commits, 2FA, allowed actions, etc.) |
| [threat-model.md](./threat-model.md) | what the SDK protects against, what's out of scope, attacker scenarios |

## Other security artifacts

| Where | What |
|---|---|
| [`SECURITY.md`](../../SECURITY.md) | vulnerability disclosure policy |
| [`SECURITY-INSIGHTS.yml`](../../SECURITY-INSIGHTS.yml) | OpenSSF security insights spec 1.0 |
| [`CODE_OF_CONDUCT.md`](../../CODE_OF_CONDUCT.md) | Contributor Covenant 2.1 |
| [`.github/workflows/scorecard.yml`](../../.github/workflows/scorecard.yml) | OpenSSF Scorecard CI |
| [`.github/workflows/osv-scanner.yml`](../../.github/workflows/osv-scanner.yml) | OSV-Scanner CI |
| [`.github/workflows/codeql.yml`](../../.github/workflows/codeql.yml) | CodeQL CI |
| [`.github/workflows/gosec.yml`](../../.github/workflows/gosec.yml) | gosec CI |
| [`.github/workflows/release.yml`](../../.github/workflows/release.yml) | SBOM + cosign + SLSA Level 3 |

## Public scorecard

Once the repo is published and Scorecard has run a few times, the public
score appears at:

[scorecard.dev/viewer/?uri=github.com/amiwrpremium/go-derive](https://scorecard.dev/viewer/?uri=github.com/amiwrpremium/go-derive)
