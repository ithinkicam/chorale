# Chorale: testing

As of 2026-11-22. This is the standard every change follows. Milestone 0 built layers 1 to 4, the scenario check and `make check`. Layers 5 and 6 arrive with the service in milestone 1 and with `chorale-listen` in milestone 2.

## The rule

All development is test first, including the interface. For each behaviour:

1. Write a test that describes it.
2. Run the test and watch it fail, for the reason you expect.
3. Write the smallest code that makes it pass.
4. Tidy the code with the test still passing.

A change is not done until `make check` passes. Tests are never skipped, disabled or loosened to get a change through. If a test is wrong, fixing it is its own small step, with the reason written in the commit.

## Where tests come from

Tests come from OpenSpec scenarios. Every `#### Scenario:` in a spec under `openspec/specs/`, or in an open change under `openspec/changes/`, has at least one automated test that passes. The test, or a subtest, is named after the scenario. For example, the scenario "Permission prompt at 60 columns" becomes a test called `TestPermissionPromptAt60Columns`.

`scripts/check-scenarios` checks this. It reads the `go test -json` results that `make test` and `make e2e` write to `bin/`. A name matches a scenario when the two are equal after removing a leading `Test`, keeping only letters and digits, and ignoring case. A scenario counts as covered only when a matching test passed. A skipped test, a failed test, or a golden file on its own doesn't count. Archived changes under `openspec/changes/archive/` are skipped, because their scenarios already live in `openspec/specs/`.

A scenario that can only be checked by hand carries this line in its spec, with the reason after it:

```
- **Checked by hand:** <reason>
```

The check lists those scenarios and doesn't fail on them.

## The test layers

There are six layers. Each requirement's design says which layer covers it.

### 1. Unit tests

Plain Go tests for logic with no screen: the client protocol, session state, vocabulary correction, contrast maths, sprite data. They are fast and run on every save.

### 2. Interface snapshot tests

The terminal interface is tested by running the real Bubble Tea program in a test. This uses `teatest` v2 from Charm, `github.com/charmbracelet/x/exp/teatest/v2`, which is experimental and pinned in `go.mod`. The helper `termtest.StartModel` feeds everything the program draws to a terminal emulator and reads back the screen as plain text.

Each test asserts its scenario's THEN clauses on that text directly, such as "the first line starts with ♪ chorale". A golden file of the same plain text records the whole screen as well, but it is never the only check. Golden files are plain text, so they can be read aloud.

Every screen is snapshotted:

1. At 60 columns by 20 rows, the narrowest supported size.
2. At 120 columns by 30 rows, where side panels appear.
3. In every shipped theme.
4. With colour turned off. This snapshot must still show every state in words, which is how the "never colour alone" rule is enforced.

Golden files live next to their tests in `testdata/`, named after the test. They are regenerated only on purpose, with `make golden`, and every change to one is read in the diff like code.

### 3. Interaction tests

Also built on `teatest`, these send key presses and mouse events to the program and check what it does. They cover the prompt editor in particular: a click places the cursor, a drag selects, a double-click selects a word, Shift-click extends, and undo restores. Each test asserts on the rendered screen and on the editor's text and selection.

### 4. Terminal-level tests

Some behaviour only shows up in the raw bytes between Chorale and a real terminal: mouse reporting in SGR mode, key releases through the Kitty keyboard protocol, copying through OSC 52, bracketed paste, Kitty graphics for the sprites, colour codes, terminal modes, and where the cursor ends up. These tests live in `e2e/` and run with `make e2e`.

Each run builds `chorale` fresh from the current code. `termtest.Start` runs it in a pseudo-terminal with a fixed environment: `TERM=xterm-256color`, a UTF-8 locale, a temporary `HOME`, and `NO_COLOR` only when a test asks for it. Everything it writes goes to the terminal emulator `github.com/charmbracelet/x/vt`. The harness reads back the screen, the cursor and every terminal mode. It tracks cursor visibility and keyboard enhancements itself, because `x/vt` keeps cursor visibility per screen and ignores keyboard enhancements.

A colour code is a Select Graphic Rendition sequence with parameter 30 to 39, 40 to 49, 90 to 97 or 100 to 107, or an extended form starting with 38 or 48. Tests with colour off check the raw bytes contain none.

These tests also enforce three accessibility rules:

1. **No motion.** Once the program has written nothing for half a second, it must write zero bytes over the next two seconds. Counting bytes catches a program that clears and redraws the same picture, which comparing screens would miss. The detector is itself tested against a fixture program that redraws on a timer.
2. **Cursor placement.** After every screen update, the real cursor is visible and sits on the focused item or the insertion point.
3. **A clean exit.** Every terminal mode turned on during a run is off at exit, and the cursor is visible.

### 5. End-to-end tests with a scripted agent

Real agents are slow and unpredictable, so tests use a fake one: a small program that speaks ACP and follows a script, such as "reply with this text, then ask permission to run this command". End-to-end tests start the real service and the real terminal app against the fake agent, and check the whole path: prompts, streamed replies, permission prompts, terminal handoff, and reconnects.

A small number of smoke tests run against the real Claude Code, Codex and Antigravity adapters. They are run by hand, not in continuous integration, because they need Cami's logins.

### 6. Speech and phone tests

The Swift helper `chorale-listen` has XCTest unit tests for everything except the microphone itself. Accuracy is measured against a small set of Cami's own recordings, which stay on her Mac and are never committed. The iOS companion will have XCUITest interface tests, run at large Dynamic Type sizes.

## Themes

This check arrives with the first fixed-colour theme in milestone 1. Every fixed-colour theme is checked by a test that computes the WCAG contrast ratio of each text colour against its background. Body text must reach 7 to 1. The default theme uses the terminal's own palette, so it is checked against a few named reference palettes instead.

## Commands

```
make build      # build bin/chorale, stamped with the git version
make test       # unit, snapshot and interaction tests
make e2e        # terminal-level tests
make lint       # golangci-lint version check, gofmt, go vet and golangci-lint
make golden     # rewrite golden files, on purpose only
make scenarios  # make test and make e2e, then the scenario check
make specs      # openspec validate --strict --all, with OpenSpec from package-lock.json
make check      # all of the above except build and golden
```

`make` pins Go to the release named on the `toolchain` line of `go.mod`. `make lint` fails unless golangci-lint is the version the `Makefile` names. OpenSpec and its dependencies come from `package-lock.json`, through `npm ci`.

Continuous integration runs `make check` on macOS and Linux, on every push to `main` and on every pull request. The workflow is `.github/workflows/check.yml`.
