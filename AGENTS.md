# Chorale: working rules for agents

Read this before doing any work in this repo. It applies to every agent: Claude Code, Codex and Antigravity CLI.

## What Chorale is

Chorale is a voice-first terminal workbench on macOS for running coding agents side by side, built for low vision. The roadmap is `docs/plan.md`. The test standard is `docs/testing.md`. The interface sketches are linked from the plan.

## All work goes through OpenSpec

Chorale uses OpenSpec for spec-driven development. Behaviour is defined in specs before it is built, and no code is written without an OpenSpec change behind it.

1. **Propose.** Every piece of work starts as a change in `openspec/changes/<change-name>/`, with a proposal, delta specs, a design where needed, and tasks. Use the `openspec-propose` skill.
2. **Review.** Before Cami reviews, two other reviews run. An outside agent from a different vendor reviews the change, for example `agy --mode plan --effort xhigh -p "<prompt>"` or `codex exec`. With `agy`, `-p` must come last, because it takes the next word as its prompt. Ask the reviewer to mark anything it couldn't verify as unconfirmed, and check its claims before passing them on. Headless `agy` can't ask for permission to read files, so put the change's files into the prompt. An adversarial reviewer from the same session reviews it too. Their findings go to Cami together, each with the proposing agent's view. Nothing a reviewer suggests goes into the change without her decision. Cami then reviews the proposal and specs, and no code is written before she approves them.
3. **Apply.** Work on a branch named after the change, and work through `tasks.md` in order with the `openspec-apply-change` skill, ticking each task as it is done. The work reaches `main` through a pull request.
4. **Archive.** When every task is done and `make check` passes, archive the change with the `openspec-archive-change` skill. This merges its specs into `openspec/specs/`, which describe how Chorale behaves now.

A change that alters no behaviour, such as tooling or a refactor, sets `skip_specs: true` in its `.openspec.yaml`. It still needs a proposal and tasks.

Check the state of things with these commands:

```
openspec list
openspec list --specs
openspec validate --strict --all
```

## Test first, always

Every behaviour is test first, including the interface. Write the test, run it, and watch it fail for the reason you expect, before writing the code that makes it pass. The test layers, golden files and the scenario check are described in `docs/testing.md`.

Never skip, disable or loosen a test to get work through. Never regenerate a golden file just to make a test pass. When a golden file changes on purpose, say in the commit what changed on screen and why.

## Code style

Follow `docs/style.md`. Code explains itself, and the reasons live in the specs, the designs and commit messages, not in comments. Doc comments are short, because they become published documentation.

## Accessibility rules

These come from `docs/plan.md` and apply to every screen:

1. Works at 60 columns by 20 rows. Text wraps to the real window width.
2. Nothing is shown by colour alone. Every state has a word or a symbol.
3. Nothing animates. No spinners, shimmer or flicker.
4. The focused item has a marker and the real terminal cursor, so macOS Zoom can follow it.
5. Fixed-colour themes reach a contrast ratio of 7 to 1 for body text.

## Plain names

The choir theme is for characters and visuals only. Keys, menus, status lines, docs and code use ordinary developer terms such as session, terminal, permission and settings. Never rename a standard thing to fit the theme.

## Pull requests

`main` is protected. Every change reaches it through a pull request that needs one approval and passing checks. Agents push as Cami's GitHub account, and GitHub doesn't let a pull request's author approve it. So for our own pull requests, the adversarial reviewer posts its review as a comment, and Cami merges using her admin bypass once she's satisfied. Bypassing skips every rule in the ruleset, so until the required checks sit in a ruleset of their own, Cami confirms both continuous integration jobs are green before merging. A contributor's pull request needs Cami's approval. Never force-push, and never push to `main` directly.

## Gates

A task is done only when these pass:

```
make check
```

For spec work on its own, `openspec validate --strict --all` is quicker. Always pass `--all` or an item name: without one, OpenSpec exits with an error when there is no terminal, and asks a question when there is one.

## Writing

Cami often listens to documents through text-to-speech. Write full sentences in short paragraphs. Put code, paths and commands in backticks. Avoid tables, arrows and symbol chains in prose. Never stamp a document with a date, such as "As of" lines or when something was decided. Git history holds that. When a fact about the outside world needs a date, such as a release date, write it as an absolute date, never "today" or "recently".

## Things never to do

1. Write code without an OpenSpec change behind it.
2. Write implementation code before its failing test.
3. Skip, disable or weaken a test.
4. Put API keys, tokens or passwords in the repo, a spec or a test fixture.
5. Commit Cami's voice recordings or transcripts.
