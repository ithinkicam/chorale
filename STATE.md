# Chorale: state

## Where things stand

Milestone 0 is built on the branch `add-project-skeleton`, in pull request 2. That covers the Go project, test layers 1 to 4, the empty session list and the too-small screen, `scripts/check-scenarios`, `make check`, continuous integration on macOS and Linux, and the code style guide in `docs/style.md`. Every task in the change is done, and `make check` passes locally and in CI. The repo is public, and `main` takes changes only through pull requests.

## Next steps

1. The outside and adversarial reviews of pull request 2 go to Cami, with any fixes she agrees to.
2. The change is archived as the last commit on its branch, and Cami merges the pull request.
3. The `harden-repository` change is proposed next. It adds security scanning, Dependabot updates and auto-merge, required checks on `main` in a ruleset nobody can bypass, `CODEOWNERS` and `SECURITY.md`. Then milestone 1 is proposed.

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
