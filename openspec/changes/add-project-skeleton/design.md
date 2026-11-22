# Design

## Context

Nothing exists yet but the plan, the testing standard and OpenSpec itself. See `proposal.md` for why the skeleton comes first. The stack was chosen in `docs/plan.md`: Go, with Bubble Tea v2 and Lip Gloss v2 for the interface.

## Goals / Non-Goals

**Goals:**

- Test layers 1 to 4 in `docs/testing.md` exist, each with at least one real test. Those are unit, interface snapshot, interaction and terminal-level tests.
- `make check` is the single gate, the same locally and in continuous integration.
- The empty session list is built test first, as the worked example for every later change.

**Non-Goals:**

- The end-to-end layer with a scripted agent, which needs the service and arrives in milestone 1.
- The speech and phone layer, which arrives with `chorale-listen` in milestone 2.
- The theme contrast check, which arrives with the first fixed-colour theme in milestone 1.
- No service, agents, themes beyond the default one, or prompt editor.
- No release packaging or Homebrew formula.

## Decisions

### Go version and module

The module is `github.com/ithinkicam/chorale`. `go.mod` declares `go 1.27` as the floor, and a `toolchain` line names the exact Go release installed when task 1.1 runs. A `toolchain` line is only a minimum, so a newer local Go would otherwise be used instead. To make it a real pin, the `Makefile` reads that line and exports `GOTOOLCHAIN` set to it. Continuous integration installs the same release through `actions/setup-go` with `go-version-file: go.mod`.

### Dependencies

The interface libraries use their current module paths: `charm.land/bubbletea/v2` and `charm.land/lipgloss/v2`. teatest v2 is `github.com/charmbracelet/x/exp/teatest/v2`, at a pinned pseudo-version, and it imports the `charm.land` path for Bubble Tea. Task 1.2 confirms this, because importing the older `github.com/charmbracelet/bubbletea` path would give types that don't match teatest's. The test-only dependencies are `github.com/creack/pty` and one terminal emulator library, chosen in task 2.2.

### Layout

```
cmd/chorale/          the command: flag parsing, then hands off to internal/app
internal/app/         the Bubble Tea model for the whole app
internal/screens/     one package per screen; this change adds sessionlist
internal/version/     the version string
internal/termtest/    the pseudo-terminal test harness and its fixture programs
e2e/                  terminal-level tests, run by make e2e
scripts/              check-scenarios and other development scripts
testdata/             shared fixtures
```

Golden files live in the `testdata/` folder of the package whose tests write them. For this change that is `internal/app`, because the screens are tested through the whole app. Built binaries go to `bin/`, which is ignored, and `.gitignore` also ignores a stray `/chorale` at the repo root.

### The version

`internal/version` holds `var Version = "dev"`. Builds that set a version do it with `-ldflags "-X github.com/ithinkicam/chorale/internal/version.Version=<version>"`. `make build` sets it from `git describe --tags --always --dirty`. The terminal-level tests build one binary with the version set to `0.0.0-test` and one with no version set, so they can check both scenarios.

### Interface tests use teatest v2, compared as plain text

Snapshot and interaction tests use teatest v2. It runs the real model in the real program loop, and sends it messages such as window sizes and keys. The alternative, calling the model's `View` directly, would skip the program loop and miss bugs in how messages are handled.

teatest captures the whole byte stream, with every frame and cursor movement in it. That makes a poor thing to compare or to review. So a test helper feeds the stream through the terminal emulator library and reads back the final screen as plain text, one line per row with trailing spaces removed. Each test asserts every THEN clause of its scenario on that text, such as "the first line starts with ♪ chorale". A golden file of the same plain text is compared as well. It records what the whole screen looks like, but it is never the only check.

Tests wait for expected text with teatest's `WaitFor`, with a time limit, instead of waiting for the program to finish. They end the program from the test side, so snapshot tests don't depend on the `q` key working. This also means a test can't hang when the model never quits.

Golden files are written only by `make golden`, which passes the `-update` flag to just the packages that have golden tests. A changed golden file is reviewed in the diff. Because golden files are plain text, they can be read aloud.

### Terminal-level tests use a real pseudo-terminal

The terminal-level tests live in `e2e/` and run with `make e2e`, which is part of `make check`. In `TestMain`, `internal/termtest` builds `chorale` with `go build -o` into a temporary folder, so every run tests the current code and never a stale binary.

The harness starts the binary in a pseudo-terminal with `creack/pty`, at a given size, with a fixed environment: `TERM=xterm-256color`, a fixed UTF-8 locale that exists on both macOS and Linux, and `NO_COLOR` set or unset as the test needs. It records every byte the program writes, and feeds them to the emulator library. From the emulator it reads the screen, the cursor position and visibility, and the state of each terminal mode. Widths are measured on the emulator's grid, so the width of `♪` is whatever the emulator says, the same on both platforms.

"Quitting restores the terminal" checks every mode that was turned on during the run, not a fixed list. The harness records each mode the program turns on and checks each one is off at exit. Its GIVEN first asserts the alternate screen is on, so a program that never changed the terminal can't pass.

### Choosing the emulator library

The candidates are `github.com/charmbracelet/x/vt` and `github.com/hinshun/vt10x`. The choice is made by feeding each one known escape sequences directly and checking it reports each mode correctly. The sequences cover the alternate screen, cursor visibility, mouse reporting in SGR mode, bracketed paste and the keyboard enhancement flags. That way the choice doesn't depend on Chorale producing those sequences yet. If a library tracks modes only through callbacks, the harness keeps its own record from them.

The harness uses `github.com/charmbracelet/x/vt` at v0.1.0. It tracks the screen, the cursor position, the alternate screen, mouse reporting, SGR mouse encoding, bracketed paste and focus events, and answers queries such as device attributes. It falls short in two places, and the wrapper in `internal/termtest` covers both:

1. It keeps cursor visibility separately for the main and alternate screens. Real terminals keep it once for the whole terminal, so a program that hides the cursor on the alternate screen and quits leaves it hidden. The wrapper follows DEC mode 25 itself, so that bug is caught.
2. It ignores the Kitty keyboard protocol and xterm's modifyOtherKeys. The wrapper parses those sequences itself.

`vt10x` was not needed. It has had no release since 2022, and it doesn't track these modes either.

### Detecting motion

The idle check counts bytes, not screens. Comparing two screens a second apart passes when a program clears and redraws the same picture, which is exactly the flicker the accessibility rules forbid. The harness waits until the program has written nothing for half a second, then asserts it writes zero bytes over the next two seconds.

The detector is tested against two small fixture programs in `internal/termtest/testdata/`. One redraws on a timer and must be caught. The other draws once and must pass. That proves the detector works before it is used on Chorale.

### What counts as a colour code

A colour code is a Select Graphic Rendition sequence with parameter 30 to 39, 40 to 49, 90 to 97 or 100 to 107, or an extended form starting with 38 or 48. The harness parses the raw bytes for these. The no-colour tests check the real bytes reaching the pseudo-terminal, because Lip Gloss decides what colour to emit only when it writes to a real output. The colour-on test checks that colour codes do appear and that all of them are palette colours, so the colour-on and colour-off cases can't pass identically.

### The scenario check

`scripts/check-scenarios` is a small Go program, run with `go run`. It reads every `#### Scenario:` heading in `openspec/specs/` and in each open change under `openspec/changes/`. It skips `openspec/changes/archive/`, because an archived change's scenarios already live in `openspec/specs/`.

It reads the `go test -json` output that `make test` and `make e2e` write to `bin/`. A scenario is covered only when a test or subtest passed whose name matches it. Names match when they are equal after removing a leading `Test`, keeping only letters and digits, and ignoring case. So "Empty session list at 60 by 20 with colour" matches `TestEmptySessionListAt60By20WithColour`. A skipped test, a failed test, a golden file alone, or a test that never ran doesn't count.

A scenario that can only be checked by hand carries this line in its spec, with the reason after it:

```
- **Checked by hand:** <reason>
```

The check lists those scenarios and doesn't fail on them. It prints every uncovered scenario and exits non-zero if there are any.

### OpenSpec in the gate

OpenSpec 2.1.1 is a development dependency in `package.json`, with `package-lock.json` committed, so its own dependencies are pinned too. The Node version is in `.nvmrc`, at 22.20.0 or later, as OpenSpec requires. `npm ci` installs it locally and in continuous integration. The `Makefile` runs `npx --no-install openspec validate --strict --all` with `OPENSPEC_TELEMETRY=0`.

The `--all` matters. Without an item name or `--all`, OpenSpec 2.1.1 exits with "Nothing to validate" when there is no terminal, which includes continuous integration and agents' shells. In a real terminal, it stops to ask a question instead.

### Linting

golangci-lint's version is set once in the `Makefile`. `make lint` checks that the installed version matches and fails with a clear message if it doesn't. Continuous integration installs that version.

### Continuous integration

`make check` runs on `macos-latest`, because that is where Chorale runs, and on `ubuntu-latest`, because it is faster and catches anything that accidentally depends on macOS. The repo is public, so minutes on both are free.

The workflow in `.github/workflows/check.yml` is set up as follows:

1. It runs on pushes to `main` and on pull requests. Branch work is checked through its pull request, so each commit is checked once rather than twice.
2. Every third-party action is pinned to a full commit hash, with the version in a comment beside it, so a moved tag can't change what runs.
3. The workflow token has `contents: read` and nothing more. Checkout sets `persist-credentials: false`. The workflow never uses `pull_request_target`, so a contributor's pull request never runs with write access or secrets.
4. Go comes from `go.mod` through `actions/setup-go`, Node from `.nvmrc` through `actions/setup-node`, and both caches are on.
5. Runs are grouped by pull request number, or by branch for pushes to `main`, and a newer run cancels the older one in its group. Every job has a time limit of 15 minutes.
6. The job names are set explicitly to `check (macos)` and `check (linux)`, because the `harden-repository` change makes them required checks.

Bypassing a ruleset skips every rule in it. So `harden-repository` adds the required checks in a second ruleset that nobody can bypass. Until then, Cami confirms both jobs are green before she merges.

### Code style

`docs/style.md` is the code style for every agent, and `AGENTS.md` points to it. Its first rule is that code explains itself through names, small functions and types, while the reasons live in specs, design documents and `AGENTS.md`. A comment is kept only for what the code can't say. The guide builds on existing guides rather than starting from nothing, and where a rule can be checked by a linter, golangci-lint checks it.

`actionlint` checks the workflow files as part of `make lint`. It is a Go tool pinned in `go.mod`, so the version is the same everywhere with no separate check. It understands GitHub Actions, so it catches invalid permissions, misspelt matrix and step references, and malformed values, which a general YAML linter would miss. It runs `shellcheck` on `run` steps where `shellcheck` is installed. It can't check the inputs of actions pinned by commit hash, so a misspelt input to a pinned action is only caught when the workflow runs.

### Archiving

Archived changes stay tracked in git, so their design decisions stay readable in the tree. Archiving is the last commit on the change's branch, before the pull request merges, so the merged pull request includes the updated `openspec/specs/`.

## Test strategy

Each requirement in `specs/app-shell/spec.md` and the layer that covers it:

1. **Empty session list on first start:** interface snapshot tests at 60 by 20 and 120 by 30, asserting each THEN clause on the plain-text screen. The cursor is checked by a terminal-level test.
2. **Readable without colour:** terminal-level tests on the real bytes, with `NO_COLOR` set and with colour on, plus a plain-text golden file with colour off.
3. **No motion:** terminal-level tests using the byte-counting idle detector, which is itself tested against fixture programs.
4. **Quit with q:** an interaction test that `q` ends the program, plus terminal-level tests of the exit status and every terminal mode, from both the session list and the too-small screen.
5. **Report the version:** terminal-level tests against binaries built with and without a version.
6. **Window too small:** interface snapshot tests at 40 by 20, 80 by 15 and 20 by 10, interaction tests that send resize messages in both directions, and a terminal-level test with colour off.

The scenario check itself has unit tests with fixture specs and fixture `go test -json` output.

## Risks / Trade-offs

- **teatest v2 is experimental and may change.** The version is pinned, and only the helpers in `internal/` call it directly, so an update touches one place.
- **Golden files can be regenerated carelessly.** They are never the only check, `make golden` is the only way to write them, and `AGENTS.md` forbids regenerating them to make a test pass.
- **The two platforms can differ in small ways, such as the width of `♪`.** Tests set a fixed environment and measure on the emulator's grid. The first failing difference is investigated, never ignored.
- **Timing in terminal-level tests can be flaky.** Waits have time limits and wait for content rather than sleeping a fixed time. The idle check waits for output to settle before it starts counting.

## Open Questions

- Whether to adopt Charm's VHS tool for recorded demo tapes. It isn't a gate, so it can be decided later.
