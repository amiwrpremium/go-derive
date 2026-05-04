# GitHub rulesets

Declarative rulesets for the `master` branch and the `v*` release tags.
GitHub's newer rulesets system (2023+) is more powerful than classic
branch protection: it supports tag protection, commit-message patterns,
and JSON import/export.

## Files

| File | Target | Purpose |
|---|---|---|
| `branch-master.json` | `refs/heads/master` | required PR + reviews + status checks + signed commits + linear history + Conventional Commits regex |
| `tags-v.json` | `refs/tags/v*` | only `release-please[bot]` and admins can create/move/delete release tags; signatures required |

## One-time bootstrap

```bash
# Install: rulesets are repo-scoped, no app to install.
gh api repos/amiwrpremium/go-derive/rulesets \
  --method POST \
  --input .github/rulesets/branch-master.json

gh api repos/amiwrpremium/go-derive/rulesets \
  --method POST \
  --input .github/rulesets/tags-v.json
```

Before importing `tags-v.json`, fill the two `actor_id` placeholders:

- The `Integration` entry's `actor_id` → `gh api /apps/release-please --jq .id`
- The `RepositoryRole` entry's `actor_id` → `5` (GitHub's built-in admin role)

## Reconciling drift

```bash
# List existing rulesets and grab the IDs.
gh api repos/amiwrpremium/go-derive/rulesets

# Update in place (replaces the entire ruleset).
gh api repos/amiwrpremium/go-derive/rulesets/<ruleset-id> \
  --method PUT \
  --input .github/rulesets/branch-master.json
```

The ruleset on the live repo is authoritative; if it drifts from the
JSON in this directory, re-run the `PUT` to bring it back in line.

## Why both rulesets *and* `settings.yml`?

`settings.yml`'s `branches:` block is the legacy branch-protection API.
It's kept as a fallback so a freshly forked repo is sane out of the box
even before rulesets are imported. Once the rulesets above are imported,
they take precedence; the legacy block remains as belt-and-suspenders.
