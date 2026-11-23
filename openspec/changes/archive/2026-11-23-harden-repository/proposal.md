# Proposal: harden the repository

## Why

The Chorale repo is public, and Cami wants it ready to share. That means the free security tools GitHub offers public repos, dependency updates that keep themselves current without becoming a chore, and required checks that nobody can merge around. This change finishes the repository setup that milestone 0 in `docs/plan.md` started, and it comes before milestone 1 so that every later change lands on a protected `main`.

## What Changes

- CodeQL scans the Go code and the workflow files on every pull request, on pushes to `main` and once a week.
- `govulncheck` runs in `make check`, pinned as a Go tool in `go.mod`, with an offline fixture that proves it still finds a vulnerability.
- The dependency review action runs on every pull request and fails when the pull request adds a package with a known vulnerability. It becomes a required check.
- OpenSSF Scorecard runs once a week and on pushes to `main`, and publishes its results to the code scanning page.
- Dependabot opens version updates once a week for Go modules, GitHub Actions and npm, with a seven-day cooldown. Patch and minor updates are grouped per ecosystem.
- A workflow approves and turns on auto-merge for Dependabot's patch and minor updates, both version updates and security updates, so they merge once every required check passes. It approves with a fine-grained token kept as a Dependabot secret.
- Some updates wait for Cami instead: major updates, any update whose size can't be told, and updates to three named dependencies. Those are the two experimental Charm pins, `teatest` and `x/exp/golden`, and `dependabot/fetch-metadata`, which the workflow itself uses. For a major update, an agent reads the release notes, runs the checks and gives Cami a short summary, so her part is the product decision, not a code review.
- A second ruleset on `main`, called "Required checks", requires `check (macos)`, `check (linux)` and `dependency-review`, with no bypass actors. Cami's bypass of "Protect main" then no longer skips the checks. It doesn't require a branch to be up to date with `main`.
- `CODEOWNERS` names Cami. `SECURITY.md` points to private vulnerability reporting. A short `CONTRIBUTING.md` points to `AGENTS.md`, `docs/testing.md` and `docs/style.md`.
- A new `docs/repository.md` records every repository setting, both the ones this change adds and the ones already set on GitHub. It also gives the way out when a required check is stuck. Both rulesets are kept as JSON in `.github/rulesets/`.
- `AGENTS.md` changes in two places. The pull request rules drop the line telling Cami to check both jobs by eye, and the gates say that `make check` needs the network.

Auto-merge covers Dependabot only. Auto-merge can't use Cami's bypass, so her own pull requests can't merge through it, and they need her judgement anyway.

## Capabilities

### New Capabilities

None. This is a tooling change, so `.openspec.yaml` sets `skip_specs: true`. Nothing about how Chorale behaves for its user changes.

### Modified Capabilities

None.

## Out of scope

- A merge queue. GitHub doesn't offer one for repos owned by a personal account.
- Requiring code owner review. Agents push as Cami's account, and she can't approve her own pull requests, so the rule could only be met by bypassing it.
- Making CodeQL or Scorecard results block merging. Their findings show on pull requests and on the security page, and a later change can add a code scanning rule once there's a record of how noisy they are.
- Checking automatically that the live rulesets match the JSON in the repo. The rules can be read with the workflow token, but the bypass list can't, and the bypass list is the part that matters most. A later change can add a partial check.
- Dependabot updates for the golangci-lint version in the `Makefile`, the Node version in `.nvmrc`, and the `toolchain` line in `go.mod`. Dependabot can't update the first two, and isn't known to update the third.
- Release signing, provenance and a Scorecard badge. Chorale has no releases yet.

## Impact

New files: `.github/workflows/codeql.yml`, `dependency-review.yml`, `scorecard.yml` and `dependabot-auto-merge.yml`, `.github/dependabot.yml`, `.github/CODEOWNERS`, `.github/rulesets/protect-main.json` and `required-checks.json`, `SECURITY.md`, `CONTRIBUTING.md`, `docs/repository.md`, a `govulncheck` fixture under `testdata/vulncheck/`, and `scripts/check-vulncheck-fixture.sh`.

Changed files: `go.mod` and `go.sum` gain `govulncheck` as a tool, the `Makefile` gains a `vulncheck` target inside `check`, and `AGENTS.md`, `README.md`, `docs/testing.md`, `docs/getting-started.md` and `STATE.md` are updated.

GitHub settings and accounts: one new ruleset, a fine-grained token that Cami creates, and a Dependabot secret holding it. Each is applied only with Cami's go-ahead at that step. The setting that lets GitHub Actions approve pull requests stays off.

`make check` now needs a network connection, because `govulncheck` downloads the Go vulnerability database.
