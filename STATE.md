# Chorale: state

As of 2026-11-22.

## Where things stand

The plan, the test standard, the agent rules and the OpenSpec setup are in place. No code exists yet. The first change, `add-project-skeleton` for milestone 0, is written and passes `openspec validate --strict`. It is waiting for Cami's review.

## Next steps

1. Cami clones the repo and sets up her Mac by following `docs/getting-started.md`.
2. Cami reviews `openspec/changes/add-project-skeleton/proposal.md` and its spec.
3. An agent applies the change with the `openspec-apply-change` skill, test first.
4. Once `make check` passes in continuous integration, the change is archived and milestone 1 is proposed as the next change.

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
