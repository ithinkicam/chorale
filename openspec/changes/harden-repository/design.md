# Design

## Context

Milestone 0 left one workflow, `.github/workflows/check.yml`, which runs `make check` on macOS and Linux with a read-only token and actions pinned to commit hashes. `actionlint` checks every workflow file in `make lint`, and a broken fixture proves it still finds problems.

The repo is public and owned by Cami's personal account. One ruleset, "Protect main", applies to the default branch. It blocks deletion and force-pushes, and requires a pull request with one approval, resolved threads, and stale approvals dismissed on push. Code owner review is not required. The admin role can bypass it, for pull requests only. All three merge methods are allowed, auto-merge is allowed, and branches are deleted on merge. Secret scanning, push protection, Dependabot alerts, Dependabot security updates and private vulnerability reporting are on. GitHub Actions can't yet approve pull requests.

Agents push as Cami's account, so they can't approve their own pull requests, and Cami merges them with her bypass. Bypassing a ruleset skips every rule in it, which is why the required checks need a ruleset of their own.

## Goals / Non-Goals

**Goals:**

- Every new check is seen to fail on a known-bad input before it is trusted, the same way the `actionlint` fixture works.
- Every new workflow follows the rules `check.yml` already follows: actions pinned to a full commit hash with the version in a comment, the least token permission each job needs, `persist-credentials: false`, no `pull_request_target`, a time limit on every job, and no path filters on anything that becomes a required check.
- Every repository setting is written down in the repo, so it can be checked and rebuilt.

**Non-Goals:**

- See the proposal's "Out of scope". The design adds no further boundaries.

## Decisions

### CodeQL uses a workflow file, not default setup

GitHub's default setup for CodeQL is a switch in the settings, with nothing in the repo. A workflow file is pinned, checked by `actionlint`, reviewed in pull requests, and visible to anyone reading the repo, so it is used instead.

`.github/workflows/codeql.yml` analyses two languages in a matrix. Go uses `build-mode: autobuild`, after `actions/setup-go` installs the toolchain from `go.mod`. The workflow files use the `actions` language with `build-mode: none`. It runs on pull requests, on pushes to `main`, and once a week, so that new CodeQL queries reach code that hasn't changed. The job needs `security-events: write` to upload results, and `contents: read` and `actions: read` beside it.

CodeQL is not a required check. Its alerts appear on the pull request and on the security page. A later change can add the ruleset rule that blocks on code scanning results.

### govulncheck runs in make check, with an offline fixture

`govulncheck` is added with `go get -tool golang.org/x/vuln/cmd/govulncheck`, so its version is pinned in `go.mod` beside `actionlint` and `gotestsum`. A new `vulncheck` target runs `go tool govulncheck ./...`, and `check` runs it after `lint`. It reports only vulnerabilities in code Chorale actually calls, which keeps it quiet enough to fail the gate.

The fixture proves it still finds problems. `testdata/vulncheck/` holds a tiny module that calls a function in a local module, and a local vulnerability database in the OSV format that marks that function as vulnerable. `scripts/check-vulncheck-fixture.sh` runs `govulncheck` on the fixture with `-db` pointing at the local database, and fails unless `govulncheck` exits with status 3, which means it found a vulnerability. The fixture needs no network and names no real advisory, so GitHub's dependency graph raises no Dependabot alert for it. Go ignores `testdata` folders, so `./...` never builds it.

The alternative was a real old package with a real advisory. It would need the network, and GitHub would raise a real alert for the fixture's `go.mod`.

`govulncheck` downloads the database from `vuln.go.dev`, so `make check` needs the network from now on. A new advisory against the Go standard library can turn `main` red without any change in the repo. That is the point of the check, and the fix is a new `toolchain` line in `go.mod`.

### Dependency review is a required check

`.github/workflows/dependency-review.yml` runs `actions/dependency-review-action` on every pull request, with no path filters, so it always reports and can be required. It fails on a vulnerability of any severity, which is the action's default. Licence checks are off. It keeps a read-only token, so it doesn't post a summary comment, and its log and the job summary show what it found. The job is named `dependency-review`, which is the name the ruleset requires.

### Scorecard publishes its results

`.github/workflows/scorecard.yml` runs `ossf/scorecard-action` once a week, on pushes to `main`, and when a branch protection rule changes. It publishes results to the OpenSSF API and uploads them to the code scanning page. Publishing places limits on the workflow: one job, no write permission at the top level, `id-token: write` and `security-events: write` on the job, and only approved actions in its steps. It runs without a personal access token, so its branch protection score is partial.

### Dependabot version updates

`.github/dependabot.yml` covers `gomod`, `github-actions` and `npm`, each in `/`, once a week. In each ecosystem, patch and minor updates are grouped into one pull request, so a busy week brings three pull requests rather than twenty. Major updates come one per pull request.

Each ecosystem sets a cooldown of seven days, so Dependabot doesn't propose a version until it has been public for a week. Most malicious releases are found and pulled within days, and a week's delay costs little for a project with no users yet. Security updates ignore the cooldown, so a fix for a known vulnerability still arrives at once.

For actions, Dependabot updates the commit hash and the version comment beside it together.

### Auto-merge for patch and minor updates

`.github/workflows/dependabot-auto-merge.yml` runs on `pull_request`. Its one job runs only when the pull request's author is `dependabot[bot]` and the pull request comes from this repo. It reads the update with `dependabot/fetch-metadata`. If the update type is `version-update:semver-patch` or `version-update:semver-minor`, it approves the pull request with `gh pr review --approve` and turns on auto-merge with `gh pr merge --auto --squash`. The job needs `contents: write` and `pull-requests: write`, and never checks out the pull request's code.

The test is a list of what may merge, not a list of what may not. A major update, a Go pseudo-version whose size Dependabot can't tell, or any type the action doesn't recognise leaves the pull request waiting for Cami.

Approval comes from `github-actions[bot]`, and auto-merge waits for every required check in both rulesets. So a Dependabot update merges only when `make check` passes on both systems and dependency review passes. Golden files and every other test guard against an update that quietly changes behaviour.

The approval needs the repo setting "Allow GitHub Actions to create and approve pull requests". With it on, any workflow in this repo with `pull-requests: write` could approve a pull request. Only this workflow asks for that permission, and pull requests from forks always get a read-only token, so an outside contributor can't use it.

The alternative was to let Dependabot pull requests skip the approval rule. Rulesets can't exempt one author from one rule without exempting it from all of them.

### The second ruleset holds the required checks

"Required checks" targets the default branch and has no bypass actors. Its one rule requires three checks: `check (macos)`, `check (linux)` and `dependency-review`, each tied to the GitHub Actions app so that no other app or status can satisfy them. It also requires the branch to be up to date with `main` before merging, so a pull request is always checked against what it will merge into. Dependabot rebases its own pull requests when `main` moves, and Cami's pull requests are merged one at a time, so the cost is small.

With no bypass, nobody can merge a red pull request, Cami included. In an emergency she can still change or turn off the ruleset itself, as the repo's admin.

Both rulesets are kept in `.github/rulesets/` as the JSON that GitHub exports, so they can be reviewed and recreated. `required-checks.json` is applied with `gh api --method POST repos/ithinkicam/chorale/rulesets --input .github/rulesets/required-checks.json`. `protect-main.json` is exported from the live ruleset, and only documents it.

### Community files

`.github/CODEOWNERS` has one line, `* @ithinkicam`. It asks Cami to review a contributor's pull request automatically. Code owner review isn't required, for the reason in the proposal.

`SECURITY.md` says that vulnerabilities are reported privately through the repo's security advisories page, with its address, and not in public issues. It says only `main` is supported, because there are no releases.

`CONTRIBUTING.md` is a few short paragraphs. It says that work starts as an OpenSpec change, that development is test first, and points to `AGENTS.md`, `docs/testing.md` and `docs/style.md`.

### docs/repository.md

This new document describes every repository setting: the security features, both rulesets and why there are two, the Actions permission, merge settings, Dependabot's schedule, groups and cooldown, and the auto-merge policy. It says what Cami does with a major update and with a CodeQL or Scorecard finding. `AGENTS.md`, `README.md` and `docs/getting-started.md` point to it.

The pull request rules in `AGENTS.md` change. The line telling Cami to confirm both jobs are green before merging goes, because the "Required checks" ruleset now enforces it. A line is added that agents never merge Dependabot's major updates on their own.

## Risks / Trade-offs

**An approval from GitHub Actions may not count toward the ruleset's one required approval.** GitHub's Dependabot documentation uses this pattern, but that was not confirmed against a ruleset. The "Protect main" ruleset also requires an extra approval for unattributed changes, whose effect on Dependabot's commits isn't documented clearly. The first Dependabot pull request after merging shows the answer. If the approval doesn't count, the fallback is for Cami to approve those pull requests by hand, and the auto-merge workflow still saves her the merge.

**The checks are only proved against fixtures and canaries before merging.** Dependabot opens pull requests only from the configuration on `main`, so auto-merge can't be seen working until this change has merged. That check is listed as a follow-up after merging, and `STATE.md` records it until it is done.

**A canary for the required-checks ruleset could merge broken code.** The canary pull request is never merged, and nobody tries `gh pr merge` on it. The ruleset is proved by reading the pull request's merge state, which must be `BLOCKED`, and the list of required checks GitHub reports for it.

**Grouped updates can hide a bad package among good ones.** If a group fails `make check`, Cami or an agent splits it with `@dependabot` commands, or ignores the package that broke it.

**`make check` needs the network.** Offline, `vulncheck` fails. Running `make lint test e2e scenarios specs` covers everything else.

## Migration Plan

The order matters, because a required check that never reports would block every pull request.

1. The new workflows, configuration and documents land on the `harden-repository` branch, and its pull request runs them all.
2. With the pull request open, and with Cami's go-ahead, the Actions setting is turned on and the "Required checks" ruleset is created. The change's own pull request is the first to meet it.
3. Canary pull requests from the branch prove that dependency review and the ruleset block what they should. Then they are closed.
4. Cami merges the change. Dependabot reads its configuration from `main` and opens its first pull requests.

To roll back, Cami turns off or deletes the "Required checks" ruleset and reverts the pull request.

## Test strategy

This change has no specs, so no scenario needs a test. Each check is instead proved to fail on a known-bad input:

1. `govulncheck`: the offline fixture and `scripts/check-vulncheck-fixture.sh`, run by `make vulncheck` inside every `make check`. It is watched failing first, before the fixture exists.
2. The new workflow files: `actionlint` in `make lint`, which already rejects a broken fixture.
3. Dependency review: a draft canary pull request that adds a package with a known advisory. The `dependency-review` job must fail. Checked by hand once, because it needs GitHub's advisory database.
4. The "Required checks" ruleset: a draft canary pull request that breaks a golden file. Both `check` jobs must fail, and the merge state must be `BLOCKED` for Cami's account. Checked by hand once, because it needs the live ruleset.
5. CodeQL and Scorecard: their first runs on the pull request, and on `main` after merging, must finish and upload results. Checked by hand, because their results live on GitHub.
6. Auto-merge: the first Dependabot patch or minor pull request after merging must be approved by `github-actions[bot]` and merge on its own once checks pass. A major update, when one arrives, must stay open. Checked by hand after merging.
