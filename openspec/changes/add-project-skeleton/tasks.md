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

- [x] 4.1 Write the plain-text screen helper for teatest, then interface snapshot tests for the scenarios "Empty session list at 60 by 20 with colour" and "Empty session list at 120 by 30 with colour", asserting each THEN clause. Run them and confirm they fail
- [x] 4.2 Write terminal-level tests for the scenarios "Empty session list with colour off" and "Colour on uses the terminal palette", checking the raw bytes. Run them and confirm they fail
- [x] 4.3 Write a terminal-level test for the scenario "Cursor on the empty message". Run it and confirm it fails
- [x] 4.4 Implement the `sessionlist` screen and the app model, and verify the tests from 4.1 to 4.3 pass
- [x] 4.5 Add golden comparisons of the plain-text screen to the tests from 4.1, plus one with colour off at 60 by 20. Write the golden files with `make golden` and read each one in the diff before committing. Then change one character in a golden file, confirm the test fails, and restore it
- [x] 4.6 Write an interaction test that `q` on the session list ends the program. Run it and confirm it fails. Implement quitting, and verify both it and the test from 2.5 pass
- [x] 4.7 Write terminal-level tests for the scenarios "Idle screen stays still" and "Startup draws the session list directly" and run them. If one fails, fix the stray redraw and verify it passes. If one passes on its first run, say in the commit that it is a regression guard, whose detector was seen to fail in task 2.3. Never add a redraw just to see it fail

## 5. Window too small

- [x] 5.1 Write interface snapshot tests for the scenarios "Started at 40 columns", "Started at 15 rows" and "Message wraps in a tiny window", and interaction tests for "Window widened to 60 columns" and "Window narrowed while running". Write terminal-level tests for "Too-small message with colour off" and "Quitting from the too-small screen". Run them all and confirm they fail
- [x] 5.2 Implement the too-small message and the switch on resize, and verify all seven tests from 5.1 pass
- [x] 5.3 Add golden comparisons for the too-small screens at 40 by 20, 80 by 15 and 20 by 10, plus one with colour off, written with `make golden` and read in the diff before committing

## 6. Scenario check

- [x] 6.1 Write unit tests for `scripts/check-scenarios`, with fixture specs and fixture `go test -json` output covering these cases: a scenario with a passing test, one with no test, one under `changes/archive/` that must be ignored, a name with digits, a matching subtest, a skipped test, a failed test, a golden file with no test, and a scenario marked "Checked by hand". Run them and confirm they fail
- [x] 6.2 Implement the scenario check and verify its tests pass. Then run it on the real repo, and verify it reports all 17 scenarios in this change as covered
- [x] 6.3 Add the scenario check and `openspec validate --strict --all` to `make check`. Verify `make check` passes from a terminal, and from a shell with `CI=1` and no terminal

## 7. Continuous integration

- [x] 7.1 Add `.github/workflows/check.yml` as described under "Continuous integration" in `design.md`. Verify both jobs pass on the pull request for this change, and that they are named `check (macos)` and `check (linux)`
- [x] 7.2 On a throwaway branch named `ci-canary`, commit a change that breaks one golden file, and open it as a draft pull request. Verify both jobs fail on it, then close the pull request and give Cami the command to delete the branch, because agents can't delete branches
- [x] 7.3 Update `docs/testing.md` so its commands, the layers this change built, the colour code definition and the "Checked by hand" line match what was built. Verify each command runs as written

## 8. Style guide and workflow checks

- [x] 8.1 Write `docs/style.md` from the research into existing style guides for agents and for Go, and point to it from `AGENTS.md`. Cami reviews it before task 8.3 starts
- [x] 8.2 Add `actionlint` as a Go tool pinned in `go.mod`, and run it in `make lint` on `.github/workflows/`. Verify it rejects a fixture workflow with an invalid permission and a misspelt matrix property, then passes on `.github/workflows/check.yml`
- [x] 8.3 Remove the comments that break `docs/style.md` from the code in this change, and add any lint settings the guide calls for. Verify `make check` passes and `make golden` leaves every golden file unchanged

## 9. Review fixes

- [x] 9.1 Add `-race` to `make test` and `make e2e`, run them and confirm the harness race fails them. Then stop the emulator's reply reader before closing the emulator, close the screens in snapshot tests, and send query replies through their own writer so a program that doesn't read its input can't hang a test. Verify both pass under `-race`
- [x] 9.2 Write screen tests that a DEC mode outside the current list, such as 2027, is reported when left on, and that a mode turned on and off within one write is still recorded. Run them and confirm they fail, then report every mode that differs from the terminal's default and record modes as they arrive
- [x] 9.3 Write a screen test that the cursor style is tracked, and a terminal-level test for the scenario "Cursor does not blink". Run them and confirm they fail with `cursor.Blink = false` removed, then restore it and reset the cursor style at exit, and verify both pass
- [x] 9.4 Write an interface snapshot test for the scenario "Message comes first in a very short window", run it and confirm it fails, then let the message take priority over the key line, and verify it passes. Fix the wrapped-message check to keep real spacing, and measure line widths on the running program instead of a fresh model
- [x] 9.5 Write scenario check tests for duplicate scenario names, a scenario heading inside a code block, and a relative OpenSpec path. Run them and confirm they fail, then fix the check and verify they pass
- [ ] 9.6 Clean up the test runs and tooling: remove the temporary folders the tests create, add `--ignore-scripts` to `npm ci`, fix `scripts/check-gofmt.sh` for deleted files, add a test that the actionlint fixture still fails, and list gotestsum in the proposal. Verify `make check` passes and the temporary folder count doesn't grow
- [ ] 9.7 Fix the docs that disagree with the code: the getting-started guide's install steps, with a pinned golangci-lint and `npm ci`, the cursor check in `docs/testing.md`, where golden files live, and the lines that say `make check` doesn't exist yet. Remove the comments that still break `docs/style.md`, and make failure messages say what was wanted. Verify each command in the getting-started guide runs as written

## Workflow follow-up

- An outside agent and our own adversarial reviewer review the change, then Cami reviews it before it is applied.
- The change is built on a branch and merged through a pull request. Cami confirms both continuous integration jobs are green before merging.
- Archiving is the last commit on the change's branch, once every task is checked and `make check` passes.
