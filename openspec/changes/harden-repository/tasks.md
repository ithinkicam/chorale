# Tasks

This change has no specs, so no task maps to a scenario. Each new local check is still built test first: its known-bad input is written first, and the check is watched failing on it before it is trusted. Every new workflow pins its actions to a full commit hash with the version in a comment, looked up when the task runs, and `make lint` must pass after each one. Steps that change settings on GitHub, create a credential or add to Cami's calendar happen only with her go-ahead at that step.

## 1. govulncheck

- [x] 1.1 Write `scripts/check-vulncheck-fixture.sh`, in the same shape as `scripts/check-actionlint-fixture.sh`, as described in `design.md` under "govulncheck runs in make check". It runs from the repo root with `-C testdata/vulncheck` and a quoted, absolute `file://` URL for the fixture database, and fails unless `govulncheck` exits with status 3. Run it and confirm it fails because neither `govulncheck` nor the fixture exists yet
- [x] 1.2 Add `govulncheck` as a Go tool with `go get -tool golang.org/x/vuln/cmd/govulncheck`, pinned in `go.mod`. Verify `go tool govulncheck -version` runs, and that `go mod tidy` leaves `go.mod` unchanged
- [x] 1.3 Build the fixture in `testdata/vulncheck/`: a module with no dependencies that calls `fmt.Println`, and a local database in the Go vulnerability database layout whose one entry marks `fmt.Println` in `stdlib` as vulnerable in every version. Verify the script from 1.1 now passes. Then change the entry to a `fmt` function the fixture doesn't call, confirm the script fails, and restore it
- [x] 1.4 Add a `vulncheck` target to the `Makefile` that runs the fixture script and then `go tool govulncheck ./...`, and add it to `check` after `lint`. Verify `make vulncheck` passes on the repo, and that `make check` passes from a shell with `CI=1` and no terminal
- [x] 1.5 Update `docs/testing.md` so its commands and its list of what `make check` runs include `make vulncheck`, and say that `make check` needs the network. Add the same to the gates section of `AGENTS.md`. Verify each command named runs as written

## 2. Workflows for scanning

- [x] 2.1 Add `.github/workflows/codeql.yml` as described under "CodeQL uses a workflow file" in `design.md`, with Go built by `go build ./...`. Verify `make lint` passes, and that both matrix jobs finish and upload results on this change's pull request
- [x] 2.2 Add `.github/workflows/dependency-review.yml` with the job named `dependency-review`, as described in `design.md`. Verify `make lint` passes, and that the job passes on this change's pull request
- [ ] 2.3 Add `.github/workflows/scorecard.yml` as described in `design.md`, checked against the publishing limits in the action's README at the pinned version. Verify `make lint` passes, and start it once by hand on the branch to confirm it runs. Publishing from `main` is checked in the follow-ups

## 3. Dependabot and auto-merge

- [ ] 3.1 Add `.github/dependabot.yml` for `gomod`, `github-actions` and `npm`, weekly, with patch and minor grouped per ecosystem, a seven-day cooldown, and the three named dependencies left out of the groups. Verify the pull request shows no Dependabot configuration error
- [x] 3.2 Add `.github/workflows/dependabot-auto-merge.yml` with the `decide` and `merge` jobs described under "Auto-merge with a token kept as a Dependabot secret" in `design.md`. Verify `make lint` passes, that `merge` uses no actions, and that both jobs are skipped, not failed, on this change's pull request, because the author isn't Dependabot
- [x] 3.3 With Cami's go-ahead, walk her through creating the fine-grained token as `design.md` describes it, and store it as the Dependabot secret `DEPENDABOT_AUTOMERGE_TOKEN` with `gh secret set --app dependabot`. Verify `gh secret list --app dependabot` lists it. The token itself never appears in the repo, a log or this conversation
- [ ] 3.4 With Cami's go-ahead, put a reminder on her calendar a week before the token expires, saying where `docs/repository.md` describes renewing it. Verify the event exists on that date
- [x] 3.5 Verify that `gh api repos/ithinkicam/chorale/actions/permissions/workflow` still reports `can_approve_pull_request_reviews` as false

## 4. Required checks

- [x] 4.1 Write `.github/rulesets/protect-main.json` with the six fields named in `design.md`, from the live "Protect main". Verify it equals the output of `gh api repos/ithinkicam/chorale/rulesets/27852514` passed through the same `jq` filter
- [ ] 4.2 Write `.github/rulesets/required-checks.json` as described under "The second ruleset holds the required checks" in `design.md`. Once this change's pull request has run all three checks once, show Cami the command that creates it, and with her go-ahead run it. Verify the live ruleset, through the same `jq` filter, equals the file, and that `gh api repos/ithinkicam/chorale/rules/branches/main` includes its required status checks
- [ ] 4.3 Verify `gh pr checks --required` on this change's pull request lists exactly `check (macos)`, `check (linux)` and `dependency-review`
- [ ] 4.4 On a throwaway branch named `required-checks-canary`, made from this branch, commit a change that breaks one golden file, and open it as a pull request, marked ready for review. Verify both `check` jobs fail, and that GitHub's `viewerCanMergeAsAdmin` is false for it and true for this change's pull request once that is green. Never try to merge the canary. Close it and give Cami the command to delete the branch
- [ ] 4.5 On a throwaway branch named `dependency-review-canary`, made from this branch, add a package with a known advisory to `go.mod`, and open it as a draft pull request. Verify the `dependency-review` job fails and names the advisory. Close it and give Cami the command to delete the branch

## 5. Community files and documents

- [ ] 5.1 Add `.github/CODEOWNERS`, `SECURITY.md` and `CONTRIBUTING.md` as described under "Community files" in `design.md`. Verify GitHub's community profile for the repo lists all three, and that the security policy link opens the private reporting page
- [x] 5.2 Write `docs/repository.md` as described under "Documents" in `design.md`, with every setting it names checked against what `gh api` reports. Verify each command in it runs as written, except the ones that turn a ruleset off, which are checked for correct syntax only
- [x] 5.3 Update the pull request rules in `AGENTS.md`, `docs/getting-started.md` and `README.md` as described under "Documents" in `design.md`. Verify every path they name exists
- [ ] 5.4 Update `STATE.md` with what this change built and with the follow-ups below, and verify `make check` and `openspec validate --strict --all` pass

## Workflow follow-up

- Archive the change once every task is done and `make check` passes, as the last commit on the branch. Cami merges the pull request once both rulesets are satisfied.
- After merging, check that CodeQL and Scorecard run on `main` and that their results appear on the security page.
- After merging, check that the first Dependabot patch or minor pull request is approved from Cami's account, merges on its own once its checks pass, and that the merge starts `check.yml` on `main`. If the extra approval rule for unattributed changes holds it, bring that to Cami.
- When the first major update, or an update to one of the three named dependencies, arrives, check that it stays open, and bring Cami the summary `design.md` describes.
- Remove these follow-ups from `STATE.md` as each one is done.
