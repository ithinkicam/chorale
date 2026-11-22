# Tasks

Every behaviour task below comes in a pair: first the failing test, then the code. Run the test after writing it and confirm it fails for the expected reason before writing the code.

## 1. Project setup

- [ ] 1.1 Create `go.mod` for `github.com/ithinkicam/chorale` with the Go 1.27 floor and a pinned toolchain, and verify `go version` inside the repo reports that toolchain
- [ ] 1.2 Add Bubble Tea v2, Lip Gloss v2 and teatest v2 at pinned versions, and verify `go mod tidy` and `go build ./...` succeed
- [ ] 1.3 Create the package layout from `design.md` with a `cmd/chorale/main.go` that does nothing yet, and verify `go build ./cmd/chorale` succeeds
- [ ] 1.4 Add `.golangci.yml` and a `Makefile` with `test`, `e2e`, `lint`, `golden` and `check` targets, and verify `make lint` and `make test` run and pass on the empty project
- [ ] 1.5 Ignore `openspec/changes/archive/` in `.gitignore`, keeping its `.gitkeep`, and verify with `git check-ignore` that a file inside it is ignored and `.gitkeep` is not

## 2. Terminal-level harness

- [ ] 2.1 Write a terminal-level test for the scenario "Quitting restores the terminal", using `creack/pty` with the first candidate emulator library, run it and confirm it fails because `chorale` doesn't open the interface yet
- [ ] 2.2 Check the chosen emulator library reports alternate screen, cursor visibility and mouse mode correctly. If it doesn't, try the second candidate, then record the choice and the reason in `design.md`
- [ ] 2.3 Move the harness into `internal/termtest` with helpers to start `chorale` at a given size, send bytes, read the screen and read the cursor, and verify the failing test from 2.1 still fails for the same reason

## 3. Version flag

- [ ] 3.1 Write a terminal-level test for the scenario "Version flag", run it and confirm it fails
- [ ] 3.2 Implement `--version` in `cmd/chorale`, and verify the test from 3.1 passes

## 4. Empty session list

- [ ] 4.1 Write interface snapshot tests for the scenarios "Empty session list at 60 by 20 with colour" and "Empty session list at 120 by 30 with colour", with empty golden files, run them and confirm they fail
- [ ] 4.2 Write an interface snapshot test for the scenario "Empty session list with colour off", plus a unit test that the output contains no colour escape codes, run them and confirm they fail
- [ ] 4.3 Implement the `sessionlist` screen and the app model, and verify the three snapshot tests pass after `make golden` writes the golden files. Review each golden file by eye in the diff before committing
- [ ] 4.4 Write an interaction test for `q` quitting, run it and confirm it fails, then implement quitting and verify both it and the terminal-level test from 2.1 pass
- [ ] 4.5 Write a terminal-level test for the scenario "Idle screen stays still" and run it. If it fails, fix any redraw that happens without input and verify it passes. If it passes on the first run, record in the commit that it is a regression guard. Never add a redraw just to see it fail

## 5. Window too narrow

- [ ] 5.1 Write interface snapshot tests for the scenarios "Started at 40 columns" and "Started at 15 rows", and an interaction test for "Window widened to 60 columns", run them and confirm they fail
- [ ] 5.2 Implement the too-narrow message and the redraw on resize, and verify both tests pass

## 6. Scenario check

- [ ] 6.1 Write a unit test for `scripts/check-scenarios` with a fixture where one scenario has a test, one doesn't, and one sits under `changes/archive/` with no test and must be ignored, run it and confirm it fails
- [ ] 6.2 Implement the scenario check and verify its test passes, then run it on the real repo and verify it reports every scenario in this change as covered
- [ ] 6.3 Add it and `openspec validate --strict` to `make check`, and verify `make check` passes

## 7. Continuous integration

- [ ] 7.1 Add `.github/workflows/check.yml` as described under "Continuous integration setup" in `design.md`, and verify both jobs pass on the pull request for this change
- [ ] 7.2 On a throwaway branch named `ci-canary`, push a commit that deliberately breaks one golden file, and verify both jobs fail on it. Then give Cami the command to delete the branch, because agents can't delete branches
- [ ] 7.3 Update `docs/testing.md` so its commands match what was built, and verify each command runs as written

## Workflow follow-up

- An outside agent and our own adversarial reviewer review the change, then Cami reviews it before it is applied.
- The change is built on a branch and merged through a pull request.
- Archive the change once every task is checked and `make check` passes.
