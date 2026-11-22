# Tasks

This change has no specs, so no task maps to a scenario. Each new check is still built test first: its known-bad input is written first, and the check is watched failing on it before it is trusted. Every new workflow pins its actions to a full commit hash with the version in a comment, looked up when the task runs, and `make lint` must pass after each one. Steps that change settings on GitHub happen only with Cami's go-ahead at that step.

## 1. govulncheck

- [ ] 1.1 Write `scripts/check-vulncheck-fixture.sh`, in the same shape as `scripts/check-actionlint-fixture.sh`, failing unless `govulncheck` exits with status 3 on `testdata/vulncheck/` with `-db` pointing at the fixture's local database. Run it and confirm it fails because neither `govulncheck` nor the fixture exists yet
- [ ] 1.2 Add `govulncheck` as a Go tool with `go get -tool golang.org/x/vuln/cmd/govulncheck`, pinned in `go.mod`. Verify `go tool govulncheck -version` runs, and that `go mod tidy` leaves `go.mod` unchanged
- [ ] 1.3 Build the fixture in `testdata/vulncheck/`: a module that calls a function in a local module, and a local OSV database that marks that function as vulnerable. Verify the script from 1.1 now passes, then point the database at a function the fixture doesn't call and confirm the script fails, and restore it
- [ ] 1.4 Add a `vulncheck` target to the `Makefile` that runs the fixture script, then `go tool govulncheck ./...`, and add it to `check` after `lint`. Verify `make vulncheck` passes on the repo, and that `make check` passes from a shell with `CI=1` and no terminal
- [ ] 1.5 Confirm the dependency graph on GitHub raises no alert for `testdata/vulncheck/go.mod` once the branch is pushed. If it does, rename the fixture's manifest so the graph can't read it, and keep 1.3 passing

## 2. Workflows for scanning

- [ ] 2.1 Add `.github/workflows/codeql.yml` as described under "CodeQL uses a workflow file" in `design.md`. Verify `make lint` passes, and that both matrix jobs finish and upload results on this change's pull request
- [ ] 2.2 Add `.github/workflows/dependency-review.yml` with the job named `dependency-review`, as described in `design.md`. Verify `make lint` passes, and that the job passes on this change's pull request
- [ ] 2.3 Add `.github/workflows/scorecard.yml` within the limits that publishing places, as described in `design.md`. Verify `make lint` passes. Its first real run happens on `main` after merging, and is checked in the follow-up below

## 3. Dependabot and auto-merge

- [ ] 3.1 Add `.github/dependabot.yml` for `gomod`, `github-actions` and `npm`, weekly, with patch and minor grouped per ecosystem and a seven-day cooldown. Verify GitHub accepts it, by checking that the pull request shows no Dependabot configuration error
- [ ] 3.2 Add `.github/workflows/dependabot-auto-merge.yml` as described under "Auto-merge for patch and minor updates" in `design.md`. Verify `make lint` passes, and that its job is skipped, not failed, on this change's pull request, because the author isn't Dependabot
- [ ] 3.3 Show Cami the command that turns on "Allow GitHub Actions to create and approve pull requests", and with her go-ahead run it. Verify `gh api repos/ithinkicam/chorale/actions/permissions/workflow` reports `can_approve_pull_request_reviews` as true

## 4. Required checks

- [ ] 4.1 Export "Protect main" to `.github/rulesets/protect-main.json`, and verify it matches what `gh api repos/ithinkicam/chorale/rulesets/27852514` reports
- [ ] 4.2 Write `.github/rulesets/required-checks.json` as described under "The second ruleset holds the required checks" in `design.md`. Once this change's pull request has run all three checks once, show Cami the command that creates it, and with her go-ahead run it. Verify `gh api repos/ithinkicam/chorale/rulesets` lists both rulesets, and that "Required checks" has no bypass actors
- [ ] 4.3 On a throwaway branch named `required-checks-canary`, made from this branch, commit a change that breaks one golden file, and open it as a draft pull request. Verify both `check` jobs fail and `gh pr view --json mergeStateStatus` reports `BLOCKED`. Never try to merge it. Close it and give Cami the command to delete the branch
- [ ] 4.4 On a throwaway branch named `dependency-review-canary`, made from this branch, add a package with a known advisory to `go.mod`, and open it as a draft pull request. Verify the `dependency-review` job fails and names the advisory. Close it and give Cami the command to delete the branch

## 5. Community files and documents

- [ ] 5.1 Add `.github/CODEOWNERS`, `SECURITY.md` and `CONTRIBUTING.md` as described under "Community files" in `design.md`. Verify GitHub's community profile for the repo lists all three, and that the security policy link opens the private reporting page
- [ ] 5.2 Write `docs/repository.md` as described in `design.md`, with every setting it names checked against what `gh api` reports. Verify each command in it runs as written
- [ ] 5.3 Update the pull request rules in `AGENTS.md`, and point `README.md` and `docs/getting-started.md` to `docs/repository.md`. Verify every path they name exists
- [ ] 5.4 Update `STATE.md` with what this change built and with the follow-ups below, and verify `make check` and `openspec validate --strict --all` pass

## Workflow follow-up

- Archive the change once every task is done and `make check` passes, as the last commit on the branch. Cami merges the pull request once both rulesets are satisfied.
- After merging, check that CodeQL and Scorecard run on `main` and that their results appear on the security page.
- After merging, check that the first Dependabot patch or minor pull request is approved by `github-actions[bot]` and merges on its own once its checks pass. If the approval doesn't count toward "Protect main", record that in `docs/repository.md` and fall back to Cami approving by hand.
- When the first major update arrives, check that it stays open.
- Remove these follow-ups from `STATE.md` as each one is done.
