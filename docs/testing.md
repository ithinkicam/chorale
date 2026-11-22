# Chorale: testing

As of 2026-11-22. This is the standard every change follows. Nothing is built yet, so the commands below are what milestone 0 sets up.

## The rule

All development is test first, including the interface. For each behaviour:

1. Write a test that describes it.
2. Run the test and watch it fail, for the reason you expect.
3. Write the smallest code that makes it pass.
4. Tidy the code with the test still passing.

A change is not done until `make check` passes. Tests are never skipped, disabled or loosened to get a change through. If a test is wrong, fixing it is its own small step, with the reason written in the commit.

## Where tests come from

Tests come from OpenSpec scenarios. Every `#### Scenario:` in a spec under `openspec/specs/`, or in an open change under `openspec/changes/`, has at least one automated test. The test's name includes the scenario's name, written in lower case with hyphens. For example, the scenario "Permission prompt at 60 columns" becomes a test called `TestPermissionPromptAt60Columns` or a golden file called `permission-prompt-at-60-columns.golden`.

A script, `scripts/check-scenarios`, lists every scenario and fails if any has no matching test. It runs as part of `make check`. A scenario that can only be checked by hand must say so in the spec, and the script then skips it by name.

## The test layers

There are six layers. Each requirement's design says which layer covers it.

### 1. Unit tests

Plain Go tests for logic with no screen: the client protocol, session state, vocabulary correction, contrast maths, sprite data. They are fast and run on every save.

### 2. Interface snapshot tests

The terminal interface is tested by rendering it and comparing the result to a stored golden file. These use `teatest` from Charm, which runs a Bubble Tea program in a test, sends it messages, and captures its final output. Its version 2 module, `github.com/charmbracelet/x/exp/teatest/v2`, exists as of 2026-11-22, though it is experimental.

Every screen is snapshotted:

1. At 60 columns by 20 rows, the narrowest supported size.
2. At 120 columns by 30 rows, where side panels appear.
3. In every shipped theme.
4. With colour turned off. This snapshot must still show every state in words, which is how the "never colour alone" rule is enforced.

Golden files live next to their tests in `testdata/`. They are regenerated only on purpose, with `make golden`, and every change to one is reviewed in the diff like code.

### 3. Interaction tests

Also built on `teatest`, these send key presses and mouse events to the program and check what it does. They cover the prompt editor in particular: a click places the cursor, a drag selects, a double-click selects a word, Shift-click extends, and undo restores. Each test asserts on the rendered screen and on the editor's text and selection.

### 4. Terminal-level tests

Some behaviour only shows up in the raw bytes between Chorale and a real terminal: mouse reporting in SGR mode, key releases through the Kitty keyboard protocol, copying through OSC 52, bracketed paste, Kitty graphics for the sprites, and where the cursor ends up. These tests run Chorale in a pseudo-terminal, send it raw escape sequences, and read the screen back through a terminal emulator library. Which library is decided in milestone 0.

These tests also enforce two accessibility rules:

1. **No motion.** Render the same state twice, a second apart. The output must be identical.
2. **Cursor placement.** After every screen update, the real cursor sits on the focused item or the insertion point.

### 5. End-to-end tests with a scripted agent

Real agents are slow and unpredictable, so tests use a fake one: a small program that speaks ACP and follows a script, such as "reply with this text, then ask permission to run this command". End-to-end tests start the real service and the real terminal app against the fake agent, and check the whole path: prompts, streamed replies, permission prompts, terminal handoff, and reconnects.

A small number of smoke tests run against the real Claude Code, Codex and Antigravity adapters. They are run by hand, not in continuous integration, because they need Cami's logins.

### 6. Speech and phone tests

The Swift helper `chorale-listen` has XCTest unit tests for everything except the microphone itself. Accuracy is measured against a small set of Cami's own recordings, which stay on her Mac and are never committed. The iOS companion will have XCUITest interface tests, run at large Dynamic Type sizes.

## Themes

Every fixed-colour theme is checked by a test that computes the WCAG contrast ratio of each text colour against its background. Body text must reach 7 to 1. The default theme uses the terminal's own palette, so it is checked against a few named reference palettes instead.

## Commands

These are the commands milestone 0 will create:

```
make test     # unit, snapshot and interaction tests
make e2e      # terminal-level and end-to-end tests
make lint     # gofmt, go vet and golangci-lint
make golden   # regenerate golden files, on purpose only
make check    # everything above, plus openspec validate --strict and the scenario check
```

Continuous integration runs `make check` on every push.
