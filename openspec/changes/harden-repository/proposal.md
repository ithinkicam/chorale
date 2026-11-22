# Proposal: harden the repository

## Why

The Chorale repo is public, and Cami wants it ready to share. That means the free security tools GitHub offers public repos, dependency updates that keep themselves current, and required checks that nobody can merge around. This change finishes the repository setup that milestone 0 in `docs/plan.md` started, and it comes before milestone 1 so that every later change lands on a protected `main`.

## What Changes

- CodeQL scans the Go code and the workflow files on every pull request, on pushes to `main` and once a week.
- `govulncheck` runs in `make check`, pinned as a Go tool in `go.mod`, with a fixture that proves it fails on a known vulnerability.
- The dependency review action runs on every pull request and fails when the pull request adds a package with a known vulnerability. It becomes a required check.
- OpenSSF Scorecard runs once a week and on pushes to `main`, and publishes its results to the code scanning page.
- Dependabot opens version updates once a week for Go modules, GitHub Actions and npm. Patch and minor updates are grouped per ecosystem.
- A workflow approves and turns on auto-merge for Dependabot's patch and minor updates, so they merge once every required check passes. Major updates, and any update whose size Dependabot can't tell, wait for Cami.
- A second ruleset on `main`, called "Required checks", requires `check (macos)`, `check (linux)` and `dependency-review`, with no bypass actors. Cami's bypass of "Protect main" then no longer skips the checks.
- The repo setting that lets GitHub Actions approve pull requests is turned on, because the auto-merge workflow needs it.
- `CODEOWNERS` names Cami. `SECURITY.md` points to private vulnerability reporting. A short `CONTRIBUTING.md` points to `AGENTS.md`, `docs/testing.md` and `docs/style.md`.
- A new `docs/repository.md` records every repository setting, both the ones this change adds and the ones already set on GitHub. Both rulesets are kept as JSON in `.github/rulesets/`.
- The pull request rules in `AGENTS.md` change, because Cami no longer needs to check by eye that both jobs are green before merging.

## Capabilities

### New Capabilities

None. This is a tooling change, so `.openspec.yaml` sets `skip_specs: true`. Nothing about how Chorale behaves for its user changes.

### Modified Capabilities

None.

## Out of scope

- A merge queue. GitHub doesn't offer one for repos owned by a personal account.
- Requiring code owner review. Agents push as Cami's account, and she can't approve her own pull requests, so the rule could only be met by bypassing it.
- Making CodeQL or Scorecard results block merging. Their findings show on pull requests and on the security page, and a later change can add a code scanning rule once there's a record of how noisy they are.
- A personal access token for Scorecard. Without one, Scorecard can't read every ruleset detail, and its branch protection score is partial.
- Checking automatically that the live rulesets match the JSON in the repo. Reading rulesets needs an admin token, which continuous integration doesn't have.
- Dependabot updates for the golangci-lint version in the `Makefile` and the Node version in `.nvmrc`. Dependabot can't update either.
- Release signing, provenance and a Scorecard badge. Chorale has no releases yet.

## Impact

New files: `.github/workflows/codeql.yml`, `dependency-review.yml`, `scorecard.yml` and `dependabot-auto-merge.yml`, `.github/dependabot.yml`, `.github/CODEOWNERS`, `.github/rulesets/protect-main.json` and `required-checks.json`, `SECURITY.md`, `CONTRIBUTING.md`, `docs/repository.md`, a `govulncheck` fixture under `testdata/`, and a script that checks the fixture.

Changed files: `go.mod` and `go.sum` gain `govulncheck` as a tool, the `Makefile` gains a `vulncheck` target inside `check`, and `AGENTS.md`, `README.md`, `docs/getting-started.md` and `STATE.md` are updated.

GitHub settings: one new ruleset, and permission for GitHub Actions to approve pull requests. Both are outward-facing, and each is applied only with Cami's go-ahead at that step.

`make check` now needs a network connection, because `govulncheck` downloads the Go vulnerability database.
