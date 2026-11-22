# Tasks

Every behaviour task below comes in a pair: first the failing test, then the code. Run the test after writing it and confirm it fails for the expected reason before writing the code. Every test asserts its scenario's THEN clauses directly. A golden file is added as a record of the screen, never as the thing that makes a test pass.

## 1. Project setup

- [x] 1.1 Create `go.mod` for `github.com/ithinkicam/chorale` with `go 1.27` and a `toolchain` line naming the installed Go release, and have the `Makefile` export `GOTOOLCHAIN` from that line. Verify that `make` reports that toolchain through `go version`
- [x] 1.2 Add `charm.land/bubbletea/v2`, `charm.land/lipgloss/v2`, teatest v2, `creack/pty` and both emulator candidates at pinned versions. Confirm teatest v2 imports the `charm.land` Bubble Tea path, then verify `go mod tidy` and `go build ./...` succeed
- [x] 1.3 Create the package layout from `design.md` with a `cmd/chorale/main.go` that does nothing yet, add `/chorale` to `.gitignore`, and verify `make build` writes `bin/chorale`
- [x] 1.4 Add `package.json` with OpenSpec 2.1.1 as a development dependency, commit `package-lock.json`, and add `.nvmrc`. Verify `npm ci` succeeds, and that `npx --no-install openspec validate --strict --all` exits 0 when run with `CI=1` and no terminal
- [x] 1.5 Add `.golangci.yml` and a `Makefile` with `build`, `test`, `e2e`, `lint`, `golden` and `check` targets, where `lint` checks the golangci-lint version. Verify `make lint` and `make test` run and pass on the empty project

## 2. Terminal-level harness

- [x] 2.1 Write a test for the screen wrapper in `internal/termtest` that feeds it known escape sequences for the alternate screen, cursor visibility, SGR mouse reporting, bracketed paste and the keyboard enhancement flags, and asserts each mode is reported on and then off. Run it and confirm it fails because the wrapper doesn't exist yet
- [x] 2.2 Implement the wrapper over `github.com/charmbracelet/x/vt` and verify the test from 2.1 passes. If it can't track a mode, try `vt10x`. Record the choice and the reason in `design.md` and remove the unused library
- [x] 2.3 Write tests for the harness: that `TestMain` builds `chorale` into a temporary folder, that a program can be started at a given size with the fixed environment, and that the screen, cursor and raw bytes can be read. Also test the idle detector against the two fixture programs, one that redraws on a timer and must be caught, and one that draws once and must pass. Run them and confirm they fail because the helpers don't exist yet
- [x] 2.4 Implement the harness helpers, the colour code parser and the idle detector, and verify the tests from 2.3 pass
- [x] 2.5 Write a terminal-level test for the scenario "Quitting restores the terminal". Run it and confirm it fails at its GIVEN, because `chorale` never switches to the alternate screen yet

## 3. Version flag

- [x] 3.1 Write terminal-level tests for the scenarios "Version flag" and "Development build version", building one binary with the version set to `0.0.0-test` and one with none. Run them and confirm they fail
- [x] 3.2 Implement `internal/version`, the `--version` flag and the `-ldflags` in `make build`, and verify the tests from 3.1 pass

## 4. Empty session list

- [ ] 4.1 Write the plain-text screen helper for teatest, then interface snapshot tests for the scenarios "Empty session list at 60 by 20 with colour" and "Empty session list at 120 by 30 with colour", asserting each THEN clause. Run them and confirm they fail
- [ ] 4.2 Write terminal-level tests for the scenarios "Empty session list with colour off" and "Colour on uses the terminal palette", checking the raw bytes. Run them and confirm they fail
- [ ] 4.3 Write a terminal-level test for the scenario "Cursor on the empty message". Run it and confirm it fails
- [ ] 4.4 Implement the `sessionlist` screen and the app model, and verify the tests from 4.1 to 4.3 pass
- [ ] 4.5 Add golden comparisons of the plain-text screen to the tests from 4.1, plus one with colour off at 60 by 20. Write the golden files with `make golden` and read each one in the diff before committing. Then change one character in a golden file, confirm the test fails, and restore it
- [ ] 4.6 Write an interaction test that `q` on the session list ends the program. Run it and confirm it fails. Implement quitting, and verify both it and the test from 2.5 pass
- [ ] 4.7 Write terminal-level tests for the scenarios "Idle screen stays still" and "Startup draws the session list directly" and run them. If one fails, fix the stray redraw and verify it passes. If one passes on its first run, say in the commit that it is a regression guard, whose detector was seen to fail in task 2.3. Never add a redraw just to see it fail

## 5. Window too small

- [ ] 5.1 Write interface snapshot tests for the scenarios "Started at 40 columns", "Started at 15 rows" and "Message wraps in a tiny window", and interaction tests for "Window widened to 60 columns" and "Window narrowed while running". Write terminal-level tests for "Too-small message with colour off" and "Quitting from the too-small screen". Run them all and confirm they fail
- [ ] 5.2 Implement the too-small message and the switch on resize, and verify all seven tests from 5.1 pass
- [ ] 5.3 Add golden comparisons for the too-small screens at 40 by 20, 80 by 15 and 20 by 10, plus one with colour off, written with `make golden` and read in the diff before committing

## 6. Scenario check

- [ ] 6.1 Write unit tests for `scripts/check-scenarios`, with fixture specs and fixture `go test -json` output covering these cases: a scenario with a passing test, one with no test, one under `changes/archive/` that must be ignored, a name with digits, a matching subtest, a skipped test, a failed test, a golden file with no test, and a scenario marked "Checked by hand". Run them and confirm they fail
- [ ] 6.2 Implement the scenario check and verify its tests pass. Then run it on the real repo, and verify it reports all 17 scenarios in this change as covered
- [ ] 6.3 Add the scenario check and `openspec validate --strict --all` to `make check`. Verify `make check` passes from a terminal, and from a shell with `CI=1` and no terminal

## 7. Continuous integration

- [ ] 7.1 Add `.github/workflows/check.yml` as described under "Continuous integration" in `design.md`. Verify both jobs pass on the pull request for this change, and that they are named `check (macos)` and `check (linux)`
- [ ] 7.2 On a throwaway branch named `ci-canary`, commit a change that breaks one golden file, and open it as a draft pull request. Verify both jobs fail on it, then close the pull request and give Cami the command to delete the branch, because agents can't delete branches
- [ ] 7.3 Update `docs/testing.md` so its commands, the layers this change built, the colour code definition and the "Checked by hand" line match what was built. Verify each command runs as written

## Workflow follow-up

- An outside agent and our own adversarial reviewer review the change, then Cami reviews it before it is applied.
- The change is built on a branch and merged through a pull request. Cami confirms both continuous integration jobs are green before merging.
- Archiving is the last commit on the change's branch, once every task is checked and `make check` passes.
