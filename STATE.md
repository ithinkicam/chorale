# Chorale: state

## Where things stand

Milestone 0 is built and archived, in pull request 2. That covers the Go project, test layers 1 to 4, the empty session list and the too-small screen, `scripts/check-scenarios`, `make check`, continuous integration on macOS and Linux, and the code style guide in `docs/style.md`. The behaviour it built is in `openspec/specs/app-shell/spec.md`.

The `harden-repository` change is merged, in pull request 4. It added CodeQL, `govulncheck` in `make check`, dependency review, OpenSSF Scorecard, Dependabot with auto-merge for patch and minor updates, the "Required checks" ruleset, and the community files. Every setting, and the reasons for it, is in `docs/repository.md`. CodeQL, Scorecard and `check` run on `main`, and Dependabot reads its configuration.

## Next steps

1. Cami decides on Dependabot's pull request 7, which moves the held-back `github.com/charmbracelet/x/exp/golden` from a pseudo-version to v0.1.0. Auto-merge correctly left it open.
2. Work through the follow-ups below.
3. Propose milestone 1: one Claude session, with the keyboard editor.

## Follow-ups from harden-repository

Remove each one from this list once it is done.

1. Check that the first Dependabot patch or minor pull request is approved from Cami's account, merges on its own once its checks pass, and that the merge starts `check.yml` on `main`. If the extra approval rule for unattributed changes holds it, bring that to Cami.
2. Dependabot doesn't update indirect Go modules, and `govulncheck` reports only code Chorale calls. When Scorecard's Vulnerabilities check names an advisory in an indirect module, raise that module by hand.

## Open decisions

The open decisions are listed in `docs/plan.md`: whether the M4 sleeps, which terminal Cami uses, the speech engine, and where Inkling 2 Mini runs.

## How to verify

```
make check
openspec validate --strict --all
```

`make check` needs the network, because `make vulncheck` downloads the Go vulnerability database.

## Context to know

1. Both of Cami's Macs have 16 GB of memory: her M4 Mac mini and plex-server, the M1.
2. The choir theme is for characters and visuals only. Features use ordinary developer names.
3. Cami often listens to documents through text-to-speech, so write in full sentences and short paragraphs.
4. GitHub's `viewerCanMergeAsAdmin` field ignores ruleset bypasses, so it can't show whether Cami can merge. The merge menu on the desktop site can.
5. The auto-merge token never expires, by Cami's choice. If it is revoked, Dependabot's pull requests wait until a new one is stored.
