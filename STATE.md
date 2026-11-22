# Chorale: state

## Where things stand

Milestone 0 is built and archived, in pull request 2. That covers the Go project, test layers 1 to 4, the empty session list and the too-small screen, `scripts/check-scenarios`, `make check`, continuous integration on macOS and Linux, and the code style guide in `docs/style.md`. The behaviour it built is in `openspec/specs/app-shell/spec.md`. The repo is public, and `main` takes changes only through pull requests.

## Next steps

1. Propose `harden-repository`, described below, with the `openspec-propose` skill. Run the outside review and the adversarial review before Cami reviews it, as `AGENTS.md` says.
2. Then propose milestone 1: one Claude session, with the keyboard editor.

## The next change: harden-repository

This is a tooling change, so its `.openspec.yaml` sets `skip_specs: true`. Cami asked for the repo to be ready to share. The free security tools for a public repo are wanted, and so is automatic merging of approved, passing pull requests. It needs:

1. CodeQL for the Go code and the workflow files.
2. `govulncheck` in `make check`, pinned as a Go tool.
3. The dependency review action on pull requests, blocking ones that add a known-vulnerable package.
4. OpenSSF Scorecard, run weekly.
5. Dependabot version updates, weekly, for Go modules, GitHub Actions and npm. A workflow using `dependabot/fetch-metadata` approves and auto-merges patch and minor updates once checks pass. Major updates wait for Cami. That workflow needs the repo setting that lets GitHub Actions approve pull requests.
6. A second ruleset on `main` that requires `check (macos)` and `check (linux)`, with no bypass actors. Bypassing a ruleset skips every rule in it, so the required checks can't sit in the ruleset Cami bypasses to merge.
7. `CODEOWNERS` naming Cami, `SECURITY.md` pointing to private vulnerability reporting, and a short `CONTRIBUTING.md` pointing to `AGENTS.md`, `docs/testing.md` and `docs/style.md`.

Already set on GitHub, so the change only needs to document them:

1. Secret scanning and push protection.
2. Dependabot alerts and security updates.
3. Private vulnerability reporting.
4. Auto-merge allowed, and branches deleted on merge.
5. The ruleset "Protect main". It blocks deletion and force-pushes, and requires a pull request with one approval and resolved threads. The admin role can bypass it, for pull requests only.

GitHub doesn't offer merge queues on repos owned by a personal account. Agents push as Cami's account, so they can't approve their own pull requests. Cami merges them with her bypass, after confirming both CI jobs are green.

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
