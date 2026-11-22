# Chorale: state

## Where things stand

Milestone 0 is built and archived, in pull request 2. That covers the Go project, test layers 1 to 4, the empty session list and the too-small screen, `scripts/check-scenarios`, `make check`, continuous integration on macOS and Linux, and the code style guide in `docs/style.md`. The behaviour it built is in `openspec/specs/app-shell/spec.md`. The repo is public, and `main` takes changes only through pull requests.

## Next steps

1. Cami merges pull request 2.
2. The `harden-repository` change is proposed next. It adds security scanning, Dependabot updates and auto-merge, required checks on `main` in a ruleset nobody can bypass, `CODEOWNERS` and `SECURITY.md`.
3. Then milestone 1 is proposed: one Claude session, with the keyboard editor.

## Open decisions

The open decisions are listed in `docs/plan.md`: whether the M4 sleeps, which terminal Cami uses, the speech engine, and where Inkling 2 Mini runs.

## How to verify

```
make check
```

## Context to know

1. Both of Cami's Macs have 16 GB of memory: her M4 Mac mini and plex-server, the M1.
2. The choir theme is for characters and visuals only. Features use ordinary developer names.
3. Cami often listens to documents through text-to-speech, so write in full sentences and short paragraphs.
