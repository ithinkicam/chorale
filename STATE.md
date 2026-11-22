# Chorale: state

As of 2026-11-22.

## Where things stand

The plan, the test standard, the agent rules and the OpenSpec setup are in place. No code exists yet. The repo is public on GitHub, and `main` takes changes only through pull requests.

The first change, `add-project-skeleton` for milestone 0, was revised on 2026-11-22. It now has a scenario for short windows, a fuller continuous integration setup, and an ignore rule for archived changes. It passes `openspec validate --strict`, and it is waiting for the outside and adversarial reviews, then Cami's review.

## Next steps

1. An outside agent and the adversarial reviewer review `add-project-skeleton`, and their findings go to Cami.
2. Cami reviews the change.
3. An agent applies it on a branch, test first, and opens a pull request.
4. Once `make check` passes in continuous integration and the pull request is merged, the change is archived.
5. The `harden-repository` change is proposed next. It adds security scanning, Dependabot updates and auto-merge, required checks on `main`, `CODEOWNERS` and `SECURITY.md`. Then milestone 1 is proposed.

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
