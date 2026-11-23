# Chorale: state

## Where things stand

Milestone 0 is built and archived, in pull request 2. That covers the Go project, test layers 1 to 4, the empty session list and the too-small screen, `scripts/check-scenarios`, `make check`, continuous integration on macOS and Linux, and the code style guide in `docs/style.md`. The behaviour it built is in `openspec/specs/app-shell/spec.md`.

The `harden-repository` change is built, archived, and in pull request 4, waiting for Cami to merge it. It added CodeQL, `govulncheck` in `make check`, dependency review, OpenSSF Scorecard, Dependabot with auto-merge for patch and minor updates, the "Required checks" ruleset, and the community files. Every setting, and the reasons for it, is in `docs/repository.md`.

## Next steps

1. Cami merges pull request 4 with her bypass. On her phone, the bypass shows only on the desktop site, in the menu beside the merge button, as "Bypass requirements".
2. Work through the follow-ups below.
3. Propose milestone 1: one Claude session, with the keyboard editor.

## Follow-ups after harden-repository merges

Remove each one from this list once it is done.

1. Check that CodeQL and Scorecard run on `main`, and that their results appear on the security page. Start Scorecard once by hand with `gh workflow run scorecard.yml` and check that it publishes.
2. Check that the Insights page for Dependabot shows no configuration error for `.github/dependabot.yml`.
3. Check that `gh api repos/ithinkicam/chorale/community/profile` lists `CODEOWNERS`, `SECURITY.md` and `CONTRIBUTING.md`.
4. Check that the first Dependabot patch or minor pull request is approved from Cami's account, merges on its own once its checks pass, and that the merge starts `check.yml` on `main`. If the extra approval rule for unattributed changes holds it, bring that to Cami.
5. When the first major update, or an update to one of the three held-back dependencies, arrives, check that it stays open, and bring Cami the summary `docs/repository.md` describes.

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
