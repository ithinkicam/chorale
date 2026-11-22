# Chorale: state

As of 2026-11-22.

## Where things stand

The plan, the test standard, the agent rules and the OpenSpec setup are in place. No code exists yet. The repo is public on GitHub, and `main` takes changes only through pull requests.

The first change, `add-project-skeleton` for milestone 0, was reviewed on 2026-11-22 by Antigravity and by an adversarial reviewer, and revised to answer both. It passes `openspec validate --strict --all`, and it is waiting for Cami's review in pull request `revise-skeleton-review`.

## Next steps

1. Cami reviews the change and merges the review pull request.
2. An agent applies it on a branch, test first, and opens a pull request.
3. Once `make check` passes in continuous integration, the change is archived as the last commit on its branch, and Cami merges the pull request.
4. The `harden-repository` change is proposed next. It adds security scanning, Dependabot updates and auto-merge, required checks on `main` in a ruleset nobody can bypass, `CODEOWNERS` and `SECURITY.md`. Then milestone 1 is proposed.

## Open decisions

The open decisions are listed in `docs/plan.md`: whether the M4 sleeps, which terminal Cami uses, the speech engine, and where Inkling 2 Mini runs.

## How to verify

```
openspec list
openspec validate --strict --all
```

## Context to know

1. Both of Cami's Macs have 16 GB of memory: her M4 Mac mini and plex-server, the M1.
2. The choir theme is for characters and visuals only. Features use ordinary developer names.
3. Cami often listens to documents through text-to-speech, so write in full sentences and short paragraphs.
