# Design

## Context

Milestone 0 left one workflow, `.github/workflows/check.yml`, which runs `make check` on macOS and Linux with a read-only token and actions pinned to commit hashes. `actionlint` checks every workflow file in `make lint`, and a broken fixture proves it still finds problems.

The repo is public and owned by Cami's personal account. One ruleset, "Protect main", applies to the default branch. It blocks deletion and force-pushes, and requires a pull request with one approval, resolved threads, and stale approvals dismissed on push. It also requires an extra approval for unattributed changes. Code owner review is not required. The admin role can bypass it, for pull requests only. All three merge methods are allowed, auto-merge is allowed, and branches are deleted on merge. Secret scanning, push protection, Dependabot alerts, Dependabot security updates and private vulnerability reporting are on. GitHub Actions can't approve pull requests, and this change keeps it that way.

Agents push as Cami's account, so they can't approve their own pull requests, and Cami merges them with her bypass. Bypassing a ruleset skips every rule in it, which is why the required checks need a ruleset of their own.

Cami makes the architectural and product decisions on this repo, and doesn't review Go code line by line. The checks, not a human reading the diff, are what stop a bad dependency update. So the design keeps the number of things that wait for her small.

## Goals / Non-Goals

**Goals:**

- Every new local check is seen to fail on a known-bad input before it is trusted, the same way the `actionlint` fixture works. The checks that run on GitHub are proved once, with a canary pull request or their first run.
- Every new workflow follows the rules `check.yml` already follows: actions pinned to a full commit hash with the version in a comment, the least token permission each job needs, `persist-credentials: false`, no `pull_request_target`, a time limit on every job, and no path filters on anything that becomes a required check.
- Every repository setting is written down in the repo, so it can be checked and rebuilt.

**Non-Goals:**

- See the proposal's "Out of scope". The design adds no further boundaries.

## Decisions

### CodeQL uses a workflow file, not default setup

GitHub's default setup for CodeQL is a switch in the settings, with nothing in the repo. A workflow file is pinned, checked by `actionlint`, reviewed in pull requests, and visible to anyone reading the repo, so it is used instead.

`.github/workflows/codeql.yml` analyses two languages in a matrix. Go uses `build-mode: manual`, with `actions/setup-go` installing the toolchain from `go.mod` and then `go build ./...`. Autobuild might run `make`, whose default target here only prints the Go version, so the build is spelt out. `./...` skips `testdata`, so the `govulncheck` fixture isn't analysed. The workflow files use the `actions` language with `build-mode: none`. It runs on pull requests, on pushes to `main`, and once a week, so that new CodeQL queries reach code that hasn't changed. The job needs `security-events: write` to upload results, and `contents: read` and `actions: read` beside it.

CodeQL is not a required check. Its alerts appear on the pull request and on the security page.

### govulncheck runs in make check, with an offline fixture

`govulncheck` is added with `go get -tool golang.org/x/vuln/cmd/govulncheck`, so its version is pinned in `go.mod` beside `actionlint` and `gotestsum`. A new `vulncheck` target runs the fixture check and then `go tool govulncheck ./...`, and `check` runs it after `lint`. It reports only vulnerabilities in code Chorale actually calls, which keeps it quiet enough to fail the gate.

The fixture proves it still finds problems, with no network and no real advisory. `testdata/vulncheck/` holds two things:

1. A tiny module with its own `go.mod` and no dependencies, whose `main` calls `fmt.Println`.
2. A local database in the Go vulnerability database layout: `index/db.json`, `index/modules.json` and one entry under `ID/`. The entry marks `fmt.Println` in the module `stdlib` as vulnerable in every version.

`scripts/check-vulncheck-fixture.sh` runs from the repo root, so `go tool` finds the pinned version, with `go tool govulncheck -C testdata/vulncheck -db "file://<absolute path to the database>" ./...`. The path is built from the script's own location and quoted, because this clone lives under `/Volumes/Drive 2`, which has a space. The script fails unless `govulncheck` exits with status 3, which means it found a vulnerability.

A local module can't stand in for the vulnerable package. A module wired in with a directory `replace` has no version, and `govulncheck` never matches an advisory against a module with no version. The standard library always has a version, so the fixture uses it. A real old package with a real advisory was the other choice, but it would need the network, and GitHub would raise a real alert for it. The fixture's module has no dependencies, so it gives GitHub's dependency graph nothing to raise an alert on.

`govulncheck` downloads the database from `vuln.go.dev`, so `make check` needs the network from now on, and the gates section of `AGENTS.md` says so. A new advisory can turn every pull request red without any change in the repo. When a fix exists, an update brings it in. A new `toolchain` line in `go.mod` covers the standard library. When no fix exists, or `vuln.go.dev` is down, nothing can merge, because nobody can bypass "Required checks". `docs/repository.md` gives the way out, with the exact commands: Cami turns off "Required checks", merges the pull request that works around it, and turns the ruleset back on. `govulncheck` has no way to ignore one finding, so the workaround is a code change, such as no longer calling the affected function.

### Dependency review is a required check

`.github/workflows/dependency-review.yml` runs `actions/dependency-review-action` on every pull request, with no path filters, so it always reports and can be required. It fails on a vulnerability of any severity, which is the action's default. Licence checks are off. It keeps a read-only token, so it doesn't post a summary comment, and its log and the job summary show what it found. The job is named `dependency-review`, which is the name the ruleset requires.

### Scorecard publishes its results

`.github/workflows/scorecard.yml` runs `ossf/scorecard-action` once a week, on pushes to `main`, and when started by hand. It has no trigger for branch protection changes, because that event fires only for classic branch protection, and this repo uses rulesets. It publishes results to the OpenSSF API and uploads them to the code scanning page. It follows the limits that publishing places on the workflow, as the action's README lists them when the task runs, including no write permission at the top level, and `id-token: write` and `security-events: write` on the job. It reads the rulesets with the default token.

### Dependabot version updates

`.github/dependabot.yml` covers `gomod`, `github-actions` and `npm`, each in `/`, once a week. In each ecosystem, patch and minor updates are grouped into one pull request, so a busy week brings three pull requests rather than twenty. Major updates come one per pull request.

Each ecosystem sets a cooldown of seven days, so Dependabot doesn't propose a version until it has been public for a week. Most malicious releases are found and pulled within days, and a week's delay costs little for a project with no users yet.

The three dependencies that wait for Cami are left out of the groups with `exclude-patterns`, so each comes in a pull request of its own, and a group is never held back by one of them. They are `github.com/charmbracelet/x/exp/teatest/v2` and `github.com/charmbracelet/x/exp/golden` in `gomod`, and `dependabot/fetch-metadata` in `github-actions`.

For actions, Dependabot updates the commit hash and the version comment beside it together.

### Security updates auto-merge too

Security updates ignore the cooldown, so a security update can arrive hours after its release. They still auto-merge under the same rules as version updates. A security update moves to the patched version named in a reviewed advisory, not to whatever is newest, so the poisoned-release risk the cooldown guards against is much smaller. Waiting for a human would leave a known vulnerability open for as long as the pull request sat, and a human reading a Go diff adds little here. Every required check still has to pass.

### Auto-merge with a token kept as a Dependabot secret

`.github/workflows/dependabot-auto-merge.yml` runs on `pull_request`, for the `opened`, `synchronize` and `reopened` types, so a rebased pull request is approved again after its old approval is dismissed. It runs only when the pull request's author is `dependabot[bot]` and the pull request comes from this repo. It never checks out the pull request's code.

It has two jobs, so the token never shares a runner with a third-party action.

1. `decide` runs `dependabot/fetch-metadata` with the default read-only token, and outputs whether the pull request may merge. It may merge only when the update type is `version-update:semver-patch` or `version-update:semver-minor`, and none of its dependencies is one of the three named in the previous section. The test is a list of what may merge, so a major update, an unknown type, or an empty type waits.
2. `merge` needs `decide`, runs only when it said yes, and uses no actions at all. It runs `gh pr review --approve` and then `gh pr merge --auto --squash`, with the token from the Dependabot secret `DEPENDABOT_AUTOMERGE_TOKEN`.

fetch-metadata calls an update between two Go pseudo-versions a patch, because it compares only the first two parts of the version. That is why the two experimental Charm pins are excluded by name. `dependabot/fetch-metadata` is excluded too, because a pull request that updates it runs the new version in its own `decide` job. The two-job split means that new version never sees the token.

The token is a fine-grained personal access token that Cami creates on her account. It is limited to this one repo, with read and write access to contents and to pull requests and nothing else. It never expires, by Cami's decision, so there is no renewal chore. If it is revoked, the `merge` job fails and Dependabot's pull requests wait. Nothing breaks, and `docs/repository.md` says how to replace it.

A token kept as a Dependabot secret solves three problems that the workflow token has:

1. The workflow token could only approve with the setting that lets GitHub Actions approve pull requests. That setting would let any pull request from a branch in this repo add a workflow step that approves itself. Dependabot secrets reach only runs that Dependabot starts, so the setting stays off.
2. An approval from Cami's account clearly counts toward "Protect main", because she isn't the pull request's author. An approval from `github-actions[bot]` might not.
3. Events caused by the workflow token don't start new workflow runs, so a merge it enabled might skip the workflows that run on pushes to `main`. A merge enabled by a personal token starts them as usual.

Approvals on Dependabot's pull requests show Cami's name, though the workflow gave them. `docs/repository.md` says so.

Auto-merge waits for every required check, so a Dependabot update merges only when `make check` passes on both systems and dependency review passes. Golden files and every other test guard against an update that quietly changes behaviour.

### Major updates

Major updates wait for Cami, but she doesn't review the code. When one arrives, an agent reads the release notes, runs `make check` on the branch, and gives her a short summary: what changes, what breaks, and what it recommends. She makes the product decision. `docs/repository.md` describes this, and `AGENTS.md` says that agents never merge a major update on their own. Updates to the three named dependencies are handled the same way.

### The second ruleset holds the required checks

"Required checks" targets the default branch and has no bypass actors. Its one rule requires three checks: `check (macos)`, `check (linux)` and `dependency-review`, each tied to the GitHub Actions app so that no other app or status can satisfy them.

It doesn't require a branch to be up to date with `main`. Dependabot rebases a pull request only when it conflicts, so with that rule its other weekly pull requests would sit behind `main` until the next week. Pressing "Update branch" would put Cami's commit at the head, which fetch-metadata refuses. Without a merge queue, the up-to-date rule mostly adds waiting, and every merge runs the checks on `main` anyway.

With no bypass, nobody can merge a red pull request, Cami included. In an emergency she can still change or turn off the ruleset itself, as the repo's admin.

Both rulesets are kept in `.github/rulesets/` as the payload GitHub accepts when a ruleset is created. That is `name`, `target`, `enforcement`, `conditions`, `rules` and `bypass_actors`, and no read-only fields such as `id`, `_links` or the timestamps. `required-checks.json` is applied with `gh api --method POST repos/ithinkicam/chorale/rulesets --input .github/rulesets/required-checks.json`. `protect-main.json` only documents the live ruleset. A file is checked against the live ruleset by putting the API's response through the same `jq` filter, which keeps only those six fields, and comparing.

### Community files

`.github/CODEOWNERS` has one line, `* @ithinkicam`. It asks Cami to review a contributor's pull request automatically. Code owner review isn't required, for the reason in the proposal.

`SECURITY.md` says that vulnerabilities are reported privately through the repo's security advisories page, with its address, and not in public issues. It says only `main` is supported, because there are no releases.

`CONTRIBUTING.md` is a few short paragraphs. It says that work starts as an OpenSpec change, that development is test first, and points to `AGENTS.md`, `docs/testing.md` and `docs/style.md`.

### Documents

`docs/repository.md` is new. It describes every repository setting: the security features, both rulesets and why there are two, merge settings, Dependabot's schedule, groups and cooldown, the auto-merge policy and its token, how to replace the token, how major updates are handled, and the way out when a required check is stuck. It says what to do with a CodeQL or Scorecard finding.

`AGENTS.md` changes in the pull request rules and the gates, as the proposal says. `docs/testing.md` gains `make vulncheck` in its commands and in what `make check` runs, and says that `make check` needs the network. `docs/getting-started.md` says `govulncheck` is a Go tool beside the others, and points to `package.json` for the OpenSpec version instead of naming one, so a Dependabot update can't leave it stale. `README.md` points to `docs/repository.md`.

## Risks / Trade-offs

**"Protect main" requires an extra approval for unattributed changes.** What GitHub counts as unattributed isn't documented clearly. If it covers Dependabot's commits, one approval isn't enough and Dependabot's pull requests wait. The first one after merging shows the answer. If they wait, the fallback is to turn off that rule in "Protect main", with Cami's decision.

**Auto-merge can't be seen working until this change has merged.** Dependabot opens pull requests only from the configuration on `main`. That check is listed as a follow-up after merging, and `STATE.md` records it until it is done.

**The token is a long-lived credential.** It is limited to one repo and two permissions, it reaches only runs Dependabot starts, and it never shares a runner with a third-party action. It has no expiry, so a leaked copy works until it is revoked. If it leaks, Cami revokes it on GitHub, and Dependabot's pull requests wait until a new one is stored.

**A canary for the required-checks ruleset could merge broken code.** The canary pull request is never merged, and nobody tries `gh pr merge` on it. The ruleset is proved by reading what GitHub reports, as the test strategy says.

**Other pull requests are blocked while the change is open.** Once "Required checks" exists, a pull request based on `main` has no `dependency-review.yml`, so that check never reports and it can't merge. That lasts until this change merges. If this change is abandoned, "Required checks" is turned off first.

**Grouped updates can hide a bad package among good ones.** If a group fails `make check`, an agent finds which package broke it and Dependabot is told to ignore that version.

**`make check` needs the network.** Offline, `vulncheck` fails. Running `make lint test e2e scenarios specs` covers everything else.

## Migration Plan

The order matters, because a required check that never reports would block every pull request.

1. The new workflows, configuration and documents land on the `harden-repository` branch, and its pull request runs them all.
2. With Cami's go-ahead, she creates the token, and it is stored as the Dependabot secret.
3. With the pull request open, its three checks having run once, and Cami's go-ahead, the "Required checks" ruleset is created. The change's own pull request is the first to meet it.
4. Canary pull requests from the branch prove that dependency review and the ruleset block what they should. Then they are closed.
5. Cami merges the change. Dependabot reads its configuration from `main` and opens its first pull requests.

To roll back, Cami turns off or deletes "Required checks" first, then reverts the pull request and revokes the token.

## Test strategy

This change has no specs, so no scenario needs a test. Each check is instead proved to fail on a known-bad input:

1. `govulncheck`: the offline fixture and `scripts/check-vulncheck-fixture.sh`, run by `make vulncheck` inside every `make check`. It is watched failing first, before the fixture exists, and again when the database marks a function the fixture doesn't call.
2. The new workflow files: `actionlint` in `make lint`, which already rejects a broken fixture.
3. Dependency review: a draft canary pull request that adds a package with a known advisory. The `dependency-review` job must fail. Checked by hand once, because it needs GitHub's advisory database.
4. The "Required checks" ruleset, checked by hand once, because it needs the live ruleset. The merge state alone proves nothing, because "Protect main" already blocks a pull request with no approval. So three things are read instead:
   1. The active rules for `main`, from `gh api repos/ithinkicam/chorale/rules/branches/main`, must include the new required status checks.
   2. `gh pr checks --required` on this change's pull request must list exactly the three checks.
   3. GitHub's `viewerCanMergeAsAdmin` field must be false for Cami on a canary pull request that breaks a golden file, and true on this change's pull request once it is green. The contrast shows that her bypass no longer skips the checks.
5. CodeQL and Scorecard: their first runs on the pull request, and on `main` after merging, must finish and upload results. Checked by hand, because their results live on GitHub.
6. Auto-merge: after merging, the first Dependabot patch or minor pull request must be approved from Cami's account and merge on its own once checks pass, and the push to `main` must start `check.yml`. A major update, or an update to one of the three named dependencies, must stay open. Checked by hand.
