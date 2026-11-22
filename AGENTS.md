# Chorale: working rules for agents

Read this before doing any work in this repo. It applies to every agent: Claude Code, Codex and Antigravity CLI.

## What Chorale is

Chorale is a voice-first terminal workbench on macOS for running coding agents side by side, built for low vision. The roadmap is `docs/plan.md`. The test standard is `docs/testing.md`. The interface sketches are linked from the plan.

## All work goes through OpenSpec

Chorale uses OpenSpec for spec-driven development. Behaviour is defined in specs before it is built, and no code is written without an OpenSpec change behind it.

1. **Propose.** Every piece of work starts as a change in `openspec/changes/<change-name>/`, with a proposal, delta specs, a design where needed, and tasks. Use the `openspec-propose` skill.
2. **Review.** Cami reviews the proposal and specs before any code is written.
3. **Apply.** Work through `tasks.md` in order with the `openspec-apply-change` skill, ticking each task as it is done.
4. **Archive.** When every task is done and `make check` passes, archive the change with the `openspec-archive-change` skill. This merges its specs into `openspec/specs/`, which describe how Chorale behaves now.

A change that alters no behaviour, such as tooling or a refactor, sets `skip_specs: true` in its `.openspec.yaml`. It still needs a proposal and tasks.

Check the state of things with these commands:

```
openspec list
openspec list --specs
openspec validate --strict
```

## Test first, always

Every behaviour is test first, including the interface. Write the test, run it, and watch it fail for the reason you expect, before writing the code that makes it pass. The test layers, golden files and the scenario check are described in `docs/testing.md`.

Never skip, disable or loosen a test to get work through. Never regenerate a golden file just to make a test pass. When a golden file changes on purpose, say in the commit what changed on screen and why.

## Accessibility rules

These come from `docs/plan.md` and apply to every screen:

1. Works at 60 columns by 20 rows. Text wraps to the real window width.
2. Nothing is shown by colour alone. Every state has a word or a symbol.
3. Nothing animates. No spinners, shimmer or flicker.
4. The focused item has a marker and the real terminal cursor, so macOS Zoom can follow it.
5. Fixed-colour themes reach a contrast ratio of 7 to 1 for body text.

## Gates

A task is done only when these pass:

```
make check
```

Until milestone 0 lands, `make check` doesn't exist yet. Use `openspec validate --strict` for spec work.

## Writing

Cami often listens to documents through text-to-speech. Write full sentences in short paragraphs. Put code, paths and commands in backticks. Avoid tables, arrows and symbol chains in prose. Use absolute dates such as 2026-11-22, never "today" or "recently".

## Things never to do

1. Write code without an OpenSpec change behind it.
2. Write implementation code before its failing test.
3. Skip, disable or weaken a test.
4. Put API keys, tokens or passwords in the repo, a spec or a test fixture.
5. Commit Cami's voice recordings or transcripts.
